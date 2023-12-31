package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	beeep "github.com/gen2brain/beeep"
)

var (
	boldStyle   = lipgloss.NewStyle().Bold(true)
	italicStyle = lipgloss.NewStyle().Italic(true)
)

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
	percent  float64

	keymap   keymap
	help     help.Model
	quitting bool
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd

		m.passed += m.timer.Interval
		pct := m.passed.Milliseconds() * 100 / m.durations[m.stateIndex].Milliseconds()
		m.percent = float64(pct) / 100

		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)

		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		if m.stateIndex == len(m.state)-1 {
			m.stateIndex = 0
		} else {
			m.stateIndex++
		}

		m.timer.Timeout = m.durations[m.stateIndex]
		m.passed = 0

		beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
		beeep.Notify(m.state[m.stateIndex], "", "")

		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

		// NOTE: Listen all Key Event
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.passed = 0
			m.stateIndex = 0
			m.timer.Timeout = m.durations[m.stateIndex]
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}

func (m model) helpView() string {
	return m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	s := "\n" + pad + italicStyle.Render(m.state[m.stateIndex]) + " - " + boldStyle.Render(m.timer.View()) +
		"\n\n" + pad + m.progress.ViewAs(m.percent) + "\n\n" + pad + m.helpView() + "\n"
	return s
}

func Start() {
	var durations = [2]time.Duration{time.Minute * 20, time.Second * 20}
	m := model{
		state:      [2]string{"Work", "Look"},
		stateIndex: 0,

		durations: durations,

		timer:    timer.NewWithInterval(durations[0], time.Second),
		progress: progress.New(progress.WithDefaultGradient()),

		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
	m.keymap.start.SetEnabled(false)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
