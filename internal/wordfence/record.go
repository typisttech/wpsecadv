package wordfence

import (
	"cmp"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/typisttech/comver"
	"github.com/typisttech/wpsecadv/internal"
	"github.com/typisttech/wpsecadv/internal/wordfence/vuln"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9](([_.]|-{1,2})?[a-z0-9]+)*$`)

func makeRecords(logger logger, v vuln.Vulnerability) []internal.Record {
	t := makeTitle(v.Title, v.Copyrights)
	l := makeLink(v.References, v.CVELink)
	sev := makeSeverity(v.CVSS.Rating)

	byKindBySlug := make(map[internal.Kind]map[string][]vuln.Software, 5)

	for _, s := range v.Software {
		warn := func(msg string, args ...any) {
			logger.Warn(msg, append([]any{"remote_id", v.ID, "type", s.Type, "slug", s.Slug}, args...)...)
		}

		slug := strings.ToLower(s.Slug)
		if slug == "" {
			warn("Skip empty slug")
			continue
		}

		if !slugPattern.MatchString(slug) {
			warn("Skip invalid slug")
			continue
		}

		kind := makeKind(s.Type, slug)

		if byKindBySlug[kind] == nil {
			byKindBySlug[kind] = make(map[string][]vuln.Software, len(v.Software))
		}
		byKindBySlug[kind][slug] = append(byKindBySlug[kind][slug], s)
	}

	rs := make([]internal.Record, 0, len(v.Software))

	for kind, bySlug := range byKindBySlug {
		for slug, ss := range bySlug {
			warn := func(msg string, args ...any) {
				logger.Warn(msg, append([]any{"remote_id", v.ID, "slug", slug}, args...)...)
			}

			var vavs []vuln.AffectedVersion
			for _, s := range ss {
				vavs = append(vavs, slices.Collect(maps.Values(s.AffectedVersions))...)
			}

			avs := makeConstraints(warn, vavs)
			if avs == "" {
				warn("Skip empty affected versions")
				continue
			}

			rs = append(rs, internal.Record{
				Kind: kind,
				Slug: slug,
				Advisory: internal.Advisory{
					ID:               "WPSECADV/WF/" + v.ID + "/" + slug,
					Title:            t,
					ReportedAt:       v.Published, // Assume timestamp formats match.
					Sources:          []internal.Source{{Name: "Wordfence", ID: v.ID}},
					Link:             l,
					CVE:              v.CVE,
					AffectedVersions: avs,
					Severity:         sev,
				},
			})
		}
	}

	return rs
}

func makeTitle(title string, cs vuln.Copyrights) string {
	var b strings.Builder
	b.WriteString(title)

	// Ensure stable output for git diffs.
	sortedCs := slices.SortedFunc(maps.Values(cs), func(a, b vuln.Copyright) int {
		return cmp.Or(
			cmp.Compare(a.Notice, b.Notice),
			cmp.Compare(a.LicenseURL, b.LicenseURL),
			cmp.Compare(a.License, b.License),
		)
	})

	for _, c := range sortedCs {
		b.WriteString("\n### " + c.Notice + "\n" + c.License + "\n" + c.LicenseURL)
	}

	return b.String()
}

func makeLink(refs []string, cveLink string) string {
	for _, r := range refs {
		if r != "" {
			return r
		}
	}

	return cveLink
}

func makeConstraints(warn func(string, ...any), avs []vuln.AffectedVersion) string { //nolint:cyclop
	cs := make([]string, 0, len(avs))

	var c strings.Builder

	for _, av := range avs {
		c.Reset()

		// Match all.
		if av.FromVersion == "*" && av.ToVersion == "*" {
			return "*"
		}

		if av.FromVersion != "*" && !validVersion(av.FromVersion) {
			warn("Skip invalid FromVersion", "from_version", av.FromVersion)
			continue
		}

		if av.ToVersion != "*" && !validVersion(av.ToVersion) {
			warn("Skip invalid ToVersion", "to_version", av.ToVersion)
			continue
		}

		// Exact version.
		if av.FromInclusive && av.ToInclusive && av.FromVersion == av.ToVersion {
			cs = append(cs, "="+av.FromVersion)
			continue
		}

		// Lower bounded.
		if av.FromVersion != "*" {
			c.WriteByte('>')
			if av.FromInclusive {
				c.WriteByte('=')
			}
			c.WriteString(av.FromVersion)
		}

		// Upper bounded.
		if av.ToVersion != "*" {
			if c.Len() > 0 {
				c.WriteByte(',')
			}

			c.WriteByte('<')
			if av.ToInclusive {
				c.WriteByte('=')
			}
			c.WriteString(av.ToVersion)
		}

		cs = append(cs, c.String())
	}

	// Ensure stable output for git diffs.
	slices.Sort(cs)

	return strings.Join(cs, "|")
}

func validVersion(s string) bool {
	switch s {
	case "*", "": // Fast path.
		return false
	default:
		_, err := comver.Parse(s)
		return err == nil
	}
}

func makeSeverity(r vuln.CVSSRating) string {
	switch r {
	case vuln.CVSSRatingNone:
		return "none"
	case vuln.CVSSRatingLow:
		return "low"
	case vuln.CVSSRatingMedium:
		return "medium"
	case vuln.CVSSRatingHigh:
		return "high"
	case vuln.CVSSRatingCritical:
		return "critical"
	default:
		return ""
	}
}

func makeKind(t vuln.SoftwareType, s string) internal.Kind {
	switch t {
	case vuln.SoftwareTypeCore:
		if s == "wpmu" {
			return internal.KindWPMU
		}

		return internal.KindCore
	case vuln.SoftwareTypePlugin:
		return internal.KindPlugin
	case vuln.SoftwareTypeTheme:
		return internal.KindTheme
	default:
		return internal.KindUnknown
	}
}
