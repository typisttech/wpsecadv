package wordfence

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/typisttech/wpsecadv/internal"
	"github.com/typisttech/wpsecadv/internal/wordfence/vuln"
)

type spyLogger struct {
	onWarn func()
}

func (l *spyLogger) Info(string, ...any) {}
func (l *spyLogger) Warn(_ string, _ ...any) {
	if l.onWarn != nil {
		l.onWarn()
	}
}

func TestMakeRecords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		vuln vuln.Vulnerability
		want []internal.Record
	}{
		{
			name: "plugin_vulnerability_produces_plugin_record",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my-plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/abc-123/my-plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "abc-123"}},
						AffectedVersions: "<=1.0",
					},
				},
			},
		},
		{
			name: "theme_vulnerability_produces_theme_record",
			vuln: vuln.Vulnerability{
				ID: "def-456",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypeTheme,
						Slug: "my-theme",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-2.0": {FromVersion: "*", FromInclusive: true, ToVersion: "2.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindTheme,
					Slug: "my-theme",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/def-456/my-theme",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "def-456"}},
						AffectedVersions: "<=2.0",
					},
				},
			},
		},
		{
			name: "core_vulnerability_produces_core_record",
			vuln: vuln.Vulnerability{
				ID: "ghi-789",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypeCore,
						Slug: "wordpress",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-6.0": {FromVersion: "*", FromInclusive: true, ToVersion: "6.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindCore,
					Slug: "wordpress",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/ghi-789/wordpress",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "ghi-789"}},
						AffectedVersions: "<=6.0",
					},
				},
			},
		},
		{
			name: "wpmu_slug_produces_wpmu_record",
			vuln: vuln.Vulnerability{
				ID: "jkl-000",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypeCore,
						Slug: "wpmu",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-3.0": {FromVersion: "*", FromInclusive: true, ToVersion: "3.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindWPMU,
					Slug: "wpmu",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/jkl-000/wpmu",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "jkl-000"}},
						AffectedVersions: "<=3.0",
					},
				},
			},
		},
		{
			name: "multiple_software_entries_produce_multiple_records",
			vuln: vuln.Vulnerability{
				ID: "mno-111",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "plugin-a",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "plugin-b",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-2.0": {FromVersion: "*", FromInclusive: true, ToVersion: "2.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypeTheme,
						Slug: "theme-a",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-3.0": {FromVersion: "*", FromInclusive: true, ToVersion: "3.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypeTheme,
						Slug: "theme-b",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-4.0": {FromVersion: "*", FromInclusive: true, ToVersion: "4.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypeCore,
						Slug: "wordpress",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-5.0": {FromVersion: "*", FromInclusive: true, ToVersion: "5.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypeCore,
						Slug: "wpmu",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-6.0": {FromVersion: "*", FromInclusive: true, ToVersion: "6.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "plugin-a",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/mno-111/plugin-a",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "mno-111"}},
						AffectedVersions: "<=1.0",
					},
				},
				{
					Kind: internal.KindPlugin,
					Slug: "plugin-b",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/mno-111/plugin-b",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "mno-111"}},
						AffectedVersions: "<=2.0",
					},
				},
				{
					Kind: internal.KindTheme,
					Slug: "theme-a",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/mno-111/theme-a",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "mno-111"}},
						AffectedVersions: "<=3.0",
					},
				},
				{
					Kind: internal.KindTheme,
					Slug: "theme-b",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/mno-111/theme-b",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "mno-111"}},
						AffectedVersions: "<=4.0",
					},
				},
				{
					Kind: internal.KindCore,
					Slug: "wordpress",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/mno-111/wordpress",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "mno-111"}},
						AffectedVersions: "<=5.0",
					},
				},
				{
					Kind: internal.KindWPMU,
					Slug: "wpmu",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/mno-111/wpmu",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "mno-111"}},
						AffectedVersions: "<=6.0",
					},
				},
			},
		},
		{
			name: "duplicate_slugs_are_combined_into_one_record",
			vuln: vuln.Vulnerability{
				ID: "dup-111",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"2.0-3.0": {FromVersion: "2.0", FromInclusive: true, ToVersion: "3.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my-plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/dup-111/my-plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "dup-111"}},
						AffectedVersions: "<=1.0|>=2.0,<=3.0",
					},
				},
			},
		},
		{
			name: "duplicate_slugs_are_combined_into_one_record_match_all_wildcard",
			vuln: vuln.Vulnerability{
				ID: "dup-match-all",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
							"bad":   {FromVersion: "not-valid", FromInclusive: true, ToVersion: "also-bad", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"2.0-3.0": {FromVersion: "2.0", FromInclusive: true, ToVersion: "3.0", ToInclusive: true},
							"*":       {FromVersion: "*", FromInclusive: true, ToVersion: "*", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my-plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/dup-match-all/my-plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "dup-match-all"}},
						AffectedVersions: "*",
					},
				},
			},
		},
		{
			name: "duplicate_slugs_with_mixed_case_are_combined_into_one_record",
			vuln: vuln.Vulnerability{
				ID: "dup-222",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "My-Plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"2.0-3.0": {FromVersion: "2.0", FromInclusive: true, ToVersion: "3.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my-plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/dup-222/my-plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "dup-222"}},
						AffectedVersions: "<=1.0|>=2.0,<=3.0",
					},
				},
			},
		},
		{
			name: "duplicate_slugs_with_invalid_entry_produces_partial_combined_record",
			vuln: vuln.Vulnerability{
				ID: "dup-333",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"bad": {FromVersion: "not-valid", FromInclusive: true, ToVersion: "also-bad", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my-plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/dup-333/my-plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "dup-333"}},
						AffectedVersions: "<=1.0",
					},
				},
			},
		},
		{
			name: "duplicate_slugs_where_all_entries_invalid_produces_no_record",
			vuln: vuln.Vulnerability{
				ID: "dup-444",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"bad": {FromVersion: "not-valid", FromInclusive: true, ToVersion: "also-bad", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"also-bad": {FromVersion: "not-valid", FromInclusive: true, ToVersion: "also-bad", ToInclusive: true},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "duplicate_slugs_with_conflicting_kind_produces_two_records",
			vuln: vuln.Vulnerability{
				ID: "dup-666",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-slug",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypeTheme,
						Slug: "my-slug",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*": {FromVersion: "*", FromInclusive: true, ToVersion: "*", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my-slug",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/dup-666/my-slug",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "dup-666"}},
						AffectedVersions: "<=1.0",
					},
				},
				{
					Kind: internal.KindTheme,
					Slug: "my-slug",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/dup-666/my-slug",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "dup-666"}},
						AffectedVersions: "*",
					},
				},
			},
		},
		{
			name: "duplicate_slugs_combined_constraints_are_sorted",
			vuln: vuln.Vulnerability{
				ID: "dup-555",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"2.0-3.0": {FromVersion: "2.0", FromInclusive: true, ToVersion: "3.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my-plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/dup-555/my-plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "dup-555"}},
						AffectedVersions: "<=1.0|>=2.0,<=3.0",
					},
				},
			},
		},
		{
			name: "software_with_only_invalid_versions_is_skipped",
			vuln: vuln.Vulnerability{
				ID: "stu-333",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "bad-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"bad": {FromVersion: "not-valid", FromInclusive: true, ToVersion: "also-bad", ToInclusive: true},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "no_slug_produces_no_records",
			vuln: vuln.Vulnerability{
				ID: "missing-slug",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "to_lowercased_slug",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "My-Uppercased-PLUGIN",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my-uppercased-plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/abc-123/my-uppercased-plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "abc-123"}},
						AffectedVersions: "<=1.0",
					},
				},
			},
		},
		{
			name: "slug_with_invalid_pattern_produces_no_records",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "invalid slug!",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "slug_starting_with_hyphen_produces_no_records",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "-invalid",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "slug_with_triple_hyphen_produces_no_records",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "invalid---slug",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "slug_with_spaces_produces_no_records",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "invalid slug",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "slug_with_double_hyphen_is_valid",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "block--editor",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "block--editor",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/abc-123/block--editor",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "abc-123"}},
						AffectedVersions: "<=1.0",
					},
				},
			},
		},
		{
			name: "slug_with_underscore_is_valid",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my_plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my_plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/abc-123/my_plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "abc-123"}},
						AffectedVersions: "<=1.0",
					},
				},
			},
		},
		{
			name: "slug_with_dot_is_valid",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my.plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			want: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my.plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/abc-123/my.plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "abc-123"}},
						AffectedVersions: "<=1.0",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeRecords(noopLogger{}, tt.vuln)

			if diff := cmp.Diff(tt.want, got, sortRecords(), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("makeRecords() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMakeRecords_Warn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		vuln      vuln.Vulnerability
		wantWarns int
	}{
		{
			name: "empty_slug_triggers_warn",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			wantWarns: 1,
		},
		{
			name: "non_empty_slug_does_not_trigger_warn",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			wantWarns: 0,
		},
		{
			name: "multiple_empty_slugs_trigger_warn_per_software",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypeTheme,
						Slug: "",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-2.0": {FromVersion: "*", FromInclusive: true, ToVersion: "2.0", ToInclusive: true},
						},
					},
				},
			},
			wantWarns: 2,
		},
		{
			name: "mixed_empty_and_non_empty_slugs_warn_only_for_empty",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "my-plugin",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-2.0": {FromVersion: "*", FromInclusive: true, ToVersion: "2.0", ToInclusive: true},
						},
					},
				},
			},
			wantWarns: 1,
		},
		{
			name: "invalid_slug_pattern_triggers_warn",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "invalid slug!",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			wantWarns: 1,
		},
		{
			name: "slug_starting_with_hyphen_triggers_warn",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "-invalid",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			wantWarns: 1,
		},
		{
			name: "slug_with_triple_hyphen_triggers_warn",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "invalid---slug",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			wantWarns: 1,
		},
		{
			name: "valid_slug_pattern_does_not_trigger_warn",
			vuln: vuln.Vulnerability{
				ID: "abc-123",
				Software: []vuln.Software{
					{
						Type: vuln.SoftwareTypePlugin,
						Slug: "block--editor",
						AffectedVersions: map[string]vuln.AffectedVersion{
							"*-1.0": {FromVersion: "*", FromInclusive: true, ToVersion: "1.0", ToInclusive: true},
						},
					},
				},
			},
			wantWarns: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var warnCount int
			l := &spyLogger{onWarn: func() { warnCount++ }}

			makeRecords(l, tt.vuln)

			if warnCount != tt.wantWarns {
				t.Errorf("makeRecords() triggered %d warnings, want %d", warnCount, tt.wantWarns)
			}
		})
	}
}

func TestMakeTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		title      string
		copyrights vuln.Copyrights
		want       string
	}{
		{
			name:       "title_without_copyrights",
			title:      "Plugin <= 1.0 - XSS",
			copyrights: nil,
			want:       "Plugin <= 1.0 - XSS",
		},
		{
			name:  "title_with_single_copyright",
			title: "Plugin <= 1.0 - XSS",
			copyrights: vuln.Copyrights{
				"ACME": {
					Notice:     "Copyright 2012 ACME Inc.",
					License:    "License text.",
					LicenseURL: "https://example.com/license",
				},
			},
			want: "Plugin <= 1.0 - XSS\n### Copyright 2012 ACME Inc.\nLicense text.\nhttps://example.com/license",
		},
		{
			name:  "title_with_multiple_copyrights_sorted_by_notice",
			title: "Plugin <= 1.0 - XSS",
			copyrights: vuln.Copyrights{
				"umbrella": {
					Notice:     "Copyright 1999 Umbrella Corporation",
					License:    "Umbrella license.",
					LicenseURL: "https://umbrella.example.com/",
				},
				"ACME": {
					Notice:     "Copyright 2012 ACME Inc.",
					License:    "ACME license.",
					LicenseURL: "https://acme.example.com/",
				},
			},
			want: "Plugin <= 1.0 - XSS\n### Copyright 1999 Umbrella Corporation\nUmbrella license.\nhttps://umbrella.example.com/\n### Copyright 2012 ACME Inc.\nACME license.\nhttps://acme.example.com/",
		},
		{
			name:       "empty_title_without_copyrights",
			title:      "",
			copyrights: nil,
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeTitle(tt.title, tt.copyrights)

			if got != tt.want {
				t.Errorf("makeTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMakeLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		refs    []string
		cveLink string
		want    string
	}{
		{
			name:    "first_non_empty_reference_is_used",
			refs:    []string{"https://example.com/ref1", "https://example.com/ref2"},
			cveLink: "https://cve.example.com/",
			want:    "https://example.com/ref1",
		},
		{
			name:    "skips_empty_refs_and_uses_first_non_empty",
			refs:    []string{"", "https://example.com/ref2"},
			cveLink: "https://cve.example.com/",
			want:    "https://example.com/ref2",
		},
		{
			name:    "falls_back_to_cve_link_when_all_refs_empty",
			refs:    []string{"", ""},
			cveLink: "https://cve.example.com/",
			want:    "https://cve.example.com/",
		},
		{
			name:    "falls_back_to_cve_link_when_refs_nil",
			refs:    nil,
			cveLink: "https://cve.example.com/",
			want:    "https://cve.example.com/",
		},
		{
			name:    "returns_empty_string_when_no_refs_and_no_cve_link",
			refs:    nil,
			cveLink: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeLink(tt.refs, tt.cveLink)

			if got != tt.want {
				t.Errorf("makeLink() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMakeConstraints(t *testing.T) {
	t.Parallel()

	noWarn := func(string, ...any) {}

	tests := []struct {
		name string
		avs  []vuln.AffectedVersion
		want string
	}{
		{
			name: "match_all_wildcard_returns_asterisk",
			avs:  []vuln.AffectedVersion{{FromVersion: "*", FromInclusive: true, ToVersion: "*", ToInclusive: true}},
			want: "*",
		},
		{
			name: "exact_version_uses_equals_prefix",
			avs:  []vuln.AffectedVersion{{FromVersion: "1.0.0", FromInclusive: true, ToVersion: "1.0.0", ToInclusive: true}},
			want: "=1.0.0",
		},
		{
			name: "upper_bounded_inclusive",
			avs:  []vuln.AffectedVersion{{FromVersion: "*", FromInclusive: true, ToVersion: "3.5", ToInclusive: true}},
			want: "<=3.5",
		},
		{
			name: "upper_bounded_exclusive",
			avs:  []vuln.AffectedVersion{{FromVersion: "*", FromInclusive: true, ToVersion: "3.5", ToInclusive: false}},
			want: "<3.5",
		},
		{
			name: "lower_bounded_inclusive",
			avs:  []vuln.AffectedVersion{{FromVersion: "1.0", FromInclusive: true, ToVersion: "*", ToInclusive: true}},
			want: ">=1.0",
		},
		{
			name: "lower_bounded_exclusive",
			avs:  []vuln.AffectedVersion{{FromVersion: "1.0", FromInclusive: false, ToVersion: "*", ToInclusive: true}},
			want: ">1.0",
		},
		{
			name: "range_inclusive_both_ends",
			avs:  []vuln.AffectedVersion{{FromVersion: "1.0.3", FromInclusive: true, ToVersion: "2.0", ToInclusive: true}},
			want: ">=1.0.3,<=2.0",
		},
		{
			name: "range_exclusive_both_ends",
			avs:  []vuln.AffectedVersion{{FromVersion: "1.0", FromInclusive: false, ToVersion: "2.0.3", ToInclusive: false}},
			want: ">1.0,<2.0.3",
		},
		{
			name: "multiple_ranges_joined_with_pipe_sorted",
			avs: []vuln.AffectedVersion{
				{FromVersion: "*", FromInclusive: true, ToVersion: "3.5", ToInclusive: true},
				{FromVersion: "2.0", FromInclusive: true, ToVersion: "4.0", ToInclusive: false},
				{FromVersion: "4.9.2", FromInclusive: true, ToVersion: "4.9.2", ToInclusive: true},
				{FromVersion: "1.0.0", FromInclusive: false, ToVersion: "3.0.8", ToInclusive: true},
			},
			want: "<=3.5|=4.9.2|>1.0.0,<=3.0.8|>=2.0,<4.0",
		},
		{
			name: "empty_slice_returns_empty_string",
			avs:  []vuln.AffectedVersion{},
			want: "",
		},
		{
			name: "invalid_from_version_skipped",
			avs:  []vuln.AffectedVersion{{FromVersion: "not-a-version", FromInclusive: true, ToVersion: "2.0", ToInclusive: true}},
			want: "",
		},
		{
			name: "invalid_to_version_skipped",
			avs:  []vuln.AffectedVersion{{FromVersion: "1.0", FromInclusive: true, ToVersion: "not-a-version", ToInclusive: true}},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeConstraints(noWarn, tt.avs)

			if got != tt.want {
				t.Errorf("makeConstraints() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMakeConstraints_Warn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		avs       []vuln.AffectedVersion
		wantWarns int
	}{
		{
			name:      "invalid_from_version_triggers_warn",
			avs:       []vuln.AffectedVersion{{FromVersion: "not-a-version", FromInclusive: true, ToVersion: "2.0.0", ToInclusive: true}},
			wantWarns: 1,
		},
		{
			name:      "invalid_to_version_triggers_warn",
			avs:       []vuln.AffectedVersion{{FromVersion: "1.0.0", FromInclusive: true, ToVersion: "not-a-version", ToInclusive: true}},
			wantWarns: 1,
		},
		{
			name:      "valid_versions_do_not_trigger_warn",
			avs:       []vuln.AffectedVersion{{FromVersion: "1.0.0", FromInclusive: true, ToVersion: "2.0.0", ToInclusive: true}},
			wantWarns: 0,
		},
		{
			name:      "empty_slice_do_not_trigger_warn",
			avs:       []vuln.AffectedVersion{},
			wantWarns: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var warnCount int
			warn := func(string, ...any) { warnCount++ }

			makeConstraints(warn, tt.avs)

			if warnCount != tt.wantWarns {
				t.Errorf("makeConstraints() triggered %d warnings, want %d", warnCount, tt.wantWarns)
			}
		})
	}
}

func TestMakeSeverity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		rating vuln.CVSSRating
		want   string
	}{
		{"none_lowercased", vuln.CVSSRatingNone, "none"},
		{"low_lowercased", vuln.CVSSRatingLow, "low"},
		{"medium_lowercased", vuln.CVSSRatingMedium, "medium"},
		{"high_lowercased", vuln.CVSSRatingHigh, "high"},
		{"critical_lowercased", vuln.CVSSRatingCritical, "critical"},
		{"empty_stays_empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeSeverity(tt.rating)

			if got != tt.want {
				t.Errorf("makeSeverity(%q) = %q, want %q", tt.rating, got, tt.want)
			}
		})
	}
}

func TestMakeKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		softType vuln.SoftwareType
		slug     string
		want     internal.Kind
	}{
		{"core", vuln.SoftwareTypeCore, "wordpress", internal.KindCore},
		{"core_any_slug", vuln.SoftwareTypeCore, "foobar", internal.KindCore},
		{"wpmu", vuln.SoftwareTypeCore, "wpmu", internal.KindWPMU},
		{"plugin", vuln.SoftwareTypePlugin, "foobar", internal.KindPlugin},
		{"theme", vuln.SoftwareTypeTheme, "foobar", internal.KindTheme},
		{"unknown", "unknown", "foobar", internal.KindUnknown},
		{"empty", "", "foobar", internal.KindUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeKind(tt.softType, tt.slug)

			if got != tt.want {
				t.Errorf("makeKind(%q, %q) = %q, want %q", tt.softType, tt.slug, got, tt.want)
			}
		})
	}
}
