package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	db "github.com/chee-zer/negentropy/database/sqlc"
	"github.com/chee-zer/negentropy/stopwatch"
	_ "github.com/mattn/go-sqlite3"
)

type model struct {
	db           *db.Queries
	tasks        map[int64]db.Task
	activeTaskId int64
	statusQuote  string
	help         string
	//tabs         TabModel
	timer          stopwatch.Model
	quitting       bool
	currentSession *db.Session
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
	taskMap, err := GetTaskMap(queries)
	if err != nil {
		log.Fatalf("couldn't not load tasks: %v", err)
	}

	dummyTimer := stopwatch.NewTimer("dummy")

	return model{
		db:             queries,
		tasks:          taskMap,
		activeTaskId:   0,
		statusQuote:    "this is status quote",
		help:           "this is help string",
		quitting:       false,
		timer:          dummyTimer,
		currentSession: nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
		if m.timer.IsRunning() {
			//TIMER RUNNING
			switch msg.String() {
			case "ctrl+c", "q":
				m.help = "Please end your session before quitting the app. Press Spacebar/enter to pause the timer"
				return m, nil
			}
			// check for no tasks here, basically nothing can be done if zero tasks
			if len(m.tasks) == 0 {
				m = m.NoTaskView()
				return m, nil
			}
		} else {
			//TIMER STOPPED
			if len(m.tasks) == 0 {
				m = m.NoTaskView()
			}
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case " ", "enter":
				// if program is here, there should not be zero tasks, so no checks required
				_, ok := m.tasks[m.activeTaskId]
				if !ok {
					m.statusQuote = "No task selected, press 'n' to create a new task"
					return m, nil
				}

				if m.timer.IsRunning() {
					return m.StopSession(), m.timer.StopCmd()
				} else {
					return m.StartSession(), m.timer.StartCmd()
				}

				return m, nil
			case "n":

			}
		}

		switch msg.String() {

		}

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

func GetTaskMap(queries *db.Queries) (map[int64]db.Task, error) {
	tasks, err := queries.GetTasks(context.Background())
	if err != nil {
		return nil, err
	}
	taskMap := make(map[int64]db.Task)
	for _, v := range tasks {
		taskMap[v.ID] = v
	}
	return taskMap, nil
}

func (m model) StartSession() model {
	taskID := m.activeTaskId
	sessionParams := db.StartSessionParams{
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
		TaskID:    taskID,
	}
	session, err := m.db.StartSession(context.Background(), sessionParams)
	if err != nil {
		m.statusQuote = "Couldn't start session: " + err.Error()
		return m
	}
	timer := stopwatch.NewTimerRunning(m.tasks[taskID].Name)
	m.timer = timer
	m.currentSession = &session
	return m
}

func (m model) StopSession() model {
	taskID := m.activeTaskId
	m.timer.StopCmd()
	endTime := time.Now().Format("2006-01-02 15:04:05")
	endSessionParams := db.EndSessionParams{
		EndTime: sql.NullString{String: endTime, Valid: true},
		TaskID:  taskID,
	}
	m.db.EndSession(context.Background(), endSessionParams)

	return m
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
