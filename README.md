# ToDoDo
ToDoDo is a TUI todolist manager that should be extinct.  It was created because I was too lazy to google for a good TUI todolist that didn't require an account or database.

![todo-in-action](https://github.com/user-attachments/assets/78d01768-acc1-47d6-8f67-761b434912e6)


## Features
- Keyboard based
- Write to file (`.tododo.md`)
- Create/Toggle Complete/Edit/Delete tasks
- Hide completed tasks
- prettier and more fun than just using a plain old markdown file


## Usage
```bash
tododo  # open up tododo and read/write default file .tododo.md
tododo ~/my-file.md  # Open up tododo with the file ~/my-file.md
tododo --help  # Help
tododo --version # Get the version, needed for opening issues
```

## Building
This is a very complex project with lots complicated dependencies.
```
git clone git@github.com:bmarse/tododo.git
cd tododo
go build -o tododo cmd/tododo/main.go
```


## Obligitory Screenshots
Included are some screenshots of tododo running on [Ghostty](https://ghostty.org/) with [Nerd Fonts](https://www.nerdfonts.com/).

### Main TUI
<img width="856" height="454" alt="tododo" src="https://github.com/user-attachments/assets/783d0a67-8199-4b96-947d-ee48cc1fe050" />

### Adding/editing Task TUI
<img width="854" height="455" alt="tododo-add-task" src="https://github.com/user-attachments/assets/37643d40-2aeb-4ae1-9d49-abe8694af181" />
