//go:build e2escripts

package main

import (
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
	t.Logf("server URL: %s", srvURL)

	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		Setup: func(e *testscript.Env) error {
			e.Vars = append(e.Vars,
				"COMPOSER_NO_INTERACTION=true",
				"COMPOSER_NO_AUDIT=true",
				"COMPOSER_NO_SECURITY_BLOCKING=true",
				"COMPOSER_ROOT_VERSION=0.0.1",
				"WPSECADV_SERVER_URL="+srvURL,
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
