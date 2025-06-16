package internal

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDBConnection() (pool *pgxpool.Pool, ctx context.Context){
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return pool, context.Background()
}

func CreateEntriesIfNotExists(dbpool *pgxpool.Pool, ctx context.Context) (bool, error) {
	// the boolean that this function returns refers to whether the table exists 
	var n int64
	err := dbpool.QueryRow(ctx, "select 1 from information_schema.tables where table_name = $1 AND table_schema = 'public'", "entries").Scan(&n)

	if (err != nil){
		// attempt to create the database
		tx, err := dbpool.Begin(ctx); 

		if err != nil {
			return false, fmt.Errorf("\x1B[31m✕\033[0m\tunable to begin transaction: %v", err)
		}

		if _, err := tx.Exec(ctx, "CREATE TABLE entries (id INTEGER primary key GENERATED ALWAYS AS IDENTITY, name TEXT, message TEXT, timestamp TIMESTAMP default current_timestamp)"); err != nil {
			return false, fmt.Errorf("\x1B[31m✕\033[0m\tunable to execute table creation: %v", err)
		}

		if err = tx.Commit(ctx); err != nil {
			return false, fmt.Errorf("✕ unable to commit table creation: %v", err) 
		}

		fmt.Println("\x1B[32m✔\033[0m\tcreated entries table for storing guestbook entries!")

		return true, nil
	}

	fmt.Println("\x1B[32m✔\033[0m\tentries table for storing guestbook entries already exists!")
	return true, nil
}