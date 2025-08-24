// Package todo manages the state of the actual todo list.
package todo

// Todo represents a collection of tasks that can be managed.
type Todo struct {
	Tasks    []*Task
	Cursor   int
	Hidden   bool
	Filename string
}

func (t *Todo) ToggleHidden() {
	t.Hidden = !t.Hidden
}

func (t *Todo) RemoveTask(index int) {
	if index < 0 || index >= len(t.Tasks) {
		return
	}

	t.Tasks = append(t.Tasks[:index], t.Tasks[index+1:]...)
}

func (t *Todo) Reposition(up bool) {
	if up {
		for t.Cursor > 0 {
			t.Cursor--
			t.Tasks[t.Cursor], t.Tasks[t.Cursor+1] = t.Tasks[t.Cursor+1], t.Tasks[t.Cursor]
			if !t.Hidden || !t.Tasks[t.Cursor+1].Checked {
				break
			}
		}
	}

	if !up {
		for t.Cursor < len(t.Tasks)-1 {
			t.Cursor++
			t.Tasks[t.Cursor], t.Tasks[t.Cursor-1] = t.Tasks[t.Cursor-1], t.Tasks[t.Cursor]
			if !t.Hidden || !t.Tasks[t.Cursor-1].Checked {
				break
			}
		}
	}
}

func (t *Todo) ModulateCursor(amount int) {
	newPosition := t.Cursor + amount
	newPosition = t.ConvertToValidCursor(newPosition)
	if amount < 0 {
		amount = -1
	} else {
		amount = 1
	}

	if t.Hidden && t.GetRemainingTaskCount() == 0 {
		t.Cursor = -1
		return
	}
	if t.Hidden && t.GetRemainingTaskCount() > 0 {
		for i := 0; i < len(t.Tasks); i++ {
			newPosition = t.ConvertToValidCursor(newPosition)
			if !t.Tasks[newPosition].Checked {
				break
			}
			newPosition += amount
		}
	}

	t.Cursor = newPosition
}

func (t *Todo) ConvertToValidCursor(index int) int {
	if index < 0 {
		for index < 0 {
			index += len(t.Tasks)
		}
		return index
	}

	if index >= len(t.Tasks) {
		return index % len(t.Tasks)
	}

	return index
}

func (t *Todo) GetRemainingTaskCount() int {
	count := 0
	for _, t := range t.Tasks {
		if !t.Checked {
			count++
		}
	}
	return count
}

func (t *Todo) AddTask(text string) {
	newTask := &Task{
		Text:    text,
		Checked: false,
	}
	t.Tasks = append(t.Tasks, newTask)
}

// Task is a single task in a Todo list.
type Task struct {
	Text    string
	Checked bool
}

func (t *Task) UpdateText(text string) {
	t.Text = text
}

func (t *Task) ToggleChecked() {
	t.Checked = !t.Checked
}
