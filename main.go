package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

var (
	term    = termenv.EnvColorProfile()
	keyword = makeFgStyle("211")
	subtle  = makeFgStyle("241")
	dot     = colorTitle(" • ", "236")
)

const content = `
# About us
We are a small group of hoomans who believe in the betterment of the hoomans around us through arts & technology. We also happened to build humane tech products.

Too much blah blah blah? We are a product studio based out of Kochi that convergence of brand, behavioural psychology, digital and art. We hate building softwares, and we build solutions.

WE'RE BASED IN KOCHI, BUT WORK WITH COMPANIES WORLDWIDE.WANT TO COLLABORATE? LET'S CHAT

`

const content1 = `
# Location

Panampilly Nagar Avenue, Kanayannur, Kerala IN
`

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func main() {
	err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initialModel() model {
	const width = 50

	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
	vp1 := viewport.New(width, 12)
	vp1.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)

	if err != nil {
		return model{}
	}

	str, err := renderer.Render(content)
	if err != nil {
		return model{}
	}

	vp.SetContent(str)

	str1, err1 := renderer.Render(content1)
	if err1 != nil {
		return model{}
	}

	vp1.SetContent(str1)

	columns := []table.Column{
		{Title: "Rank", Width: 14},
		{Title: "Name", Width: 20},
		{Title: "Role", Width: 20},
		{Title: "Fans", Width: 14},
	}

	rows := []table.Row{
		{"1", "Abid", "Co Founder", "37,274,000"},
		{"2", "Lazim", "Co FOunder", "32,065,760"},
		{"3", "Adil", "Product Engineer", "28,516,904"},
		{"4", "Anagha", "Product Designer", "22,478,116"},
		{"5", "Jesswin", "Product Engineer", "22,429,800"},
		{"6", "Jobin", "Product Engineer", "22,085,140"},
		{"7", "Suvarnesh", "Product Engineer", "21,750,020"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	l := lists{0, false, false, false, false, 0}
	items := []list.Item{
		item{title: "Abid", desc: "Co Founder"},
		item{title: "Lazim", desc: "Co FOunder"},
		item{title: "Adil", desc: "Product Engineer"},
		item{title: "Anagha", desc: "Product Designer"},
		item{title: "Jesswin", desc: "Product Engineer"},
		item{title: "Jobin", desc: "Product Engineer"},
		item{title: "Suvarnesh", desc: "Product Engineer"},
	}
	i := list.New(items, list.NewDefaultDelegate(), 0, 0)
	i.Title = "hoomans co."

	return model{
		lists:     l,
		viewport:  vp,
		viewport1: vp1,
		table:     t,
		list:      i,
	}

}

type lists struct {
	Choice   int
	Chosen   bool
	Loaded   bool
	Quitting bool
	User     bool
	UserId   int
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	viewport  viewport.Model
	viewport1 viewport.Model
	lists     lists
	table     table.Model
	list      list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.lists.Quitting = true
			return m, tea.Quit
		}
	}
	if !m.lists.Chosen {
		return updateChoices(msg, m)
	}
	if m.lists.User {
		return updateUserChosen(msg, m)
	} else {
		return updateChosen(msg, m)
	}

}

func updateUserChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			m.lists.User = false
		}
	}
	return m, nil
}

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			m.lists.Chosen = false
		case "enter":
			if m.lists.Choice == 1 {
				m.lists.User = true
			}
		default:
			if m.lists.Choice == 0 {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.lists.Choice++
			if m.lists.Choice > 2 {
				m.lists.Choice = 2
			}
		case "up":
			m.lists.Choice--
			if m.lists.Choice < 0 {
				m.lists.Choice = 0
			}
		case "enter":
			m.lists.Chosen = true
			fmt.Print(m.lists.Choice)
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	var s string
	if m.lists.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.lists.Chosen {
		s = choicesView(m)
	} else {
		if m.lists.User {
			s = chosenUserView(m)
		} else {
			s = chosenView(m)
		}
	}
	return indent.String("\n"+s+"\n\n", 2)

}

func chosenView(m model) string {

	var msg string

	switch m.lists.Choice {
	case 0:
		msg = m.viewport.View()
	case 1:
		msg = baseStyle.Render(m.table.View()) + "\n"
	case 2:
		msg = m.viewport1.View()
	}

	return msg
}

func chosenUserView(m model) string {
	var msg string

	switch m.table.SelectedRow()[0] {
	case "1":
		msg = fmt.Sprintf("This is %s", m.table.SelectedRow()[1])
	case "2":
		msg = fmt.Sprintf("This is %s", m.table.SelectedRow()[1])
	case "3":
		msg = fmt.Sprintf("This is %s", m.table.SelectedRow()[1])
	case "4":
		msg = fmt.Sprintf("This is %s", m.table.SelectedRow()[1])
	case "5":
		msg = fmt.Sprintf("This is %s", m.table.SelectedRow()[1])
	case "6":
		msg = fmt.Sprintf("This is %s", m.table.SelectedRow()[1])
	case "7":
		msg = fmt.Sprintf("This is %s", m.table.SelectedRow()[1])
	}

	return msg
}

func choicesView(m model) string {
	c := m.lists.Choice
	myFigure := figure.NewFigure("hoomans co.", "", true)
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#d1ff69"))

	tpl := style.Render(myFigure.String())
	tpl += "\n\n"
	tpl += "%s\n\n"
	tpl += subtle("up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		checkbox("About", "see what we do when we are not discussing how to change the world.", c == 0),
		checkbox("Team", "check out our crazy team and connect with us.", c == 1),
		checkbox("Location", "checkout our workspace and sync with us.", c == 2),
	)

	return fmt.Sprintf(tpl, choices)
}

func checkbox(title string, description string, checked bool) string {
	if checked {
		checkValue := colorTitle("\n"+"| "+title, "#d1ff69")
		checkValue += colorDes("\n"+"| "+description, "#d1ff69")
		return checkValue
	}
	return fmt.Sprintf("\n| %s \n| %s", title, description)
}

// Utils

// Color a string's foreground with the given value.
func colorTitle(val, color string) string {
	return termenv.String(val).Bold().Foreground(term.Color(color)).String()
}

func colorDes(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
func makeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(term.Color(fg)).
		Background(term.Color(bg)).
		Styled
}

// Generate a blend of colors.
func makeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, colorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
