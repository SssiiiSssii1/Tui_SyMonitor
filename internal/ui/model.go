// Package ui contains the Bubble Tea model, its update logic and the rendered
// view for the sysmon terminal interface.
package ui

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"sysmon-tui/internal/monitor"
	"sysmon-tui/internal/processes"
	"sysmon-tui/internal/system"
)

type model struct {
	cpu         float64
	ram         float64
	history     monitor.History
	processes   []processes.Info
	visibleRows []processes.Info
	table       table.Model
	system      system.Info
	query       string
	sortBy      string
	searching   bool
	confirmKill bool
	status      string
	width       int
}

// New builds the initial model with the system information already loaded.
func New() tea.Model {
	m := newModel()
	m.system = system.Load()
	return m
}

func newModel() model {
	columns := []table.Column{
		{Title: "PID", Width: 9},
		{Title: "PROCESS", Width: 24},
		{Title: "CPU %", Width: 10},
		{Title: "MEM %", Width: 10},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	return model{
		table:  t,
		sortBy: "cpu",
		status: "System monitor started",
	}
}

// Init starts the first data load and the periodic refresh timer.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return loadData()
		},
		tickCmd(),
	)
}

func loadData() dataMsg {
	result := dataMsg{}

	snapshot, err := monitor.Collect()
	if err != nil {
		result.err = err
		return result
	}
	result.cpu = snapshot.CPU
	result.ram = snapshot.RAM
	result.metricsOK = true

	found, err := processes.List()
	if err != nil {
		result.err = err
		return result
	}
	result.processes = found

	return result
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *model) refreshTable() {
	rows := make([]table.Row, 0)
	m.visibleRows = processes.Filter(m.processes, m.query)
	processes.Sort(m.visibleRows, m.sortBy)

	for _, p := range m.visibleRows {
		rows = append(rows, table.Row{
			fmt.Sprint(p.PID),
			p.Name,
			fmt.Sprintf("%.1f", p.CPU),
			fmt.Sprintf("%.1f", p.RAM),
		})
	}

	cursor := m.table.Cursor()
	m.table.SetRows(rows)

	if len(rows) > 0 {
		if cursor >= len(rows) {
			cursor = len(rows) - 1
		}
		m.table.SetCursor(cursor)
	}
}

func (m model) selectedProcess() (processes.Info, bool) {
	cursor := m.table.Cursor()

	if cursor < 0 || cursor >= len(m.visibleRows) {
		return processes.Info{}, false
	}

	return m.visibleRows[cursor], true
}
