package main

import (
	"context"
	"fmt"
	"io"
	"errors"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/gliderlabs/ssh"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phthallo/guestbook/internal"
	"github.com/phthallo/guestbook/api"
)

var (
	name       string
	message    string
	submitted  bool
)


func StartSSHService(dbpool *pgxpool.Pool, ctx context.Context) {
	if _, err := internal.CreateEntriesIfNotExists(dbpool, ctx); err != nil {
		fmt.Printf("Entries table creation failed. %v", err)
		return
	}

	fmt.Println("\x1B[32m✔\tguestbook is up and running!\033[0m\t\t")
	theme := huh.ThemeDracula()
	ssh.Handle(func(sess ssh.Session) {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(fmt.Sprintf("thanks for calling, %s.\nyou've reached phthallo. leave a message?", sess.User())).
					Value(&message).
					Validate( func(str string) error {
						if len(str) == 0 {
							return errors.New("speak now or forever hold your peace")
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
			name = sess.User()
			name, message = internal.Filter(name, message)
			fmt.Printf("A new message was submitted by '%s'. They said '%s'\n", name, message)
			tx, transaction_err := dbpool.Begin(context.Background())
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
			io.WriteString(sess, fmt.Sprintf("\x1B[34mty for the message, %s!\nsee it on phthallo.com/guestbook :)\n\033[0m", name))
		} else {
			io.WriteString(sess, fmt.Sprintf("\x1B[34mcome back when you've got something to say, %s\n\033[0m", name))
		}
	})
	ssh.ListenAndServe(fmt.Sprintf(":%s",os.Getenv("SSH_PORT")), nil, ssh.HostKeyFile("id_ed25519"))
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
  	dbpool, ctx := internal.CreateDBConnection()

	go api.StartAPIService(dbpool)
	go StartSSHService(dbpool, ctx)
	defer dbpool.Close()
	select {}
}