package cli

import (
	"fmt"
	"time"
)

func formatShutdownCountdown(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	var nonzeroCount uint8
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		nonzeroCount++
	}
	if minutes > 0 {
		nonzeroCount++
	}
	if seconds > 0 {
		nonzeroCount++
	}

	if nonzeroCount == 0 {
		return "00 Hours 00 Minutes 00 Seconds"
	}

	if nonzeroCount == 1 {
		if hours > 0 {
			return fmt.Sprintf("%d Hours", hours)
		} else if minutes > 0 {
			return fmt.Sprintf("%d Minutes", minutes)
		} else {
			return fmt.Sprintf("%d Seconds", seconds)
		}
	}
	return fmt.Sprintf("%d Hours %d Minutes %d Seconds",
		hours, minutes, seconds)
}
