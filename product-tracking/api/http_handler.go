package api

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"product-tracking/config"
	"product-tracking/router"
	"time"
)

type HttpHandler struct {
	Server   *http.Server
	Router   *router.Router
	Settings *config.Settings
}

func (h *HttpHandler) getIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(time.Now().String()))
}

func (h *HttpHandler) StartServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.getIndex)
	mux.HandleFunc("GET /product", h.Router.GetProduct)
	mux.HandleFunc("GET /products", h.Router.GetProducts)
	mux.HandleFunc("POST /products", h.Router.CreateProduct)
	mux.HandleFunc("PUT /product", h.Router.UpdateProduct)
	mux.HandleFunc("DELETE /product", h.Router.DeleteProduct)

	h.Server = &http.Server{
		Addr:         ":" + h.Settings.SERVER_PORT,
		Handler:      mux,
		ReadTimeout:  time.Second * h.Server.ReadTimeout,
		WriteTimeout: time.Second * h.Server.WriteTimeout,
		IdleTimeout:  time.Second * h.Server.IdleTimeout,
	}

	slog.Info("Starting server on port " + h.Server.Addr)

	if err := h.Server.ListenAndServe(); err != http.ErrServerClosed {
		log.Panicf("%s Server error: %v", time.Now().UTC(), err)
	}
}

func (h *HttpHandler) ShutdownServer() {
	slog.Info("Stopping server..")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	if err := h.Server.Shutdown(ctx); err != nil {
		slog.Error("Server graceful shutdown failed: ", "description", err.Error())
	}

	slog.Info("Server successfully stopped.")
}
