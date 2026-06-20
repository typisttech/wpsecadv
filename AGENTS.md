## Content Exclusion

**IMPORTANT:** **Do not** read, list, or use the following paths or files for analysis or suggestions. These files are machine-generated and very numerous; accessing them causes performance issues and is irrelevant to most tasks:

- internal/data/assets/*_gen.json

## Development Commands

```sh
# Test
mise run test:unit

# Lint
mise run lint

# Format & fix linting issues, including gofmt & gofumpt
mise run fix
```

## Go

### JSON v2

This repository always the Go experimental JSON v2 APIs. Import `encoding/json/v2` and `encoding/json/jsontext` instead of `encoding/json` (v1).
