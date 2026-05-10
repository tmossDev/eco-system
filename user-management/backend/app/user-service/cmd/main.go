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
	"tmossDev.github.com/eco-system/shared-components/backend/package/datastore/postgres"
	"tmossDev.github.com/eco-system/shared-components/backend/package/env"
	envConstants "tmossDev.github.com/eco-system/shared-components/backend/package/env/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/middleware"
	httpTypes "tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/types"
)

var sqlStoreConn *postgres.PostgresDataStore

var config *types.ConfigModel

var axxessLogs *accesslog.AccessLog

var irisServer *iris.Application

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

	go start()

	<-ctx.Done()
	logger.Info(constants.DefaultRequestId, "shutdown signalled...")
	shutdown()
	logger.Info(constants.DefaultRequestId, "shutdown complete->")
}

func shutdown() {
	//Close Sql Conn to DB
	writerCloseError, readerCloseError := sqlStoreConn.Close()
	if writerCloseError != nil {
		logger.Errorf(constants.CTXRequestIdKey, "Unable to close DB writer: %s", writerCloseError.Error())
	}
	if readerCloseError != nil {
		logger.Errorf(constants.CTXRequestIdKey, "Unable to close DB reader: %s", writerCloseError.Error())
	}

	//Close Accesslogs to Server
	err := axxessLogs.Close()

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

	sqlStoreConn = postgres.NewPostgresDataStore(config.Database)
	cErr := sqlStoreConn.Connect()
	if cErr != nil {
		return fmt.Errorf("unable to connect to db: %s, exiting", cErr.Error())
	}

	irisServer = iris.New()
	axxessLogs = middleware.MakeAccessLog()

	config := httpTypes.JWTConfig{
		SecretKey:     []byte("super_duper_secret_key"), // Should be loaded from environment variables
		TokenExpiry:   72 * time.Hour,
		SigningMethod: jwt.SigningMethodHS256,
		TokenPrefix:   "Bearer ",
	}
	jwtFunction := http.NewJWTMiddleware(config)

	irisServer.Use(
		axxessLogs.Handler,
		middleware.CaselessMatcherMiddleware,
		middleware.RequestIDMiddleware,
		jwtFunction([]string{"/login", "/logout", "/refresh", "/health"}),
	)

	port = env.Getenv(envConstants.Port, envConstants.DefaultPort)

	return nil
}

func start() {
	if err := irisServer.Listen(fmt.Sprintf(":%s", port)); err != nil {
		logger.Errorf("failed to start server reason: %s", err.Error())
	}
}
