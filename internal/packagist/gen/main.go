package main

import (
	"cmp"
	_ "embed"
	"encoding/json/v2"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"time"
)

var (
	//go:embed core.tmpl
	tmplTxt string
	tmpl    = template.Must(template.New("assets_gen.go").Parse(tmplTxt)) //nolint:gochecknoglobals
)

const (
	dst          = "core_gen.go"
	providersURL = "https://packagist.org/providers/wordpress/core-implementation.json"
)

func main() {
	log.Println("Generating " + dst)

	log.Printf("Fetching providers URL: %s", providersURL)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(providersURL)
	if err != nil {
		log.Fatalf("Failed to fetch packagist providers: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch packagist providers: HTTP %d", resp.StatusCode)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read packagist response: %v", err)
		return
	}

	var pr struct {
		Providers []struct {
			Name string `json:"name"`
		} `json:"providers"`
	}

	if err := json.Unmarshal(bodyBytes, &pr); err != nil {
		log.Printf("Failed to decode packagist providers JSON: %v", err)
		return
	}
	if len(pr.Providers) == 0 {
		log.Printf("packagist provider index contains no providers")
		return
	}

	packages := make([]string, 0, len(pr.Providers)+1)
	for _, p := range pr.Providers {
		packages = append(packages, p.Name)
	}

	// Hand-picked from https://packagist.org/packages/list.json?type=wordpress-core
	packages = append(packages, "pantheon-systems/wordpress-composer")

	slices.Sort(packages)
	packages = slices.Compact(packages)

	set := make(map[string][]string, len(packages))
	for _, p := range packages {
		parts := strings.Split(p, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			log.Printf("Skipping malformed package name: %q", p)
			continue
		}
		set[parts[0]] = append(set[parts[0]], parts[1])
	}

	type datum struct {
		Name     string
		Projects []string
	}
	var vData []datum

	for v, ps := range set {
		// Ensure stable output for git diffs.
		slices.Sort(ps)

		vData = append(vData, datum{
			Name:     v,
			Projects: ps,
		})
	}

	// Ensure stable output for git diffs.
	slices.SortFunc(vData, func(a, b datum) int {
		return cmp.Compare(a.Name, b.Name)
	})

	data := struct {
		Initiator string
		GoPackage string
		Vendors   []datum
	}{
		Initiator: fmt.Sprintf("go generate in %s/%s:%s", os.Getenv("GOPACKAGE"), os.Getenv("GOFILE"), os.Getenv("GOLINE")),
		GoPackage: os.Getenv("GOPACKAGE"),
		Vendors:   vData,
	}

	f, err := os.Create(dst)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return
	}

	cwd, _ := os.Getwd()

	log.Printf("Generated %s\n", filepath.Join(cwd, f.Name()))
}
