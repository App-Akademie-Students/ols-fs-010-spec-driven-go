# Data Model: Todo-Aufgaben verwalten

## Aufgabe

Represents one item in the user's personal todo list.

| Field | Type | Rules |
|---|---|---|
| `id` | Positive integer | Assigned by the repository; unique for the lifetime of the running process; immutable. |
| `title` | String | Required; leading and trailing whitespace is removed on creation; a blank result is rejected. |
| `status` | Enum: `open`, `completed` | New tasks start as `open`; either status can transition to the other. |

Tasks are independent records. There are no users, parent-child relationships, categories,
due dates, or cross-process identity guarantees in scope.

## State transitions

```text
open <-> completed
```

Deletion removes the record from the repository. A deleted ID is not a valid update target.
The in-memory repository does not promise persistence across application restarts.

## Status filter

The status filter is a request selection rather than stored task data:

| Value | Included tasks |
|---|---|
| `all` | Both `open` and `completed` tasks; used when the filter is omitted. |
| `open` | Tasks whose status is `open`. |
| `completed` | Tasks whose status is `completed`. |

Any other filter value is invalid and returns a client error.

## Validation and consistency

- Create trims the title and rejects an empty result.
- Status updates accept only `open` or `completed`.
- Repository operations are safe for concurrent callers.
- Each operation returns a detached value/snapshot so callers cannot mutate repository state
  outside the repository's synchronization boundary.
