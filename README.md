# logg

A thread-safe `zerolog` wrapper that supports global redirection.

`logg` lets you change the output destination (including multiple writers) at runtime for all existing loggers. It also enforces consistent structured fields (`pkg`, `component`) across your application.

## Quick Start

```go
import "gitea.cyphix.dev/kade/logg"

func main() {
    // Global setup
    logg.Init("info")
    
    // Usage
    logger := logg.Ctx("main", "startup")
    logger.Info().Msg("Application starting...")
}
```

## Features

- **Dynamic Redirection**: Change output writers globally. Existing loggers automatically use the new destination without being recreated.
- **Structured Context**: `Ctx(pkg, component)` ensures metadata is consistently applied.
- **Silent Mode**: Global toggle to discard all logs.
- **Console Support**: Built-in `zerolog.ConsoleWriter` integration.
- **Custom Keys**: Override default field names (`pkg`, `component`).

## Usage Examples

### Development (Pretty-Print)
```go
logg.InitConsole("debug")
```

### Multiple Destinations
```go
f, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
logg.InitWithWriters("info", os.Stderr, f)
```

### Runtime Redirection
Loggers created *before* a redirection will immediately use the new writer.
```go
logger := logg.Ctx("service", "worker")

// Logs to Stderr
logger.Info().Msg("Starting")

// Redirect all loggers to a buffer
var buf bytes.Buffer
logg.InitWithWriter("info", &buf)

logger.Info().Msg("This now goes to buf")
```

### Silent Mode
```go
logg.SetSilent(true)
// ... noise-heavy operations ...
logg.SetSilent(false)
```

### Custom Field Names
```go
logg.SetKeys("mod", "sub", "ev", "res")
logger := logg.Ctx("auth", "ldap")
// Output: {"mod": "auth", "sub": "ldap", ...}
```

## Installation

```bash
go get gitea.cyphix.dev/kade/logg
```
