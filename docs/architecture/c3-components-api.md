# C3 — Component Diagram (HTTP API)

Components inside the **HTTP API** container follow Clean Architecture.

## Diagram

```mermaid
flowchart TB
    subgraph transport["Transport Layer"]
        router[Router]
        handlers[Handlers]
        dto[DTO]
    end

    subgraph usecase["Usecase Layer"]
        ucTpl[Template Usecase]
        ucVer[Version Usecase]
        ucDoc[Document Usecase]
        ucLog[Render Log Usecase]
        ucCb[Callback Usecase]
    end

    subgraph data["Data and Infrastructure"]
        repo[Repositories]
        infra[Infrastructure]
        bootstrap[Bootstrap]
    end

    db[(PostgreSQL)]
    kafka[Kafka]

    router --> handlers --> dto
    handlers --> ucTpl & ucVer & ucDoc & ucLog & ucCb
    ucTpl & ucVer & ucDoc --> repo
    ucDoc & ucTpl --> infra
    bootstrap --> router
    repo --> db
    infra --> kafka
```

## Endpoint → Handler → Usecase

| API Group | Handler | Usecase |
|-----------|---------|---------|
| `/templates` | `TemplateHandler` | `documenttemplates.Service` |
| `/templates/:id/versions` | `TemplateVersionHandler` | `documenttemplateversions.Service` |
| `/documents` | `DocumentHandler` | `documents.Service` |
| `/documents/:id/render-logs` | `DocumentHandler` | `documentrenderlogs.Service` |
| `/callbacks/test` | `CallbackHandler` | `documentcallbackattempts.Service` |
| `/users` | `UserHandler` | `users.UserService` (example) |

## Bootstrap Wiring

```mermaid
flowchart TB
    main[cmd/app/main.go] --> RunApp[bootstrap.RunApp]
    RunApp --> LoadConfig[LoadConfig]
    RunApp --> initDB[initDB + Migrate]
    RunApp --> wireUser[wireUserService]
    RunApp --> wireDoc[wireDocumentServices]
    RunApp --> initRedis[initRedis]
    RunApp --> echo[newEcho + RegisterRoutes]
```
