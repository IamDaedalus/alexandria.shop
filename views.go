package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// the initial screen to display when the program first runs
func (m model) initialScreen() string {
	// center the ascii art within the terminal window
	return lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		cyan.Align(lipgloss.Center).Render(welcomeAscii),
	)
}

func (m model) renderInputBox(label string, index int, inputs []textinput.Model, focused bool) string {
	render := renderBoxDesc(label, index, focused, inputs)
	if focused {
		return magenta.PaddingLeft(8).Align(lipgloss.Left).Render(render)
	}
	return faded.PaddingLeft(8).Align(lipgloss.Left).Render(render)
}

func (m model) renderAuthScreen(title, helpText string, inputs []textinput.Model, curField int) string {
	var layout strings.Builder

	// render the title at the top
	titleRender := lipgloss.NewStyle().
		Foreground(cyan.GetForeground()).
		Width(55).Height(5).
		Align(lipgloss.Center).
		Render(title)

	// render input boxes
	var inputRenders []string
	for i, label := range []string{"username", "password", "password"} {
		if i >= len(inputs) {
			break
		}
		inputRenders = append(inputRenders, m.renderInputBox(label, i, inputs, curField == i))
	}

	// render the help text at the bottom
	helpBox := noBorderStyle.PaddingTop(1).Width(50).Align(lipgloss.Bottom).Render(helpText)

	// join the input fields and help text
	textFields := lipgloss.JoinVertical(lipgloss.Left, inputRenders...)
	textFields = lipgloss.JoinVertical(lipgloss.Left, textFields, helpBox)

	// combine the title and input fields
	ui := lipgloss.JoinVertical(lipgloss.Left, titleRender, textFields)

	// place the ui in the center of the screen
	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
	)

	// return the final rendered layout
	layout.WriteString(dialog)
	return layout.String()
}

func (m model) signUpScreen() string {
	return m.renderAuthScreen(
		signUpText,
		"press ctrl+l to log in | press ctrl+c to quit",
		m.signupInputs,
		m.signupCurField,
	)
}

func (m model) loginScreen() string {
	return m.renderAuthScreen(
		loginText,
		"press ctrl+s to sign up | press ctrl+c to quit",
		m.loginInputs,
		m.loginCurField,
	)
}

func (m model) infoScreen(info string) string {
	infoRender := noBorderStyle.
		PaddingTop(1).
		Width(50).Height(3).
		Align(lipgloss.Center).Render(info)

	// footer/help message render
	// helpText := "press enter"
	// helpBox := noBorderStyle.Width(50).Height(1).Align(lipgloss.Bottom).Render(helpText)

	dialog := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Width(50).Render(infoRender),
	)

	return dialog
}

// this renders the entire catalogue view
// big thanks to @its_gaurav on the Charm CLI Discord!
func (m model) catalogueScreen(curUser string) string {
	// Initialize variables
	renderWidth := (m.termWidth / 2) + 10
	if renderWidth < 0 {
		renderWidth = 0
	}

	headerRender := m.renderHeaders(curUser, "c cart [16]", renderWidth)
	footer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		Width(renderWidth).
		Align(lipgloss.Center).
		Render(catalogueHelpMsg)

	mainHeight := m.termHeight - lipgloss.Height(headerRender) - lipgloss.Height(footer)
	offset := (mainHeight - 20) / 3
	// function to determine if an item is highlighted
	isHighlighted := func(index int) bool {
		return m.curItem == index
	}

	// render the top, mid, and bot items based on current item
	var itemsRender []string
	for i := 0; i < magicNum; i++ {
		itemsRender = append(itemsRender, renderItemDisplay(renderWidth, offset, isHighlighted(i), m.curBooks[i]))
	}

	catalogueView := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderBottom(false).
		Width(renderWidth).
		Height(mainHeight - 20).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(lipgloss.Center, itemsRender...))

	seperator := "├" + strings.Repeat("─", renderWidth) + "┤"

	headerRenderModified := strings.ReplaceAll(headerRender, "└", "├")
	headerRenderModified = strings.ReplaceAll(headerRenderModified, "┘", "┤")

	catalogueRender := lipgloss.JoinVertical(
		lipgloss.Center,
		headerRenderModified,
		catalogueView,
		seperator,
		footer,
	)

	cFinalRender := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		catalogueRender,
	)

	return cFinalRender
}

// renders an item in the catalogue view
func renderItemDisplay(renderWidth, renderHeight int, selected bool, b book) string {
	selectProperties := lipgloss.NewStyle().
		Foreground(faded.GetForeground()).
		Border(lipgloss.NormalBorder())

	if selected {
		selectProperties = lipgloss.NewStyle().
			Foreground(magenta.GetForeground()).
			Border(lipgloss.ThickBorder())

		// assign the current selected book to the global selected book
		selectedBook = b
	}

	itemContent := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder(), false, true, false, true).
		Foreground(selectProperties.GetForeground()).
		Render(fmt.Sprintf("%s  |  %s  |  $%.2f", b.Title, b.Author, b.Price))

	itemDesc := lipgloss.NewStyle().
		Border(lipgloss.HiddenBorder(), false, true, false, true).
		Border(lipgloss.HiddenBorder()).
		Foreground(selectProperties.GetForeground()).
		Render(b.Description)

	contentRender := lipgloss.JoinVertical(lipgloss.Top, itemContent, itemDesc)

	return lipgloss.NewStyle().
		Border(selectProperties.GetBorder()).
		BorderForeground(selectProperties.GetForeground()).
		Width(renderWidth - 5).Height(renderHeight).
		Render(contentRender)
}

// render the headers at the top of the catalogue page
func (m model) renderHeaders(curUser, cart string, renderWidth int) string {
	tops := [][]string{
		{
			"alexandria.shop",
			fmt.Sprintf("welcome, %s", curUser),
			fmt.Sprintf("? for help and details"),
			cart,
		}, // actual headers
	}

	headerTable := table.New().
		Border(lipgloss.NormalBorder()).
		Width(renderWidth + 2).
		StyleFunc(table.StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
		})).
		Rows(tops...)

	return headerTable.Render()
}
func renderHeaderBox(s string) string {
	return headerBoxStyle.
		Width(20).Margin(-1).
		Align(lipgloss.Center).
		Render(s)
}

// function to return a nicely formatted description and input box
func renderBoxDesc(s string, idx int, focused bool, inputs []textinput.Model) string {
	desc := noBorderStyle.Bold(focused).Render(s)
	// this is the content from the input box as we type
	// side not find a way to render the textbox thicker
	inputBox := textBoxStyle.Render(inputs[idx].View())
	finalRender := lipgloss.JoinHorizontal(lipgloss.Left, desc, inputBox)

	return finalRender
}
