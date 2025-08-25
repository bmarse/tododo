package todo

import (
	"testing"
)

func newTodoWithTasks(checked []bool) *Todo {
	tasks := make([]*Task, len(checked))
	for i, c := range checked {
		tasks[i] = &Task{Text: "Task", Checked: c}
	}
	return &Todo{Tasks: tasks, Cursor: 0}
}

func TestModulateCursor_Normal(t *testing.T) {
	todo := newTodoWithTasks([]bool{false, false, false})
	todo.Cursor = 0
	todo.ModulateCursor(1)
	if todo.Cursor != 1 {
		t.Errorf("expected cursor 1, got %d", todo.Cursor)
	}
	todo.ModulateCursor(1)
	if todo.Cursor != 2 {
		t.Errorf("expected cursor 2, got %d", todo.Cursor)
	}
	todo.ModulateCursor(1)
	if todo.Cursor != 0 {
		t.Errorf("expected cursor 0 (wrap around), got %d", todo.Cursor)
	}
	todo.ModulateCursor(-1)
	if todo.Cursor != 2 {
		t.Errorf("expected cursor 2 (wrap around), got %d", todo.Cursor)
	}
}

func TestModulateCursor_HiddenAllChecked(t *testing.T) {
	todo := newTodoWithTasks([]bool{true, true, true})
	todo.Hidden = true
	todo.Cursor = 0
	todo.ModulateCursor(1)
	if todo.Cursor != -1 {
		t.Errorf("expected cursor -1 when all tasks checked and hidden, got %d", todo.Cursor)
	}
}

func TestModulateCursor_HiddenSomeUnchecked(t *testing.T) {
	todo := newTodoWithTasks([]bool{true, false, true})
	todo.Hidden = true
	todo.Cursor = 0
	todo.ModulateCursor(1)
	if todo.Cursor != 1 {
		t.Errorf("expected cursor 1 (first unchecked), got %d", todo.Cursor)
	}
	todo.ModulateCursor(1)
	if todo.Cursor != 1 {
		t.Errorf("expected cursor 1 (should stay on unchecked), got %d", todo.Cursor)
	}
}

func TestConvertToValidCursor(t *testing.T) {
	todo := newTodoWithTasks([]bool{false, false, false})

	// Test negative index
	if todo.ConvertToValidCursor(-1) != 2 {
		t.Errorf("expected cursor 2 for -1, got %d", todo.ConvertToValidCursor(-1))
	}

	// Test index within bounds
	if todo.ConvertToValidCursor(1) != 1 {
		t.Errorf("expected cursor 1 for 1, got %d", todo.ConvertToValidCursor(1))
	}

	// Test index greater than length
	if todo.ConvertToValidCursor(3) != 0 {
		t.Errorf("expected cursor 0 for 3, got %d", todo.ConvertToValidCursor(3))
	}
}

func TestGetRemainingTaskCount(t *testing.T) {
	todo := newTodoWithTasks([]bool{true, false, true})
	if todo.GetRemainingTaskCount() != 1 {
		t.Errorf("expected 1 remaining task, got %d", todo.GetRemainingTaskCount())
	}

	todo = newTodoWithTasks([]bool{true, true, true})
	if todo.GetRemainingTaskCount() != 0 {
		t.Errorf("expected 0 remaining tasks, got %d", todo.GetRemainingTaskCount())
	}

	todo = newTodoWithTasks([]bool{false, false, false})
	if todo.GetRemainingTaskCount() != 3 {
		t.Errorf("expected 3 remaining tasks, got %d", todo.GetRemainingTaskCount())
	}
}

func TestAddTask(t *testing.T) {
	todo := &Todo{Tasks: []*Task{}, Cursor: 0}
	todo.AddTask("New Task", 0)
	if len(todo.Tasks) != 1 || todo.Tasks[0].Text != "New Task" {
		t.Errorf("expected 1 task with text 'New Task', got %d tasks", len(todo.Tasks))
	}
}

func TestRemoveTodo(t *testing.T) {
	todo := &Todo{Tasks: []*Task{{Text: "Task 1", Checked: false}}, Cursor: 0}
	todo.RemoveTask(0)
	if len(todo.Tasks) != 0 {
		t.Errorf("expected 0 tasks after removal, got %d", len(todo.Tasks))
	}
}
