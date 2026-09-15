// Command namegen prints one or more generated names to stdout, for
// ad-hoc use from a shell (naming a scratch branch, a temp resource, a test
// fixture) without writing a throwaway Go program.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	namegen "github.com/abrown370/petname-generator"
)

var themesByName = map[string]namegen.Theme{
	"default": namegen.ThemeDefault,
	"space":   namegen.ThemeSpace,
	"tech":    namegen.ThemeTech,
}

var stylesByName = map[string]namegen.Style{
	"kebab": namegen.Kebab,
	"snake": namegen.Snake,
	"space": namegen.Space,
	"camel": namegen.Camel,
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "namegen:", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("namegen", flag.ContinueOnError)
	count := fs.Int("n", 1, "number of names to generate")
	unique := fs.Bool("unique", false, "generate names with no duplicates")
	theme := fs.String("theme", "default", "word list theme: default, space, tech")
	style := fs.String("style", "kebab", "output style: kebab, snake, space, camel")
	sep := fs.String("sep", "", "custom separator; overrides -style")
	words := fs.Int("words", 2, "words per name: 2 or 3")
	seed := fs.Int64("seed", 0, "seed for deterministic output (default: time-based)")
	exclude := fs.String("exclude", "", "comma-separated words to exclude, exact match")
	excludePattern := fs.String("exclude-pattern", "", "regular expression of words to exclude")
	plural := fs.Bool("plural", false, "pluralize the trailing noun, e.g. brave-falcons")
	if err := fs.Parse(args); err != nil {
		return err
	}

	t, ok := themesByName[*theme]
	if !ok {
		return fmt.Errorf("unknown theme %q (want default, space, or tech)", *theme)
	}
	s, ok := stylesByName[*style]
	if !ok {
		return fmt.Errorf("unknown style %q (want kebab, snake, space, or camel)", *style)
	}

	g, err := namegen.NewWithTheme(t)
	if err != nil {
		return err
	}

	fs.Visit(func(f *flag.Flag) {
		if f.Name == "seed" {
			g.Seed(*seed)
		}
	})

	if *sep != "" {
		g.SetCustomSeparator(*sep)
	} else {
		g.SetStyle(s)
	}

	if err := g.SetWordCount(*words); err != nil {
		return err
	}
	g.SetPluralNoun(*plural)

	if *exclude != "" {
		if err := g.ExcludeWords(strings.Split(*exclude, ",")...); err != nil {
			return err
		}
	}
	if *excludePattern != "" {
		if err := g.ExcludePattern(*excludePattern); err != nil {
			return err
		}
	}

	var names []string
	if *unique {
		names, err = g.Unique(*count)
	} else {
		names, err = g.GenerateN(*count)
	}
	if err != nil {
		return err
	}

	for _, n := range names {
		fmt.Fprintln(out, n)
	}
	return nil
}
