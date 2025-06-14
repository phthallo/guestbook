package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/gliderlabs/ssh"
	"github.com/jackc/pgx/v5"
)

var (
	name       string
	message    string
	submitted  bool
)

func main(){
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	tableExists, err := createEntriesIfNotExists(context.Background(), conn);
	fmt.Printf("Entries table exists: %t, %v", tableExists, err)

	ssh.Handle(func(sess ssh.Session) {
		io.WriteString(sess, "dfsdfsd")
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("leave a message?").
					Value(&message),
				huh.NewInput().
					Title("who are you?").
					Value(&name),
				huh.NewConfirm().
					Affirmative("submit").
					Negative("cancel").
					Value(&submitted),
			),
		).WithOutput(sess).WithInput(sess)
		err := form.Run()
		if (err != nil) {
			fmt.Println(err)
			os.Exit(1)
		} else if (submitted) {
			tx, transaction_err := conn.Begin(context.Background())
			if (transaction_err != nil){
				io.WriteString(sess, fmt.Sprintf("an error occurred starting the transaction:, %s\n", transaction_err))
			}
			if _, exec_err := tx.Exec(context.Background(), "INSERT into entries (name, message) VALUES ($1, $2)", name, message); exec_err != nil {
				fmt.Printf("%v", exec_err)
			}
			commit_err := tx.Commit(context.Background())
			if (commit_err != nil){
				io.WriteString(sess, fmt.Sprintf("an error occurred saving your message: %s\n", commit_err))
				return
			}
			io.WriteString(sess, fmt.Sprintf("thanks for the message, %s\n", name))
		} else {
			io.WriteString(sess, fmt.Sprintf("see ya, %s\n", name))
		}
	})
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

func createEntriesIfNotExists(ctx context.Context, conn *pgx.Conn) (bool, error) {
	// the boolean that this function returns refers to whether the table exists 
	var n int64
	err := conn.QueryRow(ctx, "select 1 from information_schema.tables where table_name = $1", "entries").Scan(&n)
	if (err != nil){
		// attempt to create the database
		tx, err := conn.Begin(ctx); 
		if err != nil {
			return false, fmt.Errorf("unable to begin transaction: %v", err)
		}

		if _, err := tx.Exec(ctx, "CREATE TABLE entries (name text, message text)"); err != nil {
			return false, fmt.Errorf("unable to execute table creation: %v", err)
		}

		if err = tx.Commit(ctx); err != nil {
			return false, fmt.Errorf("unable to commit table creation: %v", err) 
		}
		return true, nil
	}
	return true, nil
}