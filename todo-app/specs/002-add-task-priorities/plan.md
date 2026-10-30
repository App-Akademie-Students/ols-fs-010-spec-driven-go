# Implementation Plan: To-Do-Prioritäten festlegen und filtern

**Branch**: `002-add-task-priorities` | **Date**: 2026-10-30 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/002-add-task-priorities/spec.md`. Integrate the feature into the existing architecture and use its current components and test tools.

## Summary

Add a `low` / `medium` / `high` priority to todos through the existing Go domain, service, repository, and HTTP layers, then expose creation, editing, display, and combined status/priority filtering in the existing TypeScript UI. Keep old HTTP clients compatible: omitted priority on create defaults to `medium`, while an omitted priority on update leaves the current value unchanged. Continue treating legacy todos without a stored priority as `medium`. Preserve the existing in-memory store, REST surface, accessible HTML controls, and established Go/Vitest/Playwright tests.

## Technical Context

**Language/Version**: Go 1.25; TypeScript 5.9 (ES modules, DOM)

**Primary Dependencies**: Go standard library; existing TypeScript, Vitest 5, jsdom 26, and Playwright 1 development dependencies. No new dependencies.

**Storage**: Existing in-memory `todo.InMemoryRepository`; no persistent storage or migration.

**Testing**: `go test ./...` and `go test -race ./...` from `backend`; `npm test -- --run` and `npm run build` from `frontend`; `npm run test:browser` for Playwright accessibility and end-to-end coverage.

**Target Platform**: Existing local Go HTTP application and browser-based web UI.

**Project Type**: Go REST web service and TypeScript frontend.

**Performance Goals**: No new latency or throughput target; retain existing synchronous local-app behavior.

**Constraints**: Keep HTTP behavior backward compatible for requests that omit priority; combine status and priority filters; retain creation-order results; do not add dependencies or persistence; follow the existing service/repository and accessible-control patterns.

**Scale/Scope**: Single in-memory application instance and its existing todo list; one additive todo attribute, one additive list filter, and priority controls in existing create/list flows.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Specification-first**: Pass. The feature has an approved specification and clarified defaults, filter composition, and legacy request behavior.
- **Clear, minimal architecture**: Pass. Reuse the current domain/service/repository/HTTP boundaries and frontend; no new dependencies, storage layer, or abstraction is needed.
- **Automated behavioral tests**: Pass. Domain rules, HTTP contracts, UI behavior, and end-to-end accessibility have existing test suites to extend.
- **Explicit and safe interfaces**: Pass. Validate priority and filters, preserve consistent errors, and keep status-only updates valid.
- **Simplicity and verifiability**: Pass. Keep in-memory storage and current commands; verify with Go tests including race detection, TypeScript build/Vitest, and browser tests.

**Gate result before research**: PASS. No constitution violations or unresolved requirement questions.

## Project Structure

### Documentation (this feature)

```text
specs/002-add-task-priorities/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── openapi.yaml
└── tasks.md              # Produced by /speckit-tasks
```

### Source Code (repository root)

```text
backend/
├── internal/
│   ├── todo/
│   │   ├── model.go                 # Priority type, validation errors, Todo field
│   │   ├── service.go               # Defaults and update validation
│   │   ├── repository.go            # Priority-aware updates and combined filtering
│   │   ├── service_test.go
│   │   ├── repository_test.go
│   └── httpapi/
│       ├── handler.go               # Optional request fields and query filters
│       └── handler_test.go
└── cmd/todo/main.go
frontend/
├── index.html                       # Creation priority and separate filter controls
├── styles.css                       # Responsive priority controls/badges
├── src/
│   ├── app.ts                       # Priority-aware requests, rendering, and filtering
│   └── app.test.ts                  # Vitest/jsdom behavior coverage
└── e2e/accessibility.test.js         # Playwright keyboard/UI regression coverage
```

**Structure Decision**: Keep the existing two-part web application. Backend behavior stays in `backend/internal/todo` and `backend/internal/httpapi`; user interactions stay in the current static HTML/CSS and `frontend/src/app.ts`. Extend existing tests in place and use the feature-scoped OpenAPI contract to describe the HTTP interface.

## Implementation Design

1. Add a validated priority value to the todo domain model. The service defaults omitted create priorities to `medium`; an update may change status, priority, or both, and omitted update fields are left unchanged.
2. Extend the in-memory repository update/list operations to preserve priority and apply status plus priority criteria together. Normalize any legacy todo with no stored priority to effective `medium` without changing its ID, title, status, or position.
3. Extend the existing HTTP handlers and OpenAPI contract additively. Keep `POST /api/todos` valid without a priority and keep status-only `PATCH /api/todos/{id}` valid. Allow priority-only or combined updates, reject updates with no mutable fields, validate values, and accept optional `status` and `priority` query filters whose results satisfy both.
4. Add a default-`medium` selector to the existing create form, a labeled priority filter beside the status filter, and an accessible per-todo priority edit control. Display each todo's priority; retain existing status, deletion, feedback, and keyboard behaviors.
5. Add service/repository unit tests, HTTP contract tests, Vitest UI tests, and Playwright coverage for defaulting, valid/invalid priorities, edit omission semantics, legacy values, combined filters, and existing actions.

## Constitution Check (Post-Design)

- **Specification-first**: PASS — design maps to FR-001–FR-012 and SC-001–SC-004.
- **Clear, minimal architecture**: PASS — no boundary changes or new dependencies; each layer owns its existing responsibility.
- **Automated behavioral tests**: PASS — each changed domain, API, and UI behavior has an existing narrow test suite.
- **Explicit and safe interfaces**: PASS — omitted versus invalid priority values have distinct documented semantics; invalid values and empty updates fail explicitly.
- **Simplicity and verifiability**: PASS — only existing components and tools are used; plan verification uses the project toolchain.

**Gate result after design**: PASS. No unjustified constitution violations.
