package data

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9](([_.]|-{1,2})?[a-z0-9]+)*$`)

func TestPluginsAndThemesDisjoint(t *testing.T) {
	// TODO: Re-enable this test.
	t.Skip("Something wrong on Wordfence side. Awaiting replies.")

	t.Parallel()

	tms := themes()
	ps := plugins()

	set := make(map[string]struct{}, len(tms))
	for _, tm := range tms {
		set[tm] = struct{}{}
	}

	var conflicts []string
	for _, p := range ps {
		if _, ok := set[p]; ok {
			conflicts = append(conflicts, p)
		}
	}

	if len(conflicts) > 0 {
		slices.Sort(conflicts)

		artPath := filepath.Join(t.ArtifactDir(), "conflicts.txt")
		if err := os.WriteFile(artPath, []byte(strings.Join(conflicts, "\n")), 0o600); err != nil {
			t.Logf("failed to write conflicts artifact: %v", err)
		}

		t.Errorf("plugins() and themes() share %d slugs, want 0", len(conflicts))
	}
}

func TestPlugins(t *testing.T) {
	t.Parallel()

	got := plugins()

	if len(got) < 15000 {
		t.Errorf("plugins() len = %d, want >= 15000", len(got))
	}

	wants := []string{
		"akismet",
		"gravityforms",
		"woocommerce",
	}

	for _, want := range wants {
		if !slices.Contains(got, want) {
			t.Errorf("plugins() does not contain %q", want)
		}
	}

	for _, slug := range got {
		if strings.ToLower(slug) != slug {
			t.Errorf("plugins() contains %q, want lowercased", slug)
		}

		if !slugPattern.MatchString(slug) {
			t.Errorf("plugins() contains malformed slug %q, want %s", slug, slugPattern.String())
		}
	}
}

func TestThemes(t *testing.T) {
	t.Parallel()

	got := themes()

	if len(got) < 1700 {
		t.Errorf("themes() len = %d, want >= 1700", len(got))
	}

	wants := []string{
		"divi",
		"ithemes2",
		"twentyfifteen",
	}

	for _, want := range wants {
		if !slices.Contains(got, want) {
			t.Errorf("themes() does not contain %q", want)
		}
	}

	for _, slug := range got {
		if strings.ToLower(slug) != slug {
			t.Errorf("themes() contains %q, want lowercased", slug)
		}

		if !slugPattern.MatchString(slug) {
			t.Errorf("themes() contains malformed slug %q, want %s", slug, slugPattern.String())
		}
	}
}
