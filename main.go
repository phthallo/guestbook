package main

import (
	"context"
	"fmt"
	"io"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/gliderlabs/ssh"
	"github.com/phthallo/guestbook/internal"
	"github.com/phthallo/guestbook/api"
	"github.com/joho/godotenv"
)

var (
	name       string
	message    string
	submitted  bool
)

func ThemeGruvbox() *huh.Theme {
	t := huh.ThemeBase()
	
	var (
		background = lipgloss.AdaptiveColor{Dark: "#282828"}
		selection  = lipgloss.AdaptiveColor{Dark: "#3c3836"}
		foreground = lipgloss.AdaptiveColor{Dark: "#ebdbb2"}
		comment    = lipgloss.AdaptiveColor{Dark: "#928374"}
		green      = lipgloss.AdaptiveColor{Dark: "#b8bb26"}
		red        = lipgloss.AdaptiveColor{Dark: "#fb4934"}
		yellow     = lipgloss.AdaptiveColor{Dark: "#fabd2f"}
		orange     = lipgloss.AdaptiveColor{Dark: "#fe8019"}
		blue       = lipgloss.AdaptiveColor{Dark: "#83a598"}
	)
	t.Focused.Base = t.Focused.Base.BorderForeground(selection)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(orange)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(orange)
	t.Focused.Description = t.Focused.Description.Foreground(comment)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.Directory = t.Focused.Directory.Foreground(blue)
	t.Focused.File = t.Focused.File.Foreground(foreground)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(orange)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(yellow)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(yellow)
	t.Focused.Option = t.Focused.Option.Foreground(foreground)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(orange)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(green)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(foreground)
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(comment)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(background).Background(yellow).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(foreground).Background(selection)
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(green)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(comment)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(yellow)
	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()
	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description
	return t
}

func StartSSHService() {
	conn, ctx := internal.CreateDBConnection()
	if _, err := internal.CreateEntriesIfNotExists(conn, ctx); err != nil {
		fmt.Printf("Entries table creation failed. %v", err)
		return
	}

	fmt.Println("\x1B[32mGuestbook is up and running!\033[0m\t\t")
	theme := ThemeGruvbox()
	ssh.Handle(func(sess ssh.Session) {

	signals := make(chan ssh.Signal, 1)
	sess.Signals(signals)
	go func() {
            for {
	        <-signals
	    }
    }()
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
		).WithOutput(sess).WithInput(sess).WithTheme(theme)

		err := form.Run()
		if (err != nil) {
			fmt.Println(err)
			return
		} else if (submitted) {
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
			io.WriteString(sess, fmt.Sprintf("\x1B[32mthanks for the message, %s!\n\033[0m", name))
		} else {
			io.WriteString(sess, fmt.Sprintf("see ya, %s\n", name))
		}
	})
	defer conn.Close(context.Background())
	ssh.ListenAndServe(":2222", nil)

}

func main(){
	if err := godotenv.Load(); err != nil {
		fmt.Errorf("Error loading environment variables!")
	}
	fmt.Printf(`
   ____ _   _ _____ ____ _____ ____   ___   ___  _  __
  / ___| | | | ____/ ___|_   _| __ ) / _ \ / _ \| |/ /
 | |  _| | | |  _| \___ \ | | |  _ \| | | | | | | ' / 
 | |_| | |_| | |___ ___) || | | |_) | |_| | |_| | . \ 
  \____|\___/|_____|____/ |_| |____/ \___/ \___/|_|\_\

`)
	go api.StartAPIService()
	go StartSSHService()
	select {}
}