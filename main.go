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
	//timer stopwatch.Model
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
		db:    queries,
		tasks: tasks,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return nil, nil
}

func (m model) View() string {
	return ""
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
