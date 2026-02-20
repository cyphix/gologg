# Go Project Guidelines

This project adheres to modern Go standards (1.25+) to ensure performance and maintainability.

## Quick Reference

- **Build project**: `go build ./...`
- **Run tests**: `go test ./...`
- **Tidy modules**: `go mod tidy`

---

## Build/Configuration Instructions

This project is built using **Go 1.25** and managed with Go Modules.

### Prerequisites
- **Go**: Version 1.25 or higher.

### Setup and Development
1. **Install Dependencies**:
   ```bash
   go mod download
   ```
2. **Build**:
   ```bash
   go build ./...
   ```

### Dependencies
The project uses standard Go modules. 
- Logging: `github.com/rs/zerolog`
- Testing: `github.com/stretchr/testify`

---

## Testing Information

### Configuration
The project uses the standard library `testing` package along with `testify` for assertions.

### Running Tests
- **All tests**: `go test ./...`
- **Verbose output**: `go test -v ./...`
- **Coverage**: `go test -cover ./...`

### Adding New Tests
Tests should be placed in the same package as the code they test, using the `_test.go` suffix.

### Checking Coverage
When checking the coverage with an output file, make sure to clean the file up afterwards.

#### Example Test
```go
package logg

import (
    "testing"
	
    "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
    // ALWAYS use t.Context() when a test function needs a context
    ctx := t.Context()
    _ = ctx

    got := true
    assert.True(t, got)
}
```

---

## Core Principles

Write code that is **performant, idiomatic, and maintainable**. Focus on clarity and the standard library where possible.

### Type Safety & Explicitness
- Use specific types for domain models.
- Avoid `any` where a concrete type can be used.
- **Explicit Parameter Types**: Always state the type for each variable in function/method parameters, even if multiple consecutive parameters share the same type (e.g., use `(a string, b string)` instead of `(a, b string)`).

### Modern Go Idioms (1.25+)
- **`any` over `interface{}`**: Always use `any` for generic interfaces.
- **Integer Loops**: Use `for i := range n` for simple count-based loops.
- **Slices & Maps Packages**: Leverage `slices` and `maps` (e.g., `slices.Contains`, `maps.Keys`).
- **JSON Struct Tags**: Use `omitzero` instead of `omitempty` for types like `time.Time`, `time.Duration`, structs, slices, and maps.

### Error Handling & Debugging
- **Error Wrapping**: Use `fmt.Errorf("...: %w", err)` to wrap errors.
- **Error Checking**: Use `errors.Is` and `errors.As` for checking specific errors.
- **Early Returns**: Prefer early returns for error cases to reduce nesting.

---

## Coding Standards & Best Practices

Avoid unnecessary commenting unless asked otherwise.

### Core Patterns
- **Functional Options**: Use idiomatic Go patterns for library configuration.
- **Thread Safety**: The package uses a global mutex (`sync.RWMutex`) to ensure that configuration changes (like `Init` or `SetSilent`) are thread-safe.
- **Minimal Dependencies**: The project only depends on `zerolog` for logging and `testify` for testing.

---

## When Go Standards Can't Help

Focus your attention on:
1. **Business logic correctness** - Ensure algorithms match requirements.
2. **Concurrency Safety** - Use `sync` primitives or channels correctly to avoid race conditions.
3. **Resource Management** - Ensure database connections and file handles are closed properly (use `defer`).
4. **API Design** - Ensure the library API is intuitive and easy to use.

---

## Security & Privacy

To protect sensitive data and prevent accidental leakage of credentials to LLMs:
- **DO NOT READ `.env` files**: Never open or read the content of `.env` files or any other files containing secrets, API keys, or passwords.
- **Sensitive Data**: If you encounter files that appear to contain sensitive information (like private keys, certificates, or database credentials), avoid reading them unless absolutely necessary for the task, and even then, never include their contents in your responses.
