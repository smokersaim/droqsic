package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/smokersaim/droqsic/cmd/config"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func NewFiberApp(cfg *config.Config) *fiber.App {
	return fiber.New(fiber.Config{
		Prefork:      cfg.App.Prefork,
		AppName:      cfg.App.Name,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		BodyLimit:    5 * 1024 * 1024,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
}
