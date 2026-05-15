package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/aws"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/local"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	storePostgres "tmossDev.github.com/eco-system/shared-components/backend/package/datastore/postgres"
	"tmossDev.github.com/eco-system/shared-components/backend/package/env"
	envConstants "tmossDev.github.com/eco-system/shared-components/backend/package/env/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/middleware"
	httpTypes "tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
	"tmossDev.github.com/eco-system/user-management/backend/app/user-gateway/routes"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/repository/postgres"
	"tmossDev.github.com/eco-system/user-management/backend/package/user/service"
)

var sqlStore *storePostgres.PostgresDataStore

var config *types.ConfigModel

var axxessLogs *accesslog.AccessLog

var irisApp *iris.Application

var port string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Infof(constants.DefaultRequestId, "Starting Application Version %s", envConstants.DefaultVersion)

	err := setup()
	if err != nil {
		logger.Errorf(constants.DefaultRequestId, "unable to Setup Application: %s", err.Error())
		return
	}

	logger.Info(constants.DefaultRequestId, "Set up complete")

	go start()

	<-ctx.Done()
	logger.Info(constants.DefaultRequestId, "shutdown signalled...")
	shutdown()
	logger.Info(constants.DefaultRequestId, "shutdown complete->")
}

func shutdown() {
	//Close Sql Conn to DB
	err := sqlStore.Close()
	if err != nil {
		logger.Errorf(constants.CTXRequestIdKey, "Unable to close reader or writer DB: %s", err.Error())
	}

	//Close Accesslogs to Server
	err = axxessLogs.Close()
	if err != nil {
		logger.Errorf(constants.CTXRequestIdKey, "Unable to stop HTTP Server gracefully: %s", err.Error())
	}
}
func setup() error {
	var err error
	configManager := local.NewLocalConfigManager()
	logger.Infof(constants.DefaultRequestId, "Env: %s", env.Getenv(os.Getenv(envConstants.Env), "[not specified]"))
	if os.Getenv(envConstants.Env) != "" {
		configManager = aws.NewSecretConfigManager()
	}

	logger.Infof(constants.DefaultRequestId, "Secret Username: %s", env.Getenv(os.Getenv(envConstants.SecretName), "[not specified]"))
	config, err = configManager.GetConfig(os.Getenv("SECRET_NAME"))
	if err != nil {
		return fmt.Errorf("unable to load secret config: %s, exiting", err.Error())
	}

	sqlStore = storePostgres.NewPostgresDataStore(config.Database)
	cErr := sqlStore.Connect()
	if cErr != nil {
		return fmt.Errorf("unable to connect to db: %s, exiting", cErr.Error())
	}

	irisApp = iris.New()
	axxessLogs = middleware.MakeAccessLog()

	config := httpTypes.JWTConfig{
		SecretKey:     []byte("super_duper_secret_key"), // Should be loaded from environment variables
		TokenExpiry:   72 * time.Hour,
		SigningMethod: jwt.SigningMethodHS256,
		TokenPrefix:   "Bearer ",
	}
	jwtFunction := http.NewJWTMiddleware(config)

	irisApp.Use(
		corsMiddleware,
		axxessLogs.Handler,
		middleware.CaselessMatcherMiddleware,
		middleware.RequestIDMiddleware,
		jwtFunction([]string{"/login", "/logout", "/refresh", "/health"}),
	)
	irisApp.Options("/{path:path}", corsMiddleware)

	port = env.Getenv(envConstants.Port, envConstants.DefaultPort)

	validater := validator.NewValidator()
	logger.Info(constants.CTXRequestIdKey, "Set up validator...")

	userRepo := postgres.NewPostgresUserRepository(sqlStore)
	logger.Info(constants.CTXRequestIdKey, "Set up repositories...")

	publicUserService := service.NewPublicService(validater, userRepo)
	privateUserService := service.NewPrivateService(validater, userRepo)
	logger.Info(constants.CTXRequestIdKey, "Set up services...")

	routes.Setup(irisApp, publicUserService, privateUserService)
	logger.Info(constants.CTXRequestIdKey, "Set up routes...")

	return nil
}

func start() {
	if err := irisApp.Listen(fmt.Sprintf(":%s", port)); err != nil {
		logger.Errorf("failed to start server reason: %s", err.Error())
	}
}

func corsMiddleware(ctx iris.Context) {
	origin := os.Getenv("FRONTEND_ORIGIN")
	if origin == "" {
		origin = ctx.GetHeader("Origin")
	}

	if origin != "" {
		ctx.Header("Access-Control-Allow-Origin", origin)
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Vary", "Origin")
	}

	if ctx.Method() == iris.MethodOptions {
		ctx.StatusCode(iris.StatusNoContent)
		return
	}

	ctx.Next()
}
