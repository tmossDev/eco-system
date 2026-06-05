package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"tmossDev.github.com/eco-system/online-storefront/backend/app/cart-gateway/routes"
	cartDomainService "tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/service"
	cartPostgres "tmossDev.github.com/eco-system/online-storefront/backend/domain/cart/store/postgres"
	orderDomainService "tmossDev.github.com/eco-system/online-storefront/backend/domain/order/service"
	orderPostgres "tmossDev.github.com/eco-system/online-storefront/backend/domain/order/store/postgres"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/aws"
	"tmossDev.github.com/eco-system/shared-components/backend/package/config/local"
	configTypes "tmossDev.github.com/eco-system/shared-components/backend/package/config/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	storePostgres "tmossDev.github.com/eco-system/shared-components/backend/package/datastore/postgres"
	"tmossDev.github.com/eco-system/shared-components/backend/package/env"
	envConstants "tmossDev.github.com/eco-system/shared-components/backend/package/env/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/logger"
	"tmossDev.github.com/eco-system/shared-components/backend/package/transport/http/middleware"
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
	irisApp.Use(axxessLogs.Handler, middleware.CaselessMatcherMiddleware, middleware.RequestIDMiddleware)

	cartRepo := cartPostgres.NewPostgresCartRepository(sqlStore)
	orderRepo := orderPostgres.NewPostgresOrderRepository(sqlStore)
	cartSvc := cartDomainService.NewCartService(validator.NewValidator(), cartRepo)
	orderSvc := orderDomainService.NewOrderService(orderRepo)
	routes.Setup(irisApp, cartSvc, orderSvc)

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
