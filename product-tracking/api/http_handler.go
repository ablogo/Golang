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

func corsMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *HttpHandler) StartServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.getIndex)
	mux.HandleFunc("GET /product/{id}", h.Router.GetProduct)
	mux.HandleFunc("GET /products", h.Router.GetProducts)
	mux.HandleFunc("GET /product/prices/historical/{id}", h.Router.GetProductHistoricalLowestPrices)
	mux.HandleFunc("GET /product/prices/chart/{id}", h.Router.GetProductChartInfoById)
	mux.HandleFunc("POST /product", h.Router.CreateProduct)
	mux.HandleFunc("PUT /product/tracked/{id}", h.Router.ChangeProductTrackedStatus)
	mux.HandleFunc("PUT /product", h.Router.UpdateProduct)
	mux.HandleFunc("DELETE /product/{id}", h.Router.DeleteProduct)

	h.Server = &http.Server{
		Addr:         ":" + h.Settings.SERVER_PORT,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  time.Second * time.Duration(h.Settings.SERVER_READ_TIMEOUT),
		WriteTimeout: time.Second * time.Duration(h.Settings.SERVER_WRITE_TIMEOUT),
		IdleTimeout:  time.Second * time.Duration(h.Settings.SERVER_IDLE_TIMEOUT),
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
