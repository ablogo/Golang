package dependencies

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"native/api"
	"native/api/admin"
	"src/services"
)

type HttpHandler struct {
	Server *http.Server
}

func (h *HttpHandler) StartServer() {

	mux := http.NewServeMux()

	//mux.HandleFunc("GET /", api.GetUser)

	{
		mux.HandleFunc("POST /auth/sign-in", api.SignIn)
		mux.HandleFunc("POST /auth/sign-up", api.SignUp)
		mux.HandleFunc("POST /auth/validate-token", api.VerifyToken)
	}
	{
		mux.HandleFunc("GET /user", api.GetUser)
		mux.HandleFunc("PUT /user", api.UpdateUser)
		mux.HandleFunc("POST /user/change-password", api.ChangePassword)
		mux.HandleFunc("GET /user/img", api.GetPicture)
		mux.HandleFunc("POST /user/img", api.AddPicture)
		mux.HandleFunc("GET /user/address", api.GetAddresses)
		mux.HandleFunc("POST /user/address", api.AddAddress)
		mux.HandleFunc("DELETE /user/address", api.DeleteAddress)
		mux.HandleFunc("PUT /user/address", api.UpdateAddress)
		mux.HandleFunc("DELETE /user", api.DeleteUser)
	}
	{
		mux.HandleFunc("GET /admin/user", admin.GetUser)
		mux.HandleFunc("GET /admin/users", admin.GetUsers)
		mux.HandleFunc("DELETE /admin/user", admin.DeleteUser)
		mux.HandleFunc("GET /admin/user/address", admin.GetAddresses)
	}

	h.Server = &http.Server{
		Addr:           ":" + services.GetPort(),
		Handler:        mux,
		ReadTimeout:    time.Second * 6,
		WriteTimeout:   time.Second * 15,
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
