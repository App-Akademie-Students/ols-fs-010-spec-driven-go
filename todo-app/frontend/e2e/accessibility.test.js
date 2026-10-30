import assert from "node:assert/strict";
import { spawn, spawnSync } from "node:child_process";
import { once } from "node:events";
import { mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { after, before, it } from "node:test";
import { setTimeout as delay } from "node:timers/promises";
import { chromium } from "playwright";

const frontendDirectory = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const backendDirectory = resolve(frontendDirectory, "..", "backend");
const appUrl = "http://127.0.0.1:8080";
let appProcess;
let browser;
let temporaryDirectory;

before(async () => {
  temporaryDirectory = await mkdtemp(join(tmpdir(), "todo-accessibility-"));
  const executable = join(
    temporaryDirectory,
    process.platform === "win32" ? "todo-app.exe" : "todo-app",
  );
  const build = spawnSync("go", ["build", "-o", executable, "./cmd/todo"], {
    cwd: backendDirectory,
    encoding: "utf8",
  });
  if (build.error || build.status !== 0) {
    throw new Error(`Could not build the todo app:\n${build.stderr || build.error}`);
  }

  appProcess = spawn(executable, [], { cwd: backendDirectory, stdio: "ignore" });
  await waitForServer();
  browser = await chromium.launch({ headless: true });
});

after(async () => {
  await browser?.close();
  if (appProcess && appProcess.exitCode === null) {
    const exited = once(appProcess, "exit");
    appProcess.kill();
    await exited;
  }
  if (temporaryDirectory) {
    await rm(temporaryDirectory, { recursive: true, force: true });
  }
});

it("supports keyboard-only actions, visible focus, and accessible feedback", async () => {
  const page = await browser.newPage();
  let interceptedReopen = false;
  await page.route("**/api/todos/*", async (route) => {
    const request = route.request();
    if (
      !interceptedReopen &&
      request.method() === "PATCH" &&
      request.postDataJSON().status === "open"
    ) {
      interceptedReopen = true;
      await route.fulfill({
        status: 404,
        contentType: "application/json",
        body: JSON.stringify({ error: "Aufgabe ist nicht mehr verfügbar." }),
      });
      return;
    }
    await route.continue();
  });

  try {
    await page.goto(appUrl);
    await page.locator("#empty-state").waitFor({ state: "visible" });

    const title = page.getByRole("textbox", { name: "Aufgabentitel" });
    await page.keyboard.press("Tab");
    await expectFocusedWithVisibleOutline(title);
    await page.keyboard.type("Tastaturaufgabe");

    const createPriority = page.getByRole("combobox", { name: "Erstellungspriorität" });
    await page.keyboard.press("Tab");
    await expectFocusedWithVisibleOutline(createPriority);
    await assert.equal(await createPriority.inputValue(), "medium");

    const createButton = page.getByRole("button", { name: "Aufgabe hinzufügen" });
    await page.keyboard.press("Tab");
    await expectFocusedWithVisibleOutline(createButton);
    await page.keyboard.press("Enter");

    const success = page.getByRole("status");
    await success.getByText("Aufgabe erstellt.", { exact: true }).waitFor();
    await assertAccessibleLiveRegion(success, "polite");

    const filter = page.getByRole("combobox", { name: "Anzeigen" });
    await page.keyboard.press("Tab");
    await expectFocusedWithVisibleOutline(filter);
    await page.keyboard.press("End");
    await assert.equal(await filter.inputValue(), "completed");
    const emptyState = page.locator("#empty-state");
    await emptyState.getByText("Erledigt").waitFor();

    await page.keyboard.press("Home");
    await assert.equal(await filter.inputValue(), "all");
    await page.getByText("Tastaturaufgabe", { exact: true }).waitFor();

    const priorityEditor = page.getByRole("combobox", {
      name: "Priorität für Aufgabe ändern: Tastaturaufgabe",
    });
    await tabUntilFocused(page, priorityEditor);
    await expectFocusedWithVisibleOutline(priorityEditor);
    await page.keyboard.press("End");
    await assert.equal(await priorityEditor.inputValue(), "high");
    const savePriorityButton = page.getByRole("button", { name: "Priorität speichern" });
    await page.keyboard.press("Tab");
    await expectFocusedWithVisibleOutline(savePriorityButton);
    await page.keyboard.press("Enter");
    await success.getByText("Priorität geändert.", { exact: true }).waitFor();
    await assert.equal(await page.locator(".todo-priority").textContent(), "high");

    const priorityFilter = page.getByRole("combobox", { name: "Prioritätsfilter" });
    await tabUntilFocused(page, priorityFilter);
    await expectFocusedWithVisibleOutline(priorityFilter);
    await page.keyboard.press("End");
    await assert.equal(await priorityFilter.inputValue(), "high");
    await page.getByText("Tastaturaufgabe", { exact: true }).waitFor();

    await page.keyboard.press("Shift+Tab");
    await expectFocusedWithVisibleOutline(filter);
    await page.keyboard.press("End");
    await emptyState.getByText("Priorität „high“", { exact: false }).waitFor();
    await assert.match(await emptyState.textContent(), /Status.*Priorität/);

    await page.keyboard.press("Tab");
    await expectFocusedWithVisibleOutline(priorityFilter);
    await page.keyboard.press("Home");
    await assert.equal(await priorityFilter.inputValue(), "all");
    await assert.match(await emptyState.textContent(), /Status/);
    await page.keyboard.press("Shift+Tab");
    await page.keyboard.press("Home");
    await assert.equal(await filter.inputValue(), "all");
    await page.getByText("Tastaturaufgabe", { exact: true }).waitFor();

    const completeButton = page.getByRole("button", { name: "Als erledigt markieren" });
    await tabUntilFocused(page, completeButton);
    await expectFocusedWithVisibleOutline(completeButton);
    await page.keyboard.press("Enter");
    await success.getByText("Status geändert.", { exact: true }).waitFor();
    await assert.equal(await page.locator(".todo-status").textContent(), "Erledigt");

    const reopenButton = page.getByRole("button", { name: "Aufgabe wieder öffnen" });
    await tabUntilFocused(page, reopenButton);
    await expectFocusedWithVisibleOutline(reopenButton);
    await page.keyboard.press("Enter");

    const error = page.getByRole("alert");
    await error.getByText("Aufgabe ist nicht mehr verfügbar.", { exact: true }).waitFor();
    await assertAccessibleLiveRegion(error, "assertive");
    await assert.equal(await page.locator("#feedback").textContent(), "");
    await assert.equal(await page.locator(".todo-status").textContent(), "Erledigt");

    const deleteButton = page.getByRole("button", { name: "Aufgabe löschen: Tastaturaufgabe" });
    await page.keyboard.press("Tab");
    await expectFocusedWithVisibleOutline(deleteButton);
    await page.keyboard.press("Enter");
    await success.getByText("Aufgabe gelöscht.", { exact: true }).waitFor();
    await page.getByText("Tastaturaufgabe", { exact: true }).waitFor({ state: "detached" });
    await assert.equal(await page.locator("#error-feedback").textContent(), "");
  } finally {
    await page.close();
  }
});

async function waitForServer() {
  for (let attempt = 0; attempt < 100; attempt += 1) {
    if (appProcess.exitCode !== null) {
      throw new Error(`Todo app exited before becoming ready (code ${appProcess.exitCode}).`);
    }
    try {
      const response = await fetch(appUrl);
      if (response.ok) {
        return;
      }
    } catch {
      // The application may not have bound its port yet.
    }
    await delay(100);
  }
  throw new Error(`Todo app did not become ready at ${appUrl}.`);
}

async function expectFocusedWithVisibleOutline(locator) {
  const focus = await locator.evaluate((element) => {
    const style = getComputedStyle(element);
    return {
      focused: element === document.activeElement,
      keyboardFocus: element.matches(":focus-visible"),
      outlineStyle: style.outlineStyle,
      outlineWidth: Number.parseFloat(style.outlineWidth),
    };
  });
  assert.equal(focus.focused, true, "expected the control to receive keyboard focus");
  assert.equal(focus.keyboardFocus, true, "expected the control to match :focus-visible");
  assert.equal(focus.outlineStyle, "solid");
  assert.ok(focus.outlineWidth >= 3, "expected a visible focus outline at least 3px wide");
}

async function tabUntilFocused(page, locator) {
  for (let attempt = 0; attempt < 10; attempt += 1) {
    await page.keyboard.press("Tab");
    if (await locator.evaluate((element) => element === document.activeElement).catch(() => false)) {
      return;
    }
  }
  assert.fail("keyboard Tab navigation could not reach the expected control");
}

async function assertAccessibleLiveRegion(locator, politeness) {
  await assert.equal(await locator.getAttribute("aria-live"), politeness);
  await assert.equal(await locator.getAttribute("aria-atomic"), "true");
}
