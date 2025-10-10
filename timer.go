package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Label       string
	Running     bool
	id          int
	Interval    time.Duration
	SessionTime time.Duration
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(m.Interval, func(_ time.Time) tea.Msg {
		return TickMsg{id: m.id}
	})
}

// TODO: for Update
// [] handle play/pause
// [] reset(with prompt-reset just adds into entropy)
// [] complete task (later)

// actually i cannot reuse the timer bubble, its a countdown timer and i need to make a normal timer. Whatever thats called

// [x] msg:
// [x] startstop
// [x] tick

type StartStopMsg struct {
	id      int
	running bool
}

type TickMsg struct {
	id int
}

type ResetMsg struct {
	id int
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
		return ResetMsg{id: m.id}
	}
}

func (m Model) startStop(v bool) tea.Cmd {
	return func() tea.Msg {
		return StartStopMsg{id: m.id, running: v}
	}
}

func (m Model) reset() {
	m.Running = false
	m.SessionTime = 0
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m Model) View() string {
	return ""
}
