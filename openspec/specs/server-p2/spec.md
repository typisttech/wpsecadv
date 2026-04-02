## Purpose

Provide Composer p2 package metadata endpoints so clients can discover
security advisories for packages. This spec defines the stable and dev p2
endpoints and the required packages.json metadata for advisory discovery.

## Requirements

### Requirement: Stable p2 path returns security advisories
The server SHALL respond to `GET /p2/{vendor}/{slug}.json` with HTTP 200 and a
JSON body containing `"packages": []` and a `"security-advisories"` array
populated from `AdvisoriesMarshaler.MarshalAdvisoriesFor(vendor, slug)`.

#### Scenario: Known package with advisories
- **WHEN** a GET request is made to `/p2/{vendor}/{slug}.json` for a vendor/slug that has advisories
- **THEN** the response status SHALL be 200
- **THEN** the response body SHALL be `{"packages":[],"security-advisories":[...]}` where the array contains the advisory objects for that package

#### Scenario: Unknown package
- **WHEN** a GET request is made to `/p2/{vendor}/{slug}.json` for a vendor/slug with no known advisories
- **THEN** the response status SHALL be 404

#### Scenario: Cache-Control header on stable path
- **WHEN** a GET request is made to `/p2/{vendor}/{slug}.json`
- **THEN** the response SHALL include a `Cache-Control: max-age=86400` header

### Requirement: Dev p2 path always returns 404
The server SHALL respond to `GET /p2/{vendor}/{slug}~dev.json` with HTTP 404
regardless of vendor or slug.

#### Scenario: Dev path returns 404
- **WHEN** a GET request is made to `/p2/{vendor}/{slug}~dev.json`
- **THEN** the response status SHALL be 404

#### Scenario: Cache-Control header on dev path
- **WHEN** a GET request is made to `/p2/{vendor}/{slug}~dev.json`
- **THEN** the response SHALL include a `Cache-Control: max-age=86400` header

### Requirement: Composer repository metadata enables p2 advisory discovery
The `packages.json` repository metadata SHALL set `"metadata": true` in the
`security-advisories` block so that Composer clients look up advisories via
the `/p2/` endpoints.

#### Scenario: packages.json reports metadata enabled
- **WHEN** a client fetches `/packages.json`
- **THEN** the response body SHALL contain `"metadata": true` inside the `security-advisories` object
