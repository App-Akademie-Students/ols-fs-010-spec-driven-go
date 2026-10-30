# Feature Specification: To-Do-Prioritäten festlegen und filtern

**Feature Branch**: `Not created`

**Created**: 2026-10-30

**Status**: Draft

**Input**: User description: "Erweitere die To-Do-App Prioritäten. Eine Aufgabe kann die Priorität low, medium oder high besitzen. Benutzer können die Priorität beim Erstellen und Bearbeiten setzen. Aufgaben können nach Priorität gefiltert werden. Bestehende Aufgaben ohne explizite Priorität müssen weiterhin funktionieren."

## Clarifications

### Session 2026-10-30

- Q: Soll der Prioritätsfilter zusätzlich zum bestehenden Statusfilter gelten, sodass nur Aufgaben angezeigt werden, die beide Filter erfüllen? → A: Filter kombinieren: Nur Aufgaben, die Status- und Prioritätsfilter erfüllen, werden angezeigt.
- Q: Soll „medium“ beim Erstellen einer neuen Aufgabe bereits vorausgewählt sein, oder soll die Priorität zunächst leer bleiben und nur beim Speichern ohne Auswahl auf „medium“ gesetzt werden? → A: „medium“ ist beim Erstellen vorausgewählt und kann vor dem Speichern geändert werden.
- Q: Sollen auch ältere Aufrufer, die beim Erstellen oder Bearbeiten keine Priorität mitsenden, weiterhin akzeptiert werden und dafür „medium“ erhalten? → A: Fehlende Priorität akzeptieren; bei neuen Aufgaben „medium“ verwenden und bei bestehenden Aufgaben die bisherige Priorität erhalten.
- Q: Wenn ein älterer Aufrufer eine Aufgabe ohne Prioritätsangabe bearbeitet, soll die bereits gespeicherte Priorität erhalten bleiben oder auf „medium“ zurückgesetzt werden? → A: Die vorhandene Priorität bleibt erhalten; nur bei einer neuen Aufgabe ohne Angabe greift „medium“.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Priorität beim Erstellen und Bearbeiten festlegen (Priority: P1)

Ein Benutzer kann beim Erstellen einer Aufgabe eine der Prioritäten „low“, „medium“ oder
„high“ auswählen. Beim Bearbeiten kann er die Priorität ändern. Die gewählte Priorität ist
anschließend an der Aufgabe erkennbar.

**Why this priority**: Die Priorität gibt Benutzern direkt beim Erfassen und Pflegen ihrer
Aufgaben eine Möglichkeit, deren Wichtigkeit festzuhalten.

**Independent Test**: Aufgaben mit jeder der drei Prioritäten erstellen, anschließend eine
Priorität bearbeiten und jeweils prüfen, dass die aktuelle Auswahl angezeigt wird.

**Acceptance Scenarios**:

1. **Given** der Benutzer erstellt eine Aufgabe, **When** er „high“ auswählt und speichert,
   **Then** wird die Aufgabe mit der Priorität „high“ angezeigt.
2. **Given** eine Aufgabe besitzt die Priorität „low“, **When** der Benutzer sie bearbeitet,
   „medium“ auswählt und speichert, **Then** wird die Aufgabe mit „medium“ angezeigt.
3. **Given** der Benutzer erstellt eine Aufgabe, **When** das Formular geöffnet wird,
   **Then** ist „medium“ vorausgewählt und kann vor dem Speichern geändert werden.
4. **Given** der Benutzer wählt eine Priorität aus, **When** die verfügbaren Optionen angezeigt
   werden, **Then** kann er genau zwischen „low“, „medium“ und „high“ wählen.

---

### User Story 2 - Aufgaben nach Priorität filtern (Priority: P2)

Ein Benutzer kann die Aufgabenliste auf eine der drei Prioritäten einschränken und jederzeit
wieder alle Aufgaben anzeigen.

**Why this priority**: Der Filter hilft Benutzern, sich auf Aufgaben mit einer bestimmten
Wichtigkeit zu konzentrieren.

**Independent Test**: Aufgaben mit unterschiedlichen Prioritäten anlegen, nacheinander jede
Priorität filtern und prüfen, dass nur passende Aufgaben erscheinen; danach wieder alle
Aufgaben anzeigen.

**Acceptance Scenarios**:

1. **Given** die Liste enthält Aufgaben mit verschiedenen Prioritäten, **When** der Benutzer
   nach „high“ filtert, **Then** werden nur Aufgaben mit der Priorität „high“ angezeigt.
2. **Given** ein Prioritätsfilter ist aktiv, **When** der Benutzer „Alle“ auswählt,
   **Then** werden Aufgaben aller Prioritäten angezeigt.
3. **Given** für eine ausgewählte Priorität gibt es keine Aufgaben, **When** der Benutzer
   diese Priorität filtert, **Then** wird eine leere Ergebnisliste ohne Fehler angezeigt.
4. **Given** ein Status- und ein Prioritätsfilter sind aktiv, **When** die Aufgabenliste
   angezeigt wird, **Then** erscheinen nur Aufgaben, die beide Filter erfüllen.

---

### User Story 3 - Vorhandene Aufgaben weiterhin verwenden (Priority: P1)

Ein Benutzer kann vorhandene Aufgaben, für die bisher keine Priorität gespeichert wurde,
weiterhin anzeigen und wie gewohnt bearbeiten, erledigen oder löschen. Sie werden als
„medium“ behandelt.

**Why this priority**: Bereits vorhandene Aufgaben müssen bei der Erweiterung erhalten bleiben,
damit Benutzer ihre Arbeit nicht verlieren und die bisherigen Funktionen weiter nutzen können.

**Independent Test**: Eine Aufgabe ohne gespeicherte Priorität bereitstellen und prüfen, dass
sie sichtbar ist, als „medium“ behandelt und weiterhin bearbeitet, erledigt und gelöscht
werden kann.

**Acceptance Scenarios**:

1. **Given** eine vorhandene Aufgabe besitzt keine explizite Priorität, **When** der Benutzer
   die Aufgabenliste öffnet, **Then** bleibt die Aufgabe sichtbar und wird als „medium“
   dargestellt.
2. **Given** eine vorhandene Aufgabe besitzt keine explizite Priorität, **When** der Benutzer
   nach „medium“ filtert, **Then** erscheint diese Aufgabe in den Ergebnissen.
3. **Given** eine vorhandene Aufgabe besitzt keine explizite Priorität, **When** der Benutzer
   sie bearbeitet oder ihren Status ändert, **Then** wird die Aktion erfolgreich ausgeführt
   und die Aufgabe bleibt verfügbar.
4. **Given** eine vorhandene Aufgabe besitzt keine explizite Priorität, **When** der Benutzer
   sie löscht, **Then** wird sie wie bisher aus der Liste entfernt.
5. **Given** Aufgaben mit den Prioritäten „low“, „medium“ und „high“ sind vorhanden,
   **When** der Benutzer eine Aufgabe erledigt oder löscht, **Then** funktionieren diese
   bestehenden Aktionen unverändert.
6. **Given** eine Aufgabe besitzt eine Priorität ungleich „medium“, **When** ein älterer
   Aufrufer sie ohne Prioritätsangabe bearbeitet, **Then** bleibt ihre bisherige Priorität
   erhalten.

### Edge Cases

- Ein Filter, für den keine Aufgaben existieren, zeigt eine leere Liste statt eines Fehlers.
- Aufgaben ohne gespeicherte Priorität verhalten sich konsistent wie Aufgaben mit „medium“,
  auch beim Filtern und Bearbeiten.
- Ein Prioritätswechsel ändert nicht den Titel oder den Erledigt-Status der Aufgabe.
- Die Prioritätsauswahl lässt keine Werte außerhalb von „low“, „medium“ und „high“ zu.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Jede Aufgabe MUST genau eine Priorität aus „low“, „medium“ oder „high“ besitzen.
- **FR-002**: Benutzer MUST eine Priorität beim Erstellen einer Aufgabe auswählen können.
- **FR-003**: Beim Erstellen einer Aufgabe MUST „medium“ vorausgewählt sein; Benutzer MUST die
  Auswahl vor dem Speichern ändern können.
- **FR-004**: Benutzer MUST die Priorität einer bestehenden Aufgabe bearbeiten können.
- **FR-005**: Die Aufgabenliste MUST die aktuelle Priorität jeder Aufgabe anzeigen.
- **FR-006**: Benutzer MUST die Aufgabenliste auf „low“, „medium“ oder „high“ filtern können.
- **FR-007**: Benutzer MUST den Prioritätsfilter zurücksetzen können, um Aufgaben aller
  Prioritäten anzuzeigen.
- **FR-008**: Aufgaben ohne zuvor gespeicherte Priorität MUST weiterhin angezeigt und
  bearbeitet werden können und MUST sich wie Aufgaben mit der Priorität „medium“ verhalten.
- **FR-009**: Änderungen der Priorität MUST Titel und Erledigt-Status der Aufgabe unverändert
  lassen.
- **FR-010**: Die bestehenden Funktionen zum Anzeigen, Bearbeiten, Erledigen und Löschen von
  Aufgaben MUST für Aufgaben aller Prioritäten verfügbar bleiben.
- **FR-011**: Wenn ein Status- und ein Prioritätsfilter aktiv sind, MUST die Aufgabenliste nur
  Aufgaben anzeigen, die beide Filter erfüllen.
- **FR-012**: Die Anwendung MUST Erstellungs- und Bearbeitungsaufrufe ohne übermittelte
  Priorität weiterhin akzeptieren. Bei einer neuen Aufgabe MUST die fehlende Priorität als
  „medium“ behandelt werden; bei einer bestehenden Aufgabe MUST die bisherige Priorität
  unverändert bleiben.

### Key Entities *(include if feature involves data)*

- **Aufgabe**: Eine vom Benutzer verwaltete Aufgabe mit Titel und Erledigt-Status sowie einer
  Priorität („low“, „medium“ oder „high“). Für bereits vorhandene Aufgaben ohne explizite
  Priorität gilt „medium“.
- **Prioritätsfilter**: Eine Auswahl, die die Aufgabenliste auf Aufgaben einer Priorität
  einschränkt oder alle Prioritäten anzeigt.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In allen drei Prioritätsstufen kann ein Benutzer eine Aufgabe erstellen und die
  gespeicherte Auswahl anschließend korrekt in der Aufgabenliste erkennen.
- **SC-002**: In allen drei Prioritätsstufen zeigt der Filter ausschließlich passende Aufgaben;
  die Auswahl „Alle“ zeigt wieder Aufgaben jeder Priorität.
- **SC-003**: 100 % der bereitgestellten vorhandenen Aufgaben ohne explizite Priorität bleiben
  sichtbar, bearbeitbar und im Filter „medium“ auffindbar.
- **SC-004**: Bei einem Prioritätswechsel bleiben Titel und Erledigt-Status der Aufgabe
  unverändert.

## Assumptions

- Neue Aufgaben haben beim Erstellen standardmäßig „medium“ vorausgewählt; Benutzer können vor
  dem Speichern eine andere Priorität wählen.
- Vorhandene Aufgaben ohne explizite Priorität werden als „medium“ behandelt, ohne dass
  Benutzer sie manuell anpassen müssen.
- Ältere Aufrufer können Aufgaben ohne übermittelte Priorität erstellen und bearbeiten; neue
  Aufgaben erhalten „medium“, bestehende Aufgaben behalten ihre bisherige Priorität.
- Die Filterauswahl „Alle“ ist der Standardzustand der Aufgabenliste; Sortierung nach
  Priorität ist nicht Teil dieses Features.
