package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	reader = bufio.NewReader(os.Stdin)
)

// Colors
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	DimStyle  = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"

	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Gray   = "\033[90m"

	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
	BgPurple = "\033[45m"
	BgCyan   = "\033[46m"
)

func colorize(color, text string) string {
	return color + text + Reset
}

func Header(text string) {
	fmt.Println()
	fmt.Println(colorize(Bold+Purple, "  ╭─────────────────────────────────────────╮"))
	fmt.Println(colorize(Bold+Purple, "  │") + colorize(Bold+White, "  "+text))
	fmt.Println(colorize(Bold+Purple, "  ╰─────────────────────────────────────────╯"))
	fmt.Println()
}

func HeaderSmall(text string) {
	fmt.Println()
	fmt.Println(colorize(Bold+Cyan, "  ── "+text+" ──"))
	fmt.Println()
}

func Success(msg string) {
	fmt.Println(colorize(Green, "  ✓ ") + msg)
}

func Successf(format string, args ...interface{}) {
	fmt.Println(colorize(Green, "  ✓ ") + fmt.Sprintf(format, args...))
}

func Error(msg string) {
	fmt.Println(colorize(Red, "  ✗ ") + msg)
}

func Errorf(format string, args ...interface{}) {
	fmt.Println(colorize(Red, "  ✗ ") + fmt.Sprintf(format, args...))
}

func Warn(msg string) {
	fmt.Println(colorize(Yellow, "  ⚠ ") + msg)
}

func Warnf(format string, args ...interface{}) {
	fmt.Println(colorize(Yellow, "  ⚠ ") + fmt.Sprintf(format, args...))
}

func Info(msg string) {
	fmt.Println(colorize(Cyan, "  ℹ ") + msg)
}

func Infof(format string, args ...interface{}) {
	fmt.Println(colorize(Cyan, "  ℹ ") + fmt.Sprintf(format, args...))
}

func Step(msg string) {
	fmt.Println(colorize(Blue, "  → ") + msg)
}

func Stepf(format string, args ...interface{}) {
	fmt.Println(colorize(Blue, "  → ") + fmt.Sprintf(format, args...))
}

func Dim(msg string) {
	fmt.Println(colorize(Gray, "    "+msg))
}

func PrintTree(items []string) {
	for i, item := range items {
		prefix := "    "
		if i == len(items)-1 {
			prefix = "    "
		}
		if i == 0 {
			prefix = colorize(Gray, "  ├─ ")
		} else if i == len(items)-1 {
			prefix = colorize(Gray, "  └─ ")
		} else {
			prefix = colorize(Gray, "  ├─ ")
		}
		fmt.Println(prefix + item)
	}
}

func Divider() {
	fmt.Println(colorize(Gray, "  ─────────────────────────────────────────"))
}

func Newline() {
	fmt.Println()
}

// Input prompts
func PromptText(label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("  %s %s %s ", colorize(Bold, label), colorize(Gray, "("+defaultVal+")")+":", "")
	} else {
		fmt.Printf("  %s: ", colorize(Bold, label))
	}
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return defaultVal
	}
	return text
}

func PromptTextRequired(label string) string {
	for {
		fmt.Printf("  %s: ", colorize(Bold, label))
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text != "" {
			return text
		}
		fmt.Println(colorize(Red, "    This field is required."))
	}
}

func PromptPassword(label string) string {
	fmt.Printf("  %s: ", colorize(Bold, label))
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func PromptConfirm(label string, defaultYes bool) bool {
	defaultStr := "y/N"
	if defaultYes {
		defaultStr = "Y/n"
	}
	fmt.Printf("  %s %s: ", colorize(Bold, label), colorize(Gray, "("+defaultStr+")"))
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))

	if text == "" {
		return defaultYes
	}
	return text == "y" || text == "yes"
}

type Choice struct {
	Label       string
	Description string
	Value       string
}

func PromptSelect(label string, choices []Choice) Choice {
	fmt.Printf("  %s\n\n", colorize(Bold, label))
	for i, c := range choices {
		num := colorize(Cyan, fmt.Sprintf("%d", i+1))
		if c.Description != "" {
			fmt.Printf("    %s. %s  %s\n", num, colorize(Bold, c.Label), colorize(Gray, "- "+c.Description))
		} else {
			fmt.Printf("    %s. %s\n", num, colorize(Bold, c.Label))
		}
	}
	fmt.Println()

	for {
		fmt.Printf("  %s ", colorize(Gray, "Select"))
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" && len(choices) > 0 {
			return choices[0]
		}

		idx, err := strconv.Atoi(text)
		if err == nil && idx >= 1 && idx <= len(choices) {
			return choices[idx-1]
		}
		fmt.Println(colorize(Red, "    Invalid selection. Please enter a number."))
	}
}

func PromptMultiSelect(label string, choices []Choice) []Choice {
	fmt.Printf("  %s\n\n", colorize(Bold, label))
	for i, c := range choices {
		num := colorize(Cyan, fmt.Sprintf("%d", i+1))
		if c.Description != "" {
			fmt.Printf("    %s. %s  %s\n", num, colorize(Bold, c.Label), colorize(Gray, "- "+c.Description))
		} else {
			fmt.Printf("    %s. %s\n", num, colorize(Bold, c.Label))
		}
	}
	fmt.Printf("\n  %s\n", colorize(Gray, "Enter comma-separated numbers (e.g., 1,3,5)"))

	for {
		fmt.Printf("  %s ", colorize(Gray, "Select"))
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			return nil
		}

		var selected []Choice
		parts := strings.Split(text, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			idx, err := strconv.Atoi(p)
			if err == nil && idx >= 1 && idx <= len(choices) {
				selected = append(selected, choices[idx-1])
			}
		}
		if len(selected) > 0 {
			return selected
		}
		fmt.Println(colorize(Red, "    Invalid selection."))
	}
}

// Spinner
type Spinner struct {
	message string
	frames  []string
	done    chan bool
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		done:    make(chan bool),
	}
}

func (s *Spinner) Start() {
	i := 0
	go func() {
		for {
			select {
			case <-s.done:
				return
			default:
				fmt.Printf("\r  %s %s", colorize(Cyan, s.frames[i%len(s.frames)]), s.message)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop(success bool) {
	<-s.done
	fmt.Print("\r\033[K")
	if success {
		Success(s.message)
	} else {
		Error(s.message)
	}
}

func (s *Spinner) StopMessage(msg string, success bool) {
	<-s.done
	fmt.Print("\r\033[K")
	if success {
		Success(msg)
	} else {
		Error(msg)
	}
}

func (s *Spinner) Finish() {
	close(s.done)
}

// Box drawing
func Box(title string, lines []string) {
	maxLen := len(title)
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	maxLen += 4

	fmt.Println(colorize(Cyan, "  ╭"+strings.Repeat("─", maxLen)+"╮"))
	fmt.Println(colorize(Cyan, "  │")+colorize(Bold+White, "  "+title)+strings.Repeat(" ", maxLen-len(title)-2)+colorize(Cyan, "│"))
	if len(lines) > 0 {
		fmt.Println(colorize(Cyan, "  ├"+strings.Repeat("─", maxLen)+"┤"))
		for _, l := range lines {
			padding := maxLen - len(l) - 2
			if padding < 0 {
				padding = 0
			}
			fmt.Println(colorize(Cyan, "  │")+"  "+l+strings.Repeat(" ", padding)+colorize(Cyan, "│"))
		}
	}
	fmt.Println(colorize(Cyan, "  ╰"+strings.Repeat("─", maxLen)+"╯"))
}

// Progress bar
func ProgressBar(label string, current, total int) {
	width := 30
	pct := float64(current) / float64(total)
	filled := int(pct * float64(width))

	bar := colorize(Green, strings.Repeat("█", filled)) +
		colorize(Gray, strings.Repeat("░", width-filled))

	percentStr := fmt.Sprintf("%d%%", int(pct*100))
	fmt.Printf("\r  %s [%s] %s", label, bar, percentStr)
	if current == total {
		fmt.Println()
	}
}

// Clear screen
func Clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// IsTerminal checks if stdin is a terminal
func IsTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
