# Data Model: To-Do-Prioritäten

## Todo

The existing todo gains an effective `priority` attribute. A todo is returned in creation order as before.

| Field | Type / values | Required in responses | Rule |
|---|---|---:|---|
| `id` | Positive integer | Yes | Existing stable todo identifier. |
| `title` | Non-empty string | Yes | Existing title normalization and validation remain unchanged. |
| `status` | `open` or `completed` | Yes | Existing status lifecycle remains unchanged. |
| `priority` | `low`, `medium`, or `high` | Yes | New todos default to `medium`; a missing legacy value has effective value `medium`. |

### Validation and update rules

- Only `low`, `medium`, and `high` are valid priority values.
- Create requests may omit `priority`; the created todo receives `medium`.
- Update requests may supply `status`, `priority`, or both. At least one mutable field is required.
- An omitted update field is unchanged. In particular, legacy status-only updates do not reset priority.
- A legacy todo without a stored priority is shown and filtered as `medium`; its other attributes and list position remain unchanged.
- Changing priority does not change title or status. Changing status does not change priority.
- Existing delete and status transitions remain available regardless of priority.

## List filters

The list has two independent optional criteria:

| Filter | Values | Default |
|---|---|---|
| `status` | `all`, `open`, `completed` | `all` |
| `priority` | `all`, `low`, `medium`, `high` | `all` |

When both filters are specific, a todo is included only when it matches both. An empty match returns an empty list, not an error. Applying filters preserves the existing creation order.
