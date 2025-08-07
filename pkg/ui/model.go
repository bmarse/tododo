package ui

import (
	"log"
	"strings"
	"time"

	tl "github.com/bmarse/tododo/pkg/todo"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

type tickMsg struct{}

func tickCmd() tea.Msg {
	time.Sleep(1 * time.Second)
	return tickMsg{}
}

// Run is a blocking function that starts bubbletea.
func Run() error {
	m, err := InitialModel()
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

// InitialModel initializes the bubbletea model with the initial state of the todo list.
func InitialModel() (*Model, error) {
	todolist, err := tl.LoadTodosFromMarkdown()
	if err != nil {
		return nil, err
	}
	ti := textinput.New()
	ti.Placeholder = "Whatcha want to do?"
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80

	s := spinner.New()
	s.Style = spinnerStyle
	s.Spinner = spinner.MiniDot

	return &Model{
		todo:    todolist,
		input:   ti,
		spinner: s,
	}, nil
}

type Model struct {
	todo    tl.Todo
	input   textinput.Model
	spinner spinner.Model
	adding  bool
	saving  bool
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.saving = false
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		if m.adding {
			switch msg.Type {
			case tea.KeyEnter:
				m.adding = false
				if strings.TrimSpace(m.input.Value()) == "" {
					break
				}
				if m.todo.Cursor == -1 {
					m.todo.AddTask(m.input.Value())
					m.todo.Cursor = len(m.todo.Tasks) - 1
				} else {
					m.todo.Tasks[m.todo.Cursor].UpdateText(m.input.Value())
				}

				m.input.SetValue("")
			case tea.KeyEsc:
				m.input.SetValue("")
				m.adding = false
			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				return m, cmd

			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			m.todo.ModulateCursor(1)
		case "k", "up":
			m.todo.ModulateCursor(-1)
		case "a":
			m.adding = true
			m.todo.Cursor = -1
		case "e":
			if m.todo.Cursor == -1 {
				return m, nil
			}
			m.adding = true
			m.input.SetValue(m.todo.Tasks[m.todo.Cursor].Text)
		case "d":
			m.todo.RemoveTodoAtIndex(m.todo.Cursor)
		case "t":
			m.todo.ToggleHidden()
			m.todo.ModulateCursor(0)
		case "w":
			m.saving = true
			if err := tl.SaveTodo(m.todo); err != nil {
				log.Fatal(err)
			}
			return m, tickCmd
		case " ", "x":
			if m.todo.Cursor == -1 {
				return m, nil
			}
			m.todo.Tasks[m.todo.Cursor].ToggleChecked()
			m.todo.ModulateCursor(0)
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.adding {
		return AddingUI(m.input.View())
	}

	return MainUI(&m.todo, m.saving, m.spinner.View())
}
