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
	userClient "tmossDev.github.com/eco-system/order-management/backend/app/order-gateway/client"
	"tmossDev.github.com/eco-system/order-management/backend/app/order-gateway/routes"
	orderService "tmossDev.github.com/eco-system/order-management/backend/domain/order/service"
	orderPostgres "tmossDev.github.com/eco-system/order-management/backend/domain/order/store/postgres"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/aws"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/local"
	configTypes "tmossDev.github.com/eco-system/shared-components/backend/package/config/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	storePostgres "tmossDev.github.com/eco-system/shared-components/backend/package/datastore/postgres"
	"tmossDev.github.com/eco-system/shared-components/backend/package/env"
	envConstants "tmossDev.github.com/eco-system/shared-components/backend/package/env/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	transportHTTP "tmossDev.github.com/eco-system/shared-components/backend/package/transport/http"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/middleware"
	httpTypes "tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/types"
	userConstants "tmossDev.github.com/eco-system/shared-components/backend/package/user/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/validator"
)

var (
	axxessLogs *accesslog.AccessLog
	config     *configTypes.ConfigModel
	irisApp    *iris.Application
	port       string
	sqlStore   *storePostgres.PostgresDataStore
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Infof(constants.DefaultRequestId, "Starting Application Version %s", envConstants.DefaultVersion)
	if err := setup(); err != nil {
		logger.Errorf(constants.DefaultRequestId, "unable to setup application: %s", err.Error())
		return
	}

	go start()
	<-ctx.Done()
	shutdown()
}

func setup() error {
	configManager := local.NewLocalConfigManager()
	if os.Getenv(envConstants.Env) != "" {
		configManager = aws.NewSecretConfigManager()
	}

	var err error
	config, err = configManager.GetConfig(os.Getenv(envConstants.SecretName))
	if err != nil {
		return fmt.Errorf("unable to load secret config: %s", err.Error())
	}

	sqlStore = storePostgres.NewPostgresDataStore(config.Database)
	if err := sqlStore.Connect(); err != nil {
		return fmt.Errorf("unable to connect to db: %s", err.Error())
	}

	irisApp = iris.New()
	axxessLogs = middleware.MakeAccessLog()
	jwtMiddleware := transportHTTP.NewJWTMiddleware(httpTypes.JWTConfig{
		SecretKey:     []byte(userConstants.PASSWORD_SECRET_HASHING_KEY),
		TokenExpiry:   72 * time.Hour,
		SigningMethod: jwt.SigningMethodHS256,
		TokenPrefix:   "Bearer ",
	})

	irisApp.Use(
		corsMiddleware,
		axxessLogs.Handler,
		middleware.CaselessMatcherMiddleware,
		middleware.RequestIDMiddleware,
		jwtMiddleware([]string{"/auth/login", "/login", "/auth/logout", "/logout", "/health"}),
	)
	irisApp.Options("/{path:path}", corsMiddleware)

	orderRepo := orderPostgres.NewPostgresOrderRepository(sqlStore)
	orderSvc := orderService.NewOrderService(validator.NewValidator(), orderRepo)
	authClient := userClient.NewHTTPAuthClient(env.Getenv("USER_SERVICE_URL", "http://user-service:8080"))
	routes.Setup(irisApp, orderSvc, authClient)

	port = env.Getenv(envConstants.Port, envConstants.DefaultPort)
	return nil
}

func start() {
	if err := irisApp.Listen(fmt.Sprintf(":%s", port)); err != nil {
		logger.Errorf("failed to start server reason: %s", err.Error())
	}
}

func shutdown() {
	if err := sqlStore.Close(); err != nil {
		logger.Errorf(constants.DefaultRequestId, "Unable to close DB: %s", err.Error())
	}
	if err := axxessLogs.Close(); err != nil {
		logger.Errorf(constants.DefaultRequestId, "Unable to stop HTTP Server gracefully: %s", err.Error())
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
