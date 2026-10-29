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
      expect.objectContaining({ body: JSON.stringify({ title: "New task" }) }),
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
