// Package namegen generates readable random names by pairing an adjective
// with a noun, e.g. "brave-falcon". It exists for the cases where you need
// an identifier a person can say out loud and remember: temp resource names,
// game character defaults, sample data in a demo seed script.
package namegen

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
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
	// count exceeds the number of combinations the word lists can produce.
	ErrNotEnoughCombinations = errors.New("namegen: requested more unique names than possible combinations")
	// ErrInvalidWordCount is returned by SetWordCount for any value other
	// than 2 or 3.
	ErrInvalidWordCount = errors.New("namegen: word count must be 2 or 3")
)

// Style controls how the words in a generated name are joined.
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
	// Custom joins words with the separator set via SetCustomSeparator.
	Custom
)

// Generator produces names by combining random adjectives with a random
// noun. The zero value is not usable; construct one with New, NewSeeded, or
// NewWithWords.
type Generator struct {
	adjectives []string
	nouns      []string
	style      Style
	separator  string
	wordCount  int
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
		wordCount:  2,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// SetStyle changes how future names are joined.
func (g *Generator) SetStyle(s Style) {
	g.style = s
}

// SetCustomSeparator switches future names to the Custom style, joined with
// sep instead of one of the built-in separators, e.g. SetCustomSeparator(".")
// produces "brave.falcon".
func (g *Generator) SetCustomSeparator(sep string) {
	g.separator = sep
	g.style = Custom
}

// SetWordCount configures how many words a generated name contains. Only 2
// (adjective-noun, the default) and 3 (adjective-adjective-noun) are
// supported. The two adjectives in a three-word name are distinct whenever
// the adjective list has more than one entry.
func (g *Generator) SetWordCount(n int) error {
	if n != 2 && n != 3 {
		return fmt.Errorf("%w: got %d", ErrInvalidWordCount, n)
	}
	g.wordCount = n
	return nil
}

// ExcludeWords removes any adjective or noun that exactly matches one of the
// given words (case-insensitive, whitespace-trimmed). It leaves the word
// lists unchanged and returns ErrEmptyWordList if the removal would empty
// out the adjective or noun list.
func (g *Generator) ExcludeWords(words ...string) error {
	exclude := make(map[string]struct{}, len(words))
	for _, w := range words {
		exclude[strings.ToLower(strings.TrimSpace(w))] = struct{}{}
	}
	keep := func(w string) bool {
		_, excluded := exclude[strings.ToLower(w)]
		return !excluded
	}
	return g.applyWordFilter(keep)
}

// ExcludePattern removes any adjective or noun matching the given regular
// expression, as accepted by the regexp package. It leaves the word lists
// unchanged and returns an error if pattern fails to compile, or
// ErrEmptyWordList if the removal would empty out the adjective or noun
// list.
func (g *Generator) ExcludePattern(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("namegen: invalid exclude pattern: %w", err)
	}
	return g.applyWordFilter(func(w string) bool { return !re.MatchString(w) })
}

// applyWordFilter replaces the adjective and noun lists with the subset for
// which keep returns true. It only commits the change if both resulting
// lists are non-empty, so a filter that would exhaust one list leaves the
// generator in its previous, usable state.
func (g *Generator) applyWordFilter(keep func(string) bool) error {
	newAdj := filterWords(g.adjectives, keep)
	newNouns := filterWords(g.nouns, keep)
	if len(newAdj) == 0 || len(newNouns) == 0 {
		return ErrEmptyWordList
	}
	g.adjectives = newAdj
	g.nouns = newNouns
	return nil
}

// filterWords returns the entries of words for which keep returns true.
func filterWords(words []string, keep func(string) bool) []string {
	out := make([]string, 0, len(words))
	for _, w := range words {
		if keep(w) {
			out = append(out, w)
		}
	}
	return out
}

// Generate returns a single random name.
func (g *Generator) Generate() string {
	return format(g.style, g.separator, g.pick()...)
}

// pick returns the words for one name, in output order.
func (g *Generator) pick() []string {
	if g.wordCount == 3 {
		first := g.adjectives[g.rng.Intn(len(g.adjectives))]
		second := first
		if len(g.adjectives) > 1 {
			for second == first {
				second = g.adjectives[g.rng.Intn(len(g.adjectives))]
			}
		}
		return []string{first, second, g.nouns[g.rng.Intn(len(g.nouns))]}
	}
	return []string{g.adjectives[g.rng.Intn(len(g.adjectives))], g.nouns[g.rng.Intn(len(g.nouns))]}
}

// combinationCount returns the number of distinct names the current word
// lists and word count can produce.
func (g *Generator) combinationCount() int {
	if g.wordCount == 3 {
		pairs := 1
		if len(g.adjectives) > 1 {
			pairs = len(g.adjectives) * (len(g.adjectives) - 1)
		}
		return pairs * len(g.nouns)
	}
	return len(g.adjectives) * len(g.nouns)
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
	total := g.combinationCount()
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

// format joins words according to style, using sep only when style is
// Custom.
func format(style Style, sep string, words ...string) string {
	clean := make([]string, len(words))
	for i, w := range words {
		clean[i] = strings.ToLower(strings.TrimSpace(w))
	}
	switch style {
	case Snake:
		return strings.Join(clean, "_")
	case Space:
		return strings.Join(clean, " ")
	case Camel:
		var b strings.Builder
		b.WriteString(clean[0])
		for _, w := range clean[1:] {
			b.WriteString(strings.ToUpper(w[:1]))
			b.WriteString(w[1:])
		}
		return b.String()
	case Custom:
		return strings.Join(clean, sep)
	default: // Kebab
		return strings.Join(clean, "-")
	}
}
