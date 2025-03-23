package server

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/smokersaim/droqsic/cmd/config"
	"github.com/smokersaim/droqsic/internal/infrastructure/cache"
	"github.com/smokersaim/droqsic/internal/infrastructure/database"
	"go.uber.org/zap"
)

func InitServer(app *fiber.App, cfg *config.Config, log *zap.Logger, mongo *database.MongoDB, redis *cache.RedisCache) {
	maxProcs := cfg.Server.Instance
	cpuCount := runtime.NumCPU()
	if maxProcs > cpuCount {
		maxProcs = cpuCount
	}
	runtime.GOMAXPROCS(maxProcs)

	log.Info("Server is running", zap.Int("port", cfg.Server.Port), zap.Int("cpu_count", cpuCount), zap.Int("max_procs", maxProcs))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			log.Fatal("Startup failure", zap.Error(err))
		}
	}()

	<-quit
	log.Info("Shutdown signal received. Initiating graceful termination.")
	mongo.Disconnect(log)
	redis.Disconnect(log)

	if err := app.Shutdown(); err != nil {
		log.Error("Shutdown failure", zap.Error(err))
	} else {
		log.Info("Server shut down successfull.")
	}

}
