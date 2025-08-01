package swaggerserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-jedi/lingramm_backend/config"
	_ "github.com/go-jedi/lingramm_backend/docs"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	defaultHost            = "127.0.0.1"
	defaultPort            = 40050
	defaultShutdownTimeout = 10
	readTimeoutSec         = 40
	writeTimeoutSec        = 40
	idleTimeout            = 120
)

// ISwaggerServer defines the interface for the swagger server.
//
//go:generate mockery --name=ISwaggerServer --output=mocks --case=underscore
type ISwaggerServer interface {
	Start() error
}

type SwaggerServer struct {
	shutdownTimeout int64
	port            int
	host            string
	allowedIPs      []string
	server          *http.Server
	mux             *http.ServeMux
}

func New(cfg config.SwaggerServerConfig, allowedIPs []string) (*SwaggerServer, error) {
	ss := &SwaggerServer{
		shutdownTimeout: cfg.ShutdownTimeout,
		port:            cfg.Port,
		host:            cfg.Host,
		allowedIPs:      allowedIPs,
		mux:             http.NewServeMux(),
	}

	if err := ss.init(); err != nil {
		return nil, err
	}

	ss.routes()

	ss.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", ss.host, ss.port),
		Handler:      ss.withLogging(ss.initCORS(cfg.Cors)),
		ReadTimeout:  readTimeoutSec * time.Second,
		WriteTimeout: writeTimeoutSec * time.Second,
		IdleTimeout:  idleTimeout * time.Second,
	}

	return ss, nil
}

func (ss *SwaggerServer) init() error {
	if ss.host == "" {
		ss.host = defaultHost
	}
	if ss.shutdownTimeout == 0 {
		ss.shutdownTimeout = defaultShutdownTimeout
	}
	if ss.port == 0 {
		ss.port = defaultPort
	}

	return nil
}

// initCORS initialize cors.
func (ss *SwaggerServer) initCORS(cfg config.CorsConfig) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:      cfg.AllowOrigins,
		AllowedMethods:      cfg.AllowMethods,
		AllowedHeaders:      cfg.AllowHeaders,
		ExposedHeaders:      cfg.ExposeHeaders,
		AllowCredentials:    cfg.AllowCredentials,
		AllowPrivateNetwork: cfg.AllowPrivateNetwork,
		MaxAge:              cfg.MaxAge,
	})
	return c.Handler(ss.mux)
}

func (ss *SwaggerServer) Start() error {
	errChan := make(chan error, 1)

	go func() {
		if err := ss.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("listen: %w", err)
		}
	}()

	const waitTimeSec = 5

	select {
	case err := <-errChan:
		return err
	case <-time.After(waitTimeSec * time.Second): // wait for 5 seconds to ensure server starts.
		log.Println("server swagger started successfully")
	}

	return ss.gracefulShutdown()
}

// ping register the /ping endpoint.
func (ss *SwaggerServer) routes() {
	ss.mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		log.Println("ping endpoint called")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("pong")); err != nil {
			log.Printf("error writing pong response: %v", err)
		}
	})

	ss.mux.Handle("/swagger/", ss.allowOnlyIPs(httpSwagger.WrapHandler))
}

func (ss *SwaggerServer) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("completed in %v", time.Since(start))
	})
}

// allowOnlyIPs allow only ips.
func (ss *SwaggerServer) allowOnlyIPs(next http.Handler) http.Handler {
	allowedMap := make(map[string]struct{})

	for i := range ss.allowedIPs {
		allowedMap[ss.allowedIPs[i]] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ss.getRequestIP(r)

		if _, ok := allowedMap[ip]; !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			log.Printf("unauthorized IP access: %s", ip)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getRequestIP get request ip.
func (ss *SwaggerServer) getRequestIP(r *http.Request) string {
	headerKey := "X-Forwarded-For"

	// check X-Forwarded-For header.
	xff := r.Header.Get(headerKey)
	if xff != "" {
		parts := strings.Split(xff, ",")

		return strings.TrimSpace(parts[0])
	}

	// fall back to remote addr.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

// gracefulShutdown server with graceful shutdown.
func (ss *SwaggerServer) gracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server swagger...")

	const ctxSec = 5
	ctx, cancel := context.WithTimeout(context.Background(), ctxSec*time.Second)
	defer cancel()

	if err := ss.server.Shutdown(ctx); err != nil {
		log.Printf("error during shutdown: %v", err)
		return err
	}

	log.Println("server swagger shutting down gracefully")

	return nil
}
