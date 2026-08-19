# petname-generator

A Go library for generating readable random names, like `brave-falcon`
instead of `x7f2a9c1`. Useful anywhere you need an identifier a person can
say out loud, remember, and tell apart from the last one: temporary cloud
resources, seed data in a demo script, default names for a game character.

There's no CLI here, just a package to import.

## Install

```
go get github.com/abrown370/petname-generator
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/abrown370/petname-generator"
)

func main() {
	g := namegen.New()
	fmt.Println(g.Generate()) // e.g. "quiet-otter"

	names, err := g.GenerateN(5)
	if err != nil {
		panic(err)
	}
	fmt.Println(names)
}
```

### Deterministic output

For tests or reproducible demos, seed the generator explicitly:

```go
g := namegen.NewSeeded(42)
// g.Generate() will always return the same sequence for seed 42
```

### No duplicates

`Generate` and `GenerateN` can repeat names, same as rolling dice twice can
land on the same number. If you need a batch with no collisions, use
`Unique`:

```go
names, err := g.Unique(20)
if err != nil {
	// err is ErrNotEnoughCombinations if you asked for more unique names
	// than the word lists can produce
}
```

### Custom word lists

The default lists are plain adjective/noun pairs. Bring your own if you want
a different vocabulary or a themed set:

```go
g, err := namegen.NewWithWords(
	[]string{"solar", "lunar", "stellar"},
	[]string{"orbit", "comet", "nebula"},
)
```

### Formatting

Names default to kebab-case (`brave-falcon`). Other separators are
available via `SetStyle`:

```go
g.SetStyle(namegen.Snake) // brave_falcon
g.SetStyle(namegen.Space) // brave falcon
g.SetStyle(namegen.Camel) // braveFalcon
```

## Status

Early. The word lists are small and English-only for now, and there's no
support yet for three-word names or excluding specific words at generation
time.
