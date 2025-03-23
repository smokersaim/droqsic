//go:build wireinject
// +build wireinject

package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/smokersaim/droqsic/cmd/config"
	"github.com/smokersaim/droqsic/internal/infrastructure/cache"
	"github.com/smokersaim/droqsic/internal/infrastructure/database"
	"github.com/smokersaim/droqsic/internal/infrastructure/http/middleware"
	"github.com/smokersaim/droqsic/internal/infrastructure/logger"
	"go.uber.org/zap"
)

type AppDependencies struct {
	App   *fiber.App
	Cfg   *config.Config
	Log   *zap.Logger
	Mongo *database.MongoDB
	Redis *cache.RedisCache
}

func ProvideLogger() *zap.Logger {
	logger.InitLogger()
	return logger.GetLogger()
}

func ProvideConfig() *config.Config {
	log := logger.GetLogger()
	return config.LoadConfigs(log)
}

func ProvideMongoDB(cfg *config.Config, log *zap.Logger) (*database.MongoDB, error) {
	return database.ConnectMongoDB(cfg, log)
}

func ProvideRedisCache(cfg *config.Config, log *zap.Logger) (*cache.RedisCache, error) {
	return cache.ConnectRedisCache(cfg, log)
}

func ProvideFiberApp(cfg *config.Config) *fiber.App {
	app := NewFiberApp(cfg)
	middleware.Cors(app, cfg)
	return app
}

var AppSet = wire.NewSet(
	ProvideLogger,
	ProvideConfig,
	ProvideMongoDB,
	ProvideRedisCache,
	ProvideFiberApp,
	wire.Struct(new(AppDependencies), "App", "Cfg", "Log", "Mongo", "Redis"),
)

func InitializeApp() (*AppDependencies, error) {
	wire.Build(AppSet)
	return nil, nil
}
