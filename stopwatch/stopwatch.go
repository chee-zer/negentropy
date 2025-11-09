package stopwatch

import (
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var lastID int64

type Model struct {
	Label       string
	running     bool
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

func (m Model) StartCmd() tea.Cmd {
	return m.startStop(true)
}

func (m Model) StopCmd() tea.Cmd {
	return m.startStop(false)
}

func (m Model) ResetCmd() tea.Cmd {
	return func() tea.Msg {
		return ResetMsg{Id: m.id}
	}
}

func (m Model) startStop(v bool) tea.Cmd {
	return func() tea.Msg {
		return StartStopMsg{Id: m.id, Running: v}
	}
}

func NewTimer(label string) Model {
	return Model{
		id:       nextID(),
		Label:    label,
		running:  false,
		Interval: time.Second,
	}
}

// helper funcs
// concurrent safe incrementer because.
func nextID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

func (m *Model) reset() {
	m.running = false
	m.SessionTime = 0
	m.tag++
}

func (m *Model) IsRunning() bool {
	return m.running
}

// tag's purpose is to prevent race conditions. Here it'll be incremented everytime user resets.
// This will prevent the bug where the user maybe uses a macro to reset and start immediately.
// just for that tho, there is no functionality to change tick interval duration for now(no reason to do that in a productivity timer)
// So we'll just check the tag in tickmsg case
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StartStopMsg:
		// id 0 is master control (no use for now)
		if msg.Id != 0 && msg.Id != m.id {
			return m, nil
		}
		m.running = msg.Running
		return m, m.tick()

	case ResetMsg:
		if m.id != 0 && m.id != msg.Id {
			return m, nil
		}
		m.reset()
		return m, nil

	case TickMsg:
		if !m.running || m.tag != msg.Tag || m.id != msg.Id {
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
