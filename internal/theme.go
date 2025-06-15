package internal

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
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
