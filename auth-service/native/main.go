package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"native/handler"
	"native/router"
	"src/configuration"
	"src/postgresql"
	"src/service"
)

func main() {

	fmt.Println(time.Now().UTC(), "Starting watcher.")

	settings := &configuration.Settings{}
	settings.GetInstance()
	postgresHandler := &postgresql.DBHandler{Settings: settings}
	userSvc := &service.UserService{DB: postgresHandler}
	authRouter := &router.Router{UserSvc: userSvc}
	adminRouter := &router.AdminRouter{UserSvc: userSvc}

	httpHandler := &handler.HttpHandler{
		Settings:    settings,
		Router:      authRouter,
		AdminRouter: adminRouter,
	}

	// Create context that listens for interrupt signals from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start long-running worker in a goroutine
	go httpHandler.StartServer()

	// Block main execution until a signal is received
	<-ctx.Done()
	fmt.Println(time.Now().UTC(), "Shutdown signal received..")

	// Perform cleanup (e.g., closing database, flushing logs)
	httpHandler.ShutdownServer()

	fmt.Println(time.Now().UTC(), "Application gracefully stopped.")

}
