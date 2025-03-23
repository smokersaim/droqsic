package main

import (
	"log"

	"github.com/smokersaim/droqsic/cmd/server"
)

func main() {
	deps, err := server.InitializeApp()
	if err != nil {
		log.Fatal("Failed to initialize application: ", err)
	}

	server.InitServer(deps.App, deps.Cfg, deps.Log, deps.Mongo, deps.Redis)
}
