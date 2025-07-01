package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	currentTime string
	stopwatch   time.Duration
	running     bool
	startTime   time.Time
}

type tickMsg time.Time

func initialModel() model {
	return model{
		currentTime: time.Now().Format("15:04:05"),
		stopwatch:   0,
		running:     false,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		m.currentTime = time.Now().Format("15:04:05")
		if m.running {
			m.stopwatch = time.Since(m.startTime)
		}
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "s":
			if !m.running {
				m.startTime = time.Now().Add(-m.stopwatch)
				m.running = true
			} else {
				m.running = false
			}
		case "r":
			m.stopwatch = 0
			m.running = false
		}
	}
	return m, nil
}

func (m model) View() string {
	stopwatchStr := fmt.Sprintf("%02d:%02d:%02d",
		int(m.stopwatch.Hours()),
		int(m.stopwatch.Minutes())%60,
		int(m.stopwatch.Seconds())%60,
	)

	status := "Stopped"
	if m.running {
		status = "Running"
	}

	return fmt.Sprintf(
		"\nCurrent Time: %s\nStopwatch: %s (%s)\n\nControls:\n[s] Start/Stop Stopwatch\n[r] Reset\n[q] Quit\n",
		m.currentTime, stopwatchStr, status,
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

