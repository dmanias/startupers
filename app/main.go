package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmanias/startupers/app/conf"
	"github.com/dmanias/startupers/app/services/api/handlers"
	database "github.com/dmanias/startupers/business/sys/database/pgx"
	"github.com/dmanias/startupers/business/web/auth"
	"github.com/dmanias/startupers/foundation/keystore"
	"github.com/dmanias/startupers/foundation/logger"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var build = "develop"

type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type HTTPServer struct {
	*http.Server
}

func (s *HTTPServer) Start() error {
	return s.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func main() {
	log, err := logger.New("STARTUPERS-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func(log *zap.SugaredLogger) {
		err := log.Sync()
		if err != nil {

		}
	}(log)

	if err := run(log); err != nil {
		log.Errorw("start", "ERROR", err)
		err := log.Sync()
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// GOMAXPROCS

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "BUILD-", build)

	// -------------------------------------------------------------------------

	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout        time.Duration `conf:"default:90s"`
			WriteTimeout       time.Duration `conf:"default:60s"`
			IdleTimeout        time.Duration `conf:"default:150s"`
			ShutdownTimeout    time.Duration `conf:"default:60s"`
			CORSAllowedOrigins []string      `conf:"default:*"`
			//CORSAllowedOrigins []string `conf:"default:http://localhost:3000"`
		}
		AI struct {
			APIKey string `conf:"env:AI_API_KEY"`
			APIURL string `conf:"default:https://api.openai.com/v1/chat/completions"`
		}
		DB struct {
			User         string `conf:"env:DATABASE_USERNAME"`
			Password     string `conf:"env:DATABASE_PASSWORD"`
			Host         string `conf:"default:startupers-postgresql-primary"`
			Name         string `conf:"env:DATABASE_NAME"`
			MaxIdleConns int    `conf:"default:2"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
		Auth struct {
			KeysFolder string `conf:"default:zarf/keys/"`
			ActiveKID  string `conf:"env:ACTIVE_KID"`
			Issuer     string `conf:"default:BackEnd"`
		}
		Build struct {
			Build string `conf:"default:0.3"`
			Desc  string `conf:"default:copyright information here"`
		}
	}{}

	// Directly access environment variables for database configuration
	dbUser := os.Getenv("DATABASE_USERNAME")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	authActiveKID := os.Getenv("AUTH_ACTIVEKID")
	aiApiKey := os.Getenv("AI_API_KEY")

	cfg.DB.User = dbUser
	cfg.DB.Password = dbPassword
	cfg.DB.Name = dbName
	cfg.Auth.ActiveKID = authActiveKID
	cfg.AI.APIKey = aiApiKey

	const prefix = "STARTUP"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------

	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// -------------------------------------------------------------------------

	// Database Support

	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	db, err := database.Open(database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		err := db.Close()
		if err != nil {
			return
		}
	}()

	// -------------------------------------------------------------------------
	// Initialize authentication support

	log.Infow("startup", "status", "initializing authentication support")

	// Simple keystore versus using Vault.
	ks, err := keystore.NewFS(os.DirFS(cfg.Auth.KeysFolder))
	if err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{
		Log:       log,
		KeyLookup: ks,
		Issuer:    cfg.Auth.Issuer,
	}

	authConf, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// -------------------------------------------------------------------------
	// Serve static files from the "uploads" directory
	fs := http.FileServer(http.Dir("./uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))
	// Start API Service

	log.Infow("startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	////////
	authConfig := auth.Config{
		Log:       log,
		KeyLookup: ks,
		Issuer:    cfg.Auth.Issuer,
	}

	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown:   shutdown,
		Log:        log,
		Auth:       authConf,
		AuthConfig: &authConfig,
		DB:         db,
		APIKey:     cfg.AI.APIKey,
		Build:      cfg.Build.Build,
		ActiveKID:  cfg.Auth.ActiveKID,
		APIHost:    cfg.Web.APIHost,
	})

	corsOptions := cors.Options{
		AllowedOrigins:   cfg.Web.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}
	corsHandler := cors.New(corsOptions)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      corsHandler.Handler(apiMux),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
