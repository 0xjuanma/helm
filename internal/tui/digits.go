package tui

import "strings"

var digitPatterns = map[rune][]string{
	'0': {
		"█████",
		"█   █",
		"█   █",
		"█   █",
		"█████",
	},
	'1': {
		"  █  ",
		" ██  ",
		"  █  ",
		"  █  ",
		"█████",
	},
	'2': {
		"█████",
		"    █",
		"█████",
		"█    ",
		"█████",
	},
	'3': {
		"█████",
		"    █",
		"█████",
		"    █",
		"█████",
	},
	'4': {
		"█   █",
		"█   █",
		"█████",
		"    █",
		"    █",
	},
	'5': {
		"█████",
		"█    ",
		"█████",
		"    █",
		"█████",
	},
	'6': {
		"█████",
		"█    ",
		"█████",
		"█   █",
		"█████",
	},
	'7': {
		"█████",
		"    █",
		"   █ ",
		"  █  ",
		"  █  ",
	},
	'8': {
		"█████",
		"█   █",
		"█████",
		"█   █",
		"█████",
	},
	'9': {
		"█████",
		"█   █",
		"█████",
		"    █",
		"█████",
	},
	':': {
		"     ",
		"  █  ",
		"     ",
		"  █  ",
		"     ",
	},
}

func renderLargeTime(timeStr string) string {
	lines := make([]string, 5)
	for _, char := range timeStr {
		pattern, ok := digitPatterns[char]
		if !ok {
			continue
		}
		for i := 0; i < 5; i++ {
			if lines[i] != "" {
				lines[i] += " "
			}
			lines[i] += pattern[i]
		}
	}
	return strings.Join(lines, "\n")
}
