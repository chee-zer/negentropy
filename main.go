package main

import (
	"context"
	"database/sql"
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
	keymap         keymap
	state          appState
	pendingAction  currentAction
}
type keymap struct {
	StartStopTimer key.Binding
	Exit           key.Binding
	GoRight        key.Binding
	GoLeft         key.Binding
	DeleteTask     key.Binding
	CreateTask     key.Binding
	ResetTimer     key.Binding
	Yes            key.Binding
	No             key.Binding
}

type appState int

const (
	TimerNotRunning appState = iota
	TimerRunning
	Typing
	Confirming
)

type currentAction int

const (
	deleteTask currentAction = iota
	resetTimer
)

func NewModel(queries *db.Queries, cfg keymap, errs error) model {
	taskMap, tasks, err := GetTaskMap(queries)
	if err != nil {
		log.Fatalf("couldn't not load tasks: %v", err)
	}
	errorString := "Config loaded successfully!"
	if errs != nil {
		errorString = errs.Error()
	}

	dummyTimer := stopwatch.NewTimer("dummy")
	ti := textinput.New()
	ti.Placeholder = "Enter task name"
	ti.CharLimit = 20
	ti.Width = 20

	activeId := int64(0)
	if len(tasks) >= 1 {
		activeId = 1
	}

	tabs := NewTabModel(tasks)
	return model{
		db:             queries,
		tasks:          taskMap,
		ActiveTaskId:   activeId,
		StatusQuote:    errorString,
		help:           "this is help string",
		quitting:       false,
		Timer:          dummyTimer,
		CurrentSession: nil,
		textInput:      ti,
		keymap:         cfg,
		tabs:           tabs,
		state:          TimerNotRunning,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case stopwatch.ResetMsg, stopwatch.StartStopMsg, stopwatch.TickMsg:
		var timerCmd tea.Cmd
		m.Timer, timerCmd = m.Timer.Update(msg)
		return m, timerCmd

	case DeleteSelectedTaskMsg:
		var tabCmd tea.Cmd
		m.tabs, tabCmd = m.tabs.Update(msg)
		return m, tabCmd
	case SwitchSelectedTaskMsg:
		m.ActiveTaskId = msg.taskID
	case SwitchMsg:
		var tabCmd tea.Cmd

		m.tabs, tabCmd = m.tabs.Update(msg)
		return m, tabCmd

	case tea.KeyMsg:
		switch m.state {
		case TimerNotRunning:
			return m.updateTimerNotRunning(msg)
		case TimerRunning:
			return m.updateTimerRunning(msg)
		case Typing:
			// cursor will not blink cuz only keyMsg being passed
			return m.updateTyping(msg)
		case Confirming:
			return m.updateConfirming(msg)
		}
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
	m.state = TimerRunning
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
	m.state = TimerNotRunning
	return m
}

func (m model) ResetSession() model {
	taskID := m.ActiveTaskId
	endTime := time.Now().Format("2006-01-02 15:04:0")
	params := db.EndSessionAsEntropyParams{
		EndTime: sql.NullString{String: endTime, Valid: true},
		TaskID:  taskID,
	}
	m.db.EndSessionAsEntropy(context.Background(), params)
	m.state = TimerNotRunning
	return m
}

func (m model) updateTimerNotRunning(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	//TIMER STOPPED
	if len(m.tasks) == 0 {
		m = m.NoTaskView()
	}

	var cmd tea.Cmd
	switch {
	case key.Matches(msg, m.keymap.GoLeft):
		return m, m.tabs.SwitchLeftCmd()
	case key.Matches(msg, m.keymap.GoRight):
		return m, m.tabs.SwitchRightCmd()
	case key.Matches(msg, m.keymap.Exit):
		m.quitting = true
		return m, tea.Quit
	case key.Matches(msg, m.keymap.StartStopTimer):
		_, ok := m.tasks[m.ActiveTaskId]
		if !ok {
			// TODO: change this later(Stringer for all the keys)
			createNewHotkey := strings.Join(m.keymap.CreateTask.Keys(), "/")
			m.StatusQuote = "No task selected, press " + createNewHotkey + " to create a new task"
			return m, nil
		}
		m = m.StartSession()
		// if timer doesn't start due to db error, will return m, nil
		if m.state == TimerRunning {
			m.StatusQuote = "Session Started: " + m.tasks[m.ActiveTaskId].Name
			return m, m.Timer.StartCmd()
		}
	case key.Matches(msg, m.keymap.CreateTask):
		cmd = m.textInput.Focus()
		m.state = Typing
		return m, cmd
	case key.Matches(msg, m.keymap.DeleteTask):
		m.StatusQuote = "Delete task? y/n"
		m.pendingAction = deleteTask
		m.state = Confirming
		return m, nil
	}
	return m, nil
}

func (m model) updateTimerRunning(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keymap.Exit):
		m.help = "Please end your session before quitting the app. Press Spacebar/enter to pause the timer"
		return m, nil
	case key.Matches(msg, m.keymap.StartStopTimer):
		m.StatusQuote = "Session ended!!"
		return m.StopSession(), m.Timer.StopCmd()
	case key.Matches(msg, m.keymap.ResetTimer):
		m.StatusQuote = "Are you sure you want to reset this session? Entropy will be added."
		m.state = Confirming
		m.pendingAction = resetTimer
		return m, nil
	}
	if len(m.tasks) == 0 {
		m = m.NoTaskView()
		return m, nil
	}

	return m, nil
}

func (m model) updateTyping(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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
		m.state = TimerNotRunning

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
		m.state = TimerNotRunning
		m.StatusQuote = "Task not created -_-"
		m.textInput.Blur()
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) updateConfirming(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keymap.Exit):
		return m, tea.Quit
	case key.Matches(msg, m.keymap.No):
		m.state = TimerNotRunning
		m.StatusQuote = m.tasks[m.ActiveTaskId].Name
		return m, nil
	case key.Matches(msg, m.keymap.Yes):
		switch m.pendingAction {
		case deleteTask:
			err := m.db.DeleteTask(context.Background(), m.ActiveTaskId)
			if err != nil {
				m.StatusQuote = "Couldn't delete task: " + err.Error()
			}
			m.StatusQuote = "deleted: " + m.tasks[m.ActiveTaskId].Name
			m.state = TimerNotRunning
			return m, m.tabs.DeleteSelectedTaskCmd()
		case resetTimer:
			m = m.ResetSession()
			return m, nil

		}
	}
	return m, nil
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	// err here wont terminate the app, infact the app will launch with default keybindings
	cfg, err := GetConfig("./neg.config.json")

	defer f.Close()
	sqlitedb, err := sql.Open("sqlite3", "./database/appdb.sqlite")
	if err != nil {
		log.Fatalf("Couldn't connect to db: %v", err.Error())
	}
	defer sqlitedb.Close()

	queries := db.New(sqlitedb)

	p := tea.NewProgram(NewModel(queries, cfg, err))

	if _, err := p.Run(); err != nil {
		fmt.Printf("could'nt run program: %v", err)
		os.Exit(1)
	}

}
