# Quickstart: Todo-Aufgaben verwalten

## Prerequisites

- Go 1.25 or newer.
- Node.js and npm for compiling the TypeScript frontend.
- Playwright's Chromium browser for automated keyboard-accessibility tests
  (`npx playwright install chromium`).
- A modern browser.

## Build and test

Run frontend and backend commands from their respective directories:

```powershell
Push-Location frontend
npm install
npx playwright install chromium
npm test -- --run
npm run build
npm run test:browser
Pop-Location

Push-Location backend
go test ./...
go test -race ./...
Pop-Location
```

Expected results: the TypeScript build exits successfully, the Vitest/jsdom behavior tests pass,
the Playwright/Chromium keyboard-accessibility test passes against the Go application, Go tests
pass, and the race-enabled test run reports no data races.

## Run the application

```powershell
Push-Location backend
go run ./cmd/todo
```

Keep the process running with `backend/` as its working directory, then open
`http://localhost:8080` in a browser. The server serves static files from `../frontend`.
The process uses in-memory storage; restarting it clears the task list.

## Browser acceptance flow

1. On an empty list, confirm the empty state says that no tasks exist and points to task creation.
2. Submit a title containing only Unicode `White_Space` characters; confirm no task is created
   and feedback explicitly says a title is required and the task was not created.
3. Add two distinct task titles; confirm each shows **„Aufgabe erstellt.“** and both appear as
   open in creation order in the default **Alle** view. Add the same title twice; confirm both
   entries remain independently available.
4. Mark one task completed; confirm **„Status geändert.“** appears and it is in **Erledigt**
   but not **Offen**.
5. Reopen that task; confirm **„Status geändert.“** appears and it is in **Offen** again.
6. Delete the other task without a confirmation prompt; confirm it is absent from **Alle**,
   **Offen**, and **Erledigt**, and that **„Aufgabe gelöscht.“** is shown.
7. Select a filter with no matching tasks; confirm the empty state names the active status,
   says there are no matching tasks, and offers a way to select **Alle**. Confirm it differs
   from the globally empty-list state and from validation feedback.
8. Use only the keyboard to reach and operate every control; confirm focus is visible and
   success/error feedback is available to assistive technology.
9. Cause a request failure (for example, attempt to change or delete an ID that no longer
   exists); confirm the UI shows an error rather than success feedback and uses the specified
   not-available message.

## SC-001 usability protocol

Recruit 20 participants who have not previously used the application. Give each participant
the same task: create a task with a valid title and find that same task in **Alle**, without
assistance. Record whether both steps were completed for that task. SC-001 passes only when
all 20 participants complete both steps. Record the aggregate and per-participant outcomes in
`usability-results.md`.

## API acceptance flow

The request and response schemas, paths, status codes and error shape are defined in
[`contracts/openapi.yaml`](./contracts/openapi.yaml).

1. `GET /api/todos` returns an empty `todos` array initially.
2. `POST /api/todos` with `{"title":"Prepare demo"}` returns `201` and an `open` todo.
3. `PATCH /api/todos/{id}` with `{"status":"completed"}` returns the updated todo; setting
   status back to `open` is also supported.
4. `GET /api/todos?status=completed` and `?status=open` return only matching records.
5. `DELETE /api/todos/{id}` returns `204`; subsequent lists do not include the deleted todo.
6. Invalid JSON, blank titles, unknown fields, unsupported status/filter values and unknown
   IDs return documented client errors with a non-empty error message; unexpected server
   failures return a documented server error. The UI shows failures as errors, not successes.
7. List and filter responses preserve task creation order; duplicate titles remain distinct
   tasks by ID.
