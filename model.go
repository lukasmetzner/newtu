package main

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// -- Styles --

var borderStyle = lipgloss.NewStyle().
	Align(lipgloss.Left, lipgloss.Top).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("69"))

// -- Tea messages --

type articlesMsg struct {
	articles []Article
}

type fetchErrMsg struct {
	err error
}

// -- Model --

type model struct {
	table  table.Model
	input  textinput.Model
	config *Config
	db     *sql.DB

	allArticles       []Article
	displayedArticles []Article
	searchMode        bool
}

func newModel(cfg *Config, db *sql.DB) model {
	input := textinput.New()
	input.Focus()

	columns := []table.Column{
		{Title: "", Width: 8},
		{Title: "Source", Width: 12},
		{Title: "Date", Width: 16},
		{Title: "Title", Width: 96},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(30),
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

	return model{
		table:  t,
		input:  input,
		config: cfg,
		db:     db,
	}
}

// -- Bubble Tea interface --

func (m model) Init() tea.Cmd {
	return refreshArticles(m.config, m.db)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height - 3)
		m.table.SetWidth(msg.Width)
		borderStyle = borderStyle.Width(msg.Width)
		return m, nil

	case articlesMsg:
		m.allArticles = msg.articles
		m.setDisplayed(msg.articles)
		return m, scheduleRefresh(m.config, m.db)

	case fetchErrMsg:
		// TODO: show error in status bar instead of silently dropping
		return m, scheduleRefresh(m.config, m.db)

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return borderStyle.Render(m.input.View(), "\n", m.table.View())
}

// -- Key handling --

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "enter":
		return m.handleEnter()

	case "/":
		m.searchMode = true

	case "esc":
		if m.searchMode {
			m.searchMode = false
			m.setDisplayed(m.allArticles)
		}
		m.input.Reset()
	}

	if m.searchMode {
		m.input, _ = m.input.Update(msg)
		m.applySearch()
	} else {
		m.input, _ = m.input.Update(msg)
		m.validateInput()
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	if val := m.input.Value(); val != "" {
		if jumpIdx, err := strconv.Atoi(val); err == nil {
			m.input.Reset()
			if jumpIdx >= 1 && jumpIdx <= len(m.displayedArticles) {
				openURL(m.displayedArticles[jumpIdx-1].Link) //nolint:errcheck
			}
			return m, nil
		}
	}

	cursor := m.table.Cursor()
	if cursor >= 0 && cursor < len(m.displayedArticles) {
		openURL(m.displayedArticles[cursor].Link) //nolint:errcheck
	}
	return m, nil
}

// -- Search & input helpers --

func (m *model) applySearch() {
	query := strings.TrimPrefix(m.input.Value(), "/")
	m.displayedArticles = FilterByTitle(m.allArticles, query)
	m.table.SetRows(ArticlesToRows(m.displayedArticles))
}

func (m *model) validateInput() {
	val := m.input.Value()
	if val == "" {
		return
	}
	if strings.HasPrefix(val, "/") {
		return
	}
	if _, err := strconv.Atoi(val); err != nil {
		m.input.Reset()
	}
}

func (m *model) setDisplayed(articles []Article) {
	m.displayedArticles = articles
	m.table.SetRows(ArticlesToRows(articles))
}

// -- Commands --

func refreshArticles(cfg *Config, db *sql.DB) tea.Cmd {
	return func() tea.Msg {
		cached, err := LoadArticles(db)
		if err != nil {
			return fetchErrMsg{err: err}
		}

		fresh, err := FetchArticles(cfg.RSSFeeds)
		if err != nil {
			// Return cached articles even if fetch fails
			SortByDateDesc(cached)
			return articlesMsg{articles: cached}
		}

		if err := SaveArticles(db, fresh); err != nil {
			return fetchErrMsg{err: err}
		}

		merged := MergeArticles(fresh, cached)
		SortByDateDesc(merged)
		return articlesMsg{articles: merged}
	}
}

func scheduleRefresh(cfg *Config, db *sql.DB) tea.Cmd {
	return tea.Tick(15*time.Minute, func(time.Time) tea.Msg {
		return refreshArticles(cfg, db)()
	})
}
