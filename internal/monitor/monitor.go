// Package monitor collects CPU and memory usage and keeps a short history of
// recent samples for the sparklines shown in the UI.
package monitor

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// maxHistory is the number of samples kept for each sparkline.
const maxHistory = 30

// Snapshot is a single CPU/RAM reading, expressed as percentages.
type Snapshot struct {
	CPU float64
	RAM float64
}

// History keeps the most recent CPU and RAM samples.
type History struct {
	CPU []float64
	RAM []float64
}

// Collect reads the current CPU and memory usage. CPU usage is measured since
// the previous call, so the first reading may be zero.
func Collect() (Snapshot, error) {
	var snapshot Snapshot

	cpuValues, err := cpu.Percent(0, false)
	if err != nil {
		return snapshot, err
	}
	if len(cpuValues) > 0 {
		snapshot.CPU = cpuValues[0]
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return snapshot, err
	}
	snapshot.RAM = memory.UsedPercent

	return snapshot, nil
}

// Add appends a snapshot and drops the oldest samples once the history is full.
func (h *History) Add(snapshot Snapshot) {
	h.CPU = append(h.CPU, snapshot.CPU)
	h.RAM = append(h.RAM, snapshot.RAM)

	h.CPU = trim(h.CPU)
	h.RAM = trim(h.RAM)
}

func trim(values []float64) []float64 {
	if len(values) > maxHistory {
		return values[len(values)-maxHistory:]
	}

	return values
}
