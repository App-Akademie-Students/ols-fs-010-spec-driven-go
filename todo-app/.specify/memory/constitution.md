<!--
Sync Impact Report
Version change: unratified scaffold -> 1.0.0 (initial ratification)
Modified principles: None; five principles established from the scaffold.
Added sections: Technical Constraints and Quality Standards; Development Workflow.
Removed sections: None.
Follow-up TODOs: None.
-->

# Todo App Constitution

## Core Principles

### I. Specification-First Development
Every user-visible behavior MUST be described in an approved feature specification before
implementation begins. Plans and tasks MUST trace back to those requirements, and implementation
reviews MUST verify that observable behavior matches the specification. This keeps work focused
and prevents undocumented product behavior.

### II. Clear, Minimal Architecture
The Go application MUST keep HTTP transport, todo business rules, and persistence responsibilities
separate. Handlers MUST delegate business decisions to the service layer, and persistence MUST be
accessed through the repository boundary. New abstractions and dependencies MUST have a documented
need; choose the simplest design that satisfies specified behavior.

### III. Automated Behavioral Tests
Each new or changed behavior MUST have automated tests at the narrowest useful level. Service
rules MUST be covered by unit tests, and changes to HTTP contracts or persistence integration
MUST be covered by API or integration tests. The relevant Go test suite MUST pass before a change
is considered complete.

### IV. Explicit and Safe Interfaces
The REST API MUST validate untrusted input and return documented, consistent status codes and
error responses. The web interface MUST communicate action results and errors clearly, and
interactive controls MUST be operable using standard keyboard interaction. Invalid input and
operational errors MUST NOT be reported as successful operations.

### V. Simplicity and Verifiability
Features MUST remain no more complex than their approved requirements demand. Changes MUST
preserve unrelated existing behavior, and verification MUST use the project’s documented build
and test commands. Concurrency-sensitive Go changes MUST be checked with the race detector when
the toolchain and environment support it; any unavailable check MUST be reported.

## Technical Constraints and Quality Standards

The backend MUST use Go and expose its specified application behavior through a REST API. The
user interface MUST use the project’s simple web stack (HTML and CSS with TypeScript or
JavaScript). Persistence is in-memory for the current project scope; durable storage requires an
approved specification and plan. External dependencies MUST be justified by a concrete
requirement, and API changes MUST keep the user interface and backend contract consistent.

## Development Workflow

Work MUST follow the Spec Kit sequence: specify, plan, generate dependency-ordered tasks,
implement, and converge by checking the result against the specification. Clarification,
checklist, and analysis steps MUST be used when the feature’s requirements or artifact
consistency need them. Each implementation task MUST identify its applicable requirements and
verification. A task is complete only after its relevant tests pass and any discrepancy with
the specification is resolved or explicitly recorded for follow-up.

## Governance

This constitution governs project decisions and takes precedence over conflicting local
conventions. Amendments MUST document their rationale and impact, receive project-owner or
reviewer approval, and update the version and amendment date in the same change.

Versioning follows semantic versioning: MAJOR for incompatible changes to principles or
governance, MINOR for new principles or materially expanded requirements, and PATCH for
clarifications or non-semantic wording changes. Every specification, plan, task set, and
implementation review MUST be checked for compliance. Deviations MUST be justified and recorded;
work that violates a principle MUST be corrected or approved through a constitution amendment.

**Version**: 1.0.0 | **Ratified**: 2026-10-27 | **Last Amended**: 2026-10-27
