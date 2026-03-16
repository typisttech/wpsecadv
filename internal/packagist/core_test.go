package packagist

import (
	"regexp"
	"slices"
	"strings"
	"testing"
)

var pattern = regexp.MustCompile(`^[a-z0-9](([_.]|-{1,2})?[a-z0-9]+)*$`)

func TestIsCoreImplementation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		vendor  string
		project string
		want    bool
	}{
		{"positve_for_sanity", "roots", "wordpress-no-content", true},
		{"vendor_does_not_exist", "nonexistent", "wordpress-no-content", false},
		{"project_does_not_exist", "roots", "nonexistent", false},
		{"both_vendor_and_project_do_not_exist", "nonexistent", "nonexistent", false},
		{"empty_strings", "", "", false},
		{"vendor_is_empty_string", "", "wordpress-no-content", false},
		{"project_is_empty_string", "roots", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCoreImplementation(tt.vendor, tt.project)

			if got != tt.want {
				t.Errorf("IsCoreImplementation(%q, %q) = %t, want %t", tt.vendor, tt.project, got, tt.want)
			}
		})
	}
}

func TestCoreImplementation_AllAreCoreImplementation(t *testing.T) {
	t.Parallel()

	for _, pkg := range CoreImplementations() {
		t.Run(pkg, func(t *testing.T) {
			parts := strings.Split(pkg, "/")
			if len(parts) != 2 {
				t.Fatalf("invalid package format: %q", pkg)
			}
			vendor, project := parts[0], parts[1]

			if !IsCoreImplementation(vendor, project) {
				t.Errorf("IsCoreImplementation(%q, %q) = false, want true", vendor, project)
			}
		})
	}
}

func TestCoreImplementations(t *testing.T) {
	t.Parallel()

	got := CoreImplementations()

	if len(got) < 5 {
		t.Errorf("CoreImplementations() len = %d, want >= 5", len(got))
	}

	wants := []string{
		"johnpbloch/wordpress-core",
		"pantheon-systems/wordpress-composer", // Hardcoded.
		"roots/wordpress-no-content",
	}

	for _, want := range wants {
		if !slices.Contains(got, want) {
			t.Errorf("CoreImplementations() does not contain %q", want)
		}
	}
}

func TestCoreImplementations_Format(t *testing.T) {
	t.Parallel()

	for _, pkg := range CoreImplementations() {
		t.Run(pkg, func(t *testing.T) {
			if strings.ToLower(pkg) != pkg {
				t.Errorf("CoreImplementations() contains %q, want lowercased", pkg)
			}

			parts := strings.Split(pkg, "/")
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				t.Fatalf("CoreImplementations() contains malformed package name: %q, want exactly one slash", pkg)
			}

			if !pattern.MatchString(parts[0]) {
				t.Errorf("CoreImplementations() contains malformed vendor name %q, want %s", parts[0], pattern.String())
			}

			if !pattern.MatchString(parts[1]) {
				t.Errorf("CoreImplementations() contains malformed project name %q, want %s", parts[1], pattern.String())
			}
		})
	}
}
