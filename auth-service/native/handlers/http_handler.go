package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"native/middleware"
	"native/routers"
	"src/utils"
)

type HttpHandler struct {
	Server      *http.Server
	Settings    *utils.Settings
	Router      *routers.Router
	AdminRouter *routers.AdminRouter
}

func (h *HttpHandler) StartServer() {

	mux := http.NewServeMux()

	//mux.HandleFunc("GET /", routers.GetUser)

	{
		mux.HandleFunc("POST /auth/sign-in", h.Router.SignIn)
		mux.HandleFunc("POST /auth/sign-up", h.Router.SignUp)
		mux.HandleFunc("POST /auth/validate-token", routers.VerifyToken)
	}
	{
		mux.HandleFunc("GET /user", h.Router.GetUser)
		mux.HandleFunc("PUT /user", h.Router.UpdateUser)
		mux.HandleFunc("POST /user/change-password", h.Router.ChangePassword)
		mux.HandleFunc("GET /user/img", h.Router.GetPicture)
		mux.HandleFunc("POST /user/img", h.Router.AddPicture)
		mux.HandleFunc("GET /user/address", h.Router.GetAddresses)
		mux.HandleFunc("POST /user/address", h.Router.AddAddress)
		mux.HandleFunc("DELETE /user/address", h.Router.DeleteAddress)
		mux.HandleFunc("PUT /user/address", h.Router.UpdateAddress)
		mux.HandleFunc("DELETE /user", h.Router.DeleteUser)
	}
	{
		mux.HandleFunc("GET /admin/user", h.AdminRouter.GetUser)
		mux.HandleFunc("GET /admin/users", h.AdminRouter.GetUsers)
		mux.HandleFunc("DELETE /admin/user", h.AdminRouter.DeleteUser)
		mux.HandleFunc("GET /admin/user/address", h.AdminRouter.GetAddresses)
	}

	JWTMiddleware := middleware.JWTMiddleware(mux)

	h.Server = &http.Server{
		Addr:           ":" + h.Settings.SERVER_PORT,
		Handler:        JWTMiddleware,
		ReadTimeout:    time.Second * time.Duration(h.Settings.SERVER_READ_TIMEOUT),
		WriteTimeout:   time.Second * time.Duration(h.Settings.SERVER_WRITE_TIMEOUT),
		IdleTimeout:    time.Second * time.Duration(h.Settings.SERVER_IDLE_TIMEOUT),
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
