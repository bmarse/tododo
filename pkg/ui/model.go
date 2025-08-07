package ui

import (
	"log"
	"strings"
	"time"

	tl "github.com/bmarse/tododo/pkg/todo"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

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
	return &Model{
		Todo:  todolist,
		Input: ti,
	}, nil
}

type Model struct {
	Todo   tl.Todo
	Input  textinput.Model
	Adding bool
	Saving bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.Saving = false
	case tea.KeyMsg:
		if m.Adding {
			switch msg.Type {
			case tea.KeyEnter:
				m.Adding = false
				if strings.TrimSpace(m.Input.Value()) == "" {
					break
				}
				if m.Todo.Cursor == -1 {
					m.Todo.AddTask(m.Input.Value())
					m.Todo.Cursor = len(m.Todo.Tasks) - 1
				} else {
					m.Todo.Tasks[m.Todo.Cursor].UpdateText(m.Input.Value())
				}

				m.Input.SetValue("")
			case tea.KeyEsc:
				m.Input.SetValue("")
				m.Adding = false
			default:
				var cmd tea.Cmd
				m.Input, cmd = m.Input.Update(msg)
				return m, cmd

			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			m.Todo.ModulateCursor(1)
		case "k", "up":
			m.Todo.ModulateCursor(-1)
		case "a":
			m.Adding = true
			m.Todo.Cursor = -1
		case "e":
			if m.Todo.Cursor == -1 {
				return m, nil
			}
			m.Adding = true
			m.Input.SetValue(m.Todo.Tasks[m.Todo.Cursor].Text)
		case "d":
			m.Todo.RemoveTodoAtIndex(m.Todo.Cursor)
		case "t":
			m.Todo.ToggleHidden()
			m.Todo.ModulateCursor(0)
		case "w":
			m.Saving = true
			if err := tl.SaveTodo(m.Todo); err != nil {
				log.Fatal(err)
			}
			return m, tickCmd
		case " ", "x":
			if m.Todo.Cursor == -1 {
				return m, nil
			}
			m.Todo.Tasks[m.Todo.Cursor].ToggleChecked()
			m.Todo.ModulateCursor(0)
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.Adding {
		return AddingUI(m.Input.View())
	}

	return MainUI(&m.Todo, m.Saving)
}
