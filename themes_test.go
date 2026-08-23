package namegen

import (
	"testing"
)

func TestNewWithTheme(t *testing.T) {
	themes := []Theme{ThemeDefault, ThemeSpace, ThemeTech}
	for _, theme := range themes {
		g, err := NewWithTheme(theme)
		if err != nil {
			t.Fatalf("NewWithTheme(%d) = %v, want nil", theme, err)
		}
		if name := g.Generate(); name == "" {
			t.Fatalf("NewWithTheme(%d).Generate() returned an empty string", theme)
		}
	}
}

func TestNewWithThemeUnknown(t *testing.T) {
	if _, err := NewWithTheme(Theme(99)); err == nil {
		t.Fatal("NewWithTheme(99) = nil error, want an error for an unrecognized theme")
	}
}

func TestThemeWordListsAreValid(t *testing.T) {
	lists := map[string][2][]string{
		"space": {SpaceAdjectives, SpaceNouns},
		"tech":  {TechAdjectives, TechNouns},
	}
	for name, list := range lists {
		if len(list[0]) == 0 || len(list[1]) == 0 {
			t.Fatalf("%s theme word lists must not be empty", name)
		}
		if _, err := NewWithWords(list[0], list[1]); err != nil {
			t.Fatalf("NewWithWords(%s) = %v, want nil", name, err)
		}
	}
}
