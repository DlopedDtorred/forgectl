package ui

import (
	"testing"
)

func TestIsTerminal(t *testing.T) {
	result := IsTerminal()
	if _, ok := interface{}(result).(bool); !ok {
		t.Errorf("IsTerminal() did not return a bool, got %T", result)
	}
}

func TestColorConstants(t *testing.T) {
	consts := map[string]string{
		"Reset":     Reset,
		"Bold":      Bold,
		"DimStyle":  DimStyle,
		"Italic":    Italic,
		"Underline": Underline,
		"Red":       Red,
		"Green":     Green,
		"Yellow":    Yellow,
		"Blue":      Blue,
		"Purple":    Purple,
		"Cyan":      Cyan,
		"White":     White,
		"Gray":      Gray,
		"BgRed":     BgRed,
		"BgGreen":   BgGreen,
		"BgYellow":  BgYellow,
		"BgBlue":    BgBlue,
		"BgPurple":  BgPurple,
		"BgCyan":    BgCyan,
	}

	for name, val := range consts {
		if val == "" {
			t.Errorf("Color constant %s is empty", name)
		}
		if val[0] != '\033' {
			t.Errorf("Color constant %s = %q, expected ANSI escape sequence", name, val)
		}
	}
}
