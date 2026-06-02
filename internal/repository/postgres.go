package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (*pgxpool.Pool, error) {
	connStr := "postgres://strong_fella_admin:final_password_228@localhost:5433/strong_fellas_db?sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	fmt.Println("Successfully connected to the database!")
	

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS leads (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		phone VARCHAR(30) NOT NULL,
		moving_from TEXT NOT NULL,
		moving_to TEXT NOT NULL,
		moving_date VARCHAR(50) NULL,
		details TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		fellas_number INT NOT NULL,
	)
	`
	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to create leads table: %w", err)
	}
	fmt.Println("Database tables verified/created successfully!")

	//--------------------------------------------------------
	return pool, nil
}

func SaveLead(pool *pgxpool.Pool, lead Lead) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	insertSQL := `

	`

	_, err := pool.Exec(ctx,)
	if err != nil {
		return fmt.Errorf("unable to save lead: %w", err)
	}
	return nil