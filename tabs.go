package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	db "github.com/chee-zer/negentropy/database/sqlc"
)

type TabModel struct {
	ActiveTabIndex int
	//TasksWithProgress map[int]map[string]float32
	Tasks []db.Task
}

// msg for switching tabs/tasks.
// Direction should be set to true for switching to very left task,
// and false for right
type SwitchMsg struct {
	direction bool
}

type SwitchSelectedTaskMsg struct {
	taskID int64
}

func (m TabModel) SwitchSelectedTaskCmd() tea.Cmd {
	activeTask := m.Tasks[m.ActiveTabIndex]
	return func() tea.Msg {
		return SwitchSelectedTaskMsg{
			taskID: activeTask.ID,
		}
	}
}

func (m TabModel) SwitchLeftCmd() tea.Cmd {
	return func() tea.Msg {
		return SwitchMsg{
			direction: false,
		}
	}
}
func (m TabModel) SwitchRightCmd() tea.Cmd {
	return func() tea.Msg {
		return SwitchMsg{
			direction: true,
		}
	}
}

// returns a new Tab Model
func NewTabModel(tasks []db.Task) TabModel {
	if len(tasks) == 0 {
		return TabModel{
			ActiveTabIndex: -1,
			//TasksWithProgress: nil,
			Tasks: nil,
		}
	}

	return TabModel{
		ActiveTabIndex: 0,
		Tasks:          tasks,
	}
}

func (m TabModel) Init() tea.Cmd {
	return nil
}

func (m TabModel) Update(msg tea.Msg) (TabModel, tea.Cmd) {
	switch msg := msg.(type) {
	case SwitchMsg:
		if msg.direction {
			m.ActiveTabIndex = (len(m.Tasks) + m.ActiveTabIndex + 1) % len(m.Tasks)
		} else {
			m.ActiveTabIndex = (len(m.Tasks) + m.ActiveTabIndex - 1) % len(m.Tasks)
		}
		return m, m.SwitchSelectedTaskCmd()
	}
	return m, nil
}

func (m TabModel) View() string {
	output := ""
	for i, task := range m.Tasks {
		marker := " "
		if i == m.ActiveTabIndex {
			marker = ">"
		}
		output += fmt.Sprintf("\n%s%d: %s", marker, i, task.Name)
	}
	return output
}
