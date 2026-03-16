package internal

import "log/slog"

type Kind string

const (
	KindUnknown Kind = "unknown"
	KindCore    Kind = "core"
	KindPlugin  Kind = "plugin"
	KindTheme   Kind = "theme"
	KindWPMU    Kind = "wpmu"
)

type Record struct {
	Kind     Kind
	Slug     string
	Advisory Advisory
}

func (r Record) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", r.Advisory.ID),
		slog.String("kind", string(r.Kind)),
		slog.String("slug", r.Slug),
	)
}
