// Package ui simplifies the rendering layer of the todolist.
package ui

import (
	"fmt"
	"strings"
	"time"

	tl "github.com/bmarse/tododo/internal/todo"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	bold         lipgloss.Style = lipgloss.NewStyle().Bold(true)
	padded       lipgloss.Style = lipgloss.NewStyle().Padding(2, 4)
	faint        lipgloss.Style = lipgloss.NewStyle().Faint(true)
	border       lipgloss.Style = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true).
			BorderTop(true).
			Padding(1, 2, 1, 2)
)

func AddingUI(inputView string) string {
	msg := bold.Render("Add task:")
	return padded.Render(fmt.Sprintf("%s\n%s\n\n(Enter to save, Esc to cancel)", msg, border.Render(inputView)))
}

func TaskUI(t *tl.Task, cursor string) string {
	checked := " " // not completed

	taskText := t.Text
	if t.Checked {
		checked = bold.Render("x")
		taskText = faint.Render(taskText)
	}
	return fmt.Sprintf("%s [%s]: %s\n", cursor, checked, taskText)
}

func MenuUI() string {
	return "↑/↓: Move  a: Add  <space>: Complete  t: Toggle Hidden  e: Edit  d: Delete  w: Write  q: Quit"
}

func MainUI(todo *tl.Todo, saving bool, spinner string) string {
	s := randomMessage()
	s += "\n\n"
	tasks := ""
	if (todo.GetRemainingTaskCount() == 0 && todo.Hidden) || len(todo.Tasks) == 0 {
		tasks += "Yippee! No tasks to do..."
	}
	for i, t := range todo.Tasks {
		if todo.Hidden && t.Checked {
			continue
		}
		cursor := " " // no cursor
		if todo.Cursor == i {
			cursor = bold.Render(">")
		}
		tasks += TaskUI(t, cursor)
	}
	s += border.Render(tasks)

	if saving {
		s += fmt.Sprintf("\n %s Saving... \n", spinnerStyle.Render(spinner))
	} else {
		s += "\n\n"
	}

	s += "\n\n"
	s += MenuUI()
	return padded.Render(s)
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
	dodo := []string{
		" ..   Tododo",
		", Õ   " + randomMessages[idx],
		" //_---_ ",
		" \\  V   )",
		"  ------",
	}
	banner := bold.Render(strings.Join(dodo, "\n"))

	return faint.Render(banner)
}
