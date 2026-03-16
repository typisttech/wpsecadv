package main

import (
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/typisttech/wpsecadv/internal"
)

// jsonArray represent a temporary file that accumulates JSON array elements.
// It is safe for concurrent use.
type jsonArray struct {
	root *os.Root
	path string

	mu    sync.Mutex
	count int

	onceCreate sync.Once
}

func (j *jsonArray) append(a internal.Advisory) error {
	var err error

	j.onceCreate.Do(func() {
		f, e := j.root.OpenFile(j.path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if e != nil {
			err = fmt.Errorf("creating temporary file: %w", e)
			return
		}
		err = f.Close()
	})
	if err != nil {
		return err
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	f, err := j.root.OpenFile(j.path, os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("opening temporary file for appending: %w", err)
	}
	defer f.Close()

	if j.count > 0 {
		if _, err := f.WriteString(","); err != nil {
			return fmt.Errorf("writing comma: %w", err)
		}
	}

	j.count++

	return json.MarshalWrite(f, a)
}

func (j *jsonArray) MarshalJSON() ([]byte, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	bs, err := j.root.ReadFile(j.path)
	if err != nil {
		return nil, fmt.Errorf("reading advisories from temporary file: %w", err)
	}

	bs = append([]byte{'['}, bs...)
	bs = append(bs, byte(']'))

	var as []internal.Advisory
	if err := json.Unmarshal(bs, &as); err != nil {
		absPath := filepath.Join(j.root.Name(), j.path)
		return nil, fmt.Errorf("unmarshaling advisories %s: %w", absPath, err)
	}

	// Ensure stable output for git diffs.
	slices.SortFunc(as, func(a, b internal.Advisory) int {
		return strings.Compare(a.ID, b.ID)
	})

	return json.Marshal(as)
}
