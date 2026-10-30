import { afterEach, describe, expect, it, vi } from "vitest";
import { mountTodoApp } from "./app";

function response(body: unknown, status = 200): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

function mockFetch(...responses: Response[]): typeof fetch {
  const fetcher = vi.fn<typeof fetch>();
  for (const result of responses) {
    fetcher.mockResolvedValueOnce(result);
  }
  return fetcher;
}

function mount(fetcher: typeof fetch): void {
  document.body.innerHTML = `
    <main>
      <form id="todo-form">
        <label for="todo-title">Aufgabentitel</label>
        <input id="todo-title" name="title">
        <label for="todo-priority">Priorität</label>
        <select id="todo-priority" name="priority">
          <option value="low">low</option>
          <option value="medium" selected>medium</option>
          <option value="high">high</option>
        </select>
        <button type="submit">Aufgabe hinzufügen</button>
      </form>
      <p id="feedback" role="status" aria-live="polite"></p>
      <p id="error-feedback" role="alert" aria-live="assertive"></p>
      <label for="todo-filter">Statusfilter</label>
      <select id="todo-filter">
        <option value="all">Alle</option>
        <option value="open">Offen</option>
        <option value="completed">Erledigt</option>
      </select>
      <label for="todo-priority-filter">Prioritätsfilter</label>
      <select id="todo-priority-filter">
        <option value="all">Alle Prioritäten</option>
        <option value="low">low</option>
        <option value="medium">medium</option>
        <option value="high">high</option>
      </select>
      <p id="empty-state" hidden></p>
      <ul id="todo-list"></ul>
    </main>`;
  mountTodoApp(document, fetcher);
}

function listResponse(todos: unknown[]): Response {
  return response({ todos });
}

const firstTodo = { id: 1, title: "First task", status: "open" };

afterEach(() => {
  document.body.innerHTML = "";
  vi.restoreAllMocks();
});

describe("todo app", () => {
  it("defaults new tasks to medium and allows priority changes", async () => {
    const created = { id: 2, title: "New task", status: "open", priority: "high" };
    const fetcher = mockFetch(
      listResponse([]),
      response(created, 201),
      listResponse([created]),
    );
    mount(fetcher);
    await vi.waitFor(() => expect(fetcher).toHaveBeenCalledTimes(1));

    const createPriority = document.querySelector<HTMLSelectElement>("#todo-priority")!;
    expect(createPriority.value).toBe("medium");
    createPriority.value = "high";
    document.querySelector<HTMLInputElement>("#todo-title")!.value = "New task";
    document.querySelector("#todo-form")!.dispatchEvent(
      new Event("submit", { bubbles: true, cancelable: true }),
    );

    await vi.waitFor(() => expect(document.querySelector("#todo-list")?.textContent).toContain("high"));
    expect(fetcher).toHaveBeenNthCalledWith(
      2,
      "/api/todos",
      expect.objectContaining({ body: JSON.stringify({ title: "New task", priority: "high" }) }),
    );
  });

  it("renders a legacy task without priority as medium", async () => {
    const fetcher = mockFetch(listResponse([firstTodo]));
    mount(fetcher);

    await vi.waitFor(() => {
      expect(document.querySelector("#todo-list")?.textContent).toContain(firstTodo.title);
      expect(document.querySelector(".todo-priority")?.textContent).toBe("medium");
      expect(document.querySelector<HTMLSelectElement>(".todo-priority-edit select")?.value).toBe("medium");
    });
  });

  it("shows a helpful global empty state and accessible controls", async () => {
    mount(mockFetch(listResponse([])));
    await vi.waitFor(() => {
      expect(document.querySelector("#empty-state")?.textContent).toContain("Noch keine Aufgaben");
    });

    expect(document.querySelector("#empty-state")?.textContent).toContain("Aufgabe");
    expect(document.querySelector('[role="status"]')).not.toBeNull();
    expect(document.querySelector('[role="alert"]')).not.toBeNull();
    const title = document.querySelector<HTMLInputElement>("#todo-title")!;
    title.focus();
    expect(document.activeElement).toBe(title);
  });

  it("rejects whitespace-only titles without sending a request", async () => {
    const fetcher = mockFetch(listResponse([]));
    mount(fetcher);
    await vi.waitFor(() => expect(fetcher).toHaveBeenCalledTimes(1));

    const input = document.querySelector<HTMLInputElement>("#todo-title")!;
    input.value = "\u00a0\u2003\t";
    document.querySelector("#todo-form")!.dispatchEvent(
      new Event("submit", { bubbles: true, cancelable: true }),
    );

    expect(document.querySelector("#error-feedback")?.textContent).toContain("Titel ist erforderlich");
    expect(document.querySelector("#error-feedback")?.textContent).toContain("Keine Aufgabe erstellt");
    expect(fetcher).toHaveBeenCalledTimes(1);
  });

  it("creates a task, normalizes its title, and announces success", async () => {
    const created = { id: 2, title: "New task", status: "open" };
    const fetcher = mockFetch(listResponse([]), response(created, 201), listResponse([created]));
    mount(fetcher);
    await vi.waitFor(() => expect(fetcher).toHaveBeenCalledTimes(1));

    document.querySelector<HTMLInputElement>("#todo-title")!.value = "\u00a0New task  ";
    document.querySelector("#todo-form")!.dispatchEvent(
      new Event("submit", { bubbles: true, cancelable: true }),
    );
    await vi.waitFor(() => {
      expect(document.querySelector("#feedback")?.textContent).toBe("Aufgabe erstellt.");
      expect(document.querySelector("#todo-list")?.textContent).toContain("New task");
    });

    expect(fetcher).toHaveBeenNthCalledWith(
      2,
      "/api/todos",
      expect.objectContaining({
        body: JSON.stringify({ title: "New task", priority: "medium" }),
      }),
    );
    expect(document.querySelector("#error-feedback")?.textContent).toBe("");
  });

  it("changes task status and refreshes the list", async () => {
    const completed = { ...firstTodo, status: "completed" };
    const fetcher = mockFetch(
      listResponse([firstTodo]),
      response(completed),
      listResponse([completed]),
      response(firstTodo),
      listResponse([firstTodo]),
    );
    mount(fetcher);
    await vi.waitFor(() => expect(document.querySelector("#todo-list")?.textContent).toContain("First task"));

    document.querySelector<HTMLButtonElement>('button[aria-label="Als erledigt markieren"]')!.click();
    await vi.waitFor(() => {
      expect(document.querySelector("#feedback")?.textContent).toBe("Status geändert.");
      expect(document.querySelector("#todo-list")?.textContent).toContain("Erledigt");
    });

    document.querySelector<HTMLButtonElement>('button[aria-label="Aufgabe wieder öffnen"]')!.click();
    await vi.waitFor(() => {
      expect(document.querySelector(".todo-status")?.textContent).toBe("Offen");
    });
  });

  it("updates and displays a task priority without changing its other fields", async () => {
    const updated = { ...firstTodo, priority: "high" };
    const fetcher = mockFetch(
      listResponse([firstTodo]),
      response(updated),
      listResponse([updated]),
    );
    mount(fetcher);
    await vi.waitFor(() => expect(document.querySelector("#todo-list")?.textContent).toContain(firstTodo.title));

    const prioritySelect = document.querySelector<HTMLSelectElement>(
      'select[aria-label="Priorität für Aufgabe ändern: First task"]',
    )!;
    prioritySelect.value = "high";
    document.querySelector<HTMLButtonElement>(".todo-priority-edit button")!.click();

    await vi.waitFor(() => {
      expect(document.querySelector(".todo-priority")?.textContent).toBe("high");
      expect(document.querySelector("#feedback")?.textContent).toBe("Priorität geändert.");
    });
    expect(fetcher).toHaveBeenNthCalledWith(
      2,
      "/api/todos/1",
      expect.objectContaining({ body: JSON.stringify({ priority: "high" }) }),
    );
  });

  it("deletes immediately without confirmation and preserves remaining tasks", async () => {
    const second = { id: 2, title: "Second task", status: "open" };
    const confirm = vi.spyOn(window, "confirm").mockReturnValue(true);
    const fetcher = mockFetch(
      listResponse([firstTodo, second]),
      response(null, 204),
      listResponse([firstTodo]),
    );
    mount(fetcher);
    await vi.waitFor(() => expect(document.querySelector("#todo-list")?.textContent).toContain("Second task"));

    document.querySelector<HTMLButtonElement>('button[aria-label="Aufgabe löschen: Second task"]')!.click();
    await vi.waitFor(() => {
      expect(document.querySelector("#feedback")?.textContent).toBe("Aufgabe gelöscht.");
      expect(document.querySelector("#todo-list")?.textContent).not.toContain("Second task");
    });
    expect(document.querySelector("#todo-list")?.textContent).toContain("First task");
    expect(confirm).not.toHaveBeenCalled();
  });

  it("filters tasks and explains an empty filtered result", async () => {
    const completed = { ...firstTodo, status: "completed" };
    const fetcher = mockFetch(listResponse([firstTodo]), listResponse([]), listResponse([firstTodo]));
    mount(fetcher);
    await vi.waitFor(() => expect(document.querySelector("#todo-list")?.textContent).toContain("First task"));

    const filter = document.querySelector<HTMLSelectElement>("#todo-filter")!;
    filter.value = "completed";
    filter.dispatchEvent(new Event("change", { bubbles: true }));
    await vi.waitFor(() => {
      expect(document.querySelector("#empty-state")?.textContent).toContain("Erledigt");
      expect(document.querySelector("#empty-state")?.textContent).toContain("Alle");
    });
    expect(document.querySelector("#todo-list")?.textContent).not.toContain(completed.title);
    expect(fetcher).toHaveBeenNthCalledWith(2, expect.stringContaining("status=completed"), undefined);

    filter.value = "all";
    filter.dispatchEvent(new Event("change", { bubbles: true }));
    await vi.waitFor(() => expect(document.querySelector("#todo-list")?.textContent).toContain("First task"));
  });

  it("combines status and priority filters and resets each independently", async () => {
    const lowOpen = { id: 2, title: "Low open", status: "open", priority: "low" };
    const highCompleted = { id: 3, title: "High completed", status: "completed", priority: "high" };
    const mediumOpen = { id: 4, title: "Medium open", status: "open", priority: "medium" };
    const allTodos = [lowOpen, highCompleted, mediumOpen];
    const fetcher = mockFetch(
      listResponse(allTodos),
      listResponse([highCompleted]),
      listResponse([highCompleted]),
      listResponse([highCompleted]),
      listResponse(allTodos),
      listResponse([]),
    );
    mount(fetcher);
    await vi.waitFor(() => expect(document.querySelector("#todo-list")?.textContent).toContain("High completed"));

    const statusFilter = document.querySelector<HTMLSelectElement>("#todo-filter")!;
    const priorityFilter = document.querySelector<HTMLSelectElement>("#todo-priority-filter")!;
    statusFilter.value = "completed";
    statusFilter.dispatchEvent(new Event("change", { bubbles: true }));
    await vi.waitFor(() => expect(fetcher).toHaveBeenCalledTimes(2));

    priorityFilter.value = "high";
    priorityFilter.dispatchEvent(new Event("change", { bubbles: true }));
    await vi.waitFor(() => expect(fetcher).toHaveBeenCalledTimes(3));
    expect(fetcher).toHaveBeenNthCalledWith(
      3,
      "/api/todos?status=completed&priority=high",
      undefined,
    );
    expect(document.querySelector("#todo-list")?.textContent).toContain("High completed");

    statusFilter.value = "all";
    statusFilter.dispatchEvent(new Event("change", { bubbles: true }));
    await vi.waitFor(() => expect(fetcher).toHaveBeenCalledTimes(4));
    expect(fetcher).toHaveBeenNthCalledWith(4, "/api/todos?priority=high", undefined);

    priorityFilter.value = "all";
    priorityFilter.dispatchEvent(new Event("change", { bubbles: true }));
    await vi.waitFor(() => expect(fetcher).toHaveBeenCalledTimes(5));
    expect(fetcher).toHaveBeenNthCalledWith(5, "/api/todos", undefined);

    priorityFilter.value = "low";
    priorityFilter.dispatchEvent(new Event("change", { bubbles: true }));
    await vi.waitFor(() => {
      expect(document.querySelector("#empty-state")?.textContent).toContain("Priorität");
      expect(document.querySelector("#empty-state")?.textContent).toContain("low");
    });
    expect(fetcher).toHaveBeenNthCalledWith(6, "/api/todos?priority=low", undefined);
  });

  it("reports missing tasks as unavailable instead of reporting success", async () => {
    const fetcher = mockFetch(
      listResponse([firstTodo]),
      response({ error: "Aufgabe ist nicht mehr verfügbar." }, 404),
    );
    mount(fetcher);
    await vi.waitFor(() => expect(document.querySelector("#todo-list")?.textContent).toContain("First task"));

    document.querySelector<HTMLButtonElement>('button[aria-label="Als erledigt markieren"]')!.click();
    await vi.waitFor(() => {
      expect(document.querySelector("#error-feedback")?.textContent).toBe("Aufgabe ist nicht mehr verfügbar.");
    });
    expect(document.querySelector("#feedback")?.textContent).toBe("");
  });

  it("reports a failed create request without showing success", async () => {
    const fetcher = mockFetch(
      listResponse([]),
      response({ error: "Aktion fehlgeschlagen. Bitte erneut versuchen." }, 500),
    );
    mount(fetcher);
    await vi.waitFor(() => expect(fetcher).toHaveBeenCalledTimes(1));

    document.querySelector<HTMLInputElement>("#todo-title")!.value = "New task";
    document.querySelector("#todo-form")!.dispatchEvent(
      new Event("submit", { bubbles: true, cancelable: true }),
    );

    await vi.waitFor(() => {
      expect(document.querySelector("#error-feedback")?.textContent).toBe(
        "Aktion fehlgeschlagen. Bitte erneut versuchen.",
      );
    });
    expect(document.querySelector("#error-feedback")?.getAttribute("role")).toBe("alert");
    expect(document.querySelector("#error-feedback")?.getAttribute("aria-live")).toBe("assertive");
    expect(document.querySelector("#feedback")?.textContent).toBe("");
    expect(document.querySelector("#todo-list")?.textContent).not.toContain("New task");
  });

  it("reports a failed delete request without showing success or removing the task", async () => {
    const fetcher = mockFetch(
      listResponse([firstTodo]),
      response({ error: "Aktion fehlgeschlagen. Bitte erneut versuchen." }, 500),
    );
    mount(fetcher);
    await vi.waitFor(() =>
      expect(document.querySelector("#todo-list")?.textContent).toContain(firstTodo.title),
    );

    document.querySelector<HTMLButtonElement>(
      'button[aria-label="Aufgabe löschen: First task"]',
    )!.click();

    await vi.waitFor(() => {
      expect(document.querySelector("#error-feedback")?.textContent).toBe(
        "Aktion fehlgeschlagen. Bitte erneut versuchen.",
      );
    });
    expect(document.querySelector("#error-feedback")?.getAttribute("role")).toBe("alert");
    expect(document.querySelector("#error-feedback")?.getAttribute("aria-live")).toBe("assertive");
    expect(document.querySelector("#feedback")?.textContent).toBe("");
    expect(document.querySelector("#todo-list")?.textContent).toContain(firstTodo.title);
    expect(fetcher).toHaveBeenCalledTimes(2);
  });
});
