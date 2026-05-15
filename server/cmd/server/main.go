// Package main — точка входа Connect-сервера financial-analyst.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DanilaKorobkov/financial-analyst/app"
)

func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("конфигурация: %v", err)
	}

	handler, err := app.New(&cfg)
	if err != nil {
		log.Fatalf("сборка приложения: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}

	log.Printf("Connect-сервер слушает %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
