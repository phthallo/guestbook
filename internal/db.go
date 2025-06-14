package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func CreateDBConnection() (conn *pgx.Conn, ctx context.Context){
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	return conn, context.Background()
}

func CreateEntriesIfNotExists(conn *pgx.Conn, ctx context.Context) (bool, error) {
	// the boolean that this function returns refers to whether the table exists 
	var n int64
	err := conn.QueryRow(ctx, "select 1 from information_schema.tables where table_name = $1 AND table_schema = 'public'", "entries").Scan(&n)
	fmt.Printf("%v", n)
	if (err != nil){
		// attempt to create the database
		tx, err := conn.Begin(ctx); 
		if err != nil {
			return false, fmt.Errorf("unable to begin transaction: %v", err)
		}

		if _, err := tx.Exec(ctx, "CREATE TABLE entries (id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY, name TEXT, message TEXT)"); err != nil {
			return false, fmt.Errorf("unable to execute table creation: %v", err)
		}

		if err = tx.Commit(ctx); err != nil {
			return false, fmt.Errorf("unable to commit table creation: %v", err) 
		}
		return true, nil
	}
	return true, nil
}