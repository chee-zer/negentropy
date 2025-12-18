package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	db "github.com/chee-zer/negentropy/database/sqlc"
	"github.com/chee-zer/negentropy/stopwatch"
	_ "github.com/mattn/go-sqlite3"
)

// capitalized fields for json.Marshal, only for debugging purposes, remove laater
type model struct {
	db             *db.Queries
	tasks          map[int64]db.Task
	ActiveTaskId   int64
	StatusQuote    string
	help           string
	tabs           TabModel
	Timer          stopwatch.StopwatchModel
	quitting       bool
	CurrentSession *db.Session
	textInput      textinput.Model
	Typing         bool
	keymap         keymap
}

// global keymap
// timer - start/stop(space, enter)
// switch timer while timer not running(tab, l)
//
// exit (q)

// TODO: forgot i had these, assign these AFTER the update loop is done
type keymap struct {
	startStopTimer key.Binding
	//switchTimer    key.Binding
	exit    key.Binding
	goRight key.Binding
	goLeft  key.Binding
	//deleteTask     key.Binding
	createTask key.Binding
	//resetTimer     key.Binding
}

func NewModel(queries *db.Queries) model {
	taskMap, tasks, err := GetTaskMap(queries)
	if err != nil {
		log.Fatalf("couldn't not load tasks: %v", err)
	}

	keymapvalues := keymap{
		goRight: key.NewBinding(
			key.WithKeys("right"),
		),
		goLeft: key.NewBinding(
			key.WithKeys("left"),
		),
	}

	dummyTimer := stopwatch.NewTimer("dummy")
	ti := textinput.New()
	ti.Placeholder = "Enter task name"
	ti.CharLimit = 20
	ti.Width = 20

	tabs := NewTabModel(tasks)
	return model{
		db:             queries,
		tasks:          taskMap,
		ActiveTaskId:   0,
		StatusQuote:    "this is status quote",
		help:           "this is help string",
		quitting:       false,
		Timer:          dummyTimer,
		CurrentSession: nil,
		textInput:      ti,
		Typing:         false,
		keymap:         keymapvalues,
		tabs:           tabs,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.Typing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				newTaskName := m.textInput.Value()
				taskCreatingParams := db.CreateTaskParams{
					Name: newTaskName,
					ColorHex: sql.NullString{String: "what",
						Valid: true},
					DailyTarget: sql.NullInt64{
						Int64: 3600,
						Valid: true,
					},
				}

				task, err := m.db.CreateTask(context.Background(), taskCreatingParams)
				if err != nil {
					m.StatusQuote = "Couldn't create task: " + err.Error()
					m.textInput.Reset()

					return m, m.textInput.Focus()
				}
				m.Typing = false

				// adding the task in main model and switching to it
				m.tasks[task.ID] = task
				m.ActiveTaskId = task.ID
				m.tabs.Tasks = append(m.tabs.Tasks, task)
				m.tabs.ActiveTabIndex = len(m.tabs.Tasks) - 1
				m.textInput.Reset()
				m.textInput.Blur()
				m.StatusQuote = "Task created"
				return m, nil

			case tea.KeyEsc:
				m.textInput.Reset()
				m.Typing = false
				m.StatusQuote = "Task not created -_-"
				m.textInput.Blur()
			}

			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case stopwatch.ResetMsg, stopwatch.StartStopMsg, stopwatch.TickMsg:
		var timerCmd tea.Cmd
		m.Timer, timerCmd = m.Timer.Update(msg)
		return m, timerCmd

	case SwitchSelectedTaskMsg:
		m.ActiveTaskId = msg.taskID
	case SwitchMsg:
		var tabCmd tea.Cmd

		m.tabs, tabCmd = m.tabs.Update(msg)
		return m, tabCmd

	case tea.KeyMsg:
		if m.Timer.IsRunning() {
			//TIMER RUNNING
			switch {
			case key.Matches(msg, m.keymap.exit):
				m.help = "Please end your session before quitting the app. Press Spacebar/enter to pause the timer"
				return m, nil
			case key.Matches(msg, m.keymap.startStopTimer):
				m.StatusQuote = "Session ended!!"
				log.Printf("\n%+v\n", m)
				return m.StopSession(), m.Timer.StopCmd()
			}
			if len(m.tasks) == 0 {
				m = m.NoTaskView()
				return m, nil
			}

		} else {
			//TIMER STOPPED
			if len(m.tasks) == 0 {
				m = m.NoTaskView()
			}
			switch {
			case key.Matches(msg, m.keymap.goLeft):
				return m, m.tabs.SwitchLeftCmd()
			case key.Matches(msg, m.keymap.goRight):
				return m, m.tabs.SwitchRightCmd()
			case key.Matches(msg, m.keymap.exit):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, m.keymap.startStopTimer):
				_, ok := m.tasks[m.ActiveTaskId]
				if !ok {
					// TODO: change this later
					createNewHotkey := strings.Join(m.keymap.createTask.Keys(), "/")
					m.StatusQuote = "No task selected, press " + createNewHotkey + " to create a new task"
					return m, nil
				}
				m.StatusQuote = "Session Started!"
				return m.StartSession(), m.Timer.StartCmd()
			case key.Matches(msg, m.keymap.createTask):
				cmd = m.textInput.Focus()
				m.Typing = true
				return m, cmd
			}

		}
		// switch msg.String() {
		// case "ctrl+c", "q":
		// 	m.quitting = true
		// 	return m, tea.Quit
		// case " ", "enter":
		// 	_, ok := m.tasks[m.ActiveTaskId]
		// 	if !ok {
		// 		m.StatusQuote = "No task selected, press 'n' to create a new task"
		// 		return m, nil
		// 	}
		// 	m.StatusQuote = "Session Started!"
		// 	return m.StartSession(), m.Timer.StartCmd()
		// case "n":
		// 	cmd = m.textInput.Focus()
		// 	m.Typing = true
		// 	return m, cmd
		// }
	}
	return m, nil

}

func (m model) View() string {
	if m.quitting {
		return "quitting negetropy!"
	}
	s := fmt.Sprintf("\n\n\n\ntasks: %s\n\nActive Task ID: %d\n  %s\n\n  %s\n %s\n", m.tabs.View(), m.ActiveTaskId, m.StatusQuote, m.help, m.textInput.View())
	return s
}

func (m model) NoTaskView() model {
	m.StatusQuote = "No tasks found :( Press 'n' to create a new task!"
	return m
}

func GetTaskMap(queries *db.Queries) (map[int64]db.Task, []db.Task, error) {
	tasks, err := queries.GetTasks(context.Background())
	if err != nil {
		return nil, nil, err
	}
	taskMap := make(map[int64]db.Task)
	for _, v := range tasks {
		taskMap[v.ID] = v
	}
	return taskMap, tasks, nil
}

func (m model) StartSession() model {
	taskID := m.ActiveTaskId
	sessionParams := db.StartSessionParams{
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
		TaskID:    taskID,
	}
	session, err := m.db.StartSession(context.Background(), sessionParams)
	if err != nil {
		m.StatusQuote = "Couldn't start session: " + err.Error()
		return m
	}
	timer := stopwatch.NewTimerRunning(m.tasks[taskID].Name)
	m.Timer = timer
	m.CurrentSession = &session
	return m
}

func (m model) StopSession() model {
	taskID := m.ActiveTaskId
	endTime := time.Now().Format("2006-01-02 15:04:05")
	endSessionParams := db.EndSessionParams{
		EndTime: sql.NullString{String: endTime, Valid: true},
		TaskID:  taskID,
	}
	m.db.EndSession(context.Background(), endSessionParams)

	return m
}

func (m model) String() string {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(b)
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	log.Println("what")
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
