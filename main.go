package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/gliderlabs/ssh"
)

var (
	name       string
	message    string
	submitted  bool
)

func main(){
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
			io.WriteString(sess, fmt.Sprintf("thanks for the message, %s\n", name))
		} else {
			io.WriteString(sess, fmt.Sprintf("see ya, %s\n", name))
		}
	})
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
