# Open Lab Sessions - FS 010 Spec Dev Go


## Todo App mit Spec Kit & Go


Wir entwickeln eine kleine Todo-Anwendung nicht direkt aus einem Prompt, sondern konsequent **von der Spezifikation zum Code**




Kurzbeschreibung:
Wir entwickeln eine kleine Todo-Anwendung nicht direkt aus einem Prompt, sondern konsequent von der Spezifikation zum Code. Mit GitHub Spec Kit entstehen zunächst Requirements, Plan und Tasks. Anschließend implementiert ein AI-Coding-Agent das Go-Backend und ein einfaches Web-Frontend. Zum Abschluss prüfen wir, ob Implementierung und ursprüngliche Spec tatsächlich übereinstimmen.

Das passt sehr gut zur aktuellen Spec-Kit-Philosophie: Specify → Plan → Tasks → Implement → Converge. Für umfangreichere Projekte können clarify, checklist und analyze als zusätzliche Quality Gates eingesetzt werden.


> [!important]
> 
> Wir sind an dieser Stelle noch NICHIT im Spec-Driven- Modus
> Die volgenden Beschreibungen dienen nur zur Orientierung


## Feature Plan  

| Bereich | Funktion |
|---|---|
| Todo | Aufgabe anlegen |
| Todo | Aufgaben anzeigen |
| Todo | Aufgabe als erledigt markieren |
| Todo | Aufgabe löschen |
| Filter | alle / offen / erledigt |
| Persistenz | zunächst In-Memory |
| Backend | Go REST API |
| Frontend | HTML + CSS + TypeScript/JavaScript |
| Tests | Go Unit-/API-Tests |
| Qualität | Race Detector / Spec-Abgleich |



## Architektur
Browser
   │
   │ HTTP / JSON
   ▼
Go REST API
   │
   ▼
TodoService
   │
   ▼
TodoRepository
   │
   ▼
In-Memory Storage




## Spec-Kit Commands


- `/speckit-constitution`
- `/speckit-specify`
- `/speckit-clarify`
- `/speckit-plan`
- `/speckit-checklist`
- `/speckit-tasks`
- `/speckit-analyze`
- `/speckit-implement`
- `/speckit-converge`



## Copilot und zusätzlich Claude

* claude optional - oder auch nur claude bzw. codex

```
                Spec Kit
                   │
                   ▼
            GitHub Copilot
           Hauptentwickler
                   │
        ┌──────────┴───────────┐
        ▼                      ▼
   Spezifikation           Implementierung
        │                      │
        └──────────┬───────────┘
                   ▼
                 Claude
             Review / Kritik
```




### Installation



## Install


> [!NOTE]
> 
> ## Go installieren
> 
> Go kann von der offiziellen Webseite heruntergeladen werden:
> 
> [https://go.dev/dl/](https://go.dev/dl/)
> 
> Nach der Installation prüfen:
> 
> ```
> go version
> ```
> 
> Beispiel:
> 
> ```
> go version go1.25.3 windows/amd64
> ```
> 
> Weitere Konfiguration ist für dieses Projekt normalerweise nicht erforderlich.
> 
> 

> [!NOTE]
> ## Projekt starten
> 
> Projektverzeichnis öffnen:
> 
> ```
> cd todo-app
> ```
> 
> Abhängigkeiten laden:
> 
> ```
> go mod tidy
> ```
> 
> Anwendung starten:
> 
> ```
> go run .
> ```
> 
> Tests ausführen:
> 
> ```
> go test ./...
> ```
> 
> Optional mit Race Detector:
> 
> ```
> go test -race ./...
> ```
> 
> Unter Windows benötigt der Race Detector zusätzlich einen installierten C-Compiler, z. B. MinGW-w64.

> [!NOTE]
> ## Spec Kit installieren
> 
> Voraussetzungen:
> 
> - Python **3.11+**
>     
> - [uv](https://docs.astral.sh/uv/)
>     
> - GitHub Copilot oder ein anderer unterstützter AI-Agent
>     
> 
> Spec Kit installieren:
> 
> ```
> uv tool install specify-cli
> ```
> 
> Installation prüfen:
> 
> ```
> specify version
> ```
> 
> Ein neues Projekt mit GitHub Copilot initialisieren:
> 
> ```
> specify init <projektname> --integration copilot
> ```
> 
> Unter Windows kann explizit PowerShell als Script-Typ gewählt werden:
> 
> ```
> specify init <projektname> --integration copilot --script ps
> ```
> 
> Anschließend:
> 
> ```
> cd <projektname>
> ```
> 
> Spec Kit ist damit eingerichtet und die `/speckit-*`-Kommandos stehen im unterstützten AI-Agent zur Verfügung.





## Kickstart
```
uv tool install specify-cli
specify init todo-app --integration copilot
cd todo-app (oder besser gleich todo-app mit cs-code öffnen)
```





