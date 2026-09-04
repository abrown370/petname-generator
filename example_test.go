package namegen

import (
	"fmt"
	"sort"
)

// Single-word lists make these examples deterministic: with only one
// adjective and one noun to choose from, the RNG has nothing to decide.

func ExampleNewWithWords() {
	g, err := NewWithWords([]string{"brave"}, []string{"falcon"})
	if err != nil {
		panic(err)
	}
	fmt.Println(g.Generate())
	// Output: brave-falcon
}

func ExampleNewWithTheme() {
	g, err := NewWithTheme(ThemeSpace)
	if err != nil {
		panic(err)
	}
	name := g.Generate()
	fmt.Println(name != "")
	// Output: true
}

func ExampleGenerator_Generate() {
	g, err := NewWithWords([]string{"brave"}, []string{"falcon"})
	if err != nil {
		panic(err)
	}
	fmt.Println(g.Generate())
	// Output: brave-falcon
}

func ExampleGenerator_SetStyle() {
	g, err := NewWithWords([]string{"brave"}, []string{"falcon"})
	if err != nil {
		panic(err)
	}

	g.SetStyle(Snake)
	fmt.Println(g.Generate())
	g.SetStyle(Space)
	fmt.Println(g.Generate())
	g.SetStyle(Camel)
	fmt.Println(g.Generate())
	// Output:
	// brave_falcon
	// brave falcon
	// braveFalcon
}

func ExampleGenerator_SetCustomSeparator() {
	g, err := NewWithWords([]string{"brave"}, []string{"falcon"})
	if err != nil {
		panic(err)
	}

	g.SetCustomSeparator(".")
	fmt.Println(g.Generate())
	// Output: brave.falcon
}

func ExampleGenerator_SetWordCount() {
	g, err := NewWithWords([]string{"brave"}, []string{"falcon"})
	if err != nil {
		panic(err)
	}

	if err := g.SetWordCount(3); err != nil {
		panic(err)
	}
	fmt.Println(g.Generate())
	// Output: brave-brave-falcon
}

func ExampleGenerator_ExcludeWords() {
	g, err := NewWithWords([]string{"brave", "wry"}, []string{"falcon"})
	if err != nil {
		panic(err)
	}

	if err := g.ExcludeWords("wry"); err != nil {
		panic(err)
	}
	fmt.Println(g.Generate())
	// Output: brave-falcon
}

func ExampleGenerator_Unique() {
	g, err := NewWithWords([]string{"brave", "wry"}, []string{"falcon"})
	if err != nil {
		panic(err)
	}

	names, err := g.Unique(2)
	if err != nil {
		panic(err)
	}
	sort.Strings(names)
	fmt.Println(names)
	// Output: [brave-falcon wry-falcon]
}
