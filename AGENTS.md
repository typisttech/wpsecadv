# AGENTS.md

The role of this file is to describe common mistakes and confusion points that agents might encounter as they work in this project. If you ever encounter something in the project that surprises you, please alert the developer working with you and indicate that this is the case in this `AGENTS.md` file to help prevent future agents from having the same issue.

## Content Exclusion

**IMPORTANT:** **Do not** read, index, or use the following paths or files for analysis or suggestions. These files are machine-generated and very numerous; accessing them causes performance issues and is irrelevant to most tasks:

- internal/data/assets/*_gen.json

## Development Commands

Test: `mise run test:unit`
Lint: `mise run lint`
Fix linting issues: `mise run lint --fix`
Format source code: `mise run fmt`

## Go

### JSON v2

This repository always the Go experimental JSON v2 APIs. Import `encoding/json/v2` and `encoding/json/jsontext` instead of `encoding/json` (v1).

### Loop Variables

Capturing loop variables is no longer need. The pattern `x := x` is unnecessary.

### Table-Driven Tests

Always prefer table-driven tests: structure test cases as slices of structs with descriptive snake_cased `name` fields and use subtests (`t.Run`) for each case. This keeps tests consistent, easy to extend, and reduces duplication across similar scenarios.
