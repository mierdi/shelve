package cmd

import (
	"fmt"
	"log"
	"os"
	"shelve/git"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var cmdStashSelect = &cobra.Command{
	Use:   "list",
	Short: "l",
	Run:   selectStash,
}

func init() {
	MainCommand.AddCommand(cmdStashSelect)
}

var (
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("19"))
	quitTextStyle     = lipgloss.NewStyle().Margin(0, 0, 1, 2)
)

type model struct {
	stashList []string
	cursor    int
	selected  string
	output    string
	applied   bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch s {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.stashList)-1 {
				m.cursor++
			}
		case "enter", "L", "a":
			if o, err := git.ApplyStash(m.cursor); err != nil {
				m.output = err.Error()
			} else {
				m.output = o
			}
			m.applied = true
			m.selected = m.stashList[m.cursor]

			if s != "a" {
				return m, tea.Quit
			}
			fallthrough

		case "d":
			if _, err := git.DropStash(m.cursor); err != nil {
				m.output = err.Error()
				return m, tea.Quit
			}
			stashList, err := git.GetStashList()

			if err != nil {
				m.output = err.Error()
				return m, tea.Quit
			}
			m.stashList = stashList

			if s == "a" {
				return m, tea.Quit
			}

			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.applied {
		return quitTextStyle.Render(strings.Join([]string{
			fmt.Sprintf("Applied: %s", m.selected),
			m.output,
		}, "\n"))
	}

	if len(m.stashList) == 0 {
		return quitTextStyle.Render("No items...")
	}
	lines := make([]string, 0, len(m.stashList))

	for i, s := range m.stashList {
		line := s

		if m.cursor == i {
			line = selectedItemStyle.Render(line)
		}
		lines = append(lines, line)
	}
	lines = append(lines, "\n")

	return strings.Join(lines, "\n")
}

func selectStash(cmd *cobra.Command, args []string) {
	stashList, err := git.GetStashList()

	if err != nil {
		log.Fatal(err)
	}

	if _, err := tea.NewProgram(model{stashList: stashList}).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
