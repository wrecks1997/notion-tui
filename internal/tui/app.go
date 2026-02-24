package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/curtisepp/notion-tui/internal/cache"
	"github.com/curtisepp/notion-tui/internal/notion"
	tuiList "github.com/curtisepp/notion-tui/internal/tui/list"
	"github.com/curtisepp/notion-tui/internal/tui/search"
	"github.com/curtisepp/notion-tui/internal/tui/viewer"
)

type view int

const (
	listView view = iota
	viewerView
	searchView
)

// Messages

type workspaceLoadedMsg struct {
	results []notion.SearchResult
	err     error
}

type pageLoadedMsg struct {
	title    string
	markdown string
	err      error
}

// App is the root model.
type App struct {
	client  *notion.Client
	cache   *cache.Cache
	current view

	list    tuiList.Model
	viewer  viewer.Model
	search  search.Model
	spinner spinner.Model

	loading    bool
	err        error
	width      int
	height     int
}

func NewApp(client *notion.Client) App {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	c := cache.New(5 * time.Minute)

	return App{
		client:  client,
		cache:   c,
		current: listView,
		list:    tuiList.New(),
		viewer:  viewer.New(),
		search:  search.New(client),
		spinner: s,
		loading: true,
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.spinner.Tick,
		a.loadWorkspace(),
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle error dismiss
		if a.err != nil {
			a.err = nil
			return a, nil
		}

		// Global quit
		if key.Matches(msg, Keys.Quit) && a.current == listView && !a.list.List.SettingFilter() {
			return a, tea.Quit
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case workspaceLoadedMsg:
		a.loading = false
		if msg.err != nil {
			a.err = msg.err
			return a, nil
		}
		a.list.SetItems(msg.results)
		return a, nil

	case pageLoadedMsg:
		a.loading = false
		if msg.err != nil {
			a.err = msg.err
			return a, nil
		}
		a.viewer.SetContent(msg.title, msg.markdown, a.width)
		a.current = viewerView
		return a, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		a.spinner, cmd = a.spinner.Update(msg)
		return a, cmd
	}

	// Delegate to current view
	switch a.current {
	case listView:
		return a.updateList(msg)
	case viewerView:
		return a.updateViewer(msg)
	case searchView:
		return a.updateSearch(msg)
	}

	return a, nil
}

func (a App) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if a.list.List.SettingFilter() {
			// Let the list handle filter input
			break
		}

		switch {
		case key.Matches(msg, Keys.Search):
			a.current = searchView
			a.search.Focus()
			return a, nil

		case key.Matches(msg, Keys.Enter):
			if item := a.list.SelectedItem(); item != nil {
				a.loading = true
				return a, tea.Batch(
					a.spinner.Tick,
					a.loadPage(item.Result.ID, item.Title()),
				)
			}
		}
	}

	var cmd tea.Cmd
	a.list, cmd = a.list.Update(msg)
	return a, cmd
}

func (a App) updateViewer(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Back):
			a.current = listView
			return a, nil
		case key.Matches(msg, Keys.Quit):
			return a, tea.Quit
		}
	}

	var cmd tea.Cmd
	a.viewer, cmd = a.viewer.Update(msg)
	return a, cmd
}

func (a App) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Back):
			a.current = listView
			return a, nil
		case key.Matches(msg, Keys.Enter):
			if item := a.search.SelectedItem(); item != nil {
				a.loading = true
				a.current = listView // Temporarily show list while loading
				return a, tea.Batch(
					a.spinner.Tick,
					a.loadPage(item.Result.ID, item.Title()),
				)
			}
		}
	}

	var cmd tea.Cmd
	a.search, cmd = a.search.Update(msg)
	return a, cmd
}

func (a App) View() string {
	if a.err != nil {
		return AppStyle.Render(ErrorStyle.Render("Error: " + a.err.Error() + "\n\nPress any key to continue"))
	}

	if a.loading {
		return AppStyle.Render(a.spinner.View() + " Loading...")
	}

	var content string
	switch a.current {
	case listView:
		content = a.list.View()
	case viewerView:
		content = a.viewer.View()
	case searchView:
		content = a.search.View()
	}

	help := HelpStyle.Render("  /: search · enter: open · esc: back · q: quit")

	if a.current == viewerView {
		return AppStyle.Render(content)
	}

	return AppStyle.Render(lipgloss.JoinVertical(lipgloss.Left, content, help))
}

// Commands

func (a App) loadWorkspace() tea.Cmd {
	return func() tea.Msg {
		if cached, ok := a.cache.Get("workspace"); ok {
			return workspaceLoadedMsg{results: cached.([]notion.SearchResult)}
		}

		resp, err := a.client.Search("")
		if err != nil {
			return workspaceLoadedMsg{err: err}
		}
		a.cache.Set("workspace", resp.Results)
		return workspaceLoadedMsg{results: resp.Results}
	}
}

func (a App) loadPage(pageID, title string) tea.Cmd {
	return func() tea.Msg {
		cacheKey := "page:" + pageID
		if cached, ok := a.cache.Get(cacheKey); ok {
			md := cached.(string)
			return pageLoadedMsg{title: title, markdown: md}
		}

		blocks, err := a.client.GetBlockChildren(pageID)
		if err != nil {
			return pageLoadedMsg{err: err}
		}

		md := notion.BlocksToMarkdown(blocks)
		a.cache.Set(cacheKey, md)
		return pageLoadedMsg{title: title, markdown: md}
	}
}
