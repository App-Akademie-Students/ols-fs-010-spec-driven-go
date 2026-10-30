type TodoStatus = "open" | "completed";
type StatusFilter = "all" | TodoStatus;
type TodoPriority = "low" | "medium" | "high";
type PriorityFilter = "all" | TodoPriority;

interface Todo {
  id: number;
  title: string;
  status: TodoStatus;
  priority: TodoPriority;
}

interface TodoResponse extends Omit<Todo, "priority"> {
  priority?: TodoPriority;
}

const genericFailure = "Aktion fehlgeschlagen. Bitte erneut versuchen.";
const unavailableTodo = "Aufgabe ist nicht mehr verfügbar.";

export function mountTodoApp(
  documentRef: Document = document,
  fetcher: typeof fetch = fetch,
): void {
  const form = requiredElement<HTMLFormElement>(documentRef, "#todo-form");
  const titleInput = requiredElement<HTMLInputElement>(documentRef, "#todo-title");
  const createPrioritySelect = requiredElement<HTMLSelectElement>(documentRef, "#todo-priority");
  const filterSelect = requiredElement<HTMLSelectElement>(documentRef, "#todo-filter");
  const priorityFilterSelect = requiredElement<HTMLSelectElement>(documentRef, "#todo-priority-filter");
  const list = requiredElement<HTMLUListElement>(documentRef, "#todo-list");
  const emptyState = requiredElement<HTMLElement>(documentRef, "#empty-state");
  const feedback = requiredElement<HTMLElement>(documentRef, "#feedback");
  const errorFeedback = requiredElement<HTMLElement>(documentRef, "#error-feedback");
  let activeFilter: StatusFilter = "all";
  let activePriorityFilter: PriorityFilter = "all";
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
  priorityFilterSelect.addEventListener("change", () => {
    const filter = priorityFilterSelect.value;
    if (!isPriorityFilter(filter)) {
      showError(genericFailure);
      return;
    }
    activePriorityFilter = filter;
    void loadTodos();
  });

  async function createTodo(): Promise<void> {
    const title = titleInput.value.replace(/^\p{White_Space}+|\p{White_Space}+$/gu, "");
    if (title === "") {
      showError("Titel ist erforderlich. Keine Aufgabe erstellt.");
      return;
    }
    if (!isTodoPriority(createPrioritySelect.value)) {
      showError("Ungültige Priorität.");
      return;
    }

    try {
      await request("/api/todos", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, priority: createPrioritySelect.value }),
      });
      titleInput.value = "";
      createPrioritySelect.value = "medium";
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

  async function changePriority(todo: Todo, priority: TodoPriority): Promise<void> {
    try {
      await request(`/api/todos/${todo.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ priority }),
      });
      showSuccess("Priorität geändert.");
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
    const parameters = new URLSearchParams();
    if (activeFilter !== "all") {
      parameters.set("status", activeFilter);
    }
    if (activePriorityFilter !== "all") {
      parameters.set("priority", activePriorityFilter);
    }
    const query = parameters.size === 0 ? "" : `?${parameters.toString()}`;
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
        activeFilter === "all" && activePriorityFilter === "all"
          ? "Noch keine Aufgaben vorhanden. Erstelle deine erste Aufgabe."
          : emptyFilterMessage(activeFilter, activePriorityFilter);
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

      const priority = documentRef.createElement("span");
      priority.className = `todo-priority todo-priority-${todo.priority}`;
      priority.textContent = todo.priority;

      const priorityEdit = documentRef.createElement("div");
      priorityEdit.className = "todo-priority-edit";
      const prioritySelect = documentRef.createElement("select");
      prioritySelect.setAttribute("aria-label", `Priorität für Aufgabe ändern: ${todo.title}`);
      for (const value of ["low", "medium", "high"] as const) {
        const option = documentRef.createElement("option");
        option.value = value;
        option.textContent = value;
        option.selected = todo.priority === value;
        prioritySelect.append(option);
      }
      const priorityButton = documentRef.createElement("button");
      priorityButton.type = "button";
      priorityButton.textContent = "Priorität speichern";
      priorityButton.addEventListener("click", () => {
        if (isTodoPriority(prioritySelect.value)) {
          void changePriority(todo, prioritySelect.value);
        } else {
          showError("Ungültige Priorität.");
        }
      });
      priorityEdit.append(prioritySelect, priorityButton);

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
      row.append(title, status, priority, priorityEdit, actions);
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

function emptyFilterMessage(status: StatusFilter, priority: PriorityFilter): string {
  const filters: string[] = [];
  if (status !== "all") {
    filters.push(`Status „${statusLabel(status)}“`);
  }
  if (priority !== "all") {
    filters.push(`Priorität „${priority}“`);
  }
  return `Keine passenden Aufgaben für ${filters.join(" und ")}. Wähle „Alle“, um alle Aufgaben anzuzeigen.`;
}

function isTodoPriority(value: unknown): value is TodoPriority {
  return value === "low" || value === "medium" || value === "high";
}

function isPriorityFilter(value: string): value is PriorityFilter {
  return value === "all" || isTodoPriority(value);
}

function isTodo(value: unknown): value is TodoResponse {
  if (typeof value !== "object" || value === null) {
    return false;
  }
  const candidate = value as Record<string, unknown>;
  return (
    typeof candidate.id === "number" &&
    Number.isSafeInteger(candidate.id) &&
    candidate.id > 0 &&
    typeof candidate.title === "string" &&
    (candidate.status === "open" || candidate.status === "completed") &&
    (candidate.priority === undefined || isTodoPriority(candidate.priority))
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
  return todos.map((todo) => ({ ...todo, priority: todo.priority ?? "medium" }));
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
