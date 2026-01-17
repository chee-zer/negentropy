package main

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/key"
)

type UserConfig struct {
	Keymap               keymap
	MaxProductivityHours int
	Theme                string
	EnableAnimations     bool
}

type rootConfig struct {
	Keymap               keymapConfig `json:"keymap"`
	MaxProductivityHours int          `json:"max_productivity_hours"`
	Theme                string       `json:"theme"`
	EnableAnimations     bool         `json:"enable_animations"`
}

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

func GetConfig(path string) (UserConfig, error) {
	// INFO: initialized with default config
	fileCfg := defaultRootConfig()

	var err error
	var data []byte
	if data, err = os.ReadFile(path); err == nil {
		err = json.Unmarshal(data, &fileCfg)
	}
	//TODO: typos will be ignored, write a strict check later
	return mapToUserConfig(fileCfg), err
}

func defaultRootConfig() rootConfig {
	return rootConfig{
		Keymap:               defaultKeymapConfig(),
		MaxProductivityHours: 8,
		Theme:                "dark",
		EnableAnimations:     true,
	}
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

func mapToUserConfig(cfg rootConfig) UserConfig {
	return UserConfig{
		MaxProductivityHours: cfg.MaxProductivityHours,
		Theme:                cfg.Theme,
		EnableAnimations:     cfg.EnableAnimations,
		Keymap: keymap{
			StartStopTimer: key.NewBinding(
				key.WithKeys(cfg.Keymap.StartStopTimer...),
				key.WithHelp("space/enter", "start/stop timer"),
			),
			Exit: key.NewBinding(
				key.WithKeys(cfg.Keymap.Exit...),
				key.WithHelp("q", "quit"),
			),
			GoRight: key.NewBinding(
				key.WithKeys(cfg.Keymap.GoRight...),
				key.WithHelp("→/l", "next"),
			),
			GoLeft: key.NewBinding(
				key.WithKeys(cfg.Keymap.GoLeft...),
				key.WithHelp("←/h", "prev"),
			),
			CreateTask: key.NewBinding(
				key.WithKeys(cfg.Keymap.CreateTask...),
				key.WithHelp("n", "new task"),
			),
			DeleteTask: key.NewBinding(
				key.WithKeys(cfg.Keymap.DeleteTask...),
				key.WithHelp("x", "delete task"),
			),
			ResetTimer: key.NewBinding(
				key.WithKeys(cfg.Keymap.ResetTimer...),
				key.WithHelp("r", "reset timer"),
			),
			Yes: key.NewBinding(key.WithKeys(cfg.Keymap.Yes...)),
			No:  key.NewBinding(key.WithKeys(cfg.Keymap.No...)),
		},
	}
}
