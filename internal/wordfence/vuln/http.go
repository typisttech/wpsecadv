package vuln

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const FeedProduction = "https://www.wordfence.com/api/intelligence/v3/vulnerabilities/production"

var ErrUnexpectedResponseStatus = errors.New("unexpected response status")

// Get fetches the feed URL and returns the response body.
//
// It is the caller's responsibility to eventually close the returned body.
//
// A custom [http.Client] may be provided (for TLS settings, timeouts, transport
// tuning). If client is nil, [http.DefaultClient] is used.
func Get(ctx context.Context, client *http.Client, url string, token Token) (io.ReadCloser, error) {
	if client == nil {
		client = http.DefaultClient
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+string(token))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 8*1024))
		_ = resp.Body.Close()

		return nil, fmt.Errorf("%s %s: %w", resp.Status, string(snippet), ErrUnexpectedResponseStatus)
	}

	return resp.Body, nil
}

type Token string

func (t Token) String() string {
	if t == "" {
		return "EMPTY"
	}

	return "REDACTED"
}

// Set implements the [flag.Value] interface.
func (t *Token) Set(value string) error {
	*t = Token(value)
	return nil
}
