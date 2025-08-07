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
	hidden bool
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
				m.adding = false
				if strings.TrimSpace(m.input.Value()) == "" {
					break
				}
				if m.cursor == -1 {
					m.todos = append(m.todos, todo{task: m.input.Value()})
					m.cursor = len(m.todos) - 1

				} else {
					m.todos[m.cursor].task = m.input.Value()
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
			m.ModulateCursor(1)
		case "k", "up":
			m.ModulateCursor(-1)
		case "a":
			m.adding = true
			m.cursor = -1
		case "e":
			m.adding = true
			m.input.SetValue(m.todos[m.cursor].task)
		case "d":
			m.RemoveTodoAtIndex(m.cursor)
		case "t":
			m.cursor = 0
			m.hidden = !m.hidden
		case "w":
			m.saving = true
			if err := m.SaveTodo(); err != nil {
				log.Fatal(err)
			}
			return m, tickCmd
		case " ", "x":
			m.todos[m.cursor].checked = !m.todos[m.cursor].checked
			m.ModulateCursor(0)
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
	if (m.GetRemainingTaskCount() == 0 && m.hidden) || len(m.todos) == 0 {
		tasks += "Yippee! No tasks to do..."
	}
	for i, t := range m.todos {
		if m.hidden && t.checked {
			continue
		}
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = bold.Render(">")
		}
		checked := " " // not completed

		taskText := t.task
		if t.checked {
			checked = bold.Render("x")
			taskText = faint.Render(taskText)
		}
		tasks += fmt.Sprintf("%s [%s]: %s\n", cursor, checked, taskText)
	}
	s += border.Render(tasks)

	if m.saving {
		s += "\n Saving... \n"
	} else {
		s += "\n\n"
	}

	s += "\n\n↑/↓: Move  a: Add  <space>: Complete  t: Toggle Hidden  e: Edit  d: Delete  w: Write  q: Quit"
	return padded.Render(s)
}

func (m *model) RemoveTodoAtIndex(index int) {
	if index < 0 || index >= len(m.todos) {
		return
	}

	m.todos = append(m.todos[:index], m.todos[index+1:]...)
}

func (m *model) ModulateCursor(amount int) {
	newPosition := m.cursor + amount
	newPosition = m.ConvertToValidCursor(newPosition)
	if amount < 0 {
		amount = -1
	} else {
		amount = 1
	}

	if m.hidden && m.GetRemainingTaskCount() > 0 {
		for i := 0; i < len(m.todos); i++ {
			newPosition = m.ConvertToValidCursor(newPosition)
			if !m.todos[newPosition].checked {
				break
			}
			newPosition += amount
		}
	}

	m.cursor = newPosition
}

func (m model) ConvertToValidCursor(index int) int {
	if index < 0 {
		for index < 0 {
			index += len(m.todos)
		}
		return index
	}

	if index >= len(m.todos) {
		return index % len(m.todos)
	}

	return index
}

func (m model) GetRemainingTaskCount() int {
	count := 0
	for _, t := range m.todos {
		if !t.checked {
			count++
		}
	}
	return count
}

func (m model) GetTodos() []todo {
	if m.hidden {
		todos := make([]todo, 0, len(m.todos))
		for _, t := range m.todos {
			if !t.checked {
				todos = append(todos, t)
			}
		}

		return todos
	}
	return m.todos
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
		"your motivation doesn't also need to be extinct",
		"todolist?  more like able-list",
		"you just lost the game",
		"to task or not to task, that is the question",
		"say hello to my little task",
		"does your task spark joy",
		"no kings, only tasks",
		"help I'm trapped in a todo list factory",
		"don't forget to take breaks",
		"I'm afraid we're not in vim anymore Toto",
		"frankly my dear, I don't give a task",
		"hey, you can do this!",
		"tasks are like socks, they always seem to multiply",
		"the task really ties the list together",
		"the tasks of mice and men",
	}

	// Get a random message from the list
	idx := time.Now().Minute() % len(randomMessages)

	banner := bold.Render("Tododo ≽^•⩊ •^≼")
	s := fmt.Sprintf("%s\n~~%s~~", banner, randomMessages[idx])

	return faint.Render(s)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
