package vuln

import (
	"encoding/json/v2"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVulnerabilities(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		json string
		want []Vulnerability
	}{
		{
			name: "two_items",
			json: `{
				"0001": {"id":"0001","title":"First"},
				"0002": {"id":"0002","title":"Second"}
			}`,
			want: []Vulnerability{
				{ID: "0001", Title: "First"},
				{ID: "0002", Title: "Second"},
			},
		},
		{
			name: "single_item",
			json: `{
				"0001": {"id":"0001","title":"First"}
			}`,
			want: []Vulnerability{{ID: "0001", Title: "First"}},
		},
		{
			name: "spaces_before_single_item",
			json: `             {
				"0001": {"id":"0001","title":"First"}
			}`,
			want: []Vulnerability{{ID: "0001", Title: "First"}},
		},
		{
			name: "spaces_after_single_item",
			json: `{
				"0001": {"id":"0001","title":"First"}
			}          `,
			want: []Vulnerability{{ID: "0001", Title: "First"}},
		},
		{
			name: "spaces_around_single_item",
			json: `             {
				"0001": {"id":"0001","title":"First"}
			}          `,
			want: []Vulnerability{{ID: "0001", Title: "First"}},
		},
		{
			name: "full_object",
			json: `{
				"59ea7390-3587-4fce-b267-bf525dbe3e27": {
					"id": "59ea7390-3587-4fce-b267-bf525dbe3e27",
					"title": "Multiple Plugins <= (Various Versions) - Authenticated (Contributor+) Stored DOM-Based Cross-Site Scripting via ThickBox JavaScript Library",
					"software": [
						{
							"type": "plugin",
							"name": "Foo Bar",
							"slug": "foo-bar",
							"affected_versions": {
								"* - 3.5": {
									"from_version": "*",
									"from_inclusive": true,
									"to_version": "3.5",
									"to_inclusive": true
								}
							},
							"patched": false,
							"patched_versions": [],
							"remediation": "No known patch available. Please review the vulnerability's details in depth and employ mitigations based on your organization's risk tolerance. It may be best to uninstall the affected software and find a replacement."
						},
						{
							"type": "theme",
							"name": "baz quux \u2013 WPBaz",
							"slug": "wpbaz",
							"affected_versions": {
								"* - 4.8.9": {
									"from_version": "*",
									"from_inclusive": true,
									"to_version": "4.8.9",
									"to_inclusive": true
								},
								"4.9.2": {
									"from_version": "4.9.2",
									"from_inclusive": true,
									"to_version": "4.9.2",
									"to_inclusive": true
								}
							},
							"patched": true,
							"patched_versions": [
								"4.9.1",
								"4.9.3"
							],
							"remediation": "Update to one of the following versions, or a newer patched version: 4.9.1, 4.9.3"
						}
					],
					"informational": false,
					"description": "Multiple plugins for WordPress are vulnerable to Stored Cross-Site Scripting via the plugin's bundled ThickBox JavaScript library (version 3.1) in various versions due to insufficient input sanitization and output escaping on user supplied attributes. This makes it possible for authenticated attackers, with contributor-level access and above, to inject arbitrary web scripts in pages that will execute whenever a user accesses an injected page.",
					"references": [
						"https:\/\/www.example.com\/threat-intel\/vulnerabilities\/id\/59ea7390-3587-4fce-b267-bf525dbe3e27?source=api-prod"
					],
					"cwe": {
						"id": 79,
						"name": "Improper Neutralization of Input During Web Page Generation ('Cross-site Scripting')",
						"description": "The product does not neutralize or incorrectly neutralizes user-controllable input before it is placed in output that is used as a web page that is served to other users."
					},
					"cvss": {
						"vector": "CVSS:3.1\/AV:N\/AC:L\/PR:L\/UI:N\/S:C\/C:L\/I:L\/A:N",
						"score": 6.4,
						"rating": "Medium"
					},
					"cve": "CVE-2025-2537",
					"cve_link": "https:\/\/www.example.org\/CVERecord?id=CVE-2025-2537",
					"researchers": [
						"The hero"
					],
					"published": "2025-07-03 00:20:54",
					"updated": "2025-08-19 05:27:14",
					"copyrights": {
						"message": "This record contains material that is subject to copyright",
						"ACME": {
							"notice": "Copyright 2012-2025 ACME Inc.",
							"license": "ACME hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute this software vulnerability information. Any copy of the software vulnerability information you make for such purposes is authorized provided that you include a hyperlink to this vulnerability record and reproduce ACME's copyright designation and this license in any such copy.",
							"license_url": "https:\/\/www.example.com\/terms-and-conditions\/"
						},
						"umbrella": {
							"notice": "Copyright 1999-2025 The Umbrella Corporation",
							"license": "CVE Usage: Umbrella hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Common Vulnerabilities and Exposures (CVE\u00ae). Any copy you make for such purposes is authorized provided that you reproduce Umbrella's copyright designation and this license in any such copy.",
							"license_url": "https:\/\/www.example.org\/Legal\/TermsOfUse"
						}
					}
				}
			}`,
			want: []Vulnerability{
				{
					ID:    "59ea7390-3587-4fce-b267-bf525dbe3e27",
					Title: "Multiple Plugins <= (Various Versions) - Authenticated (Contributor+) Stored DOM-Based Cross-Site Scripting via ThickBox JavaScript Library",
					Software: []Software{
						{
							Type: SoftwareTypePlugin,
							Name: "Foo Bar",
							Slug: "foo-bar",
							AffectedVersions: map[string]AffectedVersion{
								"* - 3.5": {FromVersion: "*", FromInclusive: true, ToVersion: "3.5", ToInclusive: true},
							},
						},
						{
							Type: SoftwareTypeTheme,
							Name: "baz quux – WPBaz",
							Slug: "wpbaz",
							AffectedVersions: map[string]AffectedVersion{
								"* - 4.8.9": {FromVersion: "*", FromInclusive: true, ToVersion: "4.8.9", ToInclusive: true},
								"4.9.2":     {FromVersion: "4.9.2", FromInclusive: true, ToVersion: "4.9.2", ToInclusive: true},
							},
						},
					},
					References: []string{"https://www.example.com/threat-intel/vulnerabilities/id/59ea7390-3587-4fce-b267-bf525dbe3e27?source=api-prod"},
					CVSS:       CVSS{Rating: CVSSRatingMedium},
					CVE:        "CVE-2025-2537",
					CVELink:    "https://www.example.org/CVERecord?id=CVE-2025-2537",
					Published:  "2025-07-03 00:20:54",
					Updated:    "2025-08-19 05:27:14",
					Copyrights: Copyrights{
						"ACME": {
							Notice:     "Copyright 2012-2025 ACME Inc.",
							License:    "ACME hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute this software vulnerability information. Any copy of the software vulnerability information you make for such purposes is authorized provided that you include a hyperlink to this vulnerability record and reproduce ACME's copyright designation and this license in any such copy.",
							LicenseURL: "https://www.example.com/terms-and-conditions/",
						},
						"umbrella": {
							Notice:     "Copyright 1999-2025 The Umbrella Corporation",
							License:    "CVE Usage: Umbrella hereby grants you a perpetual, worldwide, non-exclusive, no-charge, royalty-free, irrevocable copyright license to reproduce, prepare derivative works of, publicly display, publicly perform, sublicense, and distribute Common Vulnerabilities and Exposures (CVE®). Any copy you make for such purposes is authorized provided that you reproduce Umbrella's copyright designation and this license in any such copy.",
							LicenseURL: "https://www.example.org/Legal/TermsOfUse",
						},
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := strings.NewReader(tt.json)

			var got []Vulnerability

			for v, err := range Vulnerabilities(r) {
				if err != nil {
					t.Fatalf("Vulnerabilities() unexpected error: %v", err)
				}

				got = append(got, v)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Vulnerabilities() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestVulnerabilities_Malformed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		json string
	}{
		{
			name: "missing_opening_brace",
			json: `"0001": {"id":"0001","title":"First"} }`,
		},
		{
			name: "missing_closing_brace",
			json: `{ "0001": {"id":"0001","title":"First"}`,
		},
		{
			name: "invalid_value_type",
			json: `{ "0001": "not an object" }`,
		},
		{
			name: "empty_string",
			json: ``,
		},
		{
			name: "whitespace_string",
			json: `     `,
		},
		{
			name: "null",
			json: `null`,
		},
		{
			name: "integer",
			json: `123`,
		},
		{
			name: "empty_array",
			json: `[]`,
		},
		{
			name: "partial_invalid_record",
			json: `{
				"0001": {"id":"0001","title":"First"},
				"0002": {"id":"0002",,,,,"title":"Second"}
			}`,
		},
		{
			name: "string_after_valid",
			json: `
				{ "0001": {"id":"0001","title":"First"} }
				"FOOBAR"
			`,
		},
		{
			name: "malformed_after_valid",
			json: `
				{ "0001": {"id":"0001","title":"First"} }
				FOOBAR
			`,
		},
		{
			name: "two_valid_top_level_objects",
			json: `
				{ "0001": {"id":"0001","title":"First"} }
				{ "0002": {"id":"0002","title":"Second"} }
			`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := strings.NewReader(tt.json)

			var gotErr error

			for _, err := range Vulnerabilities(r) {
				if err != nil {
					gotErr = err
					break
				}
			}

			if gotErr == nil {
				t.Error("Vulnerabilities() unexpected success")
			}
		})
	}
}

func TestCopyrights_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		json       string
		wantValues Copyrights
		wantErr    bool
	}{
		{
			name:       "empty_object",
			json:       `{}`,
			wantValues: Copyrights{},
		},
		{
			name:       "message_only_is_ignored",
			json:       `{"message": "This record contains material that is subject to copyright"}`,
			wantValues: Copyrights{},
		},
		{
			name: "single_entry",
			json: `{
				"ACME": {
					"notice": "Copyright 2012-2025 ACME Inc.",
					"license": "some license",
					"license_url": "https://www.example.com/terms-and-conditions/"
				}
			}`,
			wantValues: Copyrights{
				"ACME": {
					Notice:     "Copyright 2012-2025 ACME Inc.",
					License:    "some license",
					LicenseURL: "https://www.example.com/terms-and-conditions/",
				},
			},
		},
		{
			name: "multiple_entries",
			json: `{
				"message": "This record contains material that is subject to copyright",
				"zebra": {
					"notice": "Copyright Z",
					"license": "license Z",
					"license_url": "https://z.example.com"
				},
				"ACME": {
					"notice": "Copyright 2012-2025 ACME Inc.",
					"license": "some license",
					"license_url": "https://www.example.com/terms-and-conditions/"
				},
				"umbrella": {
					"notice": "Copyright 1999-2025 The Umbrella Corporation",
					"license": "CVE Usage: Umbrella license",
					"license_url": "https://www.example.org/Legal/TermsOfUse"
				}
			}`,
			wantValues: Copyrights{
				"ACME":     {Notice: "Copyright 2012-2025 ACME Inc.", License: "some license", LicenseURL: "https://www.example.com/terms-and-conditions/"},
				"umbrella": {Notice: "Copyright 1999-2025 The Umbrella Corporation", License: "CVE Usage: Umbrella license", LicenseURL: "https://www.example.org/Legal/TermsOfUse"},
				"zebra":    {Notice: "Copyright Z", License: "license Z", LicenseURL: "https://z.example.com"},
			},
		},
		{
			name:    "invalid_json",
			json:    `not json`,
			wantErr: true,
		},
		{
			name:    "json_array",
			json:    `[]`,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got Copyrights

			err := json.Unmarshal([]byte(tt.json), &got)

			if err == nil && tt.wantErr {
				t.Fatalf("UnmarshalJSON() unexpected success")
			}

			if err != nil && !tt.wantErr {
				t.Fatalf("UnmarshalJSON() unexpected error: %v", err)
			}

			if tt.wantErr {
				return
			}

			if diff := cmp.Diff(tt.wantValues, got); diff != "" {
				t.Errorf("UnmarshalJSON() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
