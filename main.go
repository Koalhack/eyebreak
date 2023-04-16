package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

var intKeys = key.NewBinding(key.WithKeys("ctrl+c"))

const (
	padding  = 2
	maxWidth = 80
)

type model struct {
	stateIndex int
	state      [2]string

	durations [2]time.Duration
	timer     timer.Model
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		if m.stateIndex == len(m.state)-1 {
			m.stateIndex = 0
		} else {
			m.stateIndex++
		}

		m.timer.Timeout = m.durations[m.stateIndex]

		return m, nil
		// NOTE: Listen all Key Event
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, intKeys):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	s := m.timer.View()
	s += "\n"
	s = m.state[m.stateIndex] + " - " + s
	s += "\n" + pad + "\n\n" + pad

	if m.timer.Timedout() {
		s = "All done!"
	}

	return s
}

func main() {
	var durations = [2]time.Duration{time.Second * 30, time.Second * 20}
	fmt.Println(durations)
	m := model{
		state:      [2]string{"work", "look"},
		stateIndex: 0,

		durations: durations,
		timer:     timer.NewWithInterval(durations[0], time.Second),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}

}
