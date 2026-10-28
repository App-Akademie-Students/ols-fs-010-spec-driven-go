---
description: "Task list for Todo-Aufgaben verwalten"
---

# Tasks: Todo-Aufgaben verwalten

**Input**: Design documents from `specs/001-todo-management/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/openapi.yaml`, `quickstart.md`

**Tests**: Automated service, HTTP, repository-concurrency, and frontend behavior tests are required by the implementation plan and project constitution.

**Organization**: Tasks are grouped by the three prioritized user stories. Each implementation task identifies applicable requirements and its verification criterion. US2 and US3 build on the shared task experience while remaining testable with seeded task records.

## Phase 1: Setup

**Purpose**: Initialize the Go backend module and TypeScript/Vitest frontend build and test configuration.

- [ ] T001 [P] Create the Go 1.25+ module in `backend/go.mod`. **Verification**: `go list ./...` from `backend/` recognizes the module and its packages.
- [ ] T002 [P] Configure the TypeScript compiler, Vitest with jsdom, test/build scripts, and browser output under `frontend/dist/` in `frontend/package.json`, `frontend/tsconfig.json`, `frontend/vitest.config.ts`, and a minimal test scaffold in `frontend/src/app.test.ts`. **Verification**: `npm test -- --run` and `npm run build` from `frontend/` execute successfully.

## Phase 2: Foundational

**Purpose**: Establish the shared domain model and repository boundary required by all user stories.

- [ ] T003 Define the `Todo` entity and `open`/`completed` status values in `backend/internal/todo/model.go`; preserve these data-model constraints verbatim: "`id` | Positive integer | Assigned by the repository; unique for the lifetime of the running process; immutable.", "`title` | String | Required; leading and trailing whitespace is removed on creation; a blank result is rejected.", and "`status` | Enum: `open`, `completed` | New tasks start as `open`; either status can transition to the other." **Requirements**: FR-001, FR-002, FR-003, FR-005. **Verification**: model types represent the documented ID, title, and status constraints without exposing mutable repository state.
- [ ] T004 Define the repository interface and in-memory repository type, map, `sync.RWMutex`, and constructor in `backend/internal/todo/repository.go`; keep repository state encapsulated and return detached task values or snapshots. **Requirements**: FR-001, FR-004, FR-005, FR-006, FR-009. **Verification**: interface supports create/list/update/delete and returns values that cannot mutate stored records outside the repository.

## Phase 3: User Story 1 - Aufgaben erfassen und überblicken (Priority: P1) 🎯 MVP

**Goal**: Create tasks with valid titles, start them as open, list titles and statuses, and communicate invalid-title and empty-list outcomes.

**Independent Test**: With an initially empty repository, create a valid task and list it with its trimmed title and open status; reject a whitespace-only title without creating a task; confirm the browser empty state states that no tasks exist and points to task creation.

### Tests for User Story 1

- [ ] T005 [P] [US1] Add service tests for title trimming, blank-title rejection, initial `open` status, and listing in `backend/internal/todo/service_test.go`. **Requirements**: FR-001, FR-002, FR-003, FR-004. **Verification**: tests assert normalized title, rejected blank input, zero side effects on rejection, default status, and list contents.
- [ ] T006 [P] [US1] Add HTTP contract tests for `GET /api/todos` and `POST /api/todos`, including response JSON, `200`/`201` status codes, and invalid-title errors in `backend/internal/httpapi/handler_test.go`. **Requirements**: FR-001, FR-002, FR-003, FR-004. **Verification**: `httptest` assertions cover success schemas/status codes and invalid-title client error without creating a task.
- [ ] T007 [P] [US1] Add Vitest/jsdom behavior tests for task entry, title validation feedback, and the initially empty-list state in `frontend/src/app.test.ts`. **Requirements**: FR-001, FR-002, FR-008; SC-005. **Verification**: DOM assertions confirm valid submission is offered, blank submission reports required title/no creation, and the empty-list message says no tasks exist and indicates task creation.

### Implementation for User Story 1

- [ ] T008 [US1] Implement concurrency-safe task creation and listing in `backend/internal/todo/repository.go`; assign positive, immutable IDs unique for the running process and return detached values. **Requirements**: FR-001, FR-003, FR-004. **Verification**: service/HTTP tests can create distinct tasks and list each with its ID, title, and status.
- [ ] T009 [US1] Implement task creation and listing in `backend/internal/todo/service.go`; trim leading and trailing title whitespace, reject blank results, and set accepted tasks to `open`. **Requirements**: FR-001, FR-002, FR-003, FR-004. **Verification**: T005 passes for normalized titles, blank rejection, initial status, and complete listing.
- [ ] T010 [US1] Implement `GET /api/todos` and `POST /api/todos` JSON handlers in `backend/internal/httpapi/handler.go` according to `specs/001-todo-management/contracts/openapi.yaml`. **Requirements**: FR-001, FR-002, FR-003, FR-004. **Verification**: T006 passes against the documented response schemas and status codes.
- [ ] T011 [US1] Create the semantic task-entry form, task list, status display, validation feedback, and empty-list state in `frontend/index.html`, `frontend/styles.css`, and `frontend/src/app.ts`; render titles as text and make controls keyboard-operable. **Requirements**: FR-001, FR-002, FR-003, FR-004, FR-008; SC-005. **Verification**: T007 passes and the DOM provides an accessible title input, submission control, task list, and distinct validation and empty-list feedback.
- [ ] T012 [US1] Wire the repository, service, HTTP handlers, static frontend served from `../frontend`, and startup/error handling on port `8080` in `backend/cmd/todo/main.go`. **Requirements**: FR-001, FR-002, FR-003, FR-004, FR-008. **Verification**: `go test ./...` from `backend/` passes and the quickstart can open the served page and complete the US1 flow.

**Checkpoint**: US1 supports task creation and listing with required title validation, initial status, and empty-list feedback.

## Phase 4: User Story 2 - Aufgaben erledigen und entfernen (Priority: P1)

**Goal**: Change a task between open and completed in either direction and delete tasks so they disappear from every list view.

**Independent Test**: Seed two tasks; complete and reopen one, delete the other, and confirm the status change does not alter the remaining task.

### Tests for User Story 2

- [ ] T013 [P] [US2] Add service tests for both status transitions, unknown IDs, deletion, and isolation of changes to other tasks in `backend/internal/todo/service_test.go`. **Requirements**: FR-005, FR-006, FR-009. **Verification**: tests assert both transitions, not-found behavior, deleted-task absence, and unchanged unrelated task state.
- [ ] T014 [P] [US2] Add HTTP contract tests for `PATCH /api/todos/{id}` and `DELETE /api/todos/{id}`, including invalid status, unknown ID, response codes, and response bodies in `backend/internal/httpapi/handler_test.go`. **Requirements**: FR-005, FR-006, FR-009. **Verification**: `httptest` assertions match the OpenAPI contract for valid mutations, invalid status, and unknown IDs.
- [ ] T015 [P] [US2] Add Vitest/jsdom behavior tests for completing/reopening a task, deleting a task, and preserving other task rows in `frontend/src/app.test.ts`. **Requirements**: FR-005, FR-006, FR-009; SC-002, SC-004. **Verification**: DOM assertions confirm status changes in both directions, deletion removes the targeted task, and unrelated tasks remain unchanged.

### Implementation for User Story 2

- [ ] T016 [US2] Implement concurrency-safe status updates and deletion in `backend/internal/todo/repository.go`; accept only `open` or `completed` and return not-found for absent IDs. **Requirements**: FR-005, FR-006, FR-009. **Verification**: T013 passes for both allowed statuses, missing IDs, and deletion.
- [ ] T017 [US2] Implement status transitions and deletion in `backend/internal/todo/service.go`; preserve unrelated tasks and return an explicit not-found result for missing IDs. **Requirements**: FR-005, FR-006, FR-009. **Verification**: T013 passes with unrelated task values unchanged after each operation.
- [ ] T018 [US2] Implement `PATCH /api/todos/{id}` and `DELETE /api/todos/{id}` in `backend/internal/httpapi/handler.go` according to `specs/001-todo-management/contracts/openapi.yaml`. **Requirements**: FR-005, FR-006, FR-009. **Verification**: T014 passes for response codes, response bodies, invalid status, and unknown IDs.
- [ ] T019 [US2] Add keyboard-operable complete/reopen and delete controls with clear action/error feedback in `frontend/src/app.ts`, `frontend/index.html`, and `frontend/styles.css`; refresh visible state after successful mutations. **Requirements**: FR-005, FR-006, FR-009. **Verification**: T015 passes and the updated list reflects mutations without changing unrelated task rows.

**Checkpoint**: Users can complete, reopen, and delete tasks; changes affect only the targeted task and are reflected in the list.

## Phase 5: User Story 3 - Aufgaben nach Status filtern (Priority: P2)

**Goal**: Filter the list to all, open, or completed tasks and show a useful, distinguishable empty state when a filter has no matches.

**Independent Test**: Seed open and completed tasks; verify each status filter includes only matching records, `all` includes both, an omitted filter means all, and an empty filter state names the active status and offers a switch to “Alle”.

### Tests for User Story 3

- [ ] T020 [P] [US3] Add service tests for `all`, `open`, and `completed` filtering, default-to-all behavior, and unsupported filter rejection in `backend/internal/todo/service_test.go`. **Requirements**: FR-007, FR-009. **Verification**: tests assert exact membership for each filter, default behavior, and invalid-filter errors.
- [ ] T021 [P] [US3] Add HTTP contract tests for `GET /api/todos?status=all|open|completed`, omitted status, and unsupported status values in `backend/internal/httpapi/handler_test.go`. **Requirements**: FR-007, FR-009. **Verification**: `httptest` assertions confirm contract values, filtered response contents, and client errors for unsupported values.
- [ ] T022 [P] [US3] Add Vitest/jsdom behavior tests for filter selection, no-match messaging, active-status naming, and switching to “Alle” in `frontend/src/app.test.ts`. **Requirements**: FR-007, FR-008, FR-009; SC-003, SC-005. **Verification**: DOM assertions confirm each filter’s task set and distinguish a filter-empty state from a globally empty list.

### Implementation for User Story 3

- [ ] T023 [US3] Implement status filtering in `backend/internal/todo/repository.go` and `backend/internal/todo/service.go`; support `all`, `open`, and `completed`, default an omitted filter to `all`, and reject other values. **Requirements**: FR-007, FR-009. **Verification**: T020 passes for all valid filters, omitted filter, and invalid values.
- [ ] T024 [US3] Extend `GET /api/todos` in `backend/internal/httpapi/handler.go` to parse optional `status` and return the documented client error for unsupported values. **Requirements**: FR-007, FR-009. **Verification**: T021 passes against every supported query value and invalid-filter response.
- [ ] T025 [US3] Add all/open/completed filter controls in `frontend/index.html` and `frontend/styles.css`, connect them to the API in `frontend/src/app.ts`, and display the distinct filter-empty state specified in FR-008. **Requirements**: FR-007, FR-008, FR-009; SC-003, SC-005. **Verification**: T022 passes; the no-match message names the active status and offers a path to “Alle”.

**Checkpoint**: Every filter shows its intended task group, and status changes/deletions remain consistent in filtered views.

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Verify usability outcome measurement, repository concurrency, frontend behavior/build, and documented end-to-end flows.

- [ ] T026 [P] Add concurrent create, list, status-update, and delete coverage for the in-memory repository in `backend/internal/todo/repository_test.go`. **Requirements**: FR-001, FR-004–FR-006, FR-009. **Verification**: concurrent test cases pass under `go test -race ./...` from `backend/` with no data races.
- [ ] T027 [P] Run Vitest/jsdom behavior tests and the TypeScript production build from `frontend/`; resolve issues in `frontend/src/app.test.ts`, `frontend/src/app.ts`, `frontend/vitest.config.ts`, `frontend/tsconfig.json`, and `frontend/package.json`. **Requirements**: FR-001, FR-002, FR-003, FR-004, FR-005, FR-006, FR-007, FR-008, FR-009; SC-002, SC-003, SC-004, SC-005. **Verification**: `npm test -- --run` and `npm run build` both exit successfully.
- [ ] T028 Run `go test ./...` and `go test -race ./...` from `backend/`; resolve failures in `backend/internal/todo/`, `backend/internal/httpapi/`, and `backend/cmd/todo/`. **Requirements**: FR-001, FR-002, FR-003, FR-004, FR-005, FR-006, FR-007, FR-008, FR-009; SC-002, SC-003, SC-004. **Verification**: both Go commands pass.
- [ ] T029 Run the SC-001 usability protocol and record the outcome in `specs/001-todo-management/usability-results.md`. **Requirements**: FR-001, FR-004; SC-001. **Verification**: record results for 20 people with no prior application use; pass only if all 20 independently create a valid task and find that same task in “Alle”.
- [ ] T030 Validate browser and API acceptance flows against the application and update inaccurate commands or outcomes in `specs/001-todo-management/quickstart.md`. **Requirements**: FR-001, FR-002, FR-003, FR-004, FR-005, FR-006, FR-007, FR-008, FR-009; SC-002, SC-003, SC-004, SC-005. **Verification**: every quickstart scenario has a documented expected outcome, including validation, both empty-state cases, status changes, deletion, and filters.

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies; T001 and T002 configure separate projects and can run in parallel.
- **Foundational (Phase 2)**: Depends on Setup; complete T003 before T004. Both block user-story implementation.
- **User Stories (Phases 3–5)**: Depend on the foundational model and repository boundary. US1 establishes create/list behavior; US2 integrates mutation behavior; US3 integrates filtering with both statuses and mutations.
- **Polish (Phase 6)**: Depends on implementation of all desired stories. T026 and T027 validate independent subsystems and can run in parallel.

### User Story Dependencies

- **US1 (P1)**: Starts after Phase 2; no dependency on another story.
- **US2 (P1)**: Integrates after US1; its service, API, and frontend tests can seed tasks directly.
- **US3 (P2)**: Integrates after US1 and US2; its service, API, and frontend tests can seed both statuses directly.

### Parallel Opportunities

- T001 and T002 can run in parallel because they configure independent Go and frontend manifests.
- In each story, the service tests, HTTP contract tests, and frontend behavior tests target separate files and can be authored in parallel.
- T026 and T027 can run in parallel after feature implementation stabilizes because repository concurrency tests and frontend checks exercise separate subsystems.
- User stories are sequenced for integrated delivery; independent story tests can use seeded data.

## Parallel Examples

### User Story 1

```text
T005: Service tests in backend/internal/todo/service_test.go
T006: HTTP contract tests in backend/internal/httpapi/handler_test.go
T007: Frontend behavior tests in frontend/src/app.test.ts
```

### User Story 2

```text
T013: Service tests in backend/internal/todo/service_test.go
T014: HTTP contract tests in backend/internal/httpapi/handler_test.go
T015: Frontend behavior tests in frontend/src/app.test.ts
```

### User Story 3

```text
T020: Service tests in backend/internal/todo/service_test.go
T021: HTTP contract tests in backend/internal/httpapi/handler_test.go
T022: Frontend behavior tests in frontend/src/app.test.ts
```

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 and Phase 2.
2. Complete Phase 3, including service, HTTP, and frontend tests plus create/list behavior and browser UI.
3. Validate US1 independently using its acceptance criteria, `go test ./...` from `backend/`, and `npm test -- --run` from `frontend/`.
4. Stop at the US1 checkpoint for an MVP demonstration if appropriate.

### Incremental Delivery

1. Add US2 after the US1 checkpoint; validate status changes and deletion without regressing create/list.
2. Add US3 after the US2 checkpoint; validate every filter and both empty-state cases.
3. Complete Phase 6, including the 20-person SC-001 usability protocol, race detector, frontend behavior tests/build, and quickstart acceptance flows.

## Notes

- Every task uses an unchecked checkbox and sequential task ID; story-phase tasks have their story label, and each task states applicable FR/SC IDs and an objective verification criterion.
- `[P]` marks tasks that can proceed independently without editing the same files or depending on incomplete work.
- Story traceability: US1 covers FR-001–FR-004 and the globally empty state in FR-008; US2 covers FR-005, FR-006, and FR-009; US3 covers FR-007–FR-009.
- Run Go commands from `backend/` and npm commands from `frontend/`.
- Do not add durable storage, accounts, collaboration, title editing, due dates, categories, or reminders; they are outside the approved scope.
