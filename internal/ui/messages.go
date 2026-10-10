package ui

import "sysmon-tui/internal/processes"

// dataMsg carries a freshly collected set of metrics into the model.
type dataMsg struct {
	cpu       float64
	ram       float64
	processes []processes.Info
	metricsOK bool
	err       error
}

// tickMsg triggers the next periodic refresh.
type tickMsg struct{}
