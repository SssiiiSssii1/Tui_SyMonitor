// Command sysmon-tui is a terminal system monitor built with Bubble Tea.
package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"

	"sysmon-tui/internal/ui"
)

func main() {
	program := tea.NewProgram(ui.New())
	if _, err := program.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
