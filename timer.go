package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Label   string
	Running bool
	ID      int
}

func (m Model) Init() tea.Cmd {
	return nil
}

// TODO: for Update
// [] handle play/pause
// [] reset(with prompt-reset just adds into entropy)
// [] complete task (later)

// actually i cannot reuse the timer bubble, its a countdown timer and i need to make a normal timer. Whatever thats called
// []cmd:
// []start
// [] stop
// [] reset

// [] msg:
// [] startstop
// [] tick

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m Model) View() string {
	return ""
}
