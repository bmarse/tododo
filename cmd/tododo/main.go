package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const todoFileName = ".tododo.md"

var (
	bold   lipgloss.Style = lipgloss.NewStyle().Bold(true)
	padded lipgloss.Style = lipgloss.NewStyle().Padding(2, 4)
	faint  lipgloss.Style = lipgloss.NewStyle().Faint(true)
	border lipgloss.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(true).
		Padding(1, 2, 1, 2)
)

type todo struct {
	task    string
	checked bool
}

type model struct {
	todos  []todo
	cursor int
	input  textinput.Model
	adding bool
	saving bool
}

type tickMsg struct{}

func tickCmd() tea.Msg {
	time.Sleep(1 * time.Second)
	return tickMsg{}
}

func initialModel() model {
	todos, err := LoadTodosFromMarkdown()
	if err != nil {
		log.Fatal("failed to load todos:", err)
	}
	ti := textinput.New()
	ti.Placeholder = "Whatcha want to do?"
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80
	return model{
		todos: todos,
		input: ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.saving = false
	case tea.KeyMsg:
		if m.adding {
			switch msg.Type {
			case tea.KeyEnter:
				if strings.TrimSpace(m.input.Value()) != "" {
					m.todos = append(m.todos, todo{task: m.input.Value()})
				}
				m.input.SetValue("")
				m.adding = false
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
			m.ModulateCursor(1)
		case "k", "up":
			m.ModulateCursor(-1)
		case "a":
			m.adding = true
		case "d":
			m.RemoveTodoAtIndex(m.cursor)
		case "w":
			m.saving = true
			if err := m.SaveTodo(); err != nil {
				log.Fatal(err)
			}
			return m, tickCmd
		case " ", "x":
			m.todos[m.cursor].checked = !m.todos[m.cursor].checked
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.adding {
		msg := bold.Render("Add task:")

		return padded.Render(fmt.Sprintf("%s\n%s\n\n(Enter to save, Esc to cancel)", msg, border.Render(m.input.View())))
	}
	s := randomMessage()
	s += "\n\n"
	tasks := ""
	if len(m.todos) == 0 {
		tasks += "Yippee! No tasks to do..."
	}
	for i, t := range m.todos {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = bold.Render(">")
		}
		checked := " " // not completed
		if t.checked {
			checked = bold.Render("x")
		}
		tasks += fmt.Sprintf("%s [%s]: %s\n", cursor, checked, t.task)
	}
	if m.saving {
		tasks += "\n Saving... \n"
	} else {
		tasks += "\n\n"
	}

	s += border.Render(tasks)
	s += "\n\n↑/↓: Move  a: Add  <space>: Toggle  d: Delete  w: Write  q: Quit"
	return padded.Render(s)
}

func (m *model) RemoveTodoAtIndex(index int) {
	if index < 0 || index >= len(m.todos) {
		return
	}

	m.todos = append(m.todos[:index], m.todos[index+1:]...)
}

func (m *model) ModulateCursor(amount int) {
	if amount > 0 {
		newPosition := m.cursor + amount
		for newPosition >= len(m.todos) {
			newPosition = newPosition % len(m.todos)
		}

		m.cursor = newPosition
		return
	}

	// the amount is negative
	newPosition := m.cursor + amount
	for newPosition < 0 {
		newPosition += len(m.todos)
	}

	m.cursor = newPosition
}

func (m *model) SaveTodo() error {
	// Save the current todo list to a file
	b := []byte{}
	for _, todo := range m.todos {
		check := " "
		if todo.checked {
			check = "X"
		}
		b = fmt.Appendf(b, "- [%s] %s\n", check, todo.task)
	}

	if len(b) > 0 {
		err := os.WriteFile(todoFileName, b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func LoadTodosFromMarkdown() ([]todo, error) {
	markdownContent, err := os.ReadFile(todoFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []todo{
				{
					task: "Add items to your todo list",
				},
			}, nil
		}
		return nil, err
	}

	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(markdownContent))

	var todos []todo
	checkRegex := regexp.MustCompile(`^\s*\[[ xX]\]\s*`)

	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		if n.Kind() == ast.KindListItem {
			var buf bytes.Buffer
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				if segmenter, ok := c.(interface{ Text(source []byte) []byte }); ok {
					buf.Write(segmenter.Text(markdownContent))
				}
			}
			itemText := strings.TrimSpace(buf.String())
			checked := false
			if m := checkRegex.FindString(itemText); m != "" {
				checked = strings.HasPrefix(strings.ToLower(m), "[x]")
				itemText = strings.TrimSpace(itemText[len(m):])
			}
			if itemText != "" {
				todo := todo{
					task:    itemText,
					checked: checked,
				}
				todos = append(todos, todo)
			}
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			walk(c)
		}
	}
	walk(doc)

	return todos, nil
}

func randomMessage() string {
	randomMessages := []string{
		"murr your motivation doesn't also need to be extinct",
		"mrow I will love you until I go back to being a constellation",
		"meow todolist?  more like able-list",
		"you just lost the game",
		"meow cats are capable of judgement",
	}

	// Get a random message from the list
	idx := time.Now().Minute() % len(randomMessages)

	banner := bold.Render("Tododo ≽^•⩊ •^≼")
	s := fmt.Sprintf("%s\nBeni wisdom ~~%s~~", banner, randomMessages[idx])

	return faint.Render(s)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
