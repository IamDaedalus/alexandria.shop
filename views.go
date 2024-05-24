package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// the initial screen to display when the program first runs
func (m model) initialScreen() string {
	// Center the ASCII art within the terminal window
	return lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		cyan.Align(lipgloss.Center).Render(welcomeAscii),
	)
}

func (m model) authScreen() string {
	var doc strings.Builder

	return doc.String()
}

// renders the login screen
func (m model) loginScreen() string {
	var (
		layout     strings.Builder
		login, pwd string
	)

	// the big "LOGIN" text at the top
	loginPrompt := lipgloss.NewStyle().
		Foreground(cyan.GetForeground()).
		Width(55).Height(5).
		Align(lipgloss.Center).
		Render(loginText)

		// variables for the username and password prompt boxes
	loginRender := renderBoxDesc("username", 0, m)
	pwdRender := renderBoxDesc("password", 1, m)

	// change the color of each render based on the current focus
	if m.authCurField == 0 {
		login = magenta.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(loginRender)
		pwd = faded.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(pwdRender)
	} else {
		login = faded.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(loginRender)
		pwd = magenta.PaddingLeft(8).Align(lipgloss.Left).Bold(true).Render(pwdRender)
	}

	// footer/help message render
	helpText := "press ctrl+s to sign up | press ctrl+c to quit"
	helpBox := noBorderStyle.PaddingTop(1).Width(50).Align(lipgloss.Bottom).Render(helpText)

	// join the various fields together;
	// first the input boxes and then those and the login prompt
	textFields := lipgloss.JoinVertical(lipgloss.Left, login, pwd, helpBox)
	ui := lipgloss.JoinVertical(lipgloss.Left, loginPrompt, textFields)

	// render the actual with dialogBoxStyle but this simply "puts" the render
	// in the center of the screen no matter what
	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	// everything onscreen is a string so tie it up nice with a bow and return a string
	layout.WriteString(dialog)
	return layout.String()
}

// function to return a nicely formatted description and input box
func renderBoxDesc(s string, idx int, m model) string {
	desc := noBorderStyle.Render(s)
	inputBox := textBoxStyle.Render(m.authInputs[idx].View())
	finalRender := lipgloss.JoinHorizontal(lipgloss.Left, desc, inputBox)

	return finalRender
}