package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"sysmon-tui/internal/system"
)

// View renders the whole dashboard every time the model changes.
func (m model) View() tea.View {
	var b strings.Builder

	header := titleStyle.Render("SYSMON") +
		mutedStyle.Render("  SYSTEM MONITOR")

	b.WriteString(header + "\n")
	b.WriteString(mutedStyle.Render(strings.Repeat("─", max(20, min(90, m.width)))) + "\n\n")

	cpuCard := cardStyle.Render(
		cpuColor.Render("CPU") + "\n" +
			valueStyle.Render(fmt.Sprintf("%.1f%%", m.cpu)) + "\n" +
			bar(m.cpu, 20, cpuColor) + "\n" +
			mutedStyle.Render("History ") + cpuColor.Render(sparkline(m.history.CPU)),
	)

	ramCard := cardStyle.Render(
		ramColor.Render("RAM") + "\n" +
			valueStyle.Render(fmt.Sprintf("%.1f%%", m.ram)) + "\n" +
			bar(m.ram, 20, ramColor) + "\n" +
			mutedStyle.Render("History ") + ramColor.Render(sparkline(m.history.RAM)),
	)

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cpuCard, "  ", ramCard))
	b.WriteString("\n\n")

	b.WriteString(titleStyle.Render("PROCESSES"))
	b.WriteString(mutedStyle.Render(fmt.Sprintf(
		"  %d processes  |  sorted by %s",
		len(m.visibleRows),
		strings.ToUpper(m.sortBy),
	)))
	b.WriteString("\n")

	if m.searching {
		b.WriteString(warningStyle.Render("Search: /"+m.query+"▌") + "\n")
	} else if m.query != "" {
		b.WriteString(mutedStyle.Render("Filter: "+m.query+"  (Esc to clear)") + "\n")
	}

	b.WriteString(m.table.View())
	b.WriteString("\n")

	if p, ok := m.selectedProcess(); ok {
		details := fmt.Sprintf(
			"Selected: %s  |  PID: %d  |  CPU: %.1f%%  |  MEM: %.1f%%",
			p.Name,
			p.PID,
			p.CPU,
			p.RAM,
		)
		b.WriteString(labelStyle.Render(details) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(titleStyle.Render("SYSTEM INFO"))
	b.WriteString("\n")

	systemLine := fmt.Sprintf(
		"OS: %s %s  |  Kernel: %s  |  Uptime: %s  |  CPU cores: %d",
		m.system.OS,
		m.system.Platform,
		m.system.Kernel,
		system.FormatUptime(m.system.Uptime),
		m.system.Cores,
	)

	b.WriteString(mutedStyle.Render(systemLine) + "\n")

	if m.system.CPUName != "" {
		cpuName := m.system.CPUName
		if len(cpuName) > 100 {
			cpuName = cpuName[:97] + "..."
		}
		b.WriteString(mutedStyle.Render("Processor: "+cpuName) + "\n")
	}

	b.WriteString("\n")

	if m.confirmKill {
		b.WriteString(errorStyle.Render("CONFIRM: Press y to kill this process. Any other key cancels.") + "\n")
	} else if m.searching {
		b.WriteString(warningStyle.Render("Type to search  |  Enter: finish  |  Esc: clear search") + "\n")
	} else {
		b.WriteString(mutedStyle.Render(
			"↑/↓ Navigate  / Search  c CPU sort  m RAM sort  n Name sort  k Kill  q Quit",
		) + "\n")
	}

	if m.status != "" {
		b.WriteString(mutedStyle.Render(m.status) + "\n")
	}

	return tea.NewView(b.String())
}

func bar(value float64, width int, color lipgloss.Style) string {
	if width < 1 {
		return ""
	}

	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}

	filled := int(value / 100 * float64(width))
	result := color.Render(strings.Repeat("█", filled))
	result += mutedStyle.Render(strings.Repeat("░", width-filled))

	return result
}

func sparkline(values []float64) string {
	if len(values) == 0 {
		return mutedStyle.Render("Collecting data...")
	}

	levels := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	var result strings.Builder

	for _, value := range values {
		if value < 0 {
			value = 0
		}
		if value > 100 {
			value = 100
		}

		index := int(value / 100 * float64(len(levels)-1))
		result.WriteString(levels[index])
	}

	return result.String()
}
