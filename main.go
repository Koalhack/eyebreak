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
	state      [3]string
	stateIndex int

	timer    timer.Model
	duration time.Duration
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
		m.timer.Timeout = m.duration
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, intKeys):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := m.timer.View()

	s += "\n"

	pad := strings.Repeat(" ", padding)
	s = "Work - " + s
	s += "\n" + pad + "\n\n" + pad

	if m.timer.Timedout() {
		s = "All done!"
	}

	return s
}

func main() {
	const duration time.Duration = time.Second * 20
	fmt.Println(duration)
	m := model{
		duration: duration,
		timer:    timer.NewWithInterval(duration, time.Second),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}

}
