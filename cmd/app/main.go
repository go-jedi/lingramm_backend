package main

import (
	"context"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/app"
)

// @title Lingramm API — Telegram Web App Backend
// @version 1.0
// @description This is the backend API for the Lingramm Telegram Web Application.
// @description It provides endpoints for user interactions, game logic, statistics, tools, and more.
// @description All endpoints are secured and optimized for real-time communication with Telegram Mini Apps.

// @host localhost:50050
// @BasePath /v1.
func main() {
	ctx := context.Background()

	// initialize app.
	a, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	// run application.
	if err := a.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
