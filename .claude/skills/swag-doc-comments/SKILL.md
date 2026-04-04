---
name: swag-doc-comments
description: Adds swaggo annotation comments to handler functions. Use when the user asks to document an endpoint, add swagger docs, or annotate a handler.
---

Add accurate swaggo annotation comments to one or more handler functions.

## Process

For each handler to document:

1. **Read the handler function** — understand what it does, what params it takes, what it returns.
2. **Read the DTO structs** referenced in the handler (in `internal/dtos/`) — to get the correct type for `@Param body` and `@Success`.
3. **Read `cmd/api/api.go`** — to find the registered route path and HTTP method for `@Router`.
4. **Write the annotation block** directly above the `func` line, following the format below.
5. After writing all annotations, remind the user to run `make swagger` to regenerate `docs/`.

---

## Annotation format (project standard)

Use the godoc-style format: a `// MethodName godoc` header, a blank `//` separator, then annotation lines using a **real tab** between `//` and `@`. This is the format enforced by `swag fmt` and matches the swaggo official examples.

```go
// MethodName godoc
//
// [tab]@Summary      <one-line summary>
// [tab]@Description  <longer description, same as summary if trivial>
// [tab]@Tags         <resource name, plural lowercase e.g. posts, users, comments>
// [tab]@Accept       json
// [tab]@Produce      json
// [tab]@Param        <name>  <in>  <type>  <required>  "<description>"
// [tab]@Success      <code>  {object}  <type>
// [tab]@Failure      400  {object}  map[string]string
// [tab]@Failure      404  {object}  map[string]string
// [tab]@Failure      500  {object}  map[string]string
// [tab]@Router       /<path> [<method>]
func (h *XHandler) MethodName(...) {
```

Note: `[tab]` represents a real tab character (`\t`). When writing to Go files, use actual tabs — not spaces.

---

## Rules for each annotation

### @Summary / @Description

- `@Summary`: short imperative phrase (e.g. `Get a post`, `Create a comment`)
- `@Description`: one sentence explaining behavior; repeat `@Summary` if nothing more to add

### @Tags

- Use the resource name, plural, lowercase: `posts`, `users`, `comments`, `feed`, `health`
- Derive from the handler file name (e.g. `post_handler.go` → `posts`)

### @Accept / @Produce

- Always `json` for this project unless the handler explicitly does something else

### @Param

Syntax: `@Param <name> <in> <type> <required> "<description>"`

| Source | `in` value | How to detect |
| --- | --- | --- |
| URL path `{postID}` | `path` | Handler calls `getPostIDFromContext` or similar |
| Query string | `query` | Handler reads `r.URL.Query()` |
| Request body | `body` | Handler calls `jsonDecode` — use the DTO type |
| Header | `header` | Handler reads a header directly |

- For `body` params use the full package-qualified DTO type: `dtos.CreatePostDTO`
- For path params use the Go type: `int64` or `string`
- `required` is `true` for path/body params, `false` for query params

### @Success

- Read the `respondWithData(w, http.StatusXXX, <value>)` call to get the status code and type
- Use `{object}` for single structs, `{array}` for slices
- Use domain types for responses (e.g. `domain.Post`) or DTO response types if a `FromDomain` mapper exists
- If the handler returns `w.WriteHeader(http.StatusNoContent)` with no body: `@Success 204 "No content"`

### @Failure

Never use `map[string]string` — it renders as `additionalPropN: "string"` in Swagger UI.

Use the concrete DTOs from `internal/dtos/response_dto.go`:

| Situation                     | DTO to use                  | Description string          |
| ----------------------------- | --------------------------- | --------------------------- |
| Input/validation errors (400) | `dtos.ErrorsResponseDTO`    | `"Validation errors"`       |
| Resource not found (404)      | `dtos.ErrorResponseDTO`     | `"<Resource> not found"`    |
| Server errors (500)           | `dtos.ErrorResponseDTO`     | `"Internal server error"`   |

- Use `dtos.ErrorsResponseDTO` for 400 only when the handler calls `respondWithErrors` (validation path)
- Use `dtos.ErrorResponseDTO` for 400 when the handler calls `respondWithError` (single message)
- Omit 404 if the handler doesn't look up an existing resource (e.g. CreatePost)

### @Router

- Read `cmd/api/api.go` to find the exact registered path
- Path params use `{paramName}` syntax: `/posts/{postID}/comments`
- HTTP method in lowercase brackets: `[get]`, `[post]`, `[patch]`, `[delete]`

---

## Example — annotated handler in this project

```go
// GetPost godoc
//
//  @Summary      Get a post
//  @Description  Get a post by ID with author details and comment count
//  @Tags         posts
//  @Accept       json
//  @Produce      json
//  @Param        postID  path      int64  true  "Post ID"
//  @Success      200     {object}  dtos.PostResponseDTO
//  @Failure      400     {object}  map[string]string
//  @Failure      404     {object}  map[string]string
//  @Failure      500     {object}  map[string]string
//  @Router       /posts/{postID} [get]
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
```

**Important:** When writing to actual `.go` files, use a real tab character (`\t`) after `//`, not spaces. The example above uses spaces only to satisfy the markdown linter.

---

## DTO struct examples

When documenting a handler that takes a request body, also add `example` struct tags to the referenced DTO fields if they don't already have them. This makes Swagger UI show realistic values instead of `"string"`.

Read the DTO file and add `example` tags — keep them short and generic, just enough to show the shape:

```go
type CreatePostDTO struct {
    Title   string   `json:"title"   validate:"required,min=1,max=100"      example:"My first post"`
    Content string   `json:"content" validate:"required,min=1,max=1000"     example:"This is the content of my post."`
    Tags    []string `json:"tags"    validate:"omitempty,dive,min=1,max=50"  example:"go,api"`
}
```

Rules for examples:

- Strings: short and generic — `"My first post"`, `"John Doe"`, `"johndoe"` — not marketing copy
- Timestamps: use RFC3339 format, e.g. `example:"2026-04-03T10:00:00Z"`
- IDs: small integers like `example:"1"` or `example:"42"`
- Slices of strings: comma-separated, 2 values max, e.g. `example:"go,api"`

---

## What NOT to do

- Do not guess route paths — always read `cmd/api/api.go`
- Do not guess DTO field names — always read the DTO file
- Do not add `@Security` unless the handler has an auth middleware applied to its route
- Do not write generic placeholder examples like `"string"`, `"foo"`, or `"example"` — use realistic values
