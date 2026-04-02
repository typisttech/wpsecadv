## Why

Composer clients can discover security advisories directly from Packagist v2
`/p2/{vendor}/{slug}.json` endpoints when a repository sets `"metadata": true`.
This change implements those endpoints and enables the flag, giving Composer a
standards-compliant path to per-package advisory data.

## What Changes

- Add `GET /p2/{vendor}/{slug}.json` — responds with Packagist v2 metadata:
  `packages: []` (always empty; no version metadata served) and a
  `security-advisories` array of `{ advisoryId, affectedVersions }` objects
  drawn from `AdvisoriesMarshaler.MarshalAdvisoriesFor(vendor, slug)`
- Add `GET /p2/{vendor}/{slug}~dev.json` — always returns 404 (dev/branch
  aliases carry no advisory metadata)
- Update `internal/server/static/packages.json` — flip `"metadata": false`
  to `"metadata": true` so Composer clients use the new endpoints

## Capabilities

### New Capabilities
- `p2-package-metadata`: Packagist v2 per-package endpoint that surfaces
  security advisories in the format Composer expects at `/p2/{vendor}/{slug}.json`

### Modified Capabilities

(none)

## Impact

- `internal/server/routes.go` — new route registrations for `/p2/` paths
- `internal/server/handle_p2.go` — new handler file
- `internal/server/static/packages.json` — `metadata` flag update
- `internal/server/handle_p2_test.go` — existing tests assert 404 for both
  paths; stable path assertion changes to 200
