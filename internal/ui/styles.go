package ui

import "charm.land/lipgloss/v2"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7DD3FC"))

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#475569")).
			Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#94A3B8"))

	valueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F8FAFC"))

	cpuColor = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#38BDF8"))

	ramColor = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FBBF24"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F87171"))
)
