package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/gliderlabs/ssh"
	"github.com/phthallo/guestbook/internal"
)

var (
	name       string
	message    string
	submitted  bool
)

func main(){
	conn, ctx := internal.CreateDBConnection()
	if _, err := internal.CreateEntriesIfNotExists(conn, ctx); err != nil {
		fmt.Printf("Entries table creation failed. %v", err)
		return
	}

	ssh.Handle(func(sess ssh.Session) {
		io.WriteString(sess, "dfsdfsd")
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("you've reached phthallo. leave a message?").
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
			io.WriteString(sess, fmt.Sprintf("Submmited name was %f, submitted messagw asa %f", name, message))
			tx, transaction_err := conn.Begin(context.Background())
			if (transaction_err != nil){
				io.WriteString(sess, fmt.Sprintf("an error occurred starting the transaction:, %s\n", transaction_err))
				return
			}
			if _, exec_err := tx.Exec(context.Background(), "INSERT into entries (name, message) VALUES ($1, $2)", name, message); exec_err != nil {
				io.WriteString(sess, fmt.Sprintf("%v", exec_err))
				return
			}
			if commit_err := tx.Commit(context.Background()); commit_err != nil {
				io.WriteString(sess, fmt.Sprintf("an error occurred saving your message: %s\n", commit_err))
				return
			}
			io.WriteString(sess, fmt.Sprintf("thanks for the message, %s\n", name))
		} else {
			io.WriteString(sess, fmt.Sprintf("see ya, %s\n", name))
		}
	})
	defer conn.Close(context.Background())
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}