# Quickstart: Todo-Aufgaben verwalten

## Prerequisites

- Go 1.25 or newer.
- Node.js and npm for compiling the TypeScript frontend.
- A modern browser.

## Build and test

From the application repository root:

```powershell
npm install
npm run build
go test ./...
go test -race ./...
```

Expected results: the TypeScript build exits successfully, Go tests pass, and the race-enabled
test run reports no data races.

## Run the application

```powershell
go run ./cmd/todo
```

Open `http://localhost:8080` in a browser. The process uses in-memory storage; restarting it
clears the task list.

## Browser acceptance flow

1. On an empty list, confirm a clear empty-state message and the ability to add a task.
2. Submit a title containing only spaces; confirm no task is created and a helpful validation
   message appears.
3. Add two distinct task titles; confirm both appear as open in the default **Alle** view.
4. Mark one task completed; confirm it appears in **Erledigt** and not **Offen**.
5. Reopen that task; confirm it appears in **Offen** again.
6. Delete the other task; confirm it is absent from **Alle**, **Offen**, and **Erledigt**.
7. Select a filter with no matching tasks; confirm the empty state is understandable.
8. Use the interface with the keyboard and confirm actions and status are perceivable.

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
   IDs return documented client errors without a success-shaped response.
