type TodoStatus = "open" | "completed";
type StatusFilter = "all" | TodoStatus;

interface Todo {
  id: number;
  title: string;
  status: TodoStatus;
}

const genericFailure = "Aktion fehlgeschlagen. Bitte erneut versuchen.";
const unavailableTodo = "Aufgabe ist nicht mehr verfügbar.";

export function mountTodoApp(
  documentRef: Document = document,
  fetcher: typeof fetch = fetch,
): void {
  const form = requiredElement<HTMLFormElement>(documentRef, "#todo-form");
  const titleInput = requiredElement<HTMLInputElement>(documentRef, "#todo-title");
  const filterSelect = requiredElement<HTMLSelectElement>(documentRef, "#todo-filter");
  const list = requiredElement<HTMLUListElement>(documentRef, "#todo-list");
  const emptyState = requiredElement<HTMLElement>(documentRef, "#empty-state");
  const feedback = requiredElement<HTMLElement>(documentRef, "#feedback");
  const errorFeedback = requiredElement<HTMLElement>(documentRef, "#error-feedback");
  let activeFilter: StatusFilter = "all";
  let requestSequence = 0;

  form.addEventListener("submit", (event) => {
    event.preventDefault();
    void createTodo();
  });
  filterSelect.addEventListener("change", () => {
    const filter = filterSelect.value;
    if (filter !== "all" && filter !== "open" && filter !== "completed") {
      showError(genericFailure);
      return;
    }
    activeFilter = filter;
    void loadTodos();
  });

  async function createTodo(): Promise<void> {
    const title = titleInput.value.replace(/^\p{White_Space}+|\p{White_Space}+$/gu, "");
    if (title === "") {
      showError("Titel ist erforderlich. Keine Aufgabe erstellt.");
      return;
    }

    try {
      await request("/api/todos", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title }),
      });
      titleInput.value = "";
      showSuccess("Aufgabe erstellt.");
      await loadTodos();
    } catch (error) {
      showError(errorMessage(error));
    }
  }

  async function changeStatus(todo: Todo): Promise<void> {
    const status: TodoStatus = todo.status === "open" ? "completed" : "open";
    try {
      await request(`/api/todos/${todo.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ status }),
      });
      showSuccess("Status geändert.");
      await loadTodos();
    } catch (error) {
      showError(errorMessage(error));
    }
  }

  async function deleteTodo(todo: Todo): Promise<void> {
    try {
      await request(`/api/todos/${todo.id}`, { method: "DELETE" });
      showSuccess("Aufgabe gelöscht.");
      await loadTodos();
    } catch (error) {
      showError(errorMessage(error));
    }
  }

  async function loadTodos(): Promise<void> {
    const sequence = ++requestSequence;
    const query = activeFilter === "all" ? "" : `?status=${activeFilter}`;
    try {
      const payload = await request(`/api/todos${query}`);
      if (sequence !== requestSequence) {
        return;
      }
      const todos = parseTodoList(payload);
      renderTodos(todos);
    } catch (error) {
      if (sequence === requestSequence) {
        showError(errorMessage(error));
      }
    }
  }

  function renderTodos(todos: Todo[]): void {
    list.replaceChildren();
    emptyState.hidden = todos.length > 0;
    if (todos.length === 0) {
      emptyState.textContent =
        activeFilter === "all"
          ? "Noch keine Aufgaben vorhanden. Erstelle deine erste Aufgabe."
          : `Keine passenden Aufgaben im Status „${statusLabel(activeFilter)}“. Wähle „Alle“, um alle Aufgaben anzuzeigen.`;
      return;
    }

    for (const todo of todos) {
      const row = documentRef.createElement("li");
      row.className = "todo-item";
      row.dataset.todoId = String(todo.id);

      const title = documentRef.createElement("span");
      title.className = "todo-title";
      title.textContent = todo.title;

      const status = documentRef.createElement("span");
      status.className = `todo-status todo-status-${todo.status}`;
      status.textContent = statusLabel(todo.status);

      const actions = documentRef.createElement("div");
      actions.className = "todo-actions";

      const statusButton = documentRef.createElement("button");
      statusButton.type = "button";
      statusButton.textContent = todo.status === "open" ? "Erledigt" : "Wieder öffnen";
      statusButton.setAttribute(
        "aria-label",
        todo.status === "open" ? "Als erledigt markieren" : "Aufgabe wieder öffnen",
      );
      statusButton.addEventListener("click", () => void changeStatus(todo));

      const deleteButton = documentRef.createElement("button");
      deleteButton.type = "button";
      deleteButton.className = "button-danger";
      deleteButton.textContent = "Löschen";
      deleteButton.setAttribute("aria-label", `Aufgabe löschen: ${todo.title}`);
      deleteButton.addEventListener("click", () => void deleteTodo(todo));

      actions.append(statusButton, deleteButton);
      row.append(title, status, actions);
      list.append(row);
    }
  }

  async function request(path: string, init?: RequestInit): Promise<unknown> {
    let response: Response;
    try {
      response = await fetcher(path, init);
    } catch {
      throw new Error(genericFailure);
    }
    if (!response.ok) {
      throw new Error(await responseError(response));
    }
    if (response.status === 204) {
      return undefined;
    }
    try {
      return await response.json();
    } catch {
      throw new Error(genericFailure);
    }
  }

  function showSuccess(message: string): void {
    errorFeedback.textContent = "";
    feedback.textContent = message;
  }

  function showError(message: string): void {
    feedback.textContent = "";
    errorFeedback.textContent = message;
  }

  function errorMessage(error: unknown): string {
    return error instanceof Error && error.message !== "" ? error.message : genericFailure;
  }

  void loadTodos();
}

function requiredElement<T extends Element>(documentRef: Document, selector: string): T {
  const element = documentRef.querySelector<T>(selector);
  if (!element) {
    throw new Error(`Required todo app element not found: ${selector}`);
  }
  return element;
}

function statusLabel(status: TodoStatus): string {
  return status === "open" ? "Offen" : "Erledigt";
}

function isTodo(value: unknown): value is Todo {
  if (typeof value !== "object" || value === null) {
    return false;
  }
  const candidate = value as Record<string, unknown>;
  return (
    typeof candidate.id === "number" &&
    Number.isSafeInteger(candidate.id) &&
    candidate.id > 0 &&
    typeof candidate.title === "string" &&
    (candidate.status === "open" || candidate.status === "completed")
  );
}

function parseTodoList(value: unknown): Todo[] {
  if (typeof value !== "object" || value === null || !("todos" in value)) {
    throw new Error(genericFailure);
  }
  const todos = (value as { todos: unknown }).todos;
  if (!Array.isArray(todos) || !todos.every(isTodo)) {
    throw new Error(genericFailure);
  }
  return todos;
}

async function responseError(response: Response): Promise<string> {
  if (response.status === 404) {
    return unavailableTodo;
  }
  if (response.status >= 500) {
    return genericFailure;
  }
  try {
    const payload: unknown = await response.json();
    if (typeof payload === "object" && payload !== null && "error" in payload) {
      const message = (payload as { error?: unknown }).error;
      if (typeof message === "string" && message.trim() !== "") {
        return message;
      }
    }
  } catch {
    return genericFailure;
  }
  return genericFailure;
}
