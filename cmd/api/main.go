package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kalpeshWani222/olx-api/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println("starting server.....")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)


		w.Write([]byte(`{"status" : "ok"}`))
	})


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