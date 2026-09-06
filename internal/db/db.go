package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbUrl)

	if err != nil {
		return nil, fmt.Errorf("Sql.Open: %w", err)
	}

	db.SetMaxOpenConns(25) //setting the maximum no.of open connection
	db.SetMaxIdleConns(25) 
	db.SetConnMaxLifetime(5 * time.Minute)

	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //5Sec Context timeout
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil,fmt.Errorf("db.ping: %w", err)
	}

	return db, nil
}