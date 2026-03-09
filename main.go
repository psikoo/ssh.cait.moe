package main

import (
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/muesli/termenv"
)

//
// Structs, types and consts
//

type tickMsg time.Time
type page int

type model struct {
	currentPage    page
	menuCursor     int
	projectCursor  int
	expCursor      int
	width          int
	height         int
}

type Project struct {
	Name string
	Desc string
	Link string
}

//
// String data
//

var menuItems = []string {"About", "Projects", "Contact"}

var asciiLogoLines = []string {
	` ▄████▄   ▄▄▄       ██▓▄▄▄█████▓`,
	`▒██▀ ▀█  ▒████▄    ▓██▒▓  ██▒ ▓▒`,
	`▒▓█    ▄ ▒██  ▀█▄  ▒██▒▒ ▓██░ ▒░`,
	`▒▓▓▄ ▄██▒░██▄▄▄▄██ ░██░░ ▓██▓ ░ `,
	`▒ ▓███▀ ░ ▓█   ▓██▒░██░  ▒██▒ ░ `,
	`░ ░▒ ▒  ░ ▒▒   ▓▒█░░▓    ▒ ░░   `,
	`  ░  ▒     ▒   ▒▒ ░ ▒ ░    ░    `,
	`░    ░     ░   ▒    ▒ ░  ░      `,
	`  ░            ░  ░ ░           `,
}

var aboutContent = `
Hii! name is Caitlyn, im a software developer. 

I mainly do web dev and java apps, but I also enjoy making small js scripts. I have more one off projects like this one and some others listed on my site.

You can always contact me through discord :3
`

var projects = []Project{
	{
		Name: "SSH Portfolio",
		Desc: "What you are interacting with right now!",
		Link: "https://github.com/psikoo/ssh.cait.moe",
	},
	{
		Name: "For more projects check out my website",
		Desc: "Here I keep an up to date list of my projects",
		Link: "https://www.cait.moe",
	},
}

//
// Theme
//

const (
	menuPage page = iota
	aboutPage
	projectsPage
	contactPage
)

var (
	// Colors
	white      = lipgloss.Color("#D8DEE9") // primary
	purple     = lipgloss.Color("#aa8cd7") // accent
	gray       = lipgloss.Color("#4C566A") // muted
	// Styles
	primaryStyle  = lipgloss.NewStyle().Foreground(white)
	accentStyle   = lipgloss.NewStyle().Foreground(purple).Bold(true)
	mutedStyle    = lipgloss.NewStyle().Foreground(gray)
)

//
// Utility functions
//

func tickCmd() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func initialModel() model {
	return model{
		currentPage:    menuPage,
		menuCursor:     0,
		projectCursor:  0,
		expCursor:      0,
		width:          60,
		height:         24,
	}
}

func clickableLink(label, url string) string {
	return "\x1b]8;;" + url + "\x1b\\" + label + "\x1b]8;;\x1b\\"
}

func renderAsciiLogo(width int) string {
	// Join the lines
	logoRaw := strings.Join(asciiLogoLines, "\n")
	// Style and center
	logoStyle := lipgloss.NewStyle().Foreground(purple).Bold(true).Width(width).Align(lipgloss.Center)
	return logoStyle.Render(logoRaw)
}

func (m model) renderMenu() string {
	var result strings.Builder
	// Render Tittle
	result.WriteString(renderAsciiLogo(int(float64(m.width) * 0.5)))
	result.WriteString("\n")
	result.WriteString("\n")
	result.WriteString("\n")
	// Render Content
	for i, item := range menuItems {
		// Show arrow in front of the selected item
		cursor := "  "
		if m.menuCursor == i { cursor = "→ " }
		line := cursor + item
		// Style the tittle of the selected item
		if m.menuCursor == i {
			result.WriteString(accentStyle.Render(line))
		} else {
			result.WriteString(primaryStyle.Render(line))
		}
		result.WriteString("\n")
	}
	result.WriteString("\n")
	result.WriteString("\n")
	// Render menu
	result.WriteString(mutedStyle.Render("↑/↓: navigate • enter: select • q: quit"))
	return result.String()
}

func (m model) renderAbout() string {
	var result strings.Builder
	// Render Tittle
	result.WriteString(accentStyle.Render("━━━ About Me ━━━"))
	result.WriteString("\n")
	result.WriteString("\n")
	// Render Content
	var contentStyle = lipgloss.NewStyle().Foreground(white).Width(int(float64(m.width) * 0.5)).Align(lipgloss.Left)
	result.WriteString(contentStyle.Render(aboutContent))
	result.WriteString("\n")
	result.WriteString("\n")
	// Render menu
	result.WriteString(mutedStyle.Render("esc: back to menu"))
	return result.String()
}

func (m model) renderProjects() string {
	var result strings.Builder
	// Render Tittle
	result.WriteString(accentStyle.Render("━━━ Projects ━━━"))
	result.WriteString("\n")
	result.WriteString("\n")
	result.WriteString("\n")
	// Render Content
	for i, project := range projects {
		// Show arrow in front of the selected project
		cursor := "  "
		if m.projectCursor == i { cursor = "→ " }
		name := cursor + project.Name
		// Style the tittle of the selected project
		if m.projectCursor == i {
			result.WriteString(accentStyle.Render(name))
		} else {
			result.WriteString(primaryStyle.Render(name))
		}
		result.WriteString("\n")
		// Expand selected project
		if m.projectCursor == i {
			result.WriteString("   ")
			result.WriteString(mutedStyle.Render(project.Desc))
			result.WriteString("\n   ")
			result.WriteString(accentStyle.Render(clickableLink(project.Link, project.Link)))
			result.WriteString("\n")
		}
		result.WriteString("\n")
	}
	result.WriteString("\n")
	// Render menu
	result.WriteString(mutedStyle.Render("↑/↓: navigate • ctrl+click: open link • esc: back to menu"))
	return result.String()
}

func (m model) renderContact() string {
	var result strings.Builder
	// Render Tittle
	result.WriteString(accentStyle.Render("━━━ Contact ━━━"))
	result.WriteString("\n")
	result.WriteString("\n")
	result.WriteString("\n")
	// Render Content
	result.WriteString(primaryStyle.Render(" - Website  "))
	result.WriteString(accentStyle.Render(clickableLink("https://www.cait.moe", "https://www.cait.moe")))
	result.WriteString("\n")
	result.WriteString(primaryStyle.Render(" - Github   "))
	result.WriteString(accentStyle.Render(clickableLink("https://url.cait.moe/?u=github", "https://url.cait.moe/?u=github")))
	result.WriteString("\n")
	result.WriteString(primaryStyle.Render(" - Discord  "))
	result.WriteString(accentStyle.Render(clickableLink("https://url.cait.moe/?u=discord", "https://url.cait.moe/?u=discord")))
	result.WriteString("\n")
	result.WriteString("\n")
	result.WriteString("\n")
	// Render menu
	result.WriteString(mutedStyle.Render("ctrl+click: open link • esc: back to menu"))
	return result.String()
}

//
// Renderer
//

func (m model) View() string {
	var content string
	// Render current page
	switch m.currentPage {
		case menuPage:     content = m.renderMenu()
		case aboutPage:    content = m.renderAbout()
		case projectsPage: content = m.renderProjects()
		case contactPage:  content = m.renderContact()
	}
	// Center content
	contentBox := lipgloss.NewStyle().Width(m.width).PaddingLeft(int(float64(m.width) * 0.25)).Align(lipgloss.Left).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, contentBox)
}

//
// Controls
//

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		// Handle tick
		case tickMsg:
			return m, tickCmd()
		// Handle resize
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
		// On key
		case tea.KeyMsg:
			switch msg.String() {
				// Handle quit
				case "ctrl+c", "q":
					return m, tea.Quit
				// Handle go back
				case "esc", "backspace":
					if m.currentPage != menuPage {
						m.currentPage = menuPage
					}
					return m, tickCmd()
				// Handle up
				case "up", "k":
					switch m.currentPage {
						// Menu
						case menuPage: if m.menuCursor > 0 { m.menuCursor-- }
						// Project
						case projectsPage: if m.projectCursor > 0 { m.projectCursor-- }
					}
					return m, nil
				// Handle down
				case "down", "j":
					switch m.currentPage {
						// Menu
						case menuPage:if m.menuCursor < len(menuItems)-1 { m.menuCursor++ }
						// Project
						case projectsPage:if m.projectCursor < len(projects)-1 { m.projectCursor++ }
					}
					return m, nil
				// Handle enter
				case "enter", " ":
					if m.currentPage == menuPage {
						switch m.menuCursor {
							case 0: m.currentPage = aboutPage
							case 1: m.currentPage = projectsPage
							case 2: m.currentPage = contactPage
						}
					}
					return m, nil
			}
	}
	return m, nil
}

//
// Entry
//

func main() {
	// Get port from the first argument
	port := "2222"
	if len(os.Args) > 1 { port = os.Args[1] }
	// Force to use ANSI256
	lipgloss.SetColorProfile(termenv.ANSI256)
	// Create TUI for each ssh connection
	teaHandler := func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		return initialModel(), []tea.ProgramOption{tea.WithAltScreen()}
	}
	// Start server and catch errors
	s, err := wish.NewServer(
		wish.WithAddress("0.0.0.0:"+port),
		wish.WithHostKeyPath(".ssh/host_ed25519"),
		wish.WithMiddleware(bubbletea.Middleware(teaHandler)),
	)
	if err != nil { log.Fatal(err) }
	log.Printf("Listening on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}
