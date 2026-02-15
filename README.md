# logg

A simple, thread-safe logging wrapper for `zerolog` with dynamic redirection support.

## Installation

```bash
go get github.com/cyphix/logg
```

## Usage

```go
import "github.com/cyphix/logg"

func main() {
    logg.Init("info")
    logger := logg.Ctx("main", "startup")
    logger.Info().Msg("Application starting...")
}
```
