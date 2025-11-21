package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	db "github.com/chee-zer/negentropy/database/sqlc"
	"github.com/chee-zer/negentropy/stopwatch"
	_ "github.com/mattn/go-sqlite3"
)

// type model struct {
// 	tabs  Tabs
// 	timer Model
// 	// bar Bar
// }

type model struct {
	db           *db.Queries
	tasks        []db.Task
	activeTaskId int
	timerRunning bool
	statusQuote  string
	help         string
	//tabs         TabModel
	//timer stopwatch.Model
	quitting bool
}

// global keymap
// timer - start/stop(space, enter)
// switch timer while timer not running(tab, l)
//
// exit (q)
type keymap struct {
	startStopTimer key.Binding
	switchTimer    key.Binding
	exit           key.Binding
	goRight        key.Binding
	goLeft         key.Binding
	deleteTask     key.Binding
	createTask     key.Binding
	resetTimer     key.Binding
}

func NewModel(queries *db.Queries) model {
	tasks, err := queries.GetTasks(context.Background())

	if err != nil {
		log.Fatalf("couldn't not load tasks: %v", err)
	}

	return model{
		db:           queries,
		tasks:        tasks,
		activeTaskId: 0,
		timerRunning: false,
		statusQuote:  "this is status quote",
		help:         "this is help string",
		quitting:     false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// check for no tasks here, basically nothing can be done if zero tasks

	switch msg := msg.(type) {
	case tea.KeyMsg:
		/*
			keyMsg handling
			if no task is created, prompt user to create new task, only 'n' is allowed
			 check if timer running, then split the keyMsgs
			if timer running:
			[] - only start stop allowed
			[] - reset only when timer stopped, so display a message showing so
			[] - same with navigation tasks(a, d, h, l, left, right), show message
			if timer not running:
			[] - navigation (a, d, h, l, left, right) for navigating tabs
			[] - n - opens a text prompt and creates new task
			[] - x - open a prompt to confirm deletion of task. the records wont be deleted, and the if the task with same name is created, it gets restored
			[] - start stop also allowed
			[] - quit (q)
		*/
		switch msg.String() {
		case "ctrl+c", "q":

			if m.timerRunning {
				m.help = "Please end your session before quitting the app. Press Spacebar/enter to pause the timer"
				return m, nil
			} else {
				m.quitting = true
				return m, tea.Quit
			}
		case " ", "enter":
			// if program is here, there should not be zero tasks, so no checks required
			 newTimer := 

			m.timerRunning = !m.timerRunning
			return m, nil
		}
	}

	if len(m.tasks) == 0 {
		return m.NoTaskView(), nil
	}

	return m, nil
}

func (m model) View() string {
	if len(m.tasks) == 0 {
		return fmt.Sprintf("\n  %s\n\n  %s\n", m.statusQuote, m.help)
	}
	if m.quitting {
		return "quitting negetropy!"
	}
	return fmt.Sprintf("Active Task ID: %d\n  %s\n\n  %s\n", m.activeTaskId, m.statusQuote, m.help)
}

func (m model) NoTaskView() model {
	m.statusQuote = "No tasks found :( Press 'n' to create a new task!"
	return m
}

func (m model) GetTaskMap() [int]string {
	tasks, err := m.db.GetTasks(context.Background()) 
	if err != nil {
		log.Fatalf("couldn't not load tasks: %v", err)
	}
	taskMap := make(map[int]string)
	for _, v := range tasks {
		taskMap[int(v.ID)] = v.Name
	}
	return tasks
}

func main() {
	sqlitedb, err := sql.Open("sqlite3", "./database/appdb.sqlite")
	if err != nil {
		log.Fatalf("Couldn't connect to db: %v", err.Error())
	}
	defer sqlitedb.Close()

	queries := db.New(sqlitedb)

	p := tea.NewProgram(NewModel(queries))

	if _, err := p.Run(); err != nil {
		fmt.Printf("could'nt run program: %v", err)
		os.Exit(1)
	}

}
