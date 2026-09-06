package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kalpeshWani222/olx-api/internal/config"
	"github.com/kalpeshWani222/olx-api/internal/config/handlers"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println("starting server.....")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)


	srv := http.Server{
		Addr: ":" + cfg.Port,
		Handler: mux,
		ReadTimeout: time.Second*10,
		WriteTimeout: time.Second*30,
		IdleTimeout: time.Second*60,
	}

	log.Printf("server is running on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed : %v",err)
	}
}