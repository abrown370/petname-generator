package namegen

import "fmt"

// Theme selects a built-in word list pairing for NewWithTheme, as an
// alternative to hand-supplying lists via NewWithWords.
type Theme int

const (
	// ThemeDefault pairs DefaultAdjectives with DefaultNouns.
	ThemeDefault Theme = iota
	// ThemeSpace pairs SpaceAdjectives with SpaceNouns.
	ThemeSpace
	// ThemeTech pairs TechAdjectives with TechNouns.
	ThemeTech
)

var themeWordLists = map[Theme][2][]string{
	ThemeDefault: {DefaultAdjectives, DefaultNouns},
	ThemeSpace:   {SpaceAdjectives, SpaceNouns},
	ThemeTech:    {TechAdjectives, TechNouns},
}

// NewWithTheme returns a Generator backed by the word lists for theme,
// seeded from the current time. It returns an error only for a Theme value
// that doesn't match one of the built-in constants; every built-in theme's
// word lists are valid by construction.
func NewWithTheme(theme Theme) (*Generator, error) {
	lists, ok := themeWordLists[theme]
	if !ok {
		return nil, fmt.Errorf("namegen: unknown theme %d", theme)
	}
	return NewWithWords(lists[0], lists[1])
}
