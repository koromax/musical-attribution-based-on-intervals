package main

import (
	"interval_attribution/internal/api"
)

// @title Composers API
// @version 1.0
// @description Backend API for Composers SPA (LR4)

// @host 127.0.0.1:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	api.StartServer()
}
