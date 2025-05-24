package main

import (
	"context"
	"log"

	"github.com/go-jedi/lingramm_backend/internal/app"
)

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
