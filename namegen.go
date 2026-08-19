// Package namegen generates readable random names by pairing an adjective
// with a noun, e.g. "brave-falcon". It exists for the cases where you need
// an identifier a person can say out loud and remember: temp resource names,
// game character defaults, sample data in a demo seed script.
package namegen

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	// ErrEmptyWordList is returned when either the adjective or noun list
	// passed to NewWithWords is empty.
	ErrEmptyWordList = errors.New("namegen: adjective and noun lists must not be empty")
	// ErrEmptyWord is returned when a word list contains a blank or
	// whitespace-only entry.
	ErrEmptyWord = errors.New("namegen: word lists must not contain empty strings")
	// ErrInvalidCount is returned when GenerateN or Unique is called with a
	// non-positive count.
	ErrInvalidCount = errors.New("namegen: count must be positive")
	// ErrNotEnoughCombinations is returned by Unique when the requested
	// count exceeds len(adjectives) * len(nouns).
	ErrNotEnoughCombinations = errors.New("namegen: requested more unique names than possible combinations")
)

// Style controls how the two words in a generated name are joined.
type Style int

const (
	// Kebab joins words as "brave-falcon". This is the default.
	Kebab Style = iota
	// Snake joins words as "brave_falcon".
	Snake
	// Space joins words as "brave falcon".
	Space
	// Camel joins words as "braveFalcon". It assumes ASCII word lists;
	// multi-byte runes at the start of a word are not handled specially.
	Camel
)

// Generator produces names by combining a random adjective with a random
// noun. The zero value is not usable; construct one with New, NewSeeded, or
// NewWithWords.
type Generator struct {
	adjectives []string
	nouns      []string
	style      Style
	rng        *rand.Rand
}

// New returns a Generator backed by the built-in DefaultAdjectives and
// DefaultNouns lists, seeded from the current time.
func New() *Generator {
	g, err := NewWithWords(DefaultAdjectives, DefaultNouns)
	if err != nil {
		// DefaultAdjectives/DefaultNouns are fixed and valid, so this
		// would indicate a bug in this package, not caller input.
		panic(err)
	}
	return g
}

// NewSeeded returns a Generator backed by the default word lists whose
// random sequence is fully determined by seed. Two generators created with
// the same seed produce the same sequence of names.
func NewSeeded(seed int64) *Generator {
	g := New()
	g.rng = rand.New(rand.NewSource(seed))
	return g
}

// NewWithWords returns a Generator backed by custom word lists. Both lists
// must be non-empty and free of blank entries.
func NewWithWords(adjectives, nouns []string) (*Generator, error) {
	if len(adjectives) == 0 || len(nouns) == 0 {
		return nil, ErrEmptyWordList
	}
	for _, w := range adjectives {
		if strings.TrimSpace(w) == "" {
			return nil, ErrEmptyWord
		}
	}
	for _, w := range nouns {
		if strings.TrimSpace(w) == "" {
			return nil, ErrEmptyWord
		}
	}

	adjCopy := append([]string(nil), adjectives...)
	nounCopy := append([]string(nil), nouns...)
	return &Generator{
		adjectives: adjCopy,
		nouns:      nounCopy,
		style:      Kebab,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// SetStyle changes how future names are joined.
func (g *Generator) SetStyle(s Style) {
	g.style = s
}

// Generate returns a single random name.
func (g *Generator) Generate() string {
	a := g.adjectives[g.rng.Intn(len(g.adjectives))]
	n := g.nouns[g.rng.Intn(len(g.nouns))]
	return format(g.style, a, n)
}

// GenerateN returns count random names. Names may repeat; use Unique if
// that's not acceptable.
func (g *Generator) GenerateN(count int) ([]string, error) {
	if count <= 0 {
		return nil, ErrInvalidCount
	}
	names := make([]string, count)
	for i := range names {
		names[i] = g.Generate()
	}
	return names, nil
}

// Unique returns count random names with no duplicates. It returns
// ErrNotEnoughCombinations if count exceeds the number of possible
// adjective/noun pairs.
func (g *Generator) Unique(count int) ([]string, error) {
	if count <= 0 {
		return nil, ErrInvalidCount
	}
	total := len(g.adjectives) * len(g.nouns)
	if count > total {
		return nil, fmt.Errorf("%w: requested %d, have %d", ErrNotEnoughCombinations, count, total)
	}

	seen := make(map[string]struct{}, count)
	names := make([]string, 0, count)
	for len(names) < count {
		name := g.Generate()
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	return names, nil
}

func format(style Style, a, b string) string {
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))
	switch style {
	case Snake:
		return a + "_" + b
	case Space:
		return a + " " + b
	case Camel:
		return a + strings.ToUpper(b[:1]) + b[1:]
	default: // Kebab
		return a + "-" + b
	}
}
