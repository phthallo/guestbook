package main

import (
	"context"
	"fmt"
	"io"
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/gliderlabs/ssh"
	"github.com/joho/godotenv"
	"github.com/phthallo/guestbook/internal"
	"github.com/phthallo/guestbook/api"
)

var (
	name       string
	message    string
	submitted  bool
)


func StartSSHService() {
	conn, ctx := internal.CreateDBConnection()
	if _, err := internal.CreateEntriesIfNotExists(conn, ctx); err != nil {
		fmt.Printf("Entries table creation failed. %v", err)
		return
	}

	fmt.Println("\x1B[32m✔\tguestbook is up and running!\033[0m\t\t")
	theme := internal.ThemeGruvbox()
	ssh.Handle(func(sess ssh.Session) {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("you've reached phthallo. leave a message?").
					Value(&message).
					Validate( func(str string) error {
						fmt.Println(str)
						if len(str) == 0 {
							return errors.New("speak now or forever hold your peace")
						}
						return nil
					}),
				huh.NewInput().
					Title("who are you?").
					Value(&name).
					Validate( func(str string) error {
						fmt.Println(str)
						if len(str) == 0 {
							return errors.New("you're not arya stark!")
						}
						return nil
					}),
				huh.NewConfirm().
					Affirmative("submit").
					Negative("cancel").
					Value(&submitted),
			),
		).WithOutput(sess).WithInput(sess).WithTheme(theme)

		err := form.Run()
		if (err != nil) {
			fmt.Println(err)
			return
		} else if (submitted) {
			name, message = internal.Filter(name, message)
			fmt.Println("Name:%s\nMessage:%s\n", name, message)
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
			io.WriteString(sess, fmt.Sprintf("\x1B[34mthanks for the message, %s!\nsee it on phthallo.com/guestbook :)\033[0m", name))
		} else {
			io.WriteString(sess, fmt.Sprintf("\x1B[34mcome back when you've got something to say, %s\n\033[0m", name))
		}
	})
	defer conn.Close(context.Background())
	ssh.ListenAndServe(":2222", nil)
}

func main(){
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("Error loading environment variables!")
	}
	fmt.Println(`
   ____ _   _ _____ ____ _____ ____   ___   ___  _  __
  / ___| | | | ____/ ___|_   _| __ ) / _ \ / _ \| |/ /
 | |  _| | | |  _| \___ \ | | |  _ \| | | | | | | ' / 
 | |_| | |_| | |___ ___) || | | |_) | |_| | |_| | . \ 
  \____|\___/|_____|____/ |_| |____/ \___/ \___/|_|\_\`)
	go api.StartAPIService()
	go StartSSHService()
	select {}
}