package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dmanias/startupers/app/ai"
	"github.com/dmanias/startupers/app/config"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ConfigPath = "./app/config/config.json"

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

type Question struct {
	Query string `json:"query"`
}

func handleAsk(apikey string, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	var question Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	aiResponse, err := ai.AskAI(apikey, question.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(aiResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(response)
	if err != nil {
		return
	}
}

func main() {

	cfg, err := config.LoadConfig(ConfigPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("API Host from configuration:", cfg.Web.APIHost)

	fmt.Println("Loaded configuration:", cfg)

	log, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}
	defer func() {
		if err := log.Sync(); err != nil {
			// Handle the error. Options include logging to stderr, ignoring, or panicking.
			_, _ = fmt.Fprintf(os.Stderr, "Failed to sync log: %v\n", err)
		}
	}()

	// Initialize the HTTP server with the configuration
	httpServer := &HTTPServer{
		Server: &http.Server{
			Addr:         cfg.Web.APIHost,
			ReadTimeout:  cfg.Web.ReadTimeout * time.Millisecond,
			WriteTimeout: cfg.Web.WriteTimeout * time.Millisecond,
			IdleTimeout:  cfg.Web.IdleTimeout * time.Millisecond,
		},
	}

	// Pass the correctly initialized HTTPServer implementing the Server interface to run
	if err := run(log, httpServer, cfg); err != nil {
		log.Sugar().Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(log *zap.Logger, server Server, cfg *config.Config) error {

	log.Sugar().Infow("starting application", "version", build)

	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		handleAsk(cfg.AI.APIKey, w, r)
	})

	go func() {
		log.Sugar().Infow("Server is running", "port", cfg.Web.APIHost)
		if err := server.Start(); err != http.ErrServerClosed {
			log.Sugar().Fatalw("Failed to start server", "error", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	log.Sugar().Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Sugar().Errorw("Failed to gracefully shutdown the server", "error", err)
		return err
	}

	log.Sugar().Infow("Server shutdown gracefully")
	return nil
}
