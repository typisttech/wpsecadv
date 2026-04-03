# AGENTS.md

## Content Exclusion

**IMPORTANT:** **Do not** read, list, or use the following paths or files for analysis or suggestions. These files are machine-generated and very numerous; accessing them causes performance issues and is irrelevant to most tasks:

- internal/data/assets/*_gen.json

## Development Commands

Test: `mise run test:unit`
Lint: `mise run lint`
Fix linting issues: `mise run lint --fix`
Format source code: `mise run fmt`

## Go

### JSON v2

This repository always the Go experimental JSON v2 APIs. Import `encoding/json/v2` and `encoding/json/jsontext` instead of `encoding/json` (v1).

### Table-Driven Tests

Always prefer table-driven tests: structure test cases as slices of structs with descriptive snake_cased `name` fields and use subtests (`t.Run`) for each case. This keeps tests consistent, easy to extend, and reduces duplication across similar scenarios.
