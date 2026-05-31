package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/storefront-gateway/client"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/storefront-gateway/routes"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/env"
	envConstants "tmossDev.github.com/eco-system/shared-components/backend/package/env/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/middleware"
)

var (
	axxessLogs    *accesslog.AccessLog
	irisApp       *iris.Application
	productClient client.ProductClient
	userClient    client.UserClient
	port          string
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Infof(constants.DefaultRequestId, "Starting Application Version %s", envConstants.DefaultVersion)
	setup()

	go start()

	<-ctx.Done()
	logger.Info(constants.DefaultRequestId, "shutdown signalled...")
	shutdown()
	logger.Info(constants.DefaultRequestId, "shutdown complete->")
}

func setup() {
	irisApp = iris.New()
	axxessLogs = middleware.MakeAccessLog()
	irisApp.Use(
		corsMiddleware,
		axxessLogs.Handler,
		middleware.CaselessMatcherMiddleware,
		middleware.RequestIDMiddleware,
	)
	irisApp.Options("/{path:path}", corsMiddleware)

	port = env.Getenv(envConstants.Port, envConstants.DefaultPort)
	userClient = client.NewHTTPUserClient(env.Getenv("USER_SERVICE_URL", "http://user-service:8080"))
	productClient = client.NewHTTPProductClient(
		env.Getenv("PRODUCT_SERVICE_URL", "http://product-service:8080"),
		env.Getenv("PRODUCT_GATEWAY_URL", "http://product-gateway:8080"),
	)

	routes.Setup(irisApp, userClient, productClient)
	logger.Info(constants.DefaultRequestId, "Set up complete")
}

func shutdown() {
	userClient.Shutdown()
	productClient.Shutdown()
	if err := axxessLogs.Close(); err != nil {
		logger.Errorf(constants.DefaultRequestId, "Unable to stop HTTP Server gracefully: %s", err.Error())
	}
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
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		ctx.Header("Vary", "Origin")
	}

	if ctx.Method() == iris.MethodOptions {
		ctx.StatusCode(iris.StatusNoContent)
		return
	}

	ctx.Next()
}
