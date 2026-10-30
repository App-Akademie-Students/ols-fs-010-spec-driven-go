# Tasks: To-Do-Prioritäten festlegen und filtern

**Input**: Design documents from `specs/002-add-task-priorities/`

**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/openapi.yaml`, `quickstart.md`

**Tests**: Test tasks are included because the project constitution requires automated coverage for behavior changes in the service, HTTP contract, and UI.

**Organization**: Tasks are grouped by user story and sequenced by priority. Every task has a file path; task IDs referenced in dependencies refer to this document.

## Phase 1: Setup

**Purpose**: Verify the existing project test tools before feature work; do not add dependencies.

- [X] T001 Verify the existing test commands and installed tools from `backend/go.mod` and `frontend/package.json`; report any missing tool without changing dependency manifests.

---

## Phase 2: Foundational

**Purpose**: Define the shared domain value required by all three user stories.

- [X] T002 Add the `Priority` domain type with exactly `low`, `medium`, and `high`, a `Todo.Priority` field, and a reusable validity check in `backend/internal/todo/model.go` (FR-001).

**Checkpoint**: The shared priority type is available before story-specific service, API, and UI work.

---

## Phase 3: User Story 1 - Priorität beim Erstellen und Bearbeiten festlegen (Priority: P1)

**Goal**: Users can choose a priority when creating a todo, change it later, and see the current value without changing the todo title or status.

**Independent Test**: Create todos with each of `low`, `medium`, and `high`; change a todo's priority and verify the displayed priority changes while its title and status remain unchanged.

### Tests for User Story 1

- [X] T003 [P] [US1] Add service tests for creating and changing todos with each valid priority, preserving title and status, and rejecting an invalid priority in `backend/internal/todo/service_test.go` (FR-001, FR-002, FR-004, FR-009).
- [X] T004 [P] [US1] Add HTTP tests for priority-bearing create and partial-update requests, returned priority, invalid priority rejection, and unchanged title/status in `backend/internal/httpapi/handler_test.go` (FR-001, FR-002, FR-004, FR-009).
- [X] T005 [P] [US1] Add UI tests for the default create selection, priority-bearing create/edit requests, and visible priority values in `frontend/src/app.test.ts` (FR-002, FR-003, FR-004, FR-005).

### Implementation for User Story 1

- [X] T006 [US1] Implement priority storage and partial todo updates in `backend/internal/todo/repository.go`; create and update operations must retain the todo's ID, title, status, and creation order while changing only supplied fields (FR-001, FR-004, FR-009).
- [X] T007 [US1] Implement service-level priority validation and create/update behavior in `backend/internal/todo/service.go`; accept only the data-model values `low`, `medium`, and `high` and preserve omitted update fields (FR-001, FR-002, FR-004, FR-009).
- [X] T008 [US1] Extend request decoding and response handling in `backend/internal/httpapi/handler.go` so `POST /api/todos` accepts an optional priority and `PATCH /api/todos/{id}` accepts priority, status, or both while preserving status-only requests (FR-002, FR-004, FR-010, FR-012).
- [X] T009 [P] [US1] Update `specs/002-add-task-priorities/contracts/openapi.yaml` to document the priority values, optional create/update fields, effective response priority, and status-only update compatibility (FR-001–FR-005, FR-012).
- [X] T010 [P] [US1] Add an accessible priority selector defaulting to `medium` and a per-todo priority edit control in `frontend/index.html`, with responsive control and priority-badge styling in `frontend/styles.css` (FR-002–FR-005).
- [X] T011 [US1] Implement priority-aware create/edit requests, response validation, and todo rendering in `frontend/src/app.ts`; keep the current title, status, delete, feedback, and keyboard behaviors intact (FR-002–FR-005, FR-009, FR-010).

**Checkpoint**: A user can create and edit priorities end-to-end and see the updated value without changing other todo fields.

---

## Phase 4: User Story 3 - Vorhandene Aufgaben weiterhin verwenden (Priority: P1)

**Goal**: Todos without an explicit priority and older clients remain usable; missing legacy priorities behave as `medium`, while an older status-only update preserves an existing priority.

**Independent Test**: Provide a todo with no stored priority and verify it is listed as `medium`, can be found with the `medium` filter, and remains usable with old create and status-only update requests.

### Tests for User Story 3

- [X] T012 [P] [US3] Add repository/service tests for todos with an empty legacy priority, verifying the effective value is `medium` and existing todo fields remain intact in `backend/internal/todo/repository_test.go` and `backend/internal/todo/service_test.go` (FR-008, FR-012).
- [X] T013 [P] [US3] Add HTTP regression tests showing that create requests without priority still succeed with `medium` and status-only `PATCH` requests preserve an existing non-medium priority in `backend/internal/httpapi/handler_test.go` (FR-003, FR-008, FR-012).
- [X] T014 [P] [US3] Add a frontend compatibility test for a legacy todo response without a `priority` property, verifying that it remains renderable as `medium` in `frontend/src/app.test.ts` (FR-008).

### Implementation for User Story 3

- [X] T015 [US3] Normalize empty legacy priorities to effective `medium` in repository/service/API results and preserve missing update fields in `backend/internal/todo/repository.go`, `backend/internal/todo/service.go`, and `backend/internal/httpapi/handler.go` (FR-003, FR-008, FR-012).
- [X] T016 [US3] Make `frontend/src/app.ts` accept legacy todo responses with no priority and present their effective priority as `medium` without blocking status changes or deletion (FR-008, FR-010).

**Checkpoint**: Legacy todos and old request shapes work without requiring users or older clients to supply a priority.

---

## Phase 5: User Story 2 - Aufgaben nach Priorität filtern (Priority: P2)

**Goal**: Users can filter by priority, combine that filter with the existing status filter, and return either filter to `all`.

**Independent Test**: Create todos across priorities and statuses; verify each filtered result matches both active criteria, clearing one criterion leaves the other in effect, and an empty match displays without error.

### Tests for User Story 2

- [X] T017 [P] [US2] Add service/repository tests for every priority filter, `all`, combined status-and-priority matching, empty results, and preserved creation order in `backend/internal/todo/service_test.go` and `backend/internal/todo/repository_test.go` (FR-006, FR-007, FR-011).
- [X] T018 [P] [US2] Add HTTP tests for `priority` query validation, combining it with `status`, repeated/invalid filter rejection, and empty results in `backend/internal/httpapi/handler_test.go` (FR-006, FR-007, FR-011).
- [X] T019 [P] [US2] Add UI tests for independently resetting status/priority filters, combined-filter query requests, and empty-result messaging in `frontend/src/app.test.ts` (FR-006, FR-007, FR-011).

### Implementation for User Story 2

- [X] T020 [US2] Add a list-filter value for status and priority and apply both criteria together in `backend/internal/todo/model.go`, `backend/internal/todo/service.go`, and `backend/internal/todo/repository.go`; omitted or `all` criteria must not restrict results (FR-006, FR-007, FR-011).
- [X] T021 [US2] Parse and validate the optional `status` and `priority` query parameters in `backend/internal/httpapi/handler.go`; return a bad-request error for invalid or repeated values and pass both filters to the service (FR-006, FR-007, FR-011).
- [X] T022 [US2] Add a separately labeled priority filter beside the status filter in `frontend/index.html` and update `frontend/src/app.ts` and `frontend/styles.css` to submit both active filters, reset either to `all`, and explain empty results (FR-006, FR-007, FR-011).

**Checkpoint**: Priority and status filtering work together in the UI and API, including empty and reset states.

---

## Phase 6: Polish & Cross-Cutting Validation

**Purpose**: Verify the integrated feature against the existing accessibility flow and documented validation guide.

- [X] T023 Extend `frontend/e2e/accessibility.test.js` to cover keyboard access to priority controls, priority updates, and combined filters while retaining existing status and deletion checks (FR-004–FR-007, FR-010–FR-011).
- [X] T024 Run every validation command in `specs/002-add-task-priorities/quickstart.md` and resolve any failures in the relevant source/test files before marking the feature complete.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies; verifies the existing Go and frontend test tools.
- **Foundational (Phase 2)**: Depends on Setup; blocks all user-story implementation.
- **User Story 1 (Phase 3)**: Depends on the foundational priority domain value.
- **User Story 3 (Phase 4)**: Depends on User Story 1's priority-aware create/update paths; verifies and completes legacy compatibility.
- **User Story 2 (Phase 5)**: Depends on priority-aware todo data and legacy normalization from User Stories 1 and 3.
- **Polish (Phase 6)**: Depends on all three user stories.

### User Story Dependencies

- **US1 (P1)**: Can start after Phase 2; no dependency on another story.
- **US3 (P1)**: Starts after US1 so compatibility tests can verify priority-bearing creates and updates.
- **US2 (P2)**: Starts after US1 and US3 so filters cover both explicitly prioritized and legacy todos.

### Within Each User Story

- Add the story's tests before its implementation and use them to verify behavior.
- Keep domain and repository behavior below the service and HTTP layers.
- Keep HTTP and frontend behavior aligned with `contracts/openapi.yaml`.
- Complete the independent test criteria at each checkpoint before continuing.

### Parallel Opportunities

- After T002, T003, T004, and T005 can be prepared in parallel because they touch service, handler, and frontend test files respectively.
- Backend repository/service work (T006–T008) can proceed independently from the HTML/CSS controls (T010); the API contract update (T009) can be reviewed alongside either.
- For US3, T012, T013, and T014 use separate repository/service, HTTP, and frontend test files. After those tests, backend normalization (T015) and frontend legacy handling (T016) touch separate layers.
- For US2, T017, T018, and T019 are independent test-file changes. After tests, frontend filter work (T022) can proceed in parallel with backend filter work (T020–T021), with end-to-end validation after both are integrated.

## Parallel Examples

### User Story 1

```text
Prepare the service, HTTP, and frontend tests in parallel:
- T003: backend/internal/todo/service_test.go
- T004: backend/internal/httpapi/handler_test.go
- T005: frontend/src/app.test.ts

After the tests are prepared, backend handlers and frontend markup/styles can proceed separately:
- T006–T008: backend/internal/todo/ and backend/internal/httpapi/
- T010: frontend/index.html and frontend/styles.css
```

### User Story 3

```text
Prepare independent compatibility tests in parallel:
- T012: backend/internal/todo/repository_test.go and backend/internal/todo/service_test.go
- T013: backend/internal/httpapi/handler_test.go
- T014: frontend/src/app.test.ts
```

### User Story 2

```text
Prepare filter tests in parallel:
- T017: backend/internal/todo/service_test.go and backend/internal/todo/repository_test.go
- T018: backend/internal/httpapi/handler_test.go
- T019: frontend/src/app.test.ts

Implement backend and frontend filtering against the same contract:
- T020–T021: backend/internal/todo/ and backend/internal/httpapi/
- T022: frontend/index.html, frontend/src/app.ts, and frontend/styles.css
```

## Implementation Strategy

### MVP First

1. Complete setup and the foundational priority value.
2. Complete User Story 1 to create, edit, and display each priority.
3. Validate its independent test criteria; this is the smallest useful priority-management increment.

### Incremental Delivery

1. Add User Story 3 compatibility and verify existing tasks and older request shapes.
2. Add User Story 2 combined filtering and validate all empty/reset states.
3. Run browser accessibility and all quickstart validation before completion.
