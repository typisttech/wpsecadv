package server

import (
	_ "embed"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

var (
	//go:embed static/error_packages_xor_updated_since.json
	packagesXORUpdatedSince []byte

	//go:embed static/error_updated_since_not_supported.json
	updatedSinceNotSupported []byte

	//go:embed static/error_parameters_malformed.json
	parametersMalformed []byte

	//go:embed static/error_packages_missing.json
	packagesMissing []byte
)

type AdvisoriesMarshaler interface {
	MarshalAdvisoriesFor(vendor, slug string) ([]byte, error)
}

func handleAdvisories(store AdvisoriesMarshaler) http.HandlerFunc { //nolint:cyclop
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//gosec:disable G120 -- False positive, wait for https://github.com/securego/gosec/pull/1605
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(parametersMalformed)

			return
		}

		packages := make([]string, 0, len(r.Form))
		for k, v := range r.Form {
			if !strings.HasPrefix(k, "packages[") || !strings.HasSuffix(k, "]") {
				continue
			}

			packages = append(packages, v...)
		}

		_, hasUpdatedSince := r.Form["updatedSince"]

		switch {
		case hasUpdatedSince && len(packages) > 0:
			w.WriteHeader(http.StatusBadRequest)
			w.Write(packagesXORUpdatedSince)

			return
		case hasUpdatedSince:
			w.WriteHeader(http.StatusBadRequest)
			w.Write(updatedSinceNotSupported)

			return
		case !slices.ContainsFunc(packages, func(p string) bool { return p != "" }):
			w.WriteHeader(http.StatusBadRequest)
			w.Write(packagesMissing)

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"advisories":{`))

		var written bool
		for _, p := range packages {
			pkg := strings.ToLower(p)

			parts := strings.Split(pkg, "/")
			if len(parts) != 2 {
				continue
			}

			vendor, slug := parts[0], parts[1]
			if vendor == "" || slug == "" {
				continue
			}

			advisories, err := store.MarshalAdvisoriesFor(vendor, slug)
			if err != nil {
				// TODO: Handle unexpected errors.
				continue
			}

			if written {
				w.Write([]byte(","))
			}
			written = true

			w.Write([]byte(strconv.Quote(pkg)))
			w.Write([]byte(":"))
			w.Write(advisories)
		}

		w.Write([]byte("}}"))
	}
}
