# Implementation Plan: Todo-Aufgaben verwalten

**Branch**: `001-todo-management` | **Date**: 2026-10-27 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/001-todo-management/spec.md`; technical direction
from the user: Go, standard-library `net/http`, an in-memory repository, separated HTTP,
business and persistence layers, concurrent access safety, an HTML/CSS/TypeScript browser UI,
and automated Go tests.

## Summary

Build a single-process todo web application. A Go service owns validation and todo state
transitions; a concurrency-safe in-memory repository stores todos; standard-library HTTP
handlers expose JSON creation, listing, status updates, deletion and filtering. The browser UI
is written in HTML, CSS and TypeScript, compiled to static JavaScript and served with the
application. Automated Go unit, repository-concurrency and HTTP contract tests verify the
specified flows; the race detector checks concurrent repository access.

## Technical Context

**Language/Version**: Go 1.25 or newer; the repository setup guide documents Go 1.25.3.

**Primary Dependencies**: Go standard library, including `net/http`, `encoding/json`,
`sync` and `embed`; TypeScript compiler as a frontend development dependency. No third-party
Go runtime dependency is planned.

**Storage**: Process-local in-memory repository backed by a map and guarded by
`sync.RWMutex`; records are lost when the process exits.

**Testing**: `go test ./...`, `go test -race ./...`, HTTP handler tests using `httptest`, and
`npm run build` for TypeScript type-checking and compilation.

**Target Platform**: Platforms supported by the selected Go toolchain; application is
accessed through a modern browser on the local network or machine.

**Project Type**: Single Go web application with a static browser frontend.

**Performance Goals**: No formal service-level target is specified. The design is for a
single-process personal todo list; performance targets and load testing are out of scope.

**Constraints**: Use `net/http`; keep transport, business logic and persistence separate;
make every repository operation safe under concurrent access; do not imply durable storage;
validate request bodies and provide consistent user-visible errors.

**Scale/Scope**: One personal list in one running process. No accounts, collaboration,
external services, persistent storage, import/export, or multi-instance coordination.

## Constitution Check

**Pre-research gate**

- Specification-first: PASS — the feature has an approved, clarified specification; each
  implementation behavior maps to its functional requirements and acceptance scenarios.
- Clear, minimal architecture: PASS — the requested HTTP, service and repository boundaries
  are explicit, with no additional runtime framework.
- Automated behavioral tests: PASS — service, concurrent repository and HTTP behavior will be
  covered by Go tests; the race detector is part of verification.
- Explicit and safe interfaces: PASS — the API contract defines validation and error behavior;
  browser-rendered task text must be treated as untrusted input.
- Simplicity and verifiability: PASS — the one-process in-memory scope avoids unrequested
  persistence and coordination.

**Post-design gate**

- Transport and static-asset serving remain in the HTTP boundary; business rules remain in
  the service; shared mutable state remains behind the repository lock: PASS.
- JSON and browser contracts expose only the specified task operations and status filters:
  PASS.
- Go test, race-test, TypeScript build and documented manual browser flow provide verification:
  PASS.
- No constitution violations or exceptions require complexity tracking.

## Project Structure

### Documentation (this feature)

```text
specs/001-todo-management/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── openapi.yaml
└── tasks.md              # Created by /speckit-tasks
```

### Source Code (repository root)

```text
cmd/todo/
└── main.go
internal/todo/
├── model.go
├── repository.go
├── service.go
├── repository_test.go
└── service_test.go
internal/httpapi/
├── handler.go
└── handler_test.go
web/
├── index.html
├── styles.css
├── src/
│   └── app.ts
├── dist/                 # TypeScript compiler output served by the application
├── package.json
└── tsconfig.json
```

**Structure Decision**: Use one Go module with an executable under `cmd/todo`. The
`internal/todo` package owns the domain model, service and repository interface plus
in-memory implementation. `internal/httpapi` adapts HTTP and JSON to the service. The
`web` directory contains the static page and TypeScript source; its compiled assets are served
with the application. Go tests live beside the packages they verify. The API contract is
maintained as OpenAPI under the feature documentation.

## Complexity Tracking

No constitution violations. No additional project or exception is introduced.
