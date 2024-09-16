package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"shelve/git"
	"strings"

	"github.com/charmbracelet/bubbles/list"
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

const listHeight = 10

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("19"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(2)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(string(i)))
}

type model struct {
	list     list.Model
	selected string
	output   string
	applied  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		s := msg.String()
		switch s {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter", "L", "a":
			m.selected = string(m.list.SelectedItem().(item))

			if o, err := git.ApplyStash(m.list.Index()); err != nil {
				m.output = err.Error()
			} else {
				m.output = o
			}
			m.applied = true

			if s != "a" {
				return m, tea.Quit
			}
			fallthrough

		case "d":
			if _, err := git.DropStash(m.list.Index()); err != nil {
				m.output = err.Error()
				return m, tea.Quit
			}
			m.list.RemoveItem(m.list.Index())

			if s == "a" {
				return m, tea.Quit
			}

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.applied {
		return quitTextStyle.Render(strings.Join([]string{
			fmt.Sprintf("Applied: %s", m.selected),
			"\n",
			m.output,
		}, "\n"))
	}

	return "\n" + m.list.View()
}

func selectStash(cmd *cobra.Command, args []string) {
	const defaultWidth = 20

	stashList, err := git.GetStashList()

	if err != nil {
		log.Fatal(err)
	}
	items := make([]list.Item, 0, len(stashList))

	for _, stashItem := range stashList {
		items = append(items, item(stashItem))
	}
	list := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	list.SetShowTitle(false)
	list.SetShowHelp(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.Styles.PaginationStyle = paginationStyle

	m := model{
		list: list,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
