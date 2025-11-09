package main

import tea "github.com/charmbracelet/bubbletea"

type TabModel struct {
	activeTab         int
	tasksWithProgress map[int]map[string]float32
}

//functionality needed:
// 1. Display tabs (duh)
// 2. switch tabs when timer not running
// 3. Create new tabs
// 4. Show completion percentages

func (m TabModel) Init() tea.Cmd {
	return nil
}

func (m TabModel) Update(msg tea.Msg) (tea.Cmd, tea.Model) {
	return nil, nil
}

func (m TabModel) View() string {
	return ""
}
