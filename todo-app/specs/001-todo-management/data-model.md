# Data Model: Todo-Aufgaben verwalten

## Aufgabe

Represents one item in the user's personal todo list.

| Field | Type | Rules |
|---|---|---|
| `id` | Positive integer | Assigned by the repository; unique for the lifetime of the running process; immutable. |
| `title` | String | Required; leading and trailing characters in the Unicode `White_Space` property are removed on creation and the normalized value is stored; a blank result is rejected. No product-level character or length restriction applies. |
| `status` | Enum: `open`, `completed` | New tasks start as `open`; either status can transition to the other. |

Tasks are independent records. There are no users, parent-child relationships, categories,
due dates, or cross-process identity guarantees in scope.
Duplicate titles are allowed; task identity is determined by `id`. Tasks are listed in
creation order, and filtered views retain that relative order.

## State transitions

```text
open <-> completed
```

Deletion immediately and permanently removes the record from the repository without a
confirmation step. A deleted ID is not a valid update target. The in-memory repository does
not persist tasks; they are lost when the application process exits or restarts.

## Status filter

The status filter is a request selection rather than stored task data:

| Value | Included tasks |
|---|---|
| `all` | Both `open` and `completed` tasks; used when the filter is omitted. |
| `open` | Tasks whose status is `open`. |
| `completed` | Tasks whose status is `completed`. |

Any other filter value is invalid and returns a client error.

## Validation and consistency

- Create removes leading and trailing Unicode `White_Space` characters from the title, stores
  the normalized value, and rejects an empty result.
- Status updates accept only `open` or `completed`.
- Repository operations are safe for concurrent callers.
- Each operation returns a detached value/snapshot so callers cannot mutate repository state
  outside the repository's synchronization boundary.
