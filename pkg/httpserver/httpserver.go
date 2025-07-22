package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

const (
	defaultHost     = "127.0.0.1"
	defaultPort     = 50050
	readTimeoutSec  = 40
	writeTimeoutSec = 40
	idleTimeout     = 120
)

// IHTTPServer defines the interface for the http server.
//
//go:generate mockery --name=IHTTPServer --output=mocks --case=underscore
type IHTTPServer interface {
	Start() error
}

type HTTPServer struct {
	App *fiber.App

	host              string
	port              int
	enablePrefork     bool
	enablePrintRoutes bool
}

func (hs *HTTPServer) init() error {
	if hs.host == "" {
		hs.host = defaultHost
	}

	if hs.port == 0 {
		hs.port = defaultPort
	}

	hs.App = fiber.New(fiber.Config{
		ErrorHandler: hs.errorHandler,
		AppName:      "lingvogramm_backend",
		ReadTimeout:  readTimeoutSec * time.Second,
		WriteTimeout: writeTimeoutSec * time.Second,
		IdleTimeout:  idleTimeout * time.Second,
		ProxyHeader:  fiber.HeaderXForwardedFor,
	})

	return nil
}

func New(cfg config.HTTPServerConfig) (*HTTPServer, error) {
	hs := &HTTPServer{
		host:              cfg.Host,
		port:              cfg.Port,
		enablePrefork:     cfg.EnablePrefork,
		enablePrintRoutes: cfg.EnablePrintRoutes,
	}

	if err := hs.init(); err != nil {
		return nil, err
	}

	hs.App.Use(logger.New())
	hs.initCORS(cfg.Cors)
	hs.ping()

	return hs, nil
}

// initCORS initialize cors.
func (hs *HTTPServer) initCORS(cfg config.CorsConfig) {
	hs.App.Use(cors.New(cors.Config{
		AllowOrigins:        cfg.AllowOrigins,
		AllowMethods:        cfg.AllowMethods,
		AllowHeaders:        cfg.AllowHeaders,
		ExposeHeaders:       cfg.ExposeHeaders,
		MaxAge:              cfg.MaxAge,
		AllowCredentials:    cfg.AllowCredentials,
		AllowPrivateNetwork: cfg.AllowPrivateNetwork,
	}))
}

// Start http server.
func (hs *HTTPServer) Start() error {
	listenConfig := fiber.ListenConfig{
		OnShutdownError:   hs.onShutdownError,
		OnShutdownSuccess: hs.onShutdownSuccess,
		EnablePrefork:     hs.enablePrefork,
		EnablePrintRoutes: hs.enablePrintRoutes,
	}

	errChan := make(chan error, 1)
	go func() {
		if err := hs.App.Listen(fmt.Sprintf(":%d", hs.port), listenConfig); err != nil {
			errChan <- fmt.Errorf("listen: %w", err)
		}
	}()

	const waitTimeSec = 5

	select {
	case err := <-errChan:
		return err
	case <-time.After(waitTimeSec * time.Second): // wait for 5 seconds to ensure server starts.
		log.Println("server started successfully")
	}

	return hs.gracefulStop()
}

func (hs *HTTPServer) errorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	} else {
		err = fiber.NewError(code, err.Error())
	}

	log.Printf("error handling request: %v", err)

	return c.Status(code).JSON(err)
}

// ping register the /ping endpoint.
func (hs *HTTPServer) ping() {
	hs.App.Get("/ping", func(c fiber.Ctx) error {
		log.Println("ping endpoint called")

		return c.SendString("pong")
	})
}

// gracefulStop server with graceful shutdown.
func (hs *HTTPServer) gracefulStop() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server...")

	const ctxSec = 5
	ctx, cancel := context.WithTimeout(context.Background(), ctxSec*time.Second)
	defer cancel()

	if err := hs.App.ShutdownWithContext(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
		return err
	}

	log.Println("server exiting")

	return nil
}

func (hs *HTTPServer) onShutdownError(err error) {
	log.Printf("shutdown error: %v\n", err)
}

func (hs *HTTPServer) onShutdownSuccess() {
	log.Println("shutdown successful")
}
