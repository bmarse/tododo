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

	tl "github.com/bmarse/tododo/pkg/todo"
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

type model struct {
	todo   tl.Todo
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
	todolist, err := LoadTodosFromMarkdown()
	if err != nil {
		log.Fatal("failed to load todos:", err)
	}
	ti := textinput.New()
	ti.Placeholder = "Whatcha want to do?"
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80
	return model{
		todo:  todolist,
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
			m.adding = true
			m.input.SetValue(m.todo.Tasks[m.todo.Cursor].Text)
		case "d":
			m.todo.RemoveTodoAtIndex(m.todo.Cursor)
		case "t":
			m.todo.Cursor = 0
			m.todo.ToggleHidden()
		case "w":
			m.saving = true
			if err := m.SaveTodo(); err != nil {
				log.Fatal(err)
			}
			return m, tickCmd
		case " ", "x":
			m.todo.Tasks[m.todo.Cursor].ToggleChecked()
			m.todo.ModulateCursor(0)
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
	if (m.todo.GetRemainingTaskCount() == 0 && m.todo.Hidden) || len(m.todo.Tasks) == 0 {
		tasks += "Yippee! No tasks to do..."
	}
	for i, t := range m.todo.Tasks {
		if m.todo.Hidden && t.Checked {
			continue
		}
		cursor := " " // no cursor
		if m.todo.Cursor == i {
			cursor = bold.Render(">")
		}
		checked := " " // not completed

		taskText := t.Text
		if t.Checked {
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

func (m *model) SaveTodo() error {
	// Save the current todo list to a file
	b := []byte{}
	for _, t := range m.todo.Tasks {
		check := " "
		if t.Checked {
			check = "X"
		}
		b = fmt.Appendf(b, "- [%s] %s\n", check, t.Text)
	}

	if len(b) > 0 {
		err := os.WriteFile(todoFileName, b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func LoadTodosFromMarkdown() (tl.Todo, error) {
	todolist := tl.Todo{
		Tasks: []*tl.Task{},
	}

	markdownContent, err := os.ReadFile(todoFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			todolist.Tasks = append(todolist.Tasks, &tl.Task{
				Text: "Add items to your todo list",
			})
		}
		return todolist, err
	}

	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(markdownContent))

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
				task := &tl.Task{
					Text:    itemText,
					Checked: checked,
				}
				todolist.Tasks = append(todolist.Tasks, task)
			}
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			walk(c)
		}
	}
	walk(doc)

	return todolist, nil
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
