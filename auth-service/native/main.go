package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"native/handlers"
	"native/routers"
	database "src/db"
	"src/services"
	"src/utils"
)

func main() {

	fmt.Println(time.Now().UTC(), "Starting watcher.")

	settings := &utils.Settings{}
	settings.GetInstance()
	postgresHandler := &database.DBHandler{Settings: settings}
	userSvc := &services.UserService{Postgres: postgresHandler}
	authRouter := &routers.Router{UserSvc: userSvc}
	adminRouter := &routers.AdminRouter{UserSvc: userSvc}

	httpHandler := &handlers.HttpHandler{
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
