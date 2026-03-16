package main

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/typisttech/wpsecadv/internal"
)

func newTestStore(t *testing.T) (string, string, *store) {
	t.Helper()

	tmpDir := t.TempDir()
	modDir := filepath.Join(tmpDir, "mod")
	assetsDir := "assets"

	if err := os.MkdirAll(modDir, 0o750); err != nil {
		t.Fatal(err)
	}

	logger := slog.New(slog.DiscardHandler)
	s, err := newStore(logger, 1, modDir, assetsDir)
	if err != nil {
		t.Fatalf("newStore() error = %v", err)
	}
	return modDir, assetsDir, s
}

func TestStore_Insert_Close_GenerateAssetJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kind internal.Kind
		slug string
	}{
		{
			name: "core",
			kind: internal.KindCore,
			slug: "wordpress",
		},
		{
			name: "plugin",
			kind: internal.KindPlugin,
			slug: "test-plugin",
		},
		{
			name: "theme",
			kind: internal.KindTheme,
			slug: "test-theme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			modDir, assetsDir, s := newTestStore(t)

			record := internal.Record{
				Kind: tt.kind,
				Slug: tt.slug,
				Advisory: internal.Advisory{
					ID:               "foobar",
					AffectedVersions: ">=1.0.0",
					Title:            "foobar",
					ReportedAt:       "2006-01-02 15:04:05",
					Sources:          []internal.Source{{Name: "foobar", ID: "123"}},
				},
			}

			if err := s.Insert(record); err != nil {
				t.Errorf("store.Insert() unexpected error: %v", err)
			}

			if err := s.Close(t.Context()); err != nil {
				t.Errorf("store.Close() unexpected error: %v", err)
			}

			path := filepath.Join(modDir, assetsDir, assetPath(tt.kind, tt.slug))
			if _, err := os.Stat(path); err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					t.Errorf("expected asset file %s does not exist", path)
				} else {
					t.Errorf("os.Stat(%q) unexpected error: %v", path, err)
				}
			}
		})
	}
}

func TestStore_Insert_Close_Skip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		record internal.Record
	}{
		{
			name: "unknown_kind",
			record: internal.Record{
				Kind: internal.KindUnknown,
				Slug: "foobar",
				Advisory: internal.Advisory{
					ID:               "foobar",
					AffectedVersions: ">=1.0.0",
					Title:            "foobar",
					ReportedAt:       "2006-01-02 15:04:05",
					Sources:          []internal.Source{{Name: "foobar", ID: "123"}},
				},
			},
		},
		{
			name: "wpmu",
			record: internal.Record{
				Kind: internal.KindWPMU,
				Slug: "foobar",
				Advisory: internal.Advisory{
					ID:               "foobar",
					AffectedVersions: ">=1.0.0",
					Title:            "foobar",
					ReportedAt:       "2006-01-02 15:04:05",
					Sources:          []internal.Source{{Name: "foobar", ID: "123"}},
				},
			},
		},
		{
			name: "incomplete_advisory",
			record: internal.Record{
				Kind: internal.KindPlugin,
				Slug: "foobar",
				Advisory: internal.Advisory{
					ID:               "foobar",
					AffectedVersions: ">=1.0.0",
					Title:            "foobar",
					ReportedAt:       "2006-01-02 15:04:05",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			modDir, assetsDir, s := newTestStore(t)

			if err := s.Insert(tt.record); err != nil {
				t.Fatalf("Insert() unexpected error: %v", err)
			}

			if err := s.Close(t.Context()); err != nil {
				t.Fatalf("Close() unexpected error: %v", err)
			}

			path := filepath.Join(modDir, assetsDir, assetPath(tt.record.Kind, tt.record.Slug))
			_, err := os.Stat(path)
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("os.Stat(%q) unexpected error: %v", path, err)
			}
			if err == nil {
				t.Errorf("asset JSON file %s created, want not exist", path)
			}
		})
	}
}

func TestStore_Close_GenerateGoFiles(t *testing.T) {
	t.Parallel()

	modDir, _, s := newTestStore(t)

	// No records inserted.

	if err := s.Close(t.Context()); err != nil {
		t.Errorf("store.Close() error = %v", err)
	}

	wants := []string{
		"assets_gen.go",
		"assets_gen_test.go",
	}

	for _, want := range wants {
		path := filepath.Join(modDir, want)
		if _, err := os.Stat(path); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				t.Errorf("expected %s does not exist at %s", want, path)
			} else {
				t.Errorf("os.Stat(%q) unexpected error: %v", path, err)
			}
		}
	}
}

func TestAssetPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kind internal.Kind
		slug string
		want string
	}{
		{"core", internal.KindCore, "wordpress", "core_gen.json"},
		{"core_empty_slug", internal.KindCore, "", "core_gen.json"},
		{"core_any_slug", internal.KindCore, "foobar", "core_gen.json"},
		{"core_any_slug_mixed_case", internal.KindCore, "foo_BAR", "core_gen.json"},
		{"plugin", internal.KindPlugin, "test-plugin", "plugin_746573742d706c7567696e811c9dc5_gen.json"},
		{"plugin_mixed_case", internal.KindPlugin, "Test-PLUGIN", "plugin_746573742d706c7567696e811c9dc5_gen.json"},
		{"theme", internal.KindTheme, "test-theme", "theme_746573742d7468656d65811c9dc5_gen.json"},
		{"theme_mixed_case", internal.KindTheme, "Test-THEME", "theme_746573742d7468656d65811c9dc5_gen.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assetPath(tt.kind, tt.slug)

			if got != tt.want {
				t.Errorf("assetPath(%v, %q) = %q, want %q", tt.kind, tt.slug, got, tt.want)
			}
		})
	}
}

func TestVarName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kind internal.Kind
		slug string
		want string
	}{
		{"core", internal.KindCore, "wordpress", "core"},
		{"core_empty_slug", internal.KindCore, "", "core"},
		{"core_any_slug", internal.KindCore, "foobar", "core"},
		{"core_any_slug_mixed_case", internal.KindCore, "foobar", "core"},
		{"plugin", internal.KindPlugin, "test-plugin", "plugin_746573742d706c7567696e811c9dc5"},
		{"plugin_mixed_case", internal.KindPlugin, "Test-PLUGIN", "plugin_746573742d706c7567696e811c9dc5"},
		{"theme", internal.KindTheme, "test-theme", "theme_746573742d7468656d65811c9dc5"},
		{"theme_mixed_case", internal.KindTheme, "Test-THEME", "theme_746573742d7468656d65811c9dc5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := varName(tt.kind, tt.slug)

			if got != tt.want {
				t.Errorf("varName(%v, %q) = %q, want %q", tt.kind, tt.slug, got, tt.want)
			}
		})
	}
}
