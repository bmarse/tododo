# ToDoDo
Tododo is a TUI todo manager that should be extinct.  It was created because I was too lazy to research for a good TUI todolist that didn't require an account or database.

Tododo uses a markdown file to store its todo data and is meant to be very barebones.  By default it uses `.tododo.md` but you can supply any filename you desire.  The keybindings are easy to use and always in view if you forget(See screenshots below).

![tododo](https://github.com/user-attachments/assets/fa484b9e-62d6-41c5-ac2a-1f5b390a427b)

## Features
- Keyboard based
- Write to file (`.tododo.md`)
- Create/Toggle Complete/Edit/Delete tasks
- Hide completed tasks
- prettier and more fun than just using a plain old markdown file

## Brew Install
```
brew tap bmarse/tododo
brew install tododo
```

## Usage
```bash
tododo  # open up tododo and read/write default file .tododo.md
tododo ~/my-file.md  # Open up tododo with the file ~/my-file.md
tododo --help  # Help
tododo --version # Get the version, needed for opening issues
```

## Building with Go
This is a very complex project with lots complicated dependencies.
```
git clone git@github.com:bmarse/tododo.git
cd tododo
go build -o tododo tododo.go
```

## Installation with Go
Very complex installation process:
```
go install github.com/bmarse/tododo@latest
```

## Obligatory Screenshots
Included are some screenshots of tododo running on [Ghostty](https://ghostty.org/) with [Nerd Fonts](https://www.nerdfonts.com/).

### Main TUI
<img width="912" height="657" alt="mainTUI" src="https://github.com/user-attachments/assets/3e5c1b12-95f8-4c10-b855-eee8ae3f3219" />

### Adding/editing Task TUI
<img width="912" height="657" alt="MainTUI-edit" src="https://github.com/user-attachments/assets/49a1ccfc-6ab6-4044-8f88-12cba4a8f89b" />
