package namegen

import (
	"fmt"
	"math/rand"
	"testing"
)

// wordList builds a synthetic word list of the given size, e.g.
// wordList("adj", 3) -> ["adj0", "adj1", "adj2"]. Real word lists are much
// smaller than what's needed to explore Unique's behavior near exhaustion,
// so benchmarks generate their own.
func wordList(prefix string, n int) []string {
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = fmt.Sprintf("%s%d", prefix, i)
	}
	return words
}

// benchGenerator returns a Generator over nAdj*nNoun combinations, seeded
// deterministically so benchmark runs are comparable across builds.
func benchGenerator(nAdj, nNoun int) *Generator {
	g, err := NewWithWords(wordList("adj", nAdj), wordList("noun", nNoun))
	if err != nil {
		panic(err)
	}
	g.rng = rand.New(rand.NewSource(1))
	return g
}

// BenchmarkUnique measures how Unique's rejection-sampling loop scales as
// the requested count approaches the total number of combinations. Near
// exhaustion, most draws collide with an already-picked name, so the cost
// per accepted name grows sharply rather than staying flat.
func BenchmarkUnique(b *testing.B) {
	const nAdj, nNoun = 100, 100
	total := nAdj * nNoun

	fractions := []struct {
		name string
		frac float64
	}{
		{"10pct", 0.10},
		{"50pct", 0.50},
		{"90pct", 0.90},
		{"99pct", 0.99},
		{"100pct", 1.00},
	}

	for _, f := range fractions {
		count := int(float64(total) * f.frac)
		if count < 1 {
			count = 1
		}
		b.Run(f.name, func(b *testing.B) {
			g := benchGenerator(nAdj, nNoun)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := g.Unique(count); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkUniqueThreeWords covers the same near-exhaustion behavior for
// three-word names, whose combination count and rejection pattern differ
// from the two-word case (the two adjectives must also differ from each
// other).
func BenchmarkUniqueThreeWords(b *testing.B) {
	const nAdj, nNoun = 20, 20
	g := benchGenerator(nAdj, nNoun)
	if err := g.SetWordCount(3); err != nil {
		b.Fatalf("SetWordCount(3) = %v, want nil", err)
	}
	count := int(float64(g.combinationCount()) * 0.90)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := g.Unique(count); err != nil {
			b.Fatal(err)
		}
	}
}
