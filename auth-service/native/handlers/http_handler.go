package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"native/middleware"
	"native/routers"
	"native/routers/admin"
	"src/utils"
)

type HttpHandler struct {
	Server   *http.Server
	Settings *utils.Settings
}

func (h *HttpHandler) StartServer() {

	mux := http.NewServeMux()

	//mux.HandleFunc("GET /", routers.GetUser)

	{
		mux.HandleFunc("POST /auth/sign-in", routers.SignIn)
		mux.HandleFunc("POST /auth/sign-up", routers.SignUp)
		mux.HandleFunc("POST /auth/validate-token", routers.VerifyToken)
	}
	{
		mux.HandleFunc("GET /user", routers.GetUser)
		mux.HandleFunc("PUT /user", routers.UpdateUser)
		mux.HandleFunc("POST /user/change-password", routers.ChangePassword)
		mux.HandleFunc("GET /user/img", routers.GetPicture)
		mux.HandleFunc("POST /user/img", routers.AddPicture)
		mux.HandleFunc("GET /user/address", routers.GetAddresses)
		mux.HandleFunc("POST /user/address", routers.AddAddress)
		mux.HandleFunc("DELETE /user/address", routers.DeleteAddress)
		mux.HandleFunc("PUT /user/address", routers.UpdateAddress)
		mux.HandleFunc("DELETE /user", routers.DeleteUser)
	}
	{
		mux.HandleFunc("GET /admin/user", admin.GetUser)
		mux.HandleFunc("GET /admin/users", admin.GetUsers)
		mux.HandleFunc("DELETE /admin/user", admin.DeleteUser)
		mux.HandleFunc("GET /admin/user/address", admin.GetAddresses)
	}

	JWTMiddleware := middleware.JWTMiddleware(mux)

	h.Server = &http.Server{
		Addr:           ":" + h.Settings.SERVER_PORT,
		Handler:        JWTMiddleware,
		ReadTimeout:    time.Second * 60,
		WriteTimeout:   time.Second * 60,
		IdleTimeout:    time.Second * 120,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println(time.Now().UTC(), "Starting server on port.", h.Server.Addr)

	if err := h.Server.ListenAndServe(); err != http.ErrServerClosed {
		log.Panicf("%s Server error: %v", time.Now().UTC(), err)
	}

}

func (h *HttpHandler) ShutdownServer() {
	fmt.Println(time.Now().UTC(), "Stopping server..")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	if err := h.Server.Shutdown(ctx); err != nil {
		fmt.Println(time.Now().UTC(), "Server graceful shutdown failed: ", err)
	}

	fmt.Println(time.Now().UTC(), "Server successfully stopped.")
}
