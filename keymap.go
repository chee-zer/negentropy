package main

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/key"
)

type keymapConfig struct {
	StartStopTimer []string `json:"start_stop_timer"`
	Exit           []string `json:"exit"`
	GoRight        []string `json:"go_right"`
	GoLeft         []string `json:"go_left"`
	DeleteTask     []string `json:"delete_task"`
	CreateTask     []string `json:"create_task"`
	ResetTimer     []string `json:"reset_timer"`
	Yes            []string `json:"yes"`
	No             []string `json:"no"`
}

func GetConfig(path string) (keymap, error) {
	config := defaultKeymapConfig()
	var err error
	var data []byte
	if data, err = os.ReadFile(path); err == nil {
		err = json.Unmarshal(data, &config)
	}
	return mapConfigToKeymap(config), err
}

func defaultKeymapConfig() keymapConfig {
	return keymapConfig{
		StartStopTimer: []string{" ", "enter"},
		Exit:           []string{"q", "ctrl+c"},
		GoRight:        []string{"right", "l"},
		GoLeft:         []string{"left", "h"},
		CreateTask:     []string{"n"},
		DeleteTask:     []string{"x"},
		ResetTimer:     []string{"r"},
		Yes:            []string{"y"},
		No:             []string{"n"},
	}
}

func mapConfigToKeymap(cfg keymapConfig) keymap {
	return keymap{
		StartStopTimer: key.NewBinding(
			key.WithKeys(cfg.StartStopTimer...),
			key.WithHelp("space/enter", "start/stop timer"),
		),
		Exit: key.NewBinding(
			key.WithKeys(cfg.Exit...),
			key.WithHelp("q", "quit"),
		),
		GoRight: key.NewBinding(
			key.WithKeys(cfg.GoRight...),
			key.WithHelp("→/l", "next"),
		),
		GoLeft: key.NewBinding(
			key.WithKeys(cfg.GoLeft...),
			key.WithHelp("←/h", "prev"),
		),
		CreateTask: key.NewBinding(
			key.WithKeys(cfg.CreateTask...),
			key.WithHelp("n", "new task"),
		),
		DeleteTask: key.NewBinding(
			key.WithKeys(cfg.DeleteTask...),
			key.WithHelp("x", "delete task"),
		),
		ResetTimer: key.NewBinding(
			key.WithKeys(cfg.ResetTimer...),
			key.WithHelp("r", "reset timer"),
		),
		Yes: key.NewBinding(key.WithKeys(cfg.Yes...)),
		No:  key.NewBinding(key.WithKeys(cfg.No...)),
	}
}
