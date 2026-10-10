package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"sysmon-tui/internal/monitor"
	"sysmon-tui/internal/processes"
)

// Update routes every message (keystrokes, window resizes, ticks and data
// loads) to the matching piece of state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.table.SetHeight(max(5, msg.Height-19))
		m.refreshTable()
		return m, nil

	case tickMsg:
		return m, tea.Batch(
			func() tea.Msg {
				return loadData()
			},
			tickCmd(),
		)

	case dataMsg:
		if msg.err != nil {
			m.status = "Refresh error: " + msg.err.Error()

			if msg.metricsOK {
				m.cpu = msg.cpu
				m.ram = msg.ram
				m.history.Add(monitor.Snapshot{CPU: m.cpu, RAM: m.ram})
			}

			return m, nil
		}

		m.cpu = msg.cpu
		m.ram = msg.ram
		m.processes = msg.processes
		m.history.Add(monitor.Snapshot{CPU: m.cpu, RAM: m.ram})
		m.refreshTable()

		return m, nil

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if m.confirmKill {
		return m.handleKillConfirmation(key)
	}

	if m.searching {
		return m.handleSearchKey(msg)
	}

	switch key {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "down", "left", "right", "home", "end", "pgup", "pgdown":
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd

	case "/":
		m.searching = true
		m.query = ""
		m.status = "Type a process name or PID"
		return m, nil

	case "esc":
		m.query = ""
		m.searching = false
		m.status = "Search cleared"
		m.refreshTable()
		return m, nil

	case "c":
		m.sortBy = "cpu"
		m.status = "Sorted by CPU usage"
		m.refreshTable()
		return m, nil

	case "m":
		m.sortBy = "ram"
		m.status = "Sorted by memory usage"
		m.refreshTable()
		return m, nil

	case "n":
		m.sortBy = "name"
		m.status = "Sorted by process name"
		m.refreshTable()
		return m, nil

	case "k":
		p, ok := m.selectedProcess()
		if !ok {
			m.status = "Select a process first"
			return m, nil
		}

		if p.PID == 1 {
			m.status = "Refusing to kill PID 1"
			return m, nil
		}

		m.confirmKill = true
		m.status = fmt.Sprintf("Kill PID %d (%s)? Press y to confirm, any other key to cancel", p.PID, p.Name)
		return m, nil
	}

	return m, nil
}

func (m model) handleKillConfirmation(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "y", "Y":
		p, ok := m.selectedProcess()
		if !ok {
			m.status = "No process selected"
		} else if p.PID == 1 {
			m.status = "Refusing to kill PID 1"
		} else if err := processes.Kill(p.PID); err != nil {
			m.status = "Kill failed: " + err.Error()
		} else {
			m.status = fmt.Sprintf("Kill signal sent to PID %d", p.PID)
		}

	default:
		m.status = "Kill cancelled"
	}

	m.confirmKill = false
	return m, nil
}

func (m model) handleSearchKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "enter":
		m.searching = false
	case "esc":
		m.query = ""
		m.searching = false
		m.refreshTable()
	case "backspace":
		if len(m.query) > 0 {
			m.query = m.query[:len(m.query)-1]
			m.refreshTable()
		}
	default:
		if text := msg.Key().Text; text != "" {
			m.query += text
			m.refreshTable()
		}
	}

	return m, nil
}
