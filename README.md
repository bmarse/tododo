# ToDoDo
Tododo is a TUI todo manager that should be extinct.  It was created because I was too lazy to research for a good TUI todolist that didn't require an account or database.

Tododo uses a markdown file to store its todo data and is meant to be very barebones.  By default it uses `.tododo.md` but you can supply any filename you desire.  The keybindings are easy to use and always in view if you forget(See screenshots below).

![tododo](https://github.com/user-attachments/assets/a9d1d60f-5b94-4628-9729-92c78f6de7db)


## Features
- Keyboard based
- Write to file (`.tododo.md`)
- Create/Toggle Complete/Edit/Delete tasks
- Hide completed tasks
- prettier and more fun than just using a plain old markdown file

## Brew Install
```
# Install
brew install bmarse/tododo/tododo

# Upgrade
brew upgrade bmarse/tododo/tododo
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
$ tododo --help

 ..   Tododo                                 
, Õ   help I'm trapped in a todo list factory
 //_---_                                     
 \  V   )                                    
  ------                                     

NAME:
   tododo - The todo manager that should be extinct

USAGE:
   tododo [options] FILE

   FILE is the file we will use to store and load todos.

VERSION:
   brew-v0.6.0-stable

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

KEY COMMANDS:
    ↑/↓ (j/k): Move the cursor up and down to the next task
    a: Add a new task to your todo list
    <space> (x): Mark the selected task as completed or not completed
    n/m: Move the selected task up or down in the list
    d: Delete the selected task from your todo list
    w (ctrl+s): Save your current todo list to the provided file
    e: Edit the text of the selected task
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
<img width="1096" height="755" alt="tododo-main" src="https://github.com/user-attachments/assets/79309923-566c-4d2f-9394-7e8648193502" />

### Adding/editing Task TUI
<img width="1052" height="711" alt="tododo-add-task" src="https://github.com/user-attachments/assets/df040f81-a595-4a30-b804-42af3952639b" />


