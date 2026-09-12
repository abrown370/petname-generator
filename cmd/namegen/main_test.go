package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunDefaults(t *testing.T) {
	var buf bytes.Buffer
	if err := run(nil, &buf); err != nil {
		t.Fatalf("run: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}
	if !strings.Contains(lines[0], "-") {
		t.Errorf("name %q does not look kebab-cased", lines[0])
	}
}

func TestRunCount(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"-n", "5"}, &buf); err != nil {
		t.Fatalf("run: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 5 {
		t.Fatalf("got %d lines, want 5", len(lines))
	}
}

func TestRunDeterministicWithSeed(t *testing.T) {
	var a, b bytes.Buffer
	if err := run([]string{"-seed", "42", "-n", "3"}, &a); err != nil {
		t.Fatalf("run: %v", err)
	}
	if err := run([]string{"-seed", "42", "-n", "3"}, &b); err != nil {
		t.Fatalf("run: %v", err)
	}
	if a.String() != b.String() {
		t.Errorf("same seed produced different output:\n%s\nvs\n%s", a.String(), b.String())
	}
}

func TestRunStyles(t *testing.T) {
	cases := []struct {
		style string
		want  string
	}{
		{"kebab", "-"},
		{"snake", "_"},
		{"space", " "},
	}
	for _, c := range cases {
		var buf bytes.Buffer
		if err := run([]string{"-style", c.style, "-seed", "1"}, &buf); err != nil {
			t.Fatalf("run(%s): %v", c.style, err)
		}
		if !strings.Contains(buf.String(), c.want) {
			t.Errorf("style %s: output %q does not contain %q", c.style, buf.String(), c.want)
		}
	}
}

func TestRunCustomSeparatorOverridesStyle(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"-style", "snake", "-sep", "."}, &buf); err != nil {
		t.Fatalf("run: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if !strings.Contains(out, ".") || strings.Contains(out, "_") {
		t.Errorf("got %q, want dot-separated name with no underscore", out)
	}
}

func TestRunThreeWords(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"-words", "3", "-seed", "7"}, &buf); err != nil {
		t.Fatalf("run: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if len(strings.Split(out, "-")) != 3 {
		t.Errorf("got %q, want three hyphen-separated words", out)
	}
}

func TestRunUnknownTheme(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"-theme", "nope"}, &buf); err == nil {
		t.Fatal("expected error for unknown theme, got nil")
	}
}

func TestRunUnknownStyle(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"-style", "nope"}, &buf); err == nil {
		t.Fatal("expected error for unknown style, got nil")
	}
}

func TestRunExcludeWords(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"-exclude", "not-a-real-word,also-not-real"}, &buf); err != nil {
		t.Fatalf("run: %v", err)
	}
	if strings.TrimSpace(buf.String()) == "" {
		t.Fatal("expected a name, got empty output")
	}
}

func TestRunUniqueNoDuplicates(t *testing.T) {
	var buf bytes.Buffer
	if err := run([]string{"-unique", "-n", "10", "-seed", "9"}, &buf); err != nil {
		t.Fatalf("run: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	seen := make(map[string]bool)
	for _, l := range lines {
		if seen[l] {
			t.Errorf("duplicate name %q in unique output", l)
		}
		seen[l] = true
	}
}
