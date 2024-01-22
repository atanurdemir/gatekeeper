package main

import (
	"github.com/atanurdemir/gatekeeper/src/config"
	"github.com/atanurdemir/gatekeeper/src/server"
	"github.com/atanurdemir/gatekeeper/src/store"
)

func main() {

	// Setup Config
	config.SetupConfig()

	// Store Cleanup
	go store.CleanupExpiredBans()

	// HTTP Server
	server.NewServer().Start("3000")
}
