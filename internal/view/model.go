package view

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	content  string
	ready    bool
	viewport viewport.Model
	p        *tea.Program
}

func New() *model {
	return &model{}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}
	case contentUpdateMsg:
		m.viewport.SetContent(m.content)
	}

	return m, nil
}

func (m *model) View() string {
	if !m.ready {
		return "\n> Initializing..."
	}

	return m.viewport.View()
}

func (m *model) Run() error {
	m.p = tea.NewProgram(m, tea.WithAltScreen())

	_, err := m.p.Run()

	return err
}
