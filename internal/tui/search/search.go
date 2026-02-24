package search

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/curtisepp/notion-tui/internal/notion"
	tuiList "github.com/curtisepp/notion-tui/internal/tui/list"
)

const debounceInterval = 300 * time.Millisecond

// Messages

type SearchResultsMsg struct {
	Results []notion.SearchResult
	Err     error
}

type debounceTickMsg struct{}

// Model holds the search view state.
type Model struct {
	Input      textinput.Model
	Results    list.Model
	client     *notion.Client
	lastQuery  string
	debouncing bool
	width      int
	height     int
}

func New(client *notion.Client) Model {
	ti := textinput.New()
	ti.Placeholder = "Search Notion..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6C5CE7"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5"))

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#6C5CE7")).
		BorderForeground(lipgloss.Color("#6C5CE7"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#A389F4")).
		BorderForeground(lipgloss.Color("#6C5CE7"))

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	l.SetShowHelp(false)

	return Model{
		Input:   ti,
		Results: l,
		client:  client,
	}
}

func (m *Model) Focus() {
	m.Input.Focus()
	m.Input.SetValue("")
	m.lastQuery = ""
	m.Results.SetItems([]list.Item{})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Input.Width = msg.Width - 8
		m.Results.SetSize(msg.Width-4, msg.Height-8)

	case debounceTickMsg:
		m.debouncing = false
		query := m.Input.Value()
		if query != "" && query != m.lastQuery {
			m.lastQuery = query
			return m, m.doSearch(query)
		}

	case SearchResultsMsg:
		if msg.Err != nil {
			return m, nil
		}
		items := make([]list.Item, len(msg.Results))
		for i, r := range msg.Results {
			items[i] = tuiList.PageItem{Result: r}
		}
		m.Results.SetItems(items)
		return m, nil
	}

	var cmds []tea.Cmd

	prevValue := m.Input.Value()
	var inputCmd tea.Cmd
	m.Input, inputCmd = m.Input.Update(msg)
	cmds = append(cmds, inputCmd)

	// Start debounce if input changed
	if m.Input.Value() != prevValue && !m.debouncing {
		m.debouncing = true
		cmds = append(cmds, tea.Tick(debounceInterval, func(time.Time) tea.Msg {
			return debounceTickMsg{}
		}))
	}

	var listCmd tea.Cmd
	m.Results, listCmd = m.Results.Update(msg)
	cmds = append(cmds, listCmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#6C5CE7")).
		Padding(0, 1).
		Render("Search")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("  esc: back · enter: open")

	return header + "\n\n  " + m.Input.View() + "\n\n" + m.Results.View() + "\n" + help
}

func (m Model) SelectedItem() *tuiList.PageItem {
	item, ok := m.Results.SelectedItem().(tuiList.PageItem)
	if !ok {
		return nil
	}
	return &item
}

func (m Model) doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := m.client.Search(query)
		if err != nil {
			return SearchResultsMsg{Err: err}
		}
		return SearchResultsMsg{Results: results.Results}
	}
}
