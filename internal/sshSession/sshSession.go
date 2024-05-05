package sshSession

import (
	"sort"
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

const heartAscii = `   &&&&&&&   &&&&&&&   
  &&&&&&&&&&&&&&&&&&&& 
 &&&&&&&&&&&&&&&&&&&&& 
 &&&&&&&&&&&&&&&&&&&&& 
 &&&&&&&&&&&&&&&&&&&&& 
  &&&&&&&&&&&&&&&&&&&  
    &&&&&&&&&&&&&&&    
      &&&&&&&&&&&      
        &&&&&&&        
          &&&          `

const paddingInline = 2

func CreateHandler(sshSession ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := sshSession.Pty()
	const hiddenTab = ""

	tabs := []string{"About (a)", "Contact (c)", hiddenTab}
	defaultStyle := lipgloss.NewStyle().MaxWidth(lipgloss.Width(header)).Padding(1, paddingInline)

	aboutTabContent := map[int]tabContent{}
	aboutTabContent[0] = tabContent{style: defaultStyle.Copy().Bold(true).Foreground(lipgloss.Color("1")), content: "# Hey there! So cool that you're SSH'd in! 🚀"}
	aboutTabContent[1] = tabContent{style: defaultStyle, content: "My name's Danny and this is a fun little project to play around with the Go programming language\nand make my little portfolio site (https://dannyisaac.com) a bit more interesting."}
	aboutTabContent[2] = tabContent{style: defaultStyle, content: "I'm a fullstack software engineer currently working at a company called Klarna, where I'm \npart of a team that's working on making https://www.klarna.com a smoother and more enjoyable experience \nfor our +150M global users."}
	aboutTabContent[3] = tabContent{style: defaultStyle, content: "I really love what I do and I'm always looking for new challenges and ways to grow as an enginner. \nI'm also a musician and have been playing the piano professionally for over 15 years now."}
	aboutTabContent[4] = tabContent{style: defaultStyle, content: "With that said, thank you again for SSHing in and please feel free to reach out\nto me on LinkedIn or via email (press c) if you have any questions or just want to chat. \nI'm always up for a good conversation!"}
	aboutTabContent[5] = tabContent{style: defaultStyle.Copy().BorderBottom(true), content: "I hope you have a great day and that you enjoy the rest of your time on my site. Take care!"}
	aboutTabContent[6] = tabContent{style: defaultStyle.Copy().Foreground(lipgloss.Color("1")), content: heartAscii}

	contactTabContent := map[int]tabContent{}
	contactTabContent[0] = tabContent{style: defaultStyle.Copy().Bold(true).Foreground(lipgloss.Color("1")), content: "# Here's how you can reach me:"}
	contactTabContent[1] = tabContent{style: defaultStyle, content: "LinkedIn: https://www.linkedin.com/in/danny-isaac/"}
	contactTabContent[2] = tabContent{style: defaultStyle, content: "Github: https://github.com/Dan7n"}
	contactTabContent[3] = tabContent{style: defaultStyle, content: "Email: mailto://danny95.nbl@gmail.com"}

	content := []map[int]tabContent{aboutTabContent, contactTabContent}

	m := model{
		Tabs:         tabs,
		TabContent:   content,
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

type tabContent struct {
	style   lipgloss.Style
	content string
}

// model is the main model for the SSH session.
type model struct {
	Tabs         []string
	TabContent   []map[int]tabContent
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
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#fff"}
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

		// we do len(mod.Tabs)-2 to exclude the hidden tab because that's always the last tab
		isFirst, isLast, isActive := i == 0, i == len(mod.Tabs)-2, i == mod.activeTab
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

		if isActive {
			style.Bold(true)
		} else {
			style.Bold(false)
			style.Faint(true)
		}

		if t == "" {
			// create a hidden tab that's as wide as the remaining space with a bottom border
			style.Padding(0)
			style.PaddingRight(lipgloss.Width(header) - (lipgloss.Width(mod.Tabs[0]) * 4))
			style.BorderLeft(false)
			style.BorderRight(false)
			style.BorderTop(false)
			style.AlignVertical(lipgloss.Bottom)
			style.AlignHorizontal(lipgloss.Bottom)
			style.Height(2)
		}

		renderedTabs = append(renderedTabs, style.Render(t))

	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	activeTabContent := mod.TabContent[mod.activeTab]

	// convert the map keys to a slice so we can sort them
	slice := make([]int, 0, len(activeTabContent))
	for i := range activeTabContent {
		slice = append(slice, i)
	}
	sort.Slice(slice, func(i, j int) bool {
		return j < i
	})

	// now slice is sorted and we can render the content in the correct order
	for idx := range slice {
		style := activeTabContent[idx].style
		content := activeTabContent[idx].content

		doc.WriteString(style.Render(content))
		doc.WriteString("\n")
	}

	footerTxt := "Press `q` or `Ctrl+C` to quit."
	footerBorder := lipgloss.NormalBorder()
	footerBorder.Right = ""
	footerBorder.TopRight = ""
	footerBorder.BottomRight = ""
	footerBorder.Left = ""
	footerBorder.TopLeft = ""
	footerBorder.BottomLeft = ""
	footerBorder.Bottom = ""
	doc.WriteString(lipgloss.NewStyle().Padding(0, paddingInline).Border(footerBorder, true).MarginTop(2).Width(lipgloss.Width(header)).Render(footerTxt))

	// doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(mod.TabContent[mod.activeTab]))
	return docStyle.Render(doc.String())
}
