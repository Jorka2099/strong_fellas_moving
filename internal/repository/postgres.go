package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Lead struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	MovingFrom   string    `json:"moving_from"`
	MovingTo     string    `json:"moving_to"`
	MovingDate   string    `json:"moving_date"`
	FellasNumber int       `json:"fellas_number"`
	Details      string    `json:"details"`
	CreatedAt    time.Time `json:"created_at"`
}

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
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		phone VARCHAR(30) NOT NULL,
		moving_from TEXT NOT NULL,
		moving_to TEXT NOT NULL,
		moving_date VARCHAR(50) NULL,
		fellas_number INT NOT NULL,
		details TEXT,
		created_at TIMESTAMPTZ DEFAULT NOW()
	)
	`
	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to create leads table: %w", err)
	}

	fmt.Println("Database tables verified/created successfully!")
	return pool, nil
}

func SaveLead(pool *pgxpool.Pool, lead Lead) error {
	if pool == nil {
		return fmt.Errorf("database pool is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO leads (name, phone, moving_from, moving_to, moving_date, fellas_number, details) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := pool.Exec(ctx, query,
		lead.Name,
		lead.Phone,
		lead.MovingFrom,
		lead.MovingTo,
		lead.MovingDate,
		lead.FellasNumber,
		lead.Details,
	)
	if err != nil {
		return fmt.Errorf("unable to save lead: %w", err)
	}
	return nil
}

func GetAllLeads(pool *pgxpool.Pool) ([]Lead, error) {
	if pool == nil {
		return nil, fmt.Errorf("database pool is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	    SELECT id, name, phone, moving_from, moving_to, moving_date, fellas_number, details, created_at FROM leads ORDER BY created_at DESC;`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to query leads: %w", err)
	}
	defer rows.Close()

	var leads []Lead

	for rows.Next() {
		var lead Lead
		err := rows.Scan(&lead.ID, &lead.Name, &lead.Phone, &lead.MovingFrom, &lead.MovingTo, &lead.MovingDate, &lead.FellasNumber, &lead.Details, &lead.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("unable to scan lead: %w", err)
		}
		leads = append(leads, lead)
	}

	return leads, nil
}
