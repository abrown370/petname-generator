package namegen

import (
	"errors"
	"strings"
	"testing"
)

func TestNewWithWords(t *testing.T) {
	cases := []struct {
		name       string
		adjectives []string
		nouns      []string
		wantErr    error
	}{
		{"nil adjectives", nil, []string{"cat"}, ErrEmptyWordList},
		{"nil nouns", []string{"red"}, nil, ErrEmptyWordList},
		{"both nil", nil, nil, ErrEmptyWordList},
		{"empty slice adjectives", []string{}, []string{"cat"}, ErrEmptyWordList},
		{"blank adjective", []string{"red", "  "}, []string{"cat"}, ErrEmptyWord},
		{"tab-only noun", []string{"red"}, []string{"\t"}, ErrEmptyWord},
		{"single word each", []string{"red"}, []string{"cat"}, nil},
		{"normal lists", []string{"red", "blue"}, []string{"cat", "dog"}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewWithWords(tc.adjectives, tc.nouns)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("NewWithWords() = %v, want nil", err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewWithWords() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	cases := []struct {
		name  string
		style Style
		a, b  string
		want  string
	}{
		{"kebab basic", Kebab, "red", "fox", "red-fox"},
		{"kebab trims whitespace", Kebab, " red ", " fox ", "red-fox"},
		{"kebab lowercases", Kebab, "Red", "FOX", "red-fox"},
		{"snake basic", Snake, "RED", "Fox", "red_fox"},
		{"space basic", Space, "Red", "Fox", "red fox"},
		{"camel basic", Camel, "red", "fox", "redFox"},
		{"camel single-letter noun", Camel, "red", "a", "redA"},
		{"camel single-letter adjective", Camel, "a", "fox", "aFox"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := format(tc.style, "", tc.a, tc.b)
			if got != tc.want {
				t.Fatalf("format(%v, %q, %q) = %q, want %q", tc.style, tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestFormatCustomSeparator(t *testing.T) {
	got := format(Custom, ".", "red", "fox")
	if want := "red.fox"; got != want {
		t.Fatalf("format(Custom, %q, ...) = %q, want %q", ".", got, want)
	}
}

func TestFormatThreeWords(t *testing.T) {
	cases := []struct {
		name  string
		style Style
		sep   string
		words []string
		want  string
	}{
		{"kebab", Kebab, "", []string{"brave", "swift", "falcon"}, "brave-swift-falcon"},
		{"snake", Snake, "", []string{"brave", "swift", "falcon"}, "brave_swift_falcon"},
		{"camel", Camel, "", []string{"brave", "swift", "falcon"}, "braveSwiftFalcon"},
		{"custom", Custom, "_", []string{"brave", "swift", "falcon"}, "brave_swift_falcon"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := format(tc.style, tc.sep, tc.words...)
			if got != tc.want {
				t.Fatalf("format(%v, %q, %v) = %q, want %q", tc.style, tc.sep, tc.words, got, tc.want)
			}
		})
	}
}

func TestSetWordCount(t *testing.T) {
	g, err := NewWithWords([]string{"red", "blue", "green"}, []string{"fox"})
	if err != nil {
		t.Fatalf("NewWithWords() = %v", err)
	}

	if err := g.SetWordCount(4); !errors.Is(err, ErrInvalidWordCount) {
		t.Fatalf("SetWordCount(4) = %v, want %v", err, ErrInvalidWordCount)
	}

	if err := g.SetWordCount(3); err != nil {
		t.Fatalf("SetWordCount(3) = %v, want nil", err)
	}

	for i := 0; i < 50; i++ {
		name := g.Generate()
		parts := strings.Split(name, "-")
		if len(parts) != 3 {
			t.Fatalf("Generate() = %q, want 3 hyphen-separated words", name)
		}
		if parts[0] == parts[1] {
			t.Fatalf("Generate() = %q, want the two adjectives to differ", name)
		}
	}
}

func TestSetWordCountThreeWithSingleAdjective(t *testing.T) {
	g, err := NewWithWords([]string{"red"}, []string{"fox", "owl"})
	if err != nil {
		t.Fatalf("NewWithWords() = %v", err)
	}
	if err := g.SetWordCount(3); err != nil {
		t.Fatalf("SetWordCount(3) = %v, want nil", err)
	}

	name := g.Generate()
	if name != "red-red-fox" && name != "red-red-owl" {
		t.Fatalf("Generate() = %q, want red-red-fox or red-red-owl", name)
	}
}

func TestSetCustomSeparator(t *testing.T) {
	g, err := NewWithWords([]string{"red"}, []string{"fox"})
	if err != nil {
		t.Fatalf("NewWithWords() = %v", err)
	}
	g.SetCustomSeparator("::")

	if got, want := g.Generate(), "red::fox"; got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}

func TestUniqueWithThreeWords(t *testing.T) {
	g, err := NewWithWords([]string{"red", "blue"}, []string{"fox", "owl"})
	if err != nil {
		t.Fatalf("NewWithWords() = %v", err)
	}
	if err := g.SetWordCount(3); err != nil {
		t.Fatalf("SetWordCount(3) = %v, want nil", err)
	}

	// 2 adjectives give 2*1=2 ordered distinct pairs, times 2 nouns = 4.
	if _, err := g.Unique(5); !errors.Is(err, ErrNotEnoughCombinations) {
		t.Fatalf("Unique(5) = %v, want %v", err, ErrNotEnoughCombinations)
	}

	names, err := g.Unique(4)
	if err != nil {
		t.Fatalf("Unique(4) = %v, want nil", err)
	}
	if len(names) != 4 {
		t.Fatalf("Unique(4) returned %d names, want 4", len(names))
	}
}

func TestGenerateNValidatesCount(t *testing.T) {
	g, err := NewWithWords([]string{"red"}, []string{"fox"})
	if err != nil {
		t.Fatalf("NewWithWords() = %v", err)
	}

	cases := []struct {
		name  string
		count int
	}{
		{"zero", 0},
		{"negative", -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := g.GenerateN(tc.count); !errors.Is(err, ErrInvalidCount) {
				t.Fatalf("GenerateN(%d) = %v, want %v", tc.count, err, ErrInvalidCount)
			}
			if _, err := g.Unique(tc.count); !errors.Is(err, ErrInvalidCount) {
				t.Fatalf("Unique(%d) = %v, want %v", tc.count, err, ErrInvalidCount)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	cases := []struct {
		name       string
		adjectives []string
		nouns      []string
		count      int
		wantErr    error
	}{
		{"exact single combination", []string{"red"}, []string{"fox"}, 1, nil},
		{"asking for more than the one combination", []string{"red"}, []string{"fox"}, 2, ErrNotEnoughCombinations},
		{"exhausts all combinations", []string{"red", "blue"}, []string{"fox", "owl"}, 4, nil},
		{"one past exhaustion", []string{"red", "blue"}, []string{"fox", "owl"}, 5, ErrNotEnoughCombinations},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g, err := NewWithWords(tc.adjectives, tc.nouns)
			if err != nil {
				t.Fatalf("NewWithWords() = %v", err)
			}

			names, err := g.Unique(tc.count)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("Unique(%d) = %v, want %v", tc.count, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unique(%d) = %v, want nil", tc.count, err)
			}
			if len(names) != tc.count {
				t.Fatalf("Unique(%d) returned %d names, want %d", tc.count, len(names), tc.count)
			}
			seen := make(map[string]bool, len(names))
			for _, n := range names {
				if seen[n] {
					t.Fatalf("Unique(%d) returned duplicate name %q", tc.count, n)
				}
				seen[n] = true
			}
		})
	}
}

func TestSeededGeneratorsAreReproducible(t *testing.T) {
	a := NewSeeded(42)
	b := NewSeeded(42)

	gotA, err := a.GenerateN(10)
	if err != nil {
		t.Fatalf("GenerateN() = %v", err)
	}
	gotB, err := b.GenerateN(10)
	if err != nil {
		t.Fatalf("GenerateN() = %v", err)
	}

	if len(gotA) != len(gotB) {
		t.Fatalf("got %d and %d names, want equal lengths", len(gotA), len(gotB))
	}
	for i := range gotA {
		if gotA[i] != gotB[i] {
			t.Fatalf("name %d differs between same-seed generators: %q != %q", i, gotA[i], gotB[i])
		}
	}
}

func TestDefaultWordListsAreValid(t *testing.T) {
	if len(DefaultAdjectives) == 0 || len(DefaultNouns) == 0 {
		t.Fatal("default word lists must not be empty")
	}
	if _, err := NewWithWords(DefaultAdjectives, DefaultNouns); err != nil {
		t.Fatalf("NewWithWords(DefaultAdjectives, DefaultNouns) = %v", err)
	}
}
