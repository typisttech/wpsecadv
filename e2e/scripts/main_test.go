//go:build e2escripts

package main

import (
	"encoding/json/v2"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"golang.org/x/mod/semver"
)

func TestScripts(t *testing.T) {
	srvURL := os.Getenv("WPSECADV_SERVER_URL")
	if srvURL == "" {
		t.Fatal("WPSECADV_SERVER_URL environment variable is not set")
	}
	t.Logf("server URL: %q", srvURL)

	caFile := os.Getenv("TESTSCRIPT_COMPOSER_CAFILE")
	t.Logf("ca file: %q", caFile)

	type repo struct {
		Type    string `json:"type"`
		URL     string `json:"url"`
		Options struct {
			SSL struct {
				CAFile string `json:"cafile,omitzero"`
			} `json:"ssl,omitzero"`
		} `json:"options,omitzero"`
	}

	r := repo{
		Type: "composer",
		URL:  srvURL,
	}
	if caFile != "" {
		r.Options.SSL.CAFile = caFile
	}

	rb, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("json.Marshal error: %v", err)
	}
	t.Logf("repo: %s", rb)

	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Setup: func(e *testscript.Env) error {
			e.Vars = append(e.Vars,
				"COMPOSER_NO_INTERACTION=true",
				"COMPOSER_NO_AUDIT=true",
				"COMPOSER_NO_SECURITY_BLOCKING=true",
				"COMPOSER_ROOT_VERSION=0.0.1",
				"REPO="+string(rb),
			)

			dir := os.Getenv("TESTSCRIPT_COMPOSER_CACHE_DIR")
			dir = strings.TrimSpace(dir)
			if dir != "" {
				e.Vars = append(e.Vars, "COMPOSER_CACHE_DIR="+dir)
			}

			return nil
		},
		Condition: func(cond string) (bool, error) {
			if !strings.HasPrefix(cond, "composer:") {
				return false, fmt.Errorf("unknown condition: %s", cond)
			}

			// This is set by the image.
			cur := os.Getenv("COMPOSER_VERSION")
			cur = "v" + cur
			if cur == "" || !semver.IsValid(cur) {
				// Assume latest composer version when running outside containers.
				return true, nil
			}

			targ := strings.TrimPrefix(cond, "composer:")
			targ = "v" + targ + ".0"

			if !semver.IsValid(targ) {
				return false, fmt.Errorf("invalid composer version: %q", targ)
			}

			return semver.Compare(cur, targ) >= 0, nil
		},
	})
}
