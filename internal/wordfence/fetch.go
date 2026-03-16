package wordfence

import (
	"context"
	"net/http"

	"github.com/typisttech/wpsecadv/internal"
	"github.com/typisttech/wpsecadv/internal/wordfence/vuln"
)

type logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
}

func Fetch(ctx context.Context, logger logger, client *http.Client, url string, token vuln.Token) (<-chan internal.Record, <-chan error) {
	logger.Info("Fetching Wordfence vulnerability data feed", "feed", url, "token", token)

	errc := make(chan error, 1)

	f, err := vuln.Get(ctx, client, url, token)
	if err != nil {
		rc := make(chan internal.Record)
		close(rc)
		errc <- err
		close(errc)
		return rc, errc
	}

	rc := make(chan internal.Record)

	go func() {
		defer func() {
			_ = f.Close()
			close(rc)
			close(errc)
		}()

		for v, err := range vuln.Vulnerabilities(f) {
			if err != nil {
				errc <- err
				return
			}

			for _, r := range makeRecords(logger, v) {
				select {
				case rc <- r:
				case <-ctx.Done():
					errc <- ctx.Err()
					return
				}
			}
		}
	}()

	return rc, errc
}
