# Testing Guide

This document outlines our testing strategy for the Go Social project.

## Table of Contents

- [Testing Philosophy](#testing-philosophy)
- [Testing Strategy by Layer](#testing-strategy-by-layer)
- [Writing a Test](#writing-a-test)
- [Test Naming](#test-naming)
- [Mocking](#mocking)
- [Running Tests](#running-tests)
- [Best Practices](#best-practices)

## Testing Philosophy

We follow these principles:

- **Simplicity** - Tests should be easy to read
- **Isolation** - Tests should be independent
- **Speed** - Tests should run quickly
- **Coverage** - Test all critical paths and errors
- **Confidence** - Tests enable safe refactoring

## Testing Strategy by Layer

Our application has different testing strategies per layer:

### Handlers (Unit Tests)

- **Location**: `internal/handlers/*_test.go`
- **Purpose**: Test HTTP handling, validation, error mapping
- **Mocks**: Service layer
- **Database**: No (mocked)
- **Parallel**: Yes

### Services (Unit Tests)

- **Location**: `internal/services/*_test.go`
- **Purpose**: Test business logic
- **Mocks**: Repository layer
- **Database**: No (mocked)
- **Parallel**: Yes

### Repositories (Integration Tests)

- **Location**: `internal/repositories/*_integration_test.go`
- **Purpose**: Test database operations
- **Mocks**: None
- **Database**: Yes (real test DB)
- **Parallel**: No

### Packages (Unit Tests)

- **Location**: `packages/*_test.go`
- **Purpose**: Test utilities
- **Mocks**: None
- **Database**: No
- **Parallel**: Yes

## Writing a Test

Every test follows the **AAA pattern**: Arrange, Act, Assert

### Step 1: Arrange

Set up test data, mocks, and dependencies.

### Step 2: Act

Execute the code under test.

### Step 3: Assert

Verify the results.

### Example

```go
func TestPostHandler_CreatePost(t *testing.T) {
    t.Parallel()
    
    t.Run("returns 201 when post is created successfully", func(t *testing.T) {
        t.Parallel()
        
        // ARRANGE
        mockService := &mockPostService{
            createPostFunc: func(ctx context.Context, post *domain.Post) error {
                post.ID = 1
                return nil
            },
        }
        handler := NewPostHandler(mockService, NewValidator())
        
        // ACT
        req := testutils.MakeJSONRequest(t, http.MethodPost, "/posts", payload)
        w := httptest.NewRecorder()
        handler.CreatePost(w, req)
        
        // ASSERT
        assert.Equal(t, http.StatusCreated, w.Code)
    })
}
```

## Test Naming

We use **BDD-style** naming:

```txt
returns <status> when <condition>
```

### Examples

**Handlers:**

- `returns 201 when post is created successfully`
- `returns 400 when JSON is invalid`
- `returns 404 when post is not found`

**Services:**

- `creates post successfully`
- `returns error when user not found`

**Repositories:**

- `creates post in database`
- `returns error for non-existing user`

## Mocking

We use **manual mocks** for simplicity.

### Mock Structure

```go
type mockPostService struct {
    createPostFunc      func(ctx context.Context, post *domain.Post) error
    createPostCallCount int
}

func (m *mockPostService) CreatePost(ctx context.Context, post *domain.Post) error {
    m.createPostCallCount++
    if m.createPostFunc != nil {
        return m.createPostFunc(ctx, post)
    }
    return nil
}
```

### Benefits

- Simple and explicit
- No external dependencies
- Easy to debug
- Type-safe

## Running Tests

### Commands

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests only
make test-integration

# With coverage
make test-coverage

# Specific test
go test -v -run TestPostHandler_CreatePost ./internal/handlers/...
```

### Flags

- `-v` - Verbose output
- `-race` - Race detection
- `-count=1` - Disable caching
- `-timeout=30s` - Timeout

## Best Practices

### DO

- Use `t.Parallel()` for unit tests
- Follow BDD naming
- Test all error paths
- Keep tests isolated
- Use meaningful test data
- Assert specific errors

### DON'T

- Share state between tests
- Use global variables
- Skip error testing
- Mock database in integration tests
- Run integration tests in parallel

### Time Comparisons

```go
// Wrong
assert.Equal(t, expectedTime, actualTime)

// Correct
assert.True(t, expectedTime.Equal(actualTime))
```

## Summary

| Layer        | Test Type   | Real DB | Mocks        | Parallel |
|--------------|-------------|---------|--------------|----------|
| Handlers     | Unit        | No      | Services     | Yes      |
| Services     | Unit        | No      | Repositories | Yes      |
| Repositories | Integration | Yes     | None         | No       |
| Packages     | Unit        | No      | None         | Yes      |

### Key Points

1. **AAA Pattern** - Arrange, Act, Assert
2. **BDD Naming** - returns X when Y
3. **Manual Mocks** - Simple and explicit
4. **Parallel Unit Tests** - Fast execution
5. **Real DB for Integration** - Test actual queries
6. **Full Coverage** - Success and error cases

---

For questions, discuss with the team.
