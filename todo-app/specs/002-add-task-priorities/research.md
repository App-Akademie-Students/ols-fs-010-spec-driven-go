# Research: To-Do-Prioritäten

Research was limited to the existing repository architecture and the clarified feature specification. The active application already separates Go domain/service/repository/HTTP responsibilities and uses static HTML/CSS with a TypeScript DOM controller. No technology or external-service unknown remains.

## Decisions

### Store priority as a validated todo domain value

- **Decision**: Add `low`, `medium`, and `high` to the existing todo model and set `medium` as the effective priority for new or legacy todos with no value.
- **Rationale**: Priority is intrinsic todo state, so it belongs alongside title and status. A single domain default keeps API responses, list filters, and UI rendering consistent.
- **Alternatives considered**: Store priority only in frontend state (not durable even within the server process); allow arbitrary strings (breaks filtering and validation); require a data migration (there is no persistent store).

### Extend the existing partial-update HTTP route

- **Decision**: Keep `PATCH /api/todos/{id}`; accept optional `status` and `priority`, with at least one field required. Omitted update fields remain unchanged. Preserve requests containing only `status`.
- **Rationale**: An additive partial update supports priority-only edits and existing clients without introducing another endpoint. The clarified compatibility rule avoids accidentally resetting a saved priority on an older status-only request.
- **Alternatives considered**: Add a separate priority route (duplicates update behavior); require clients to resend both fields (risks overwriting state and breaks old clients).

### Combine status and priority filters

- **Decision**: Extend `GET /api/todos` with an optional `priority` parameter and apply it in addition to the existing `status` parameter. An omitted or `all` criterion imposes no restriction.
- **Rationale**: This implements the clarified AND behavior without changing the existing status-only URL or default result set. Preserve creation order and existing empty-result behavior.
- **Alternatives considered**: Replace status filtering (breaks existing usage); allow only one active filter (cannot express the confirmed behavior).

### Keep UI and validation within the established stack

- **Decision**: Extend the existing DOM-based form/list, styles, and accessible native controls; verify with Go tests, Vitest/jsdom, TypeScript build, and the existing Playwright browser suite.
- **Rationale**: These are already the application’s interaction and verification paths, including keyboard navigation and live feedback.
- **Alternatives considered**: Add a component framework, a new test runner, or a dedicated service (unneeded for this bounded feature and contrary to minimal-architecture principles).

### No persistence migration

- **Decision**: Keep current in-memory persistence; map a missing legacy priority to effective `medium` when handling or exposing a todo.
- **Rationale**: The repository has no durable schema, and the requirements are satisfied by normalizing missing values at the existing boundary.
- **Alternatives considered**: Add a database or migration (out of scope and not supported by the current architecture).
