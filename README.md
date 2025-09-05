# ToDoDo
Tododo is a TUI todo manager that should be extinct.  It was created because I was too lazy to research for a good TUI todolist that didn't require an account or database.

Tododo uses a markdown file to store its todo data and is meant to be very barebones.  By default it uses `.tododo.md` but you can supply any filename you desire.  The keybindings are easy to use and always in view if you forget(See screenshots below).

![tododo](https://github.com/user-attachments/assets/e43e7467-4b20-40bd-92f5-64d33d626421)



## Features
- Keyboard based
- Write to file (`.tododo.md`)
- Create/Toggle Complete/Edit/Delete tasks
- Hide completed tasks
- prettier and more fun than just using a plain old markdown file

## Brew Install
```
# Install
brew install bmarse/tap/tododo

# Upgrade
brew upgrade bmarse/tap/tododo
```

## Usage
```bash
tododo  # open up tododo and read/write default file .tododo.md
tododo ~/my-file.md  # Open up tododo with the file ~/my-file.md
tododo --help  # Help
tododo --version # Get the version, needed for opening issues
```

## CLI Help
```
╰─$ tododo --help
 ..   Tododo
, Õ   tasks are like socks, they always seem to multiply
 //_---_
 \  V   )
  ------

NAME:
   tododo - The todo manager that should be extinct

USAGE:
   tododo [options] FILE

   FILE is the file we will use to store and load todos.

VERSION:
   brew-v0.7.0-stable

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

KEY COMMANDS:
    ↑/↓ (k/j): Move the cursor up and down to the next task
    a: Add a new task to your todo list
    <space> (x): Mark the selected task as completed or not completed
    d: Delete the selected task from your todo list
    w (ctrl+s): Save your current todo list to the provided file
    <tab>: Toggle the indentation level of the selected task
    e: Edit the text of the selected task
    m/n: Move the selected task up or down in the list
    t: Show or hide completed tasks in your todo list
    q (ctrl+c): Exit the application
    ?: Show or hide this help menu
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
<img width="1023" height="847" alt="tododo-main" src="https://github.com/user-attachments/assets/5252ec77-3e36-49ee-ae8c-699d3a935aa0" />


### Adding/editing Task TUI
<img width="1023" height="847" alt="tododo-add-task" src="https://github.com/user-attachments/assets/81f69262-2bc2-4e1c-805e-2921a82b82ef" />



