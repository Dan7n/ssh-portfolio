package sshSession

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

const header = `	______  ________  ___   __   ___   __   __  __     ________  ______  ________  ________  ______
/_____/\/_______/\/__/\ /__/\/__/\ /__/\/_/\/_/\   /_______/\/_____/\/_______/\/_______/\/_____/\
\:::_ \ \::: _  \ \::\_\\  \ \::\_\\  \ \ \ \ \ \  \__.::._\/\::::_\/\::: _  \ \::: _  \ \:::__\/
 \:\ \ \ \::(_)  \ \:. '-\  \ \:. '-\  \ \:\_\ \ \    \::\ \  \:\/___/\::(_)  \ \::(_)  \ \:\ \  __
  \:\ \ \ \:: __  \ \:. _    \ \:. _    \ \::::_\/    _\::\ \__\_::._\:\:: __  \ \:: __  \ \:\ \/_/\
   \:\/.:| \:.\ \  \ \. \'-\  \ \. \'-\  \ \\::\ \   /__\::\__/\ /____\:\:.\ \  \ \:.\ \  \ \:\_\ \ \
    \____/_/\__\/\__\/\__\/ \__\/\__\/ \__\/ \__\/   \________\/ \_____\/\__\/\__\/\__\/\__\/\_____\/ `

func CreateHandler(sshSession ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := sshSession.Pty()

	tabs := []string{"About (a)", "Contact (c)"}
	defaultStyle := lipgloss.NewStyle().MaxWidth(lipgloss.Width(header)).Padding(1, 2)

	aboutTabContent := map[string]lipgloss.Style{}
	aboutTabContent["# Hey there! So cool that you're SSH'd in! 🎉"] = defaultStyle
	aboutTabContent["My name's Danny and this is a fun little project to play around with the Go programming language and make my little portfolio site a bit more interesting."] = defaultStyle
	aboutTabContent["I'm currently working at a company called Klarna, where I'm part of the team that's working on making klarna.com a better place for our customers"] = defaultStyle
	aboutTabContent["Please feel free to reach out to me on LinkedIn or via email if you have any questions or just want to chat. I'm always up for a good conversation!"] = defaultStyle
	aboutTabContent["I hope you have a great day and that you enjoy the rest of your time on my site. Take care!"] = defaultStyle

	contactTabContent := map[string]lipgloss.Style{}
	contactTabContent["LinkedIn: https://www.linkedin.com/in/danny-bergman/"] = defaultStyle
	contactTabContent["Github: https://github.com/Dan7n"] = defaultStyle
	contactTabContent["Email: danny95.nbl@gmail.com"] = defaultStyle

	tabContent := []map[string]lipgloss.Style{aboutTabContent, contactTabContent}

	m := model{
		Tabs:         tabs,
		TabContent:   tabContent,
		activeTab:    0,
		windowWidth:  pty.Window.Width,
		windowHeight: pty.Window.Height,
	}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

// Here is where the model, Update, and View functions are defined.
// The model is a struct that holds the terminal information and styles.
// The Update function handles messages and updates the model accordingly.
// The View function renders the terminal information to the screen.

// model is the main model for the SSH session.
type model struct {
	Tabs         []string
	TabContent   []map[string]lipgloss.Style
	activeTab    int
	windowWidth  int
	windowHeight int
}

func (mod model) Init() tea.Cmd {
	return nil
}

func (mod model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// if a user quits the program
		case "q", "ctrl+c":
			return mod, tea.Quit

		// if a user presses the "a" key to go to the "about" tab
		case "tab", "a":
			mod.activeTab = 0
			return mod, nil

		// if a user presses the "c" key to go to the "contact" tab
		case "c":
			mod.activeTab = 1
			return mod, nil
		}
	}
	return mod, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", "─", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2).AlignHorizontal(lipgloss.Left)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 3)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	// windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Left)
)

func (mod model) View() string {
	doc := strings.Builder{}
	var renderedTabs []string
	doc.WriteString(docStyle.Render(header))
	doc.WriteString("\n\n")

	for i, t := range mod.Tabs {
		var style lipgloss.Style
		style.Width(lipgloss.Width(header))

		isFirst, isLast, isActive := i == 0, i == len(mod.Tabs)-1, i == mod.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}

		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		if isActive {
			style.PaddingRight(lipgloss.Width(header) - (lipgloss.Width(t) * len(mod.Tabs)) - 1)
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	activeTabContent := mod.TabContent[mod.activeTab]
	// todo: keys are not sorted, so the order of the content will be random - fix this
	for line, styles := range activeTabContent {

		doc.WriteString(styles.Render(line))
		doc.WriteString("\n")

	}

	// doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(mod.TabContent[mod.activeTab]))
	return docStyle.Render(doc.String())
}
