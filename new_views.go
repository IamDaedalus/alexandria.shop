package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) checkoutScreen() string {
	headers := []string{
		"esc to go back",
		fmt.Sprint("welcome, ", m.curUser.username),
		"? for details/help",
		"",
	}

	footerMsg := "ctrl+l to logout  |  up/down to navigate  "
	m.content = mainRenderContent{
		headerContents: headers,
		footerMessage:  footerMsg,
	}

	render := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(cyan.GetForeground()).
		Margin(1, 1, 0, 1).
		Padding(5).
		Height(m.spatials.innerH - (m.spatials.innerH % m.itemsDispCount)).
		Width(m.spatials.bookDeetsW + m.spatials.listSectionW + 3).
		Render(checkoutMsg)
		// Render(checkoutMsg, "\n\n", m.renderInputBox("email", 0, charLimitEmail+2, m.checkoutInputs, true))

	return m.mainBorderRender(render)
}

func (m *model) cartScreen() string {
	var details string

	headers := []string{
		"esc to go back",
		fmt.Sprint("welcome, ", m.curUser.username),
		"? for details/help",
		"c to checkout",
	}

	footerMsg := "ctrl+l to logout  |  up/down to navigate  "
	m.content = mainRenderContent{
		headerContents: headers,
		footerMessage:  footerMsg,
	}

	isHighlighted := func(index int) bool {
		return m.cartItemIter == index
	}
	items := make([]string, 0)
	count := m.itemsDispCount
	if len(m.curCartItems) > 0 {
		// section to render the items in the list
		if count > len(m.curCartItems) {
			count = len(m.curCartItems)
		}
		for i := 0; i < count; i++ {
			items = append(items, renderItem(m.spatials.listSectionW-4, m.curCartItems[i], isHighlighted(i)))
		}
	}

	itemsRender := lipgloss.JoinVertical(lipgloss.Center, items...)

	if len(m.c.items) == 0 {
		details = "\nthere are no items in your cart right now!\nyou can add items by pressing + or = on a selected book!"
	} else {
		details = "here you'll find a list of all your books in the cart. you can scroll through and remove those you don't"
		details += " need in the catalogue section."
		details += "\nonce you're ready to complete your purchase, press c to go to the checkout"
		details += fmt.Sprintf("\n\nTOTAL: $%.2f", m.c.booksTotal())
	}

	midSectionJoin := renderMidSections(
		m.spatials,
		m.itemsDispCount,
		itemsRender,
		details,
	)

	return m.mainBorderRender(midSectionJoin)
}

func (m *model) mainScreen() string {
	headers := []string{
		"alexandria.shop",
		fmt.Sprint("welcome, ", m.curUser.username),
		"? for details/help",
		fmt.Sprintf("c cart [%d] %.2f", len(m.c.items), m.c.booksTotal()),
	}

	if m.itemsDispCount > magicNum {
		m.curBooks, _ = getBooksForPage(m.db, m.itemsDispCount, m.mainOffset)
		if m.curBooks == nil {
			log.Fatalf("no books found")
		}
	}

	// titles := extractTitles(m.curBooks)

	footerMsg := "ctrl+l to logout  |  up/down to navigate  "
	m.content = mainRenderContent{
		headerContents: headers,
		// bookItems:      titles,
		footerMessage: footerMsg,
	}

	isHighlighted := func(index int) bool {
		return m.mainItemsIter == index
	}
	// section to render the items in the list
	items := make([]string, 0)
	for i := 0; i < len(m.curBooks); i++ {
		text := m.curBooks[i].Title
		// purpose of the asterisk is to indicate the item is already in the cart
		// previously it was a big "[in cart]" text but we need to save space with
		// the new UI
		if m.c.bookInCart(m.curBooks[i]) {
			text = fmt.Sprintf("* %s", m.curBooks[i].Title)
		}

		text = truncateText(text, m.spatials.listSectionW-4)

		items = append(items, renderItem(m.spatials.listSectionW-4, text, isHighlighted(i)))
		// set the selectedBook global var when an option is highlighted
		if isHighlighted(i) {
			selectedBook = m.curBooks[i]
		}
	}

	itemsRender := lipgloss.JoinVertical(lipgloss.Center, items...)

	midSectionJoin := renderMidSections(
		m.spatials,
		m.itemsDispCount,
		itemsRender,
		styleBookDeets(),
	)

	return m.mainBorderRender(midSectionJoin)
}

func renderHeader(w int, content string, border bool, margin ...int) string {
	border = !border
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false). //, false, true).//, false, false, false, false).
		Margin(margin...).
		Width(w).
		Align(lipgloss.Center).
		Render(content)
}

func renderItem(w int, content string, selected bool) string {
	selectedCol := cyan.GetForeground()

	s := lipgloss.NewStyle().
		Foreground(white.GetBackground()).
		Border(lipgloss.NormalBorder())

	if selected {
		s = lipgloss.NewStyle().
			Foreground(magenta.GetForeground()).
			Border(lipgloss.ThickBorder())

		selectedCol = magenta.GetForeground()
	}

	return lipgloss.NewStyle().
		Border(s.GetBorder()).
		BorderForeground(selectedCol).
		Foreground(selectedCol).
		Padding(0, 1).
		Width(w).
		Align(lipgloss.Left).
		Render(content)
}

func setupDimensions(termHeight, termWidth int) dimensions {
	// in order to achieve a common ground for dynamic responsiveness, set a maximum size
	// this is barely a scratch though
	dynRenderWidth := termWidth - (termWidth / 6)            // calculates the best width
	dynRenderHeight := termHeight - (termHeight / 4)         // calculates the best height
	actualRenderW := dynRenderWidth - 21                     // that's the width of the main border
	actualRenderH := dynRenderHeight + 5                     // that's the height of the main border
	innerW := actualRenderW - 5                              // fake padding subtracted from th main border width
	innerH := (actualRenderH / 4) + (actualRenderH / 2)      // fake padding subtracted from th main border height
	listSectionW := innerW / 3                               // listsection width
	bookDeetW := (innerW / 3) + (innerW / 4) + (innerW / 13) // bookdetails width

	return dimensions{
		dynRenderWidth,
		dynRenderHeight,
		actualRenderW,
		actualRenderH,
		innerW,
		innerH,
		listSectionW,
		bookDeetW,
	}
}

func styleBookDeets() string {
	return lipgloss.NewStyle().
		Render(fmt.Sprintf(
			"%s by %s\nPrice: $%.2f\nTags: %s\n\n\n%s",
			selectedBook.Title, selectedBook.Author, selectedBook.Price,
			selectedBook.Genre, selectedBook.LongDesc))
}

func (m *model) mainBorderRender(midSectionJoin string) string {
	var (
		headers       = make([]string, 0)
		customMargins = [][]int{
			{0, 2, 0, 0},
			{0, 2, 0, 2},
			{0, 2, 0, 2},
			{0, 0, 0, 2},
		}
	)

	// renders headers
	for i := 0; i < 4; i++ {
		headers = append(headers, renderHeader((m.spatials.actualRenderW-4)/5, m.content.headerContents[i], false, customMargins[i]...))
	}
	// joins all header related text
	innerHeaderRender := lipgloss.JoinHorizontal(lipgloss.Left, headers...)

	// renders the header section complete with header text
	header := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(cyan.GetForeground()).
		Margin(0, 1, 0, 1).
		Width(m.spatials.actualRenderW - 4).
		Align(lipgloss.Center).
		Render(innerHeaderRender)

	footer := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(cyan.GetForeground()).
		Align(lipgloss.Center).
		Margin(1, 0, 0, 1).
		Width(m.spatials.actualRenderW - 4).
		Render(m.content.footerMessage)

	// this is the best render height across 5 different terminal emulators!
	h := 30
	if m.spatials.actualRenderH > 33 {
		h = 33
	}

	mainBorder := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(cyan.GetForeground()).
		Width(m.spatials.actualRenderW).
		Height(h-(m.spatials.innerH%m.itemsDispCount)).
		Render(header, midSectionJoin, footer)

	finalRender := lipgloss.Place(
		m.termWidth, m.termHeight,
		lipgloss.Center, lipgloss.Center,
		mainBorder,
	)

	return finalRender
}

func renderMidSections(spatials dimensions, itemCount int, listContent, deetContent string) string {
	listSection := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(cyan.GetForeground()).
		Margin(1, 1, 0, 1).
		Padding(0, 1).
		Width(spatials.listSectionW).
		Height(spatials.innerH - (spatials.innerH % itemCount)).
		Render(listContent)

	bookSection := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(cyan.GetForeground()).
		Padding(0, 1, 0, 1).
		Width(spatials.bookDeetsW).
		Height(spatials.innerH - (spatials.innerH % itemCount)).
		Render(deetContent)

	return lipgloss.JoinHorizontal(lipgloss.Center, listSection, bookSection)
}
