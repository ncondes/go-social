---
name: add-endpoint
description: Scaffolds a new REST endpoint across all clean architecture layers.
---

Given arguments `<resource> <action>` (e.g. `post follow` or `user block`), scaffold all clean architecture layers.

## Steps

1. Ask the user for the HTTP method and route path if not obvious from the action name.
2. Create or update files in this order: domain → repository → service → handler → router.
3. Register the new handler in `internal/handlers/handlers.go`, `internal/services/services.go`, `internal/repositories/repositories.go`, and `cmd/api/api.go`.
4. Only read files you need to update (router, aggregators). Do NOT re-read files for patterns — use the templates below.

---

## Goal: scaffolding only

Generate placeholder skeletons — correct structure, correct wiring, no guessed business logic.

- Leave SQL queries as `// TODO: write query` — never infer table names or schema
- Leave DTO fields as `// TODO: add fields` — never guess request/response shape
- Leave method signatures with `// TODO: define params` where params are unclear
- The code must compile (correct types, imports, method names) but implementations are stubs

---

## Layer templates

### domain/`<resource>`.go — add to existing file

```go
// Add method to the Repository interface
<Action>(ctx context.Context /* TODO: define params */) error

// Add method to the Service interface
<Action><Resource>(ctx context.Context /* TODO: define params */) error
```

### repositories/`<resource>`_repository.go — add method

```go
func (r *<Resource>Repository) <Action>(ctx context.Context /* TODO: define params */) error {
    ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
    defer cancel()

    // TODO: write query
    query := ``

    // TODO: exec query and handle result
    _ = query

    return nil
}
```

### services/`<resource>`_service.go — add method

```go
func (s *<Resource>Service) <Action><Resource>(ctx context.Context /* TODO: define params */) error {
    // TODO: add any business logic before/after repo call
    return s.<resource>Repository.<Action>(ctx /* TODO: pass params */)
}
```

### handlers/`<resource>`_handler.go — add method

```go
func (h *<Resource>Handler) <Action><Resource>(w http.ResponseWriter, r *http.Request) {
    // TODO: extract path params from context if needed
    // e.g. resourceID := get<Resource>IDFromContext(r.Context())

    // TODO: decode and validate request body if needed
    // e.g. var dto *dtos.<Action><Resource>DTO
    // jsonDecode / h.validator.validateStruct

    if err := h.<resource>Service.<Action><Resource>(r.Context() /* TODO: pass params */); err != nil {
        handleError(w, r, err)
        return
    }

    w.WriteHeader(http.StatusNoContent) // TODO: change status / use respondWithData if response body needed
}
```

### dtos/`<resource>`_dto.go — add if request/response body needed

```go
type <Action><Resource>DTO struct {
    // TODO: add fields with validate tags
}
```

### cmd/api/api.go — register route

```go
r.Post("/<action>", app.handlers.<Resource>Handler.<Action><Resource>)
```

Add inside the existing `/{<resource>ID}` route block. Only read this file to find the right place to insert.

### Aggregator files to update

- `internal/handlers/handlers.go` — add field and wire in `New()`
- `internal/services/services.go` — add field and wire in `New()`
- `internal/repositories/repositories.go` — add field and wire in `New()`

Only update these if the resource handler/service/repository is **new**. Skip if it already exists.

---

## Conventions to follow

- All repository methods use `context.WithTimeout(ctx, queryTimeoutDuration)`
- Interfaces live in `internal/domain/`, not in the implementing package
- DTOs use `go-playground/validator` struct tags
- Use `handleError(w, r, err)` for service errors, `respondWithError` for input errors
- Path param middleware (e.g. `PostIDMiddleware`) lives in `internal/handlers/`
