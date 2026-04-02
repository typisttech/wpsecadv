## Context

The server uses stdlib `net/http.ServeMux` with Go 1.22+ method+path patterns.
All response JSON is produced as pre-serialized `[]byte` — handlers write raw
bytes directly rather than marshalling structs at request time. Advisory data
is embedded in the binary at compile time as per-slug JSON files.

The `AdvisoriesMarshaler` interface is the seam between the HTTP layer and the
data layer. The `data.Store` already knows how to map `(vendor, slug)` to a
JSON array of full `Advisory` objects.

Currently both `/p2/{vendor}/{slug}.json` and `/p2/{vendor}/{slug}~dev.json`
fall through to the default 404, which is wrapped by the global
`withCacheControl("max-age=86400")` middleware.

## Goals / Non-Goals

**Goals:**
- Serve `{ "packages": [], "security-advisories": [...] }` at the stable p2 path
- Return 404 for the dev p2 path
- Enable `"metadata": true` in `packages.json`

**Non-Goals:**
- Serving actual package version metadata (`packages` stays always-empty)
- Supporting `~dev.json` advisory lookups
- Changing how the existing `/api/security-advisories/` endpoint works

## Decisions

### Route registration: single handler, not two routes

Go 1.22 ServeMux wildcards match whole path segments. The last segment of
`/p2/foo/bar.json` is `bar.json` and of `/p2/foo/bar~dev.json` is
`bar~dev.json` — there is no way to write two patterns with a per-segment
wildcard that distinguish these at the routing layer alone.

**Decision:** Register one route — `GET /p2/{vendor}/{file}` — and dispatch
inside the handler on whether `{file}` ends in `~dev.json` (→ 404), ends in
`.json` (→ serve advisories, slug = file without `.json` suffix), or anything
else (→ 404).

Considered: registering `GET /p2/{vendor}/{slug}~dev.json` as a literal-suffix
pattern. Rejected because Go ServeMux does not support literal suffixes within
a wildcard segment.

### Reuse `AdvisoriesMarshaler`

The p2 handler accepts the existing `AdvisoriesMarshaler` interface.
`MarshalAdvisoriesFor(vendor, slug)` returns a JSON array of full `Advisory`
objects; that array is used verbatim as the `security-advisories` value in the
p2 response — no transformation needed.

**Decision:** No new interface. The p2 handler takes `AdvisoriesMarshaler`
directly and writes the returned bytes straight into the response body.

Considered: a separate `P2AdvisoriesMarshaler` interface returning only
`advisoryId` + `affectedVersions`. Rejected: unnecessary complexity; Composer
ignores unknown fields, so serving the full advisory objects is harmless and
keeps the data layer unchanged.

### Response assembly: inline byte writes

Consistent with `handleAdvisories`, the p2 handler assembles the response with
direct `w.Write` calls using pre-built byte literals and the `[]byte` from
`MarshalAdvisoriesFor`, avoiding any per-request struct marshalling.

```
{"packages":[],"security-advisories": <bytes from store> }
```

### Cache-Control

The global `withCacheControl("max-age=86400")` wrapper already covers all
responses including 404s, matching the existing test expectation. No per-route
override needed.

## Risks / Trade-offs

- **Test breakage** → `handle_p2_test.go` currently asserts `StatusNotFound`
  for the stable path. That assertion must change to `StatusOK`. Mitigation:
  the test file is small and the change is mechanical.
