package main

import (
	"cmp"
	"context"
	_ "embed"
	"encoding/json/v2"
	"fmt"
	"hash/fnv"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"text/template"

	"github.com/typisttech/wpsecadv/internal"
	"golang.org/x/sync/errgroup"
)

var (
	//go:embed assets.tmpl
	assetsTxt  string
	assetsTmpl = template.Must(template.New("assets_gen.go").Parse(assetsTxt)) //nolint:gochecknoglobals

	//go:embed assets_test.tmpl
	assetsTestTxt  string
	assetsTestTmpl = template.Must(template.New("assets_gen_test.go").Parse(assetsTestTxt)) //nolint:gochecknoglobals
)

type store struct {
	logger   *slog.Logger
	parallel int

	moduleRoot *os.Root
	assetsRoot *os.Root
	tempRoot   *os.Root

	coreJSON    *jsonArray
	pluginJSONs sync.Map // key: string (slug), value: *jsonArray
	themeJSONs  sync.Map // key: string (slug), value: *jsonArray
}

func newStore(logger *slog.Logger, parallel int, modDir, assetsDir string) (*store, error) {
	modR, err := os.OpenRoot(modDir)
	if err != nil {
		return nil, err
	}

	tempDir, err := os.MkdirTemp("", "*-wpsecadv-assetsgen")
	if err != nil {
		return nil, err
	}
	tempR, err := os.OpenRoot(tempDir)
	if err != nil {
		return nil, err
	}

	if err := modR.RemoveAll(assetsDir); err != nil {
		return nil, err
	}
	if err := modR.MkdirAll(assetsDir, 0o744); err != nil {
		return nil, err
	}
	assetsR, err := modR.OpenRoot(assetsDir)
	if err != nil {
		return nil, err
	}

	return &store{
		logger:   logger,
		parallel: parallel,

		moduleRoot: modR,
		assetsRoot: assetsR,
		tempRoot:   tempR,

		coreJSON: &jsonArray{
			root: tempR,
			path: assetPath(internal.KindCore, ""),
		},
	}, nil
}

func (s *store) Insert(r internal.Record) error {
	if r.Kind == internal.KindUnknown {
		s.logger.Warn("Skipping unknown kind record", "record", r)
		return nil
	}

	if r.Kind == internal.KindWPMU {
		s.logger.Warn("Skipping WPMU record", "record", r)
		return nil
	}

	if !r.Advisory.Complete() {
		s.logger.Warn("Skipping record with incomplete advisory data", "record", r)
		return nil
	}

	s.logger.Debug("Appending record to JSON array", "record", r)

	j, err := s.jsonArray(r.Kind, r.Slug)
	if err != nil {
		return fmt.Errorf("loading JSON array for %s (%s): %w", r.Slug, r.Kind, err)
	}

	if err := j.append(r.Advisory); err != nil {
		return fmt.Errorf("appending advisory %s for %s (%s) to JSON array: %w", r.Advisory.ID, r.Slug, r.Kind, err)
	}

	s.logger.Debug("Appended advisory to JSON array", "record", r, "path", j.path)
	return nil
}

func (s *store) jsonArray(k internal.Kind, slug string) (*jsonArray, error) {
	switch k {
	case internal.KindCore:
		return s.coreJSON, nil
	case internal.KindPlugin:
		v, _ := s.pluginJSONs.LoadOrStore(slug, &jsonArray{
			root: s.tempRoot,
			path: assetPath(k, slug),
		})

		return v.(*jsonArray), nil //nolint:forcetypeassert
	case internal.KindTheme:
		v, _ := s.themeJSONs.LoadOrStore(slug, &jsonArray{
			root: s.tempRoot,
			path: assetPath(k, slug),
		})

		return v.(*jsonArray), nil //nolint:forcetypeassert
	default:
		return nil, fmt.Errorf("unsupported kind: %s", k)
	}
}

func assetPath(k internal.Kind, slug string) string {
	// Any changes in this function must be mirrored to pluginAdvisories in assets.tmpl.
	return varName(k, slug) + "_gen.json"
}

func varName(k internal.Kind, slug string) string {
	// Any changes in this function must be mirrored to pluginAdvisories in assets.tmpl.
	var prefix string

	switch k {
	case internal.KindCore:
		return "core"
	case internal.KindPlugin:
		prefix = "plugin"
	case internal.KindTheme:
		prefix = "theme"
	default:
		prefix = "unknown"
	}

	slug = strings.ToLower(slug)

	// Use a hash of the slug to ensure the variable name is a valid Go identifier
	// and safe to be used as a file name. The chance of collision is negligible for our use case.
	return fmt.Sprintf("%s_%x", prefix, fnv.New32a().Sum([]byte(slug)))
}

func (s *store) Close(ctx context.Context) error {
	s.logger.Debug("Finalizing code generation")

	defer s.moduleRoot.Close()
	defer s.tempRoot.Close()
	defer s.assetsRoot.Close()

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(max(s.parallel, 4))

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return s.genAssetsJSON(s.coreJSON, internal.KindCore, "")
		}
	})

	var (
		pCount    int
		pAdvCount int
	)
	s.pluginJSONs.Range(func(k, v any) bool {
		slug := k.(string)  //nolint:forcetypeassert
		j := v.(*jsonArray) //nolint:forcetypeassert

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return s.genAssetsJSON(j, internal.KindPlugin, slug)
			}
		})

		pCount++
		pAdvCount += j.count

		return true
	})

	var (
		tCount    int
		tAdvCount int
	)
	s.themeJSONs.Range(func(k, v any) bool {
		slug := k.(string)  //nolint:forcetypeassert
		j := v.(*jsonArray) //nolint:forcetypeassert

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return s.genAssetsJSON(j, internal.KindTheme, slug)
			}
		})

		tCount++
		tAdvCount += j.count

		return true
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return s.genGoFiles()
		}
	})

	if err := g.Wait(); err != nil {
		return err
	}

	s.logger.Debug("Removing temporary directory", "path", s.tempRoot.Name())
	if err := os.RemoveAll(s.tempRoot.Name()); err != nil {
		return fmt.Errorf("removing temporary directory: %w", err)
	}

	s.logger.Info(
		"Generated assets",
		"store", s,
		"core_advisories", s.coreJSON.count,
		"plugins", pCount,
		"plugin_advisories", pAdvCount,
		"themes", tCount,
		"theme_advisories", tAdvCount,
	)

	return nil
}

func (s *store) genAssetsJSON(j *jsonArray, k internal.Kind, slug string) error {
	if j.count == 0 {
		// This should never happen.
		s.logger.Error("Skip generating asset JSON because no advisories found", "kind", k, "slug", slug)
		return nil
	}
	s.logger.Debug("Generating asset JSON", "kind", k, "slug", slug)

	dst, err := s.assetsRoot.OpenFile(assetPath(k, slug), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("creating asset JSON for %s (%s): %w", k, slug, err)
	}
	defer dst.Close()

	if err := json.MarshalWrite(dst, j); err != nil {
		return fmt.Errorf("marshaling JSON array to asset JSON for %s (%s): %w", k, slug, err)
	}

	s.logger.Debug("Generated asset JSON", "kind", k, "slug", slug, "path", dst.Name(), "advisories", j.count)
	return nil
}

func (s *store) genGoFiles() error {
	s.logger.Debug("Generating assets_gen.go and assets_gen_test.go")

	type asset struct {
		Slug    string
		VarName string
		Path    string
	}

	pAssets := make([]asset, 0, 16000)
	s.pluginJSONs.Range(func(k, _ any) bool {
		slug := k.(string) //nolint:forcetypeassert

		pAssets = append(pAssets, asset{
			Slug: slug,
			Path: assetPath(internal.KindPlugin, slug),
		})

		return true
	})
	slices.SortFunc(pAssets, func(a, b asset) int {
		return cmp.Compare(a.Slug, b.Slug)
	})

	tAssets := make([]asset, 0, 2000)
	s.themeJSONs.Range(func(k, _ any) bool {
		slug := k.(string) //nolint:forcetypeassert

		tAssets = append(tAssets, asset{
			Slug:    slug,
			VarName: varName(internal.KindTheme, slug),
			Path:    assetPath(internal.KindTheme, slug),
		})

		return true
	})
	slices.SortFunc(tAssets, func(a, b asset) int {
		return cmp.Compare(a.Slug, b.Slug)
	})

	d := struct {
		Initiator string
		GoPackage string

		AssetsDir string

		CoreAsset    asset
		PluginAssets []asset
		ThemeAssets  []asset
	}{
		Initiator: fmt.Sprintf("go generate in %s/%s:%s", os.Getenv("GOPACKAGE"), os.Getenv("GOFILE"), os.Getenv("GOLINE")),
		GoPackage: os.Getenv("GOPACKAGE"),

		AssetsDir: filepath.Base(s.assetsRoot.Name()),

		CoreAsset: asset{
			Path: assetPath(internal.KindCore, ""),
		},
		PluginAssets: pAssets,
		ThemeAssets:  tAssets,
	}

	tmpls := map[string]*template.Template{
		"assets_gen.go":      assetsTmpl,
		"assets_gen_test.go": assetsTestTmpl,
	}
	for name, tmpl := range tmpls {
		f, err := s.moduleRoot.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			return fmt.Errorf("opening %s: %w", name, err)
		}
		defer f.Close()

		if err := tmpl.Execute(f, d); err != nil {
			return err
		}

		s.logger.Debug("Generated Go file", "name", name, "path", f.Name())
	}

	return nil
}

func (s *store) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("module_dir", s.moduleRoot.Name()),
		slog.String("asset_dir", s.assetsRoot.Name()),
		slog.String("temp_dir", s.tempRoot.Name()),
	)
}
