// Package processes collects the running processes from the operating system
// and provides filtering, sorting and termination helpers.
package processes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v4/process"
)

// Info describes a single running process.
type Info struct {
	PID  int32
	Name string
	CPU  float64
	RAM  float32
}

// List returns every process that exposes a readable name.
func List() ([]Info, error) {
	running, err := process.Processes()
	if err != nil {
		return nil, err
	}

	result := make([]Info, 0, len(running))

	for _, p := range running {
		name, err := p.Name()
		if err != nil {
			continue
		}

		processCPU, err := p.Percent(0)
		if err != nil {
			processCPU = 0
		}

		processRAM, err := p.MemoryPercent()
		if err != nil {
			processRAM = 0
		}

		result = append(result, Info{
			PID:  p.Pid,
			Name: name,
			CPU:  processCPU,
			RAM:  processRAM,
		})
	}

	return result, nil
}

// Filter returns a new slice containing only the processes whose name or PID
// matches the query. An empty query returns a copy of every process.
func Filter(all []Info, query string) []Info {
	filtered := make([]Info, 0, len(all))

	if query == "" {
		filtered = append(filtered, all...)
		return filtered
	}

	lowerQuery := strings.ToLower(query)

	for _, p := range all {
		if !strings.Contains(strings.ToLower(p.Name), lowerQuery) &&
			!strings.Contains(fmt.Sprint(p.PID), query) {
			continue
		}

		filtered = append(filtered, p)
	}

	return filtered
}

// Sort orders the slice in place. Valid keys are "cpu", "ram" and "name";
// any other value leaves the order unchanged.
func Sort(all []Info, by string) {
	switch by {
	case "cpu":
		sort.SliceStable(all, func(i, j int) bool {
			return all[i].CPU > all[j].CPU
		})
	case "ram":
		sort.SliceStable(all, func(i, j int) bool {
			return all[i].RAM > all[j].RAM
		})
	case "name":
		sort.SliceStable(all, func(i, j int) bool {
			return strings.ToLower(all[i].Name) < strings.ToLower(all[j].Name)
		})
	}
}

// Kill sends a termination signal to the process with the given PID. It is
// only ever called after the user has explicitly confirmed.
func Kill(pid int32) error {
	target, err := process.NewProcess(pid)
	if err != nil {
		return err
	}

	return target.Kill()
}
