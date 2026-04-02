## 1. Configuration

- [x] 1.1 Update `internal/server/static/packages.json`: change `"metadata": false` to `"metadata": true`

## 2. Handler

- [x] 2.1 Create `internal/server/handle_p2.go` with a `handleP2(store AdvisoriesMarshaler) http.HandlerFunc`
- [x] 2.2 Dispatch on the `{file}` path value: suffix `~dev.json` → 404; suffix `.json` → extract slug and continue; anything else → 404
- [x] 2.3 Call `store.MarshalAdvisoriesFor(vendor, slug)`; on error return 404
- [x] 2.4 On success write `{"packages":[],"security-advisories":<bytes>}` with status 200 and `Content-Type: application/json`

## 3. Routing

- [x] 3.1 Register `GET /p2/{vendor}/{file}` in `internal/server/routes.go` wired to `handleP2(store)`

## 4. Tests

- [x] 4.1 Update `handle_p2_test.go`: known-package stable path now expects 200 with correct JSON body
- [x] 4.2 Add test case: unknown-package stable path expects 404
- [x] 4.3 Confirm dev path case still expects 404 (no change needed, verify only)
- [x] 4.4 Update `stubStore` in `server_test.go` if needed to support p2 test data

## 5. Verification

- [x] 5.1 Run `mise run test:unit` — all tests pass
- [x] 5.2 Run `mise run lint` — no lint errors
