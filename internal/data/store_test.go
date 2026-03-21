package data

import (
	"encoding/json/v2"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/typisttech/wpsecadv/internal"
	"github.com/typisttech/wpsecadv/internal/packagist"
)

func TestStore_MarshalAdvisoriesFor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		vendor   string
		slug     string
		wantCVEs []string
	}{
		{
			name:   "core_exists_roots_wordpress_full",
			vendor: "roots",
			slug:   "wordpress-full",
			wantCVEs: []string{
				"CVE-2012-0287",
				"CVE-2018-5776",
				"CVE-2022-21664",
			},
		},
		{
			name:   "core_exists_roots_wordpress_no_content",
			vendor: "roots",
			slug:   "wordpress-no-content",
			wantCVEs: []string{
				"CVE-2012-0287",
				"CVE-2018-5776",
				"CVE-2022-21664",
			},
		},
		{
			name:   "core_exists_johnpbloch_wordpress_core",
			vendor: "johnpbloch",
			slug:   "wordpress-core",
			wantCVEs: []string{
				"CVE-2012-0287",
				"CVE-2018-5776",
				"CVE-2022-21664",
			},
		},
		{
			name:   "core_exists_wp_core_wordpress",
			vendor: "wp-core",
			slug:   "wordpress",
			wantCVEs: []string{
				"CVE-2012-0287",
				"CVE-2018-5776",
				"CVE-2022-21664",
			},
		},
		{
			name:   "core_exists_wp_core_wordpress_no_content",
			vendor: "wp-core",
			slug:   "wordpress-no-content",
			wantCVEs: []string{
				"CVE-2012-0287",
				"CVE-2018-5776",
				"CVE-2022-21664",
			},
		},
		{
			name:   "any_vendor_plugin_exists",
			vendor: "any-vendor",
			slug:   "woocommerce",
			wantCVEs: []string{
				"CVE-2025-15033",
				"CVE-2025-26762",
				"CVE-2025-49042",
			},
		},
		{
			name:   "wpackagist_plugin_exists",
			vendor: "wpackagist-plugin",
			slug:   "woocommerce",
			wantCVEs: []string{
				"CVE-2025-15033",
				"CVE-2025-26762",
				"CVE-2025-49042",
			},
		},
		{
			name:   "any_vendor_theme_exists",
			vendor: "any-vendor",
			slug:   "twentyfifteen",
			wantCVEs: []string{
				"CVE-2015-3429",
			},
		},
		{
			name:   "wpackagist_theme_exists",
			vendor: "wpackagist-theme",
			slug:   "twentyfifteen",
			wantCVEs: []string{
				"CVE-2015-3429",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &Store{}

			got, err := s.MarshalAdvisoriesFor(tt.vendor, tt.slug)
			if err != nil {
				t.Fatalf("MarshalAdvisoriesFor() unexpected error: %v", err)
			}
			if len(got) == 0 {
				t.Fatalf("MarshalAdvisoriesFor() len = 0, want > 0")
			}

			var advs []internal.Advisory
			if err := json.Unmarshal(got, &advs); err != nil {
				t.Fatalf("json.Unmarshal() unexpected error: %v", err)
			}

			for _, wantCVE := range tt.wantCVEs {
				if !slices.ContainsFunc(advs, func(a internal.Advisory) bool { return a.CVE == wantCVE }) {
					t.Errorf("MarshalAdvisoriesFor(%q, %q) does not contain advisory %s", tt.vendor, tt.slug, wantCVE)
				}
			}
		})
	}
}

func TestStore_MarshalAdvisoriesFor_NotExist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		vendor string
		slug   string
	}{
		{name: "wpackagist_plugin_not_exist", vendor: "wpackagist-plugin", slug: "non-existent-plugin"},
		{name: "wpackagist_plugin_exists_bad_vendor", vendor: "wpackagist-theme", slug: "woocommerce"},
		{name: "wpackagist_theme_not_exist", vendor: "wpackagist-theme", slug: "non-existent-theme"},
		{name: "wpackagist_theme_exists_bad_vendor", vendor: "wpackagist-plugin", slug: "twentyfifteen"},
		{name: "any_vendor_not_exist", vendor: "any-vendor", slug: "non-existent"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &Store{}

			_, err := s.MarshalAdvisoriesFor(tt.vendor, tt.slug)
			if err == nil {
				t.Fatalf("MarshalAdvisoriesFor(%q, %q) unexpected success", tt.vendor, tt.slug)
			}
			if !errors.Is(err, errPackageNotFound) {
				t.Fatalf("MarshalAdvisoriesFor(%q, %q) error = %[3]T %[3]q, want errPackageNotFound", tt.vendor, tt.slug, err)
			}
		})
	}
}

func TestStore_MarshalAdvisoriesFor_Core(t *testing.T) {
	t.Parallel()

	for _, pkg := range packagist.CoreImplementations() {
		t.Run(pkg, func(t *testing.T) {
			parts := strings.Split(pkg, "/")
			if len(parts) != 2 {
				t.Fatalf("malformed CoreImplementations entry: %q", pkg)
			}
			vendor, project := parts[0], parts[1]

			s := &Store{}
			got, err := s.MarshalAdvisoriesFor(vendor, project)
			if err != nil {
				t.Fatalf("MarshalAdvisoriesFor(%q, %q) unexpected error: %v", vendor, project, err)
			}
			if len(got) == 0 {
				t.Fatalf("MarshalAdvisoriesFor(%q, %q) len = 0, want > 0", vendor, project)
			}

			var advs []internal.Advisory
			if err := json.Unmarshal(got, &advs); err != nil {
				t.Fatalf("json.Unmarshal() unexpected error: %v", err)
			}
			if len(advs) < 300 {
				t.Errorf("MarshalAdvisoriesFor(%q, %q) len = %d advisories, want >= 300", vendor, project, len(advs))
			}
		})
	}
}

func TestStore_MarshalAdvisoriesFor_NonCore(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		slugs func() []string
	}{
		{"plugins", plugins},
		{"themes", themes},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, slug := range tt.slugs() {
				t.Run(slug, func(t *testing.T) {
					s := &Store{}

					got, err := s.MarshalAdvisoriesFor("any-vendor", slug)
					if err != nil {
						t.Fatalf("MarshalAdvisoriesFor(%q, %q)unexpected error: %v", "any-vendor", slug, err)
					}
					if len(got) == 0 {
						t.Fatalf("MarshalAdvisoriesFor(%q, %q)len = 0, want > 0", "any-vendor", slug)
					}

					var advs []internal.Advisory
					if err := json.Unmarshal(got, &advs); err != nil {
						t.Fatalf("json.Unmarshal() unexpected error: %v", err)
					}
					if len(advs) == 0 {
						t.Errorf("MarshalAdvisoriesFor(%q, %q)len = 0, want > 0", "any-vendor", slug)
					}
				})
			}
		})
	}
}
