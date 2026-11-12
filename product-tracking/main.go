package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"product-tracking/api"
	"product-tracking/config"
	"product-tracking/kafka"
	"product-tracking/postgresql"
	"product-tracking/repository"
	"product-tracking/router"
	"product-tracking/service"
	"product-tracking/webscraping"
)

func main() {

	fmt.Println(time.Now().UTC(), "Starting product tracking.")

	// Initialize log handler pointing to standard output
	var logger *slog.Logger
	if os.Getenv("ENV") == "production" {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(handler)
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)

	// Create context that listens for interrupt signals from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	settings := &config.Settings{}
	settings.GetInstance()
	postgresHandler := &postgresql.DBHandler{Settings: settings}
	postgresHandler.Initialize()
	productRepo := repository.NewProductRepository(postgresHandler.DB)
	productSvc := &service.ProductService{Repo: productRepo, Ctx: ctx}
	router := &router.Router{ProductSvc: productSvc}
	kafkaSvc := &kafka.KafkaSvc{Settings: settings}
	devTools := &webscraping.DevTools{ProductSvc: productSvc, KafkaSvc: kafkaSvc}

	httpHandler := &api.HttpHandler{
		Settings: settings,
		Router:   router,
	}

	// Start long-running worker in a goroutine
	go httpHandler.StartServer()
	go backgroundWork(ctx, devTools, logger)

	// Block main execution until a signal is received
	<-ctx.Done()
	slog.Info("Shutdown signal received..")

	// Perform cleanup (e.g., closing database, flushing logs)
	httpHandler.ShutdownServer()
	postgresHandler.Close()

	slog.Info("Application gracefully stopped.")
}

func backgroundWork(ctx context.Context, devTools *webscraping.DevTools, log *slog.Logger) {

	log.Info("Background worker started..")
	defer log.Info("Background worker finished.")

	for {
		select {
		case <-ctx.Done():
			log.Info("Background shutdown signal received.")
			return
		default:
			log.Info("Get product information starting..")
			devTools.Execute()
			log.Info("Get product information ended.")
			timeOff := rand.Intn(28) + 6
			log.Info("Next execution in " + strconv.Itoa(timeOff) + " hours.")
			time.Sleep(time.Hour * time.Duration(timeOff))
		}
	}
}
