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
   und die Aufgabe anlegt, **Then** erscheint die Aufgabe in der Gesamtliste als offen und
   eine sichtbare Rückmeldung bestätigt die Erstellung.
2. **Given** es gibt mehrere Aufgaben, **When** der Benutzer die Gesamtliste öffnet,
   **Then** werden alle Aufgaben mit Titel und Status in Erstellungsreihenfolge angezeigt.
3. **Given** der eingegebene Titel besteht nur aus Leerraum, **When** der Benutzer die Aufgabe
   anlegen möchte, **Then** wird keine Aufgabe erstellt und die Rückmeldung sagt, dass ein
   Titel erforderlich ist und keine Aufgabe erstellt wurde.

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
   **Then** wird ihr angezeigter Status erledigt und eine sichtbare Rückmeldung bestätigt die
   Änderung.
2. **Given** eine Aufgabe ist erledigt, **When** der Benutzer die Gesamtliste betrachtet,
   **Then** wird sie weiterhin mit dem Status erledigt angezeigt.
3. **Given** eine Aufgabe wird in der Liste angezeigt, **When** der Benutzer sie löscht,
   **Then** verschwindet sie ohne vorherige Bestätigungsabfrage aus der Liste, wird auch in
   keiner Statusfilter-Ansicht angezeigt und eine sichtbare Rückmeldung bestätigt das Löschen.
4. **Given** eine Aufgabe ist erledigt, **When** der Benutzer sie wieder als offen markiert,
   **Then** wird ihr Status offen, sie erscheint in der Ansicht „Offen“ und eine sichtbare
   Rückmeldung bestätigt die Änderung.

Für fehlgeschlagene Erstellungs-, Statusänderungs- oder Löschanfragen zeigt die Anwendung eine
Fehlerrückmeldung statt einer Erfolgsmeldung. Ist die Aufgabe nicht mehr vorhanden, sagt die
Rückmeldung, dass sie nicht mehr verfügbar ist.

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
   auswählt, **Then** werden nur offene Aufgaben in ihrer bisherigen Erstellungsreihenfolge
   angezeigt.
2. **Given** die Liste enthält offene und erledigte Aufgaben, **When** der Benutzer „Erledigt“
   auswählt, **Then** werden nur erledigte Aufgaben in ihrer bisherigen Erstellungsreihenfolge
   angezeigt.
3. **Given** eine Statusfilter-Ansicht ist aktiv, **When** der Benutzer „Alle“ auswählt,
   **Then** werden offene und erledigte Aufgaben angezeigt.
4. **Given** der ausgewählte Filter hat keine passenden Aufgaben, **When** die Liste angezeigt
   wird, **Then** benennt der Leerzustand den aktiven Status, teilt mit, dass keine passenden
   Aufgaben vorhanden sind, und ermöglicht den Wechsel zu „Alle“.

### Edge Cases

- Ein leerer oder nur aus Leerraum bestehender Titel darf keine Aufgabe erzeugen.
- Führender und nachgestellter Leerraum gemäß der Unicode-Eigenschaft `White_Space` wird entfernt;
  der normalisierte Titel wird gespeichert und angezeigt. Titel dürfen ansonsten beliebige
  Unicode-Zeichen und eine beliebige Produktlänge enthalten.
- Mehrere Aufgaben dürfen denselben Titel haben; sie bleiben anhand ihrer jeweiligen Aufgabe
  unterscheidbar.
- Aufgaben werden in Erstellungsreihenfolge angezeigt; Filter erhalten diese Reihenfolge.
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
  nachgestellter Zeichen der Unicode-Eigenschaft `White_Space` keinen Inhalt enthalten. Es MUST
  den so normalisierten Titel speichern. Andere Zeichen und Titel-Längen sind auf Produktebene
  nicht eingeschränkt.
  Bei Ablehnung MUST die Rückmeldung sagen, dass ein Titel erforderlich ist und keine Aufgabe
  erstellt wurde.
- **FR-003**: Das System MUST neu angelegte Aufgaben zunächst als offen führen.
- **FR-004**: Benutzer MUST alle vorhandenen Aufgaben mit Titel und Status anzeigen können.
  Aufgaben MUST in Erstellungsreihenfolge erscheinen; gefilterte Ansichten MUST diese relative
  Reihenfolge erhalten.
- **FR-005**: Benutzer MUST den Status einer Aufgabe zwischen offen und erledigt wechseln
  können.
- **FR-006**: Benutzer MUST eine angezeigte Aufgabe ohne vorherige Bestätigungsabfrage löschen
  können; das Löschen wirkt sofort und ist nicht rückgängig zu machen. Gelöschte Aufgaben MUST
  in keiner Ansicht mehr erscheinen.
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
- **FR-010**: Die Anwendung MUST erfolgreiche Erstellung, Statusänderung und Löschung durch
  sichtbares Feedback mit den Meldungen „Aufgabe erstellt.“, „Status geändert.“ und
  „Aufgabe gelöscht.“ bestätigen. Bei fehlgeschlagenen API-Anfragen MUST sie „Aktion
  fehlgeschlagen. Bitte erneut versuchen.“ oder eine genauere Fehlerrückmeldung statt einer
  Erfolgsmeldung anzeigen. Eine nicht mehr vorhandene Aufgabe MUST als „Aufgabe ist nicht mehr
  verfügbar.“ gemeldet werden.
- **FR-011**: Alle interaktiven Bedienelemente MUST mit der Tastatur erreichbar und mit ihren
  üblichen Tastaturaktionen bedienbar sein. Der Tastaturfokus MUST sichtbar sein. Erfolgs- und
  Fehlerrückmeldungen MUST für assistive Technologien wahrnehmbar sein.

### Key Entities

- **Aufgabe**: Ein vom Benutzer verwalteter Eintrag mit einem Titel und einem Status „offen“
  oder „erledigt“. Titel dürfen doppelt vorkommen; Aufgaben sind anhand ihrer jeweiligen
  Identität unterscheidbar.
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
- **SC-006**: In allen Browser-Abnahmetests können sämtliche interaktiven Elemente ausschließlich
  per Tastatur erreicht und bedient werden; der Fokus ist sichtbar und Erfolgs- sowie
  Fehlerrückmeldungen sind für assistive Technologien verfügbar.

## Assumptions

- Die Anwendung ist für eine einzelne Person und eine persönliche Aufgabenliste vorgesehen;
  Konten, Zusammenarbeit und gemeinsames Bearbeiten sind nicht Teil dieses Umfangs.
- Ein Aufgabentitel ist erforderlich. Führender und nachgestellter Leerraum gemäß der Unicode-
  Eigenschaft `White_Space` wird beim Anlegen entfernt und der normalisierte Titel gespeichert;
  darüber hinaus gibt es keine produktseitige Längen- oder Zeichenbeschränkung.
- Doppelte Aufgabentitel sind zulässig.
- Die Aufgabenliste und alle gefilterten Ansichten behalten die Erstellungsreihenfolge bei.
- Die Ansicht „Alle“ ist die Standardansicht beim Öffnen der Aufgabenliste.
- Das Löschen entfernt eine Aufgabe unmittelbar und dauerhaft aus der aktuellen Liste, ohne
  Bestätigungsabfrage; ein Wiederherstellen ist nicht Bestandteil dieses Umfangs.
- Aufgaben werden nur im Arbeitsspeicher des laufenden Prozesses gehalten und gehen verloren,
  wenn die Anwendung beendet oder neu gestartet wird. Dauerhafte Speicherung ist nicht
  Bestandteil dieses Funktionsumfangs.
- Bearbeiten von Titeln, Fälligkeitstermine, Kategorien und Erinnerungen sind nicht Bestandteil
  dieses Funktionsumfangs.
