// Package system gathers static information about the host machine: the
// operating system, kernel version, uptime and CPU details.
package system

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
)

// Info holds static information about the machine running the monitor.
type Info struct {
	OS       string
	Platform string
	Kernel   string
	Uptime   uint64
	CPUName  string
	Cores    int32
}

// Load collects the current system information. Missing pieces of information
// are left at their zero value instead of failing the whole call.
func Load() Info {
	info := Info{}

	hostInfo, err := host.Info()
	if err == nil {
		info.OS = hostInfo.Platform
		info.Platform = hostInfo.PlatformVersion
		info.Kernel = hostInfo.KernelVersion
		info.Uptime = hostInfo.Uptime
	}

	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		info.CPUName = cpuInfo[0].ModelName
		info.Cores = cpuInfo[0].Cores
	}

	return info
}

// FormatUptime turns a number of seconds into a compact "1d 2h 3m" string.
func FormatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}
