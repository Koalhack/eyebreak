package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
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
	passed    time.Duration

	timer    timer.Model
	progress progress.Model
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmds []tea.Cmd
		var cmd tea.Cmd

		m.passed += m.timer.Interval
		pct := m.passed.Milliseconds() * 100 / m.durations[m.stateIndex].Milliseconds()
		cmds = append(cmds, m.progress.SetPercent(float64(pct)/100))

		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

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
		m.passed = 0

		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

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
	s := "\n" + pad + m.state[m.stateIndex] + " - " + m.timer.View() +
		"\n" + pad + m.progress.View() + "\n\n" + pad

	if m.timer.Timedout() {
		s = "All done!"
	}

	return s
}

func main() {
	var durations = [2]time.Duration{time.Second * 30, time.Second * 20}
	m := model{
		state:      [2]string{"work", "look"},
		stateIndex: 0,

		durations: durations,

		timer:    timer.NewWithInterval(durations[0], time.Second),
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}

}
