package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"watcher/api"
	repository "watcher/db"
	"watcher/services"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Println(time.Now().UTC(), "Starting watcher.")

	godotenv.Load()

	timeout, err := strconv.Atoi(os.Getenv("TIME_OUT"))
	if err == nil {
		time.Sleep(time.Second * time.Duration(timeout))
	}

	// Connect to local SQLite database file
	sqlyte := &repository.Sqlyte{}
	httpHandler := &api.HttpHandler{Repo: sqlyte}

	err = sqlyte.Connect("./watcher.db")
	if err != nil {
		log.Fatal(err)
	}
	sqlyte.Initialize()

	// Configure the SQLite custom logger options
	opts := slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	sqliteHandler := repository.NewSQLiteHandler(sqlyte.Db, opts)
	logger := slog.New(sqliteHandler)

	slog.SetDefault(logger)

	// Create context that listens for interrupt signals from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start long-running worker in a goroutine
	go httpHandler.StartServer()
	go backgroundWork(ctx, sqlyte, logger)

	// Block main "thread" until a stop signal is received
	<-ctx.Done()
	fmt.Println(time.Now().UTC(), "Shutdown signal received. Cleaning up..")

	// Perform cleanup (closing database, flushing logs)
	httpHandler.ShutdownServer()
	sqlyte.Close()

	fmt.Println(time.Now().UTC(), "Application gracefully stopped.")
}

func backgroundWork(ctx context.Context, repo *repository.Sqlyte, log *slog.Logger) {

	fmt.Println(time.Now().UTC(), "Background worker started..")

	url, k8Cmd, k8Params, timeout, files := services.GetEnvValues()

	for {
		select {
		case <-ctx.Done():
			fmt.Println(time.Now().UTC(), "Background shutdown signal received.")
			return
		default:
			for _, item := range files {
				services.VerifyFileNotChanged(repo, item.Type, item.Path, item.Content, item.Change)
			}
			services.CheckWebSiteStatus(repo, log, url)
			services.CheckK8sStatus(repo, log, k8Cmd, k8Params...)
			time.Sleep(time.Minute * time.Duration(timeout))
		}
	}
}
