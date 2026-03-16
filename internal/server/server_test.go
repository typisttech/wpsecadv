package server

import (
	"errors"
	"net/http"
	"time"
)

type stubStore struct {
	data map[string][]byte
}

func (s *stubStore) MarshalAdvisoriesFor(vendor, slug string) ([]byte, error) {
	if s.data == nil {
		return nil, errors.New("not found")
	}
	key := vendor + "/" + slug
	if b, ok := s.data[key]; ok {
		return b, nil
	}

	return nil, errors.New("not found")
}

func newTestServer() http.Handler {
	return newTestServerWithData(nil)
}

func newTestServerWithData(data map[string][]byte) http.Handler {
	store := &stubStore{data: data}
	modTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	return New(store, modTime)
}
