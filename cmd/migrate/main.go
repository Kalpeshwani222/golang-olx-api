package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kalpeshWani222/olx-api/internal/config"
)

func main() {

	if (len(os.Args) < 2) {
		log.Fatal("Usage make migrate <up | down>")
	}

	cfg := config.MustLoad()

	m, err := migrate.New("file://migrations", cfg.DbUrl)

	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}
	
	//handling the up and down migration
	switch os.Args[1]{
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}

	case "down" :
	if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}

	default : 
	log.Fatalf("unknow comman: %s", os.Args[1])

	}
}