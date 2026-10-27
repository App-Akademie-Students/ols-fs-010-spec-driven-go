# Research: Todo-Aufgaben verwalten

## Go HTTP routing and API behavior

- **Decision**: Use Go 1.22+ method-aware `http.ServeMux` patterns, with the selected
  project baseline of Go 1.25. Register collection and item routes by HTTP method and keep
  handlers as adapters over the service.
- **Rationale**: The standard library provides method and path-variable routing without
  adding a router dependency. It fits the requested `net/http` implementation and keeps
  transport concerns separate from todo rules.
- **Alternatives considered**: A third-party router or a hand-written dispatcher. Neither is
  needed for this small route set.

## In-memory repository and concurrent access

- **Decision**: Store records in a map guarded by `sync.RWMutex`. Guard every map read and
  mutation, return values or copies rather than exposing shared mutable records, and keep ID
  allocation within the write-locked create operation.
- **Rationale**: The data structure has clear read/write semantics for a small personal list.
  Correct lock coverage prevents concurrent map access races and protects state transitions.
- **Alternatives considered**: `sync.Map` and an unguarded map. `sync.Map` is optimized for
  specialized access patterns and obscures this repository's simple operations; an unguarded
  map is unsafe when HTTP requests run concurrently.

## Request decoding and errors

- **Decision**: Decode JSON from the request body with `json.Decoder`, reject unknown fields,
  limit body size, require exactly one JSON object, validate values in the service, and return
  a consistent JSON error object. Use `400` for malformed or invalid input, `404` for an
  unknown todo ID, `201` for creation, `200` for reads/updates, and `204` for deletion.
- **Rationale**: Streaming decoding with a size limit avoids unbounded request bodies.
  Strict input and explicit errors prevent accidental acceptance of malformed requests or
  success-shaped failures.
- **Alternatives considered**: Reading the entire body without a limit or silently ignoring
  unknown fields; both make the API less predictable.

## Frontend safety and TypeScript

- **Decision**: Keep the UI framework-free. Compile TypeScript to browser JavaScript and
  construct todo rows with DOM APIs and `textContent`; do not interpolate task text into
  `innerHTML`.
- **Rationale**: This meets the requested HTML/CSS/TypeScript stack while avoiding an
  unnecessary runtime framework. Text-node rendering safely displays arbitrary user task
  titles.
- **Alternatives considered**: Directly inserting task titles as HTML or adding a UI framework.
  The first risks markup injection; the second adds dependencies and architecture not needed
  by the feature.

## Concurrency verification

- **Decision**: Include repository tests that exercise concurrent create, list, status change
  and delete operations, and run `go test -race ./...`.
- **Rationale**: Functional tests establish repository behavior; the Go race detector checks
  for unsynchronized memory access during concurrent workloads.
- **Alternatives considered**: Sequential tests only, which do not exercise the explicit
  concurrent-access requirement.

## Sources

- Go `net/http` package and `ServeMux`: https://pkg.go.dev/net/http
- Go `encoding/json` package: https://pkg.go.dev/encoding/json
- Go `sync.RWMutex`: https://pkg.go.dev/sync#RWMutex
- Go race detector: https://go.dev/doc/articles/race_detector
- MDN `Node.textContent`: https://developer.mozilla.org/docs/Web/API/Node/textContent

## Unresolved Decisions

None. Product behavior is specified and clarified; technical choices follow the user's
explicit stack constraints and the project constitution.
