// Package logg provides a simple, thread-safe logging wrapper for zerolog with dynamic redirection support.
//
// It allows for global initialization of logging levels and output destinations,
// which can be changed at runtime even for loggers that have already been created.
// It also encourages structured logging by providing package and component context
// through the Ctx function.
package logg

import (
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	isSilent     bool
	mu           sync.RWMutex
	output       io.Writer = os.Stderr
	storedOutput io.Writer = os.Stderr

	KeyPkg       = "pkg"
	KeyComponent = "component"
	KeyEvent     = "event"
	KeyResult    = "result"
)

type proxyWriter struct{}

func (p *proxyWriter) Write(b []byte) (n int, err error) {
	mu.RLock()
	defer mu.RUnlock()
	return output.Write(b)
}

func init() {
	log.Logger = log.Output(&proxyWriter{})
}

// SetKeys allows customizing the keys used for structured logging.
func SetKeys(pkg string, component string, event string, result string) {
	mu.Lock()
	defer mu.Unlock()
	KeyPkg = pkg
	KeyComponent = component
	KeyEvent = event
	KeyResult = result
}

// Init initializes the logger with the specified level and logs to os.Stderr.
func Init(level string) {
	InitWithWriters(level, os.Stderr)
}

// InitConsole initializes the logger with a ConsoleWriter for pretty-printed, colored output.
func InitConsole(level string) {
	InitWithWriters(level, zerolog.ConsoleWriter{Out: os.Stderr})
}

// InitWithWriter initializes the logger with the specified level and output writer.
func InitWithWriter(level string, w io.Writer) {
	InitWithWriters(level, w)
}

// InitWithWriters initializes the logger with the specified level and multiple output writers.
func InitWithWriters(level string, writers ...io.Writer) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	mu.Lock()
	defer mu.Unlock()

	if len(writers) == 0 {
		output = os.Stderr
	} else if len(writers) == 1 {
		output = writers[0]
	} else {
		output = io.MultiWriter(writers...)
	}
	storedOutput = output

	if isSilent {
		output = io.Discard
	}
}

// SetSilent enables or disables silent mode.
// In silent mode, all logs are discarded. When disabled, the previous writer is restored.
func SetSilent(silent bool) {
	mu.Lock()
	defer mu.Unlock()

	isSilent = silent
	if silent {
		output = io.Discard
	} else {
		output = storedOutput
	}
}

// Ctx returns a zerolog.Logger pre-configured with package and component context.
// This preserves the full zerolog API while ensuring structured metadata is always present.
func Ctx(pkg string, component string) zerolog.Logger {
	return log.With().
		Str(KeyPkg, pkg).
		Str(KeyComponent, component).
		Logger()
}
