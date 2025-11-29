package stopwatch

import (
	"log"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var lastID int64

type Model struct {
	Label       string
	Running     bool
	id          int
	tag         int
	Interval    time.Duration
	SessionTime time.Duration
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(m.Interval, func(_ time.Time) tea.Msg {
		return TickMsg{Id: m.id, Tag: m.tag}
	})
}

// TODO: for Update
// [x] handle play/pause
// [x] reset(with prompt-reset just adds into entropy)
// [] complete task (later)

// actually i cannot reuse the timer bubble, its a countdown timer and i need to make a normal timer. Whatever thats called

// [x] msg:
// [x] startstop
// [x] tick

type StartStopMsg struct {
	Id      int
	Running bool
}

type TickMsg struct {
	Id  int
	Tag int
}

type ResetMsg struct {
	Id int
}

// [x]cmd:
// [x]start
// [x] stop
// [x] reset

func (m Model) ID() int {
	return m.id
}

// Starts the timer
func (m Model) StartCmd() tea.Cmd {
	return m.startStop(true)
}

// Stops the timer
func (m Model) StopCmd() tea.Cmd {
	return m.startStop(false)
}

// Sends a ResetMsg
func (m Model) ResetCmd() tea.Cmd {
	return func() tea.Msg {
		return ResetMsg{Id: m.id}
	}
}

func (m Model) startStop(v bool) tea.Cmd {
	return func() tea.Msg {
		log.Println("value passed: ", v)
		log.Println("current value m.Running: ", m.Running)
		return StartStopMsg{Id: m.id, Running: v}
	}
}

func NewTimer(label string) Model {
	return Model{
		id:       nextID(),
		Label:    label,
		Running:  false,
		Interval: time.Second,
	}
}

func NewTimerRunning(label string) Model {
	m := NewTimer(label)
	m.Running = true
	return m
}

// helper funcs
// concurrent safe incrementer because.
func nextID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

func (m *Model) reset() {
	m.Running = false
	m.SessionTime = 0
	m.tag++
}

func (m *Model) IsRunning() bool {
	return m.Running
}

// tag's purpose is to prevent race conditions. Here it'll be incremented everytime user resets.
// This will prevent the bug where the user maybe uses a macro to reset and start immediately.
// just for that tho, there is no functionality to change tick interval duration for now(no reason to do that in a productivity timer)
// So we'll just check the tag in tickmsg case
// Note: this is a submodel, so it can return Model instead of tea.Model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StartStopMsg:
		log.Println("\n\n entered startstopmsg update")
		// id 0 is master control (no use for now)
		if msg.Id != 0 && msg.Id != m.id {
			log.Println("this shouldnt be printed")
			return m, nil
		}
		m.Running = msg.Running
		log.Printf("\nmrunning: %v\tmsgrunning: %v\n", m.Running, msg.Running)
		if m.Running {
			return m, m.tick()
		}
		return m, nil

	case ResetMsg:
		if m.id != 0 && m.id != msg.Id {
			return m, nil
		}
		m.reset()
		return m, nil

	case TickMsg:
		if !m.Running || m.tag != msg.Tag || m.id != msg.Id {
			return m, nil
		}
		m.SessionTime += m.Interval
		return m, m.tick()
	}
	return m, nil
}

func (m Model) View() string {

	return m.SessionTime.String()
}
