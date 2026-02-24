package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/curtisepp/notion-tui/internal/notion"
)

// PageItem implements list.Item for Notion search results.
type PageItem struct {
	Result notion.SearchResult
}

func (i PageItem) Title() string {
	title := extractTitle(i.Result)
	icon := ""
	if i.Result.Icon != nil && i.Result.Icon.Emoji != "" {
		icon = i.Result.Icon.Emoji + " "
	} else if i.Result.Object == "database" {
		icon = "🗃 "
	} else {
		icon = "📄 "
	}
	return icon + title
}

func (i PageItem) Description() string {
	edited := i.Result.LastEditedTime.Format("Jan 2, 2006")
	kind := i.Result.Object
	return fmt.Sprintf("%s · edited %s", kind, edited)
}

func (i PageItem) FilterValue() string {
	return extractTitle(i.Result)
}

func extractTitle(r notion.SearchResult) string {
	for _, prop := range r.Properties {
		if prop.Type == "title" && len(prop.Title) > 0 {
			return prop.Title[0].PlainText
		}
	}
	return "Untitled"
}

// Model wraps the bubbles list component.
type Model struct {
	List   list.Model
	width  int
	height int
}

func New() Model {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#6C5CE7")).
		BorderForeground(lipgloss.Color("#6C5CE7"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#A389F4")).
		BorderForeground(lipgloss.Color("#6C5CE7"))

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Notion Workspace"
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#6C5CE7")).
		Padding(0, 1)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.DisableQuitKeybindings()

	return Model{List: l}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.List.SetSize(msg.Width-4, msg.Height-4)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.List.View()
}

func (m *Model) SetItems(results []notion.SearchResult) {
	items := make([]list.Item, len(results))
	for i, r := range results {
		items[i] = PageItem{Result: r}
	}
	m.List.SetItems(items)
}

func (m Model) SelectedItem() *PageItem {
	item, ok := m.List.SelectedItem().(PageItem)
	if !ok {
		return nil
	}
	return &item
}
