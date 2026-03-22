package data

import (
	"errors"
	"fmt"
	"hash/fnv"
	"path/filepath"

	"github.com/typisttech/wpsecadv/internal/packagist"
)

//go:generate go run ./gen -out assets -token $WORDFENCE_INTELLIGENCE_API_KEY -url $WORDFENCE_FEED_URL

var errPackageNotFound = errors.New("package not found")

type Store struct{}

func (s *Store) MarshalAdvisoriesFor(vendor, slug string) ([]byte, error) {
	switch {
	case vendor == "wp-plugin" || vendor == "wpackagist-plugin":
		return pluginAdvisories(slug)
	case vendor == "wp-theme" || vendor == "wpackagist-theme":
		return themeAdvisories(slug)
	case vendor == "wp-core" || packagist.IsCoreImplementation(vendor, slug):
		return coreAdvisories(slug)
	default:
		b, err := themeAdvisories(slug)
		if err != nil {
			return pluginAdvisories(slug)
		}
		return b, nil
	}
}

func coreAdvisories(_ string) ([]byte, error) {
	return coreAdvs, nil
}

func pluginAdvisories(slug string) ([]byte, error) {
	n := fmt.Sprintf("plugin_%x_gen.json", fnv.New32a().Sum([]byte(slug)))
	p := filepath.Join(assetsDir, n)

	b, err := pluginAdvsFS.ReadFile(p)
	if err != nil {
		return nil, errPackageNotFound
	}

	return b, nil
}
