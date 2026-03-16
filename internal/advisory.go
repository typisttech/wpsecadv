package internal

type Advisory struct {
	// ID is the unique identifier of the advisory in the WP Sec Adv database.
	// For example, "WPSECADV/AAA/BBB/CCC".
	ID string `json:"advisoryId"`

	Title string `json:"title"`

	// ReportedAt is the date the issue was reported in UTC.
	// For example, "2006-01-02 15:04:05".
	ReportedAt string `json:"reportedAt"`

	Sources []Source `json:"sources"`

	// Link is a URL to issue disclosure.
	Link string `json:"link,omitempty"`

	CVE string `json:"cve,omitempty"`

	// AffectedVersions is a string in form of a composer constraint.
	// For example, ">=1.0.0,<2.0.0|>3.0.0,<=3.4.0|=5.0.0".
	AffectedVersions string `json:"affectedVersions"`

	// Severity is the lowercased CVSS3 severity rating scale.
	// Possible values: "none", "low", "medium", "high", "critical".
	Severity string `json:"severity,omitempty"`
}

type Source struct {
	// Name is the name of the source.
	// For example, "Wordfence".
	Name string `json:"name"`
	ID   string `json:"remoteId"`
}

// Complete reports whether the advisory has [all required fields] to be
// considered a [SecurityAdvisory] by Composer.
//
// [all required fields]: https://github.com/composer/composer/blob/8851f9df3a9d16a38e27e33002678daaefc0cc30/src/Composer/Advisory/PartialSecurityAdvisory.php#L59-L61
// [SecurityAdvisory]: https://github.com/composer/composer/blob/8851f9df3a9d16a38e27e33002678daaefc0cc30/src/Composer/Advisory/SecurityAdvisory.php
func (a Advisory) Complete() bool {
	return a.partial() &&
		a.Title != "" &&
		a.ReportedAt != "" &&
		len(a.Sources) > 0 &&
		a.Sources[0].ID != ""
}

// partial reports whether the advisory has [all required fields] to be
// considered a [PartialSecurityAdvisory] by Composer.
//
// [all required fields]: https://github.com/composer/composer/blob/8851f9df3a9d16a38e27e33002678daaefc0cc30/src/Composer/Advisory/PartialSecurityAdvisory.php#L63
// [PartialSecurityAdvisory]: https://github.com/composer/composer/blob/8851f9df3a9d16a38e27e33002678daaefc0cc30/src/Composer/Advisory/PartialSecurityAdvisory.php
func (a Advisory) partial() bool {
	return a.ID != "" && a.AffectedVersions != ""
}
