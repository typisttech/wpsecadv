package vuln

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
)

var ErrMalformedJSON = errors.New("malformed JSON")

// Vulnerabilities reads from r and decodes it into [Vulnerability] and error
// pairs.
//
// It returns a single-use iterator.
//
// Usage:
//
//	r := strings.NewReader(someJSONString)
//
//	for v, err := range Vulnerabilities(r) {
//		if err != nil {
//		// Handle error.
//			break // Must break after the first non-nil error.
//		}
//		// Use v.
//	}
func Vulnerabilities(r io.Reader) func(yield func(Vulnerability, error) bool) { //nolint:cyclop
	return func(yield func(Vulnerability, error) bool) {
		dec := jsontext.NewDecoder(r)

		t, err := dec.ReadToken()
		if err != nil {
			yield(Vulnerability{}, err)

			return
		}

		if k := t.Kind(); k != jsontext.KindBeginObject {
			yield(Vulnerability{}, fmt.Errorf("expected JSON object starts with `{`, got %v: %w", k, ErrMalformedJSON))

			return
		}

		for {
			// Consume object keys.
			if t, err = dec.ReadToken(); err != nil {
				yield(Vulnerability{}, err)

				return
			}

			if k := t.Kind(); k == jsontext.KindEndObject {
				break
			}

			var v Vulnerability

			err := json.UnmarshalDecode(dec, &v)
			if err != nil {
				yield(Vulnerability{}, err)

				return
			}

			if !yield(v, nil) {
				return
			}
		}

		t, err = dec.ReadToken()
		if err != nil && !errors.Is(err, io.EOF) {
			yield(Vulnerability{}, err)

			return
		}

		if err == nil {
			yield(Vulnerability{}, fmt.Errorf("expected nothing after JSON object trailing `}`, got %v: %w", t, ErrMalformedJSON))

			return
		}
	}
}

func (c *Copyrights) UnmarshalJSON(b []byte) error {
	var raw map[string]any

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}
	// Ignore the 'message' field.
	delete(raw, "message")

	m := make(Copyrights, len(raw))
	var cp Copyright

	for k, v := range raw {
		bs, err := json.Marshal(v)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bs, &cp); err != nil {
			return err
		}

		m[k] = cp
	}

	*c = m

	return nil
}
