package main

import "github.com/charmbracelet/bubbles/key"

type model struct {
	tabs  Tabs
	timer Model
	// bar Bar
}

// global keymap
// timer - start/stop(space, enter)
// switch timer while timer not running(tab, l)
//
// exit (q)
type keymap struct {
	startStop   key.Binding
	switchTimer key.Binding
	exit        key.Binding
}

func main() {

}
