package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kalpeshWani222/olx-api/internal/config"
	"github.com/kalpeshWani222/olx-api/internal/db"
	"github.com/kalpeshWani222/olx-api/internal/handlers"
)

func main() {
	//loading the configs
	cfg := config.MustLoad()

	//connecting to the db
	db,err := db.Connect(cfg.DbUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v",err)
	}

	hanlder := slog.NewJSONHandler(os.Stdout,&slog.HandlerOptions{
		AddSource: true,
		Level: slog.LevelInfo, //default LevelInfo
	})
	logger := slog.New(hanlder)
	slog.SetDefault(logger)

	fmt.Println("database connected.....")
	fmt.Println("starting server.....")

	lh := handlers.NewListingHandler(db,logger)

	mux := http.NewServeMux()

	//endpoints
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
    mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	//server
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