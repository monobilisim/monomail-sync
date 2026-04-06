package internal

import (
	"fmt"
	"runtime"
	"time"
)

var startTime = time.Now()

func GetUptime() string {
	duration := time.Since(startTime)
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func GetSystemInfo() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"version":      "1.0.0",
		"uptime":       GetUptime(),
		"goroutines":   runtime.NumGoroutine(),
		"memory_alloc": m.Alloc / 1024 / 1024,
		"memory_total": m.TotalAlloc / 1024 / 1024,
		"memory_sys":   m.Sys / 1024 / 1024,
		"num_cpu":      runtime.NumCPU(),
	}
}

func CheckDB() error {
	return db.Ping()
}
