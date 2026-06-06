package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/sadra/tui-dashboard/internal/theme"
)

func RenderGauge(label string, percent float64, width int) string {
	filled := int(float64(width) * percent / 100.0)
	if filled > width {
		filled = width
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	color := colorForPercent(percent)
	styledBar := lipgloss.NewStyle().Foreground(color).Render(bar)

	return fmt.Sprintf("%s [%s] %.1f%%", label, styledBar, percent)
}

func RenderSparkline(label string, history []float64, maxValue float64) string {
	if len(history) == 0 {
		return fmt.Sprintf("%s [no data]", label)
	}

	sparkChars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	sparkline := ""

	for _, val := range history {
		index := int((val / maxValue) * float64(len(sparkChars)-1))
		if index < 0 {
			index = 0
		}
		if index >= len(sparkChars) {
			index = len(sparkChars) - 1
		}
		sparkline += sparkChars[index]
	}

	avg := 0.0
	for _, v := range history {
		avg += v
	}
	avg /= float64(len(history))

	color := colorForPercent(avg)
	styledSparkline := lipgloss.NewStyle().Foreground(color).Render(sparkline)

	return fmt.Sprintf("%s [%s] avg: %.1f%%", label, styledSparkline, avg)
}

func colorForPercent(percent float64) lipgloss.Color {
	if percent >= 90 {
		return theme.Color.Error
	} else if percent >= 75 {
		return theme.Color.Warning
	} else if percent >= 50 {
		return theme.Color.Accent
	}
	return theme.Color.Success
}

func FormatBytes(bytes float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	value := bytes

	for _, unit := range units {
		if value < 1024 {
			return fmt.Sprintf("%.1f %s", value, unit)
		}
		value /= 1024
	}

	return fmt.Sprintf("%.1f PB", value)
}

func FormatBytesPerSec(bytesPerSec float64) string {
	bits := bytesPerSec * 8
	units := []string{"b/s", "Kb/s", "Mb/s", "Gb/s"}
	value := bits

	for _, unit := range units {
		if value < 1000 {
			return fmt.Sprintf("%.1f %s", value, unit)
		}
		value /= 1000
	}

	return fmt.Sprintf("%.1f Tb/s", value)
}
