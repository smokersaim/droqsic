//go:build wireinject
// +build wireinject

package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/smokersaim/droqsic/cmd/config"
	"github.com/smokersaim/droqsic/internal/infrastructure/http/middleware"
	"github.com/smokersaim/droqsic/internal/infrastructure/logger"
	"go.uber.org/zap"
)

type AppDependencies struct {
	App *fiber.App
	Cfg *config.Config
	Log *zap.Logger
}

func ProvideLogger() *zap.Logger {
	logger.InitLogger()
	return logger.GetLogger()
}

func ProvideConfig() *config.Config {
	log := logger.GetLogger()
	return config.LoadConfigs(log)
}

func ProvideFiberApp(cfg *config.Config) *fiber.App {
	app := NewFiberApp(cfg)
	middleware.Cors(app, cfg)
	return app
}

var AppSet = wire.NewSet(
	ProvideLogger,
	ProvideConfig,
	ProvideFiberApp,
	wire.Struct(new(AppDependencies), "App", "Cfg", "Log"),
)

func InitializeApp() (*AppDependencies, error) {
	wire.Build(AppSet)
	return nil, nil
}
