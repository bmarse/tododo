package ui

type KeyCommand struct {
	Key         string
	AliasKey    string
	Title       string
	Description string
}

func GetKeys() []KeyCommand {
	return []KeyCommand{
		{
			"↑/↓",
			"k/j",
			"Move",
			"Move the cursor up and down to the next task",
		},
		{
			"a",
			"",
			"Add Task",
			"Add a new task to your todo list",
		},
		{
			"<space>",
			"x",
			"Toggle Complete",
			"Mark the selected task as completed or not completed",
		},
		{
			"d",
			"",
			"Delete Task",
			"Delete the selected task from your todo list",
		},
		{
			"w",
			"ctrl+s",
			"Write to file",
			"Save your current todo list to the provided file",
		},
		{
			"<tab>",
			"",
			"Toggle Indent",
			"Toggle the indentation level of the selected task",
		},
		{
			"e",
			"",
			"Edit Task",
			"Edit the text of the selected task",
		},
		{
			"m/n",
			"",
			"Reposition Task",
			"Move the selected task up or down in the list",
		},
		{
			"t",
			"",
			"Toggle Hidden",
			"Show or hide completed tasks in your todo list",
		},
		{
			"q",
			"ctrl+c",
			"Quit",
			"Exit the application",
		},
		{
			"?",
			"",
			"Toggle Help",
			"Show or hide this help menu",
		},
	}
}
