// Package ui simplifies the rendering layer of the todolist.
package ui

import (
	"fmt"
	"strings"
	"time"

	tl "github.com/bmarse/tododo/internal/todo"
	"github.com/charmbracelet/lipgloss"
)

var rainbowColors = []string{
	"#ffadad", // Red
	"#ffd6a5", // Orange
	"#fdffb6", // Yellow
	"#caffbf", // Green
	"#9bf6ff", // Blue
	"#a0c4ff", // Indigo
	"#bdb2ff", // Violet
}

var (
	standard     lipgloss.Style = lipgloss.NewStyle()
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
		checked = "x"
		taskText = faint.Render(taskText)
	}
	checkbox := standard.Render(fmt.Sprintf("[%s]:", checked))
	return fmt.Sprintf("%s %s %s\n", cursor, checkbox, taskText)
}

func MenuUI(hideCommandMenu bool) string {
	if hideCommandMenu {
		return faint.Render("Press '?' to toggle command menu...")
	}
	menu := strings.Builder{}
	lineLength := 0
	for k := range GetKeys() {
		menuAppend := fmt.Sprintf("%s: %s | ", bold.Render(GetKeys()[k].Key), faint.Render(GetKeys()[k].Title))
		lineLength += len(menuAppend)
		if lineLength > windowWidth-5 {
			// remove last " | "
			menuAppend = strings.TrimSuffix(menuAppend, " | ")
			menuAppend += "\n"
			lineLength = 0
		}
		menu.WriteString(menuAppend)
	}
	return strings.TrimSuffix(menu.String(), " | ")
}

func MainUI(todo *tl.Todo, saving bool, spinner string, hideCommandMenu bool) string {
	s := RandomMessage()
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
	s += MenuUI(hideCommandMenu)
	return padded.Render(s)
}

func updateColors(dancing bool, index int) {
	if dancing {
		border = border.Foreground(lipgloss.Color(rainbowColors[index]))
		bold = bold.Foreground(lipgloss.Color(rainbowColors[index]))
		faint = faint.Foreground(lipgloss.Color(rainbowColors[index]))
		standard = standard.Foreground(lipgloss.Color(rainbowColors[index]))
	} else {
		border = border.Foreground(lipgloss.Color(""))
		bold = bold.Foreground(lipgloss.Color(""))
		faint = faint.Foreground(lipgloss.Color(""))
		standard = standard.Foreground(lipgloss.Color(""))
	}
}

func GetKeyHelp() string {
	s := strings.Builder{}
	s.WriteString("KEY COMMANDS:\n")
	for _, k := range GetKeys() {
		a := ""
		if k.AliasKey != "" {
			a = fmt.Sprintf(" (%s)", bold.Render(k.AliasKey))
		}
		s.WriteString(fmt.Sprintf("    %s%s: %s\n", bold.Render(k.Key), faint.Render(a), k.Description))
	}
	return s.String()
}

func RandomMessage() string {
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
