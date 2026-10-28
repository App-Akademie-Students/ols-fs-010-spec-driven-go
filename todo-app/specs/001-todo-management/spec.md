# Feature Specification: Todo-Aufgaben verwalten

**Feature Branch**: `Not created`

**Created**: 2026-10-27

**Status**: Ready for implementation

**Input**: User description: "Erstelle eine To-Do-Anwendung. Benutzer können Aufgaben erstellen, alle Aufgaben anzeigen, Aufgaben als erledigt markieren und Aufgaben löschen. Benutzer können Aufgaben nach offenen und erledigten Status filtern."

## Clarifications

### Session 2026-10-27

- Q: Soll eine erledigte Aufgabe wieder auf „offen“ gesetzt werden können? → A: Ja, erledigte Aufgaben können wieder auf offen gesetzt werden.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Aufgaben erfassen und überblicken (Priority: P1)

Ein Benutzer erstellt Aufgaben mit einem Titel und kann anschließend alle Aufgaben samt ihrem
Status in einer gemeinsamen Liste überblicken. Neue Aufgaben sind zunächst offen.

**Why this priority**: Das Erfassen und Wiederfinden von Aufgaben bildet den Kernnutzen der
Anwendung.

**Independent Test**: Eine Aufgabe anlegen und prüfen, dass sie in der Gesamtliste mit ihrem
Titel und dem Status „offen“ erscheint.

**Acceptance Scenarios**:

1. **Given** die Aufgabenliste ist leer, **When** der Benutzer einen nicht-leeren Titel eingibt
   und die Aufgabe anlegt, **Then** erscheint die Aufgabe in der Gesamtliste als offen.
2. **Given** es gibt mehrere Aufgaben, **When** der Benutzer die Gesamtliste öffnet,
   **Then** werden alle Aufgaben mit Titel und Status angezeigt.
3. **Given** der eingegebene Titel besteht nur aus Leerraum, **When** der Benutzer die Aufgabe
   anlegen möchte, **Then** wird keine Aufgabe erstellt und der Benutzer erhält einen Hinweis.

---

### User Story 2 - Aufgaben erledigen und entfernen (Priority: P1)

Ein Benutzer wechselt den Status einer Aufgabe zwischen offen und erledigt oder löscht eine
nicht mehr benötigte Aufgabe aus seiner Liste.

**Why this priority**: Aufgabenstatus und das Entfernen überholter Einträge halten die Liste
nützlich und aktuell.

**Independent Test**: Mit einer vorhandenen Aufgabe beide Statuswechsel prüfen; mit einer
weiteren vorhandenen Aufgabe prüfen, dass Löschen sie aus der Liste entfernt.

**Acceptance Scenarios**:

1. **Given** eine Aufgabe ist offen, **When** der Benutzer sie als erledigt markiert,
   **Then** wird ihr angezeigter Status erledigt.
2. **Given** eine Aufgabe ist erledigt, **When** der Benutzer die Gesamtliste betrachtet,
   **Then** wird sie weiterhin mit dem Status erledigt angezeigt.
3. **Given** eine Aufgabe wird in der Liste angezeigt, **When** der Benutzer sie löscht,
   **Then** verschwindet sie aus der Liste und wird auch in keiner Statusfilter-Ansicht angezeigt.
4. **Given** eine Aufgabe ist erledigt, **When** der Benutzer sie wieder als offen markiert,
   **Then** wird ihr Status offen und sie erscheint in der Ansicht „Offen“.

---

### User Story 3 - Aufgaben nach Status filtern (Priority: P2)

Ein Benutzer schränkt die sichtbare Liste auf offene oder erledigte Aufgaben ein oder wechselt
zur Ansicht aller Aufgaben.

**Why this priority**: Filter helfen dabei, relevante Aufgaben schneller zu finden, nachdem die
Liste mehrere Einträge enthält.

**Independent Test**: Eine Liste mit offenen und erledigten Aufgaben bereitstellen und prüfen,
dass jede Filterauswahl genau die passenden Aufgaben zeigt.

**Acceptance Scenarios**:

1. **Given** die Liste enthält offene und erledigte Aufgaben, **When** der Benutzer „Offen“
   auswählt, **Then** werden nur offene Aufgaben angezeigt.
2. **Given** die Liste enthält offene und erledigte Aufgaben, **When** der Benutzer „Erledigt“
   auswählt, **Then** werden nur erledigte Aufgaben angezeigt.
3. **Given** eine Statusfilter-Ansicht ist aktiv, **When** der Benutzer „Alle“ auswählt,
   **Then** werden offene und erledigte Aufgaben angezeigt.
4. **Given** der ausgewählte Filter hat keine passenden Aufgaben, **When** die Liste angezeigt
   wird, **Then** benennt der Leerzustand den aktiven Status, teilt mit, dass keine passenden
   Aufgaben vorhanden sind, und ermöglicht den Wechsel zu „Alle“.

### Edge Cases

- Ein leerer oder nur aus Leerraum bestehender Titel darf keine Aufgabe erzeugen.
- Eine Filteransicht ohne passende Aufgaben zeigt den aktiven Status und unterscheidet den
  Filter-Leerzustand vom Zustand einer insgesamt leeren Aufgabenliste und von einer Fehlermeldung.
- Eine gelöschte Aufgabe darf nach einem Wechsel zwischen „Alle“, „Offen“ und „Erledigt“ nicht
  erneut erscheinen.
- Das Markieren als erledigt darf andere Aufgaben und deren Status nicht verändern.
- Eine leere Aufgabenliste teilt ausdrücklich mit, dass noch keine Aufgaben vorhanden sind,
  und weist auf die Möglichkeit zum Anlegen einer Aufgabe hin.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Benutzer MUST eine Aufgabe mit einem nicht-leeren Titel anlegen können.
- **FR-002**: Das System MUST Titel ablehnen, die leer sind oder nach Entfernung führender und
  nachgestellter Leerzeichen keinen Inhalt enthalten; es MUST den Benutzer über die Ablehnung
  informieren.
- **FR-003**: Das System MUST neu angelegte Aufgaben zunächst als offen führen.
- **FR-004**: Benutzer MUST alle vorhandenen Aufgaben mit Titel und Status anzeigen können.
- **FR-005**: Benutzer MUST den Status einer Aufgabe zwischen offen und erledigt wechseln
  können.
- **FR-006**: Benutzer MUST eine angezeigte Aufgabe löschen können; gelöschte Aufgaben MUST in
  keiner Ansicht mehr erscheinen.
- **FR-007**: Benutzer MUST die Aufgabenliste auf alle, offene oder erledigte Aufgaben filtern
  können.
- **FR-008**: Das System MUST für eine vollständig leere Aufgabenliste einen Leerzustand
  anzeigen, der ausdrücklich mitteilt, dass noch keine Aufgaben vorhanden sind, und auf die
  Möglichkeit zum Anlegen einer Aufgabe hinweist. Hat eine aktive Statusfilter-Ansicht keine
  Treffer, MUST der Leerzustand ausdrücklich mitteilen, dass keine Aufgaben dem aktiven Filter
  entsprechen, den aktiven Status benennen und den Wechsel zu „Alle“ ermöglichen. Der
  Leerzustand MUST vom Feedback für ungültige Eingaben unterscheidbar sein.
- **FR-009**: Statusänderungen und Löschungen MUST in allen Ansichten konsistent berücksichtigt
  werden.

### Key Entities

- **Aufgabe**: Ein vom Benutzer verwalteter Eintrag mit einem Titel und einem Status „offen“
  oder „erledigt“. Aufgaben sind voneinander unterscheidbar, damit einzelne Einträge geändert
  oder gelöscht werden können.
- **Statusfilter**: Eine Auswahl, die bestimmt, ob alle Aufgaben, nur offene oder nur erledigte
  Aufgaben in der Liste erscheinen.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In einer moderierten Nutzungsprüfung mit 20 Personen, die die Anwendung zuvor
  noch nicht verwendet haben, können alle 20 ohne Hilfe eine Aufgabe mit einem gültigen Titel
  anlegen und dieselbe Aufgabe anschließend in der Ansicht „Alle“ wiederfinden. Als Erfolg
  zählt ausschließlich, wenn beide Schritte für dieselbe Aufgabe abgeschlossen sind.
- **SC-002**: In allen Abnahmetests wird eine neu angelegte Aufgabe als offen angezeigt, nach
  dem Erledigt-Markieren als erledigt und nach dem Wiederöffnen wieder als offen.
- **SC-003**: In 100 % der Filterprüfungen enthält die Ergebnisliste ausschließlich Aufgaben
  des gewählten Status; „Alle“ zeigt beide Statusgruppen.
- **SC-004**: In 100 % der Löschprüfungen erscheint eine gelöschte Aufgabe weder in der
  Gesamtliste noch in einer Statusfilter-Ansicht.
- **SC-005**: Bei einem leeren oder nur aus Leerraum bestehenden Titel teilt das Feedback
  ausdrücklich mit, dass der Titel erforderlich ist und keine Aufgabe angelegt wurde. Bei einer
  vollständig leeren Liste teilt der Leerzustand mit, dass noch keine Aufgaben vorhanden sind,
  und weist auf die Aufgabenerstellung hin. Bei einer leeren Filteransicht nennt der
  Leerzustand den aktiven Status, teilt mit, dass keine passenden Aufgaben vorhanden sind, und
  ermöglicht den Wechsel zu „Alle“. Jeder dieser Zustände ist vom jeweils anderen sowie von
  einer Fehlermeldung unterscheidbar.

## Assumptions

- Die Anwendung ist für eine einzelne Person und eine persönliche Aufgabenliste vorgesehen;
  Konten, Zusammenarbeit und gemeinsames Bearbeiten sind nicht Teil dieses Umfangs.
- Ein Aufgabentitel ist erforderlich. Führende und nachgestellte Leerzeichen werden beim
  Anlegen ignoriert.
- Die Ansicht „Alle“ ist die Standardansicht beim Öffnen der Aufgabenliste.
- Das Löschen entfernt eine Aufgabe unmittelbar und dauerhaft aus der aktuellen Liste; ein
  Wiederherstellen ist nicht Bestandteil dieses Umfangs.
- Aufgaben werden während der aktiven Nutzung der Anwendung verwaltet. Verlässliche Aufbewahrung
  nach einem Neustart ist nicht Bestandteil dieses Funktionsumfangs.
- Bearbeiten von Titeln, Fälligkeitstermine, Kategorien und Erinnerungen sind nicht Bestandteil
  dieses Funktionsumfangs.
