// Package main — точка входа Connect-сервера financial-analyst.
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DanilaKorobkov/financial-analyst/app"
)

const readHeaderTimeout = 5 * time.Second

func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("конфигурация: %v", err)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	server := &http.Server{
		Addr:              addr,
		Handler:           app.New(cfg),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	log.Printf("Connect-сервер слушает %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
