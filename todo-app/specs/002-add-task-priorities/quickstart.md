# Quickstart: To-Do-Prioritäten validieren

## Prerequisites

- Go 1.25 or newer available as `go`.
- Node.js and npm available as `node` and `npm`.
- Frontend dependencies installed in `frontend/` (`npm install`) for Vitest, TypeScript, and Playwright.
- Playwright Chromium available for the browser test.

## Automated validation

From the repository root, run the Go tests and race detector:

```powershell
Push-Location backend
go test ./...
go test -race ./...
Pop-Location
```

Run the TypeScript build and frontend unit tests:

```powershell
Push-Location frontend
npm run build
npm test -- --run
Pop-Location
```

Run the existing Playwright accessibility/end-to-end test:

```powershell
Push-Location frontend
npm run test:browser
Pop-Location
```

## Behavior to verify

1. Create a todo with no priority and confirm its returned and displayed priority is `medium`.
2. Create todos with `low`, `medium`, and `high`; change one priority without altering its title or status.
3. Combine `status` and `priority` filters and confirm every result matches both; reset either filter to `all` and confirm only the remaining criterion applies.
4. Send an existing status-only `PATCH` request and confirm the todo's priority is preserved.
5. Submit an invalid priority or an update with no mutable fields and confirm the API rejects it without changing the todo.
6. Confirm a todo with no legacy priority is returned and filtered as `medium`.
7. Confirm empty filter results, status changes, deletion, keyboard access, and live feedback continue to work.

The API request and response shapes are specified in [contracts/openapi.yaml](./contracts/openapi.yaml); entity and omission semantics are in [data-model.md](./data-model.md).
