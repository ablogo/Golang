package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"
	repository "watcher/db"
	"watcher/services"
)

type HttpHandler struct {
	Repo   *repository.Sqlyte
	Server *http.Server
}

func (h *HttpHandler) getRecords(w http.ResponseWriter, r *http.Request) {

	records, err := h.Repo.GetRecords()
	status := http.StatusOK
	content := "application/json"
	if err != nil || records == nil {
		slog.Error(err.Error())
		status = http.StatusNotFound
		content = "text/html; charset=UTF-8"
		w.Header().Set("Cache-Control", "no-cache, no-store, or max-age=0")
	}

	w.Header().Set("Content-Type", content)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(records)
}

func (h *HttpHandler) getLogs(w http.ResponseWriter, r *http.Request) {

	logs, err := h.Repo.GetLogs()
	status := http.StatusOK
	content := "application/json"
	if err != nil || logs == nil {
		slog.Error(err.Error())
		status = http.StatusNotFound
		content = "text/html; charset=UTF-8"
		w.Header().Set("Cache-Control", "no-cache, no-store, or max-age=0")
	}

	w.Header().Set("Content-Type", content)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(logs)
}

func (h *HttpHandler) getIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(time.Now().String()))
}

func (h *HttpHandler) StartServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", h.getIndex)
	mux.HandleFunc("GET /records", h.getRecords)
	mux.HandleFunc("GET /logs", h.getLogs)

	h.Server = &http.Server{
		Addr:         ":" + services.GetPort(),
		Handler:      mux,
		ReadTimeout:  time.Second * 6,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 120,
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
