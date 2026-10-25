package logg

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, "pkg", KeyPkg)
	assert.Equal(t, "component", KeyComponent)
	assert.Equal(t, "event", KeyEvent)
	assert.Equal(t, "result", KeyResult)
}

func TestCtx(t *testing.T) {
	pkg := "test-pkg"
	component := "test-component"

	logger := Ctx(pkg, component)

	assert.NotNil(t, logger)
	assert.IsType(t, zerolog.Logger{}, logger)
}

func TestInitializationAndRedirection(t *testing.T) {
	// Ensure we reset to a known state after the test
	defer SetSilent(false)
	defer Init("info")

	var buf bytes.Buffer
	InitWithWriter("debug", &buf)

	logger := Ctx("test", "init")
	logger.Debug().Msg("hello debug")

	var result map[string]any
	err := json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "debug", result["level"])
	assert.Equal(t, "hello debug", result["message"])
	assert.Equal(t, "test", result["pkg"])
	assert.Equal(t, "init", result["component"])
}

func TestMultiWriters(t *testing.T) {
	defer SetSilent(false)

	var buf1, buf2 bytes.Buffer
	InitWithWriters("info", &buf1, &buf2)

	logger := Ctx("test", "multi")
	logger.Info().Msg("multi log")

	assert.Contains(t, buf1.String(), "multi log")
	assert.Contains(t, buf2.String(), "multi log")
}

func TestSetSilent(t *testing.T) {
	defer SetSilent(false)

	var buf bytes.Buffer
	InitWithWriter("info", &buf)

	SetSilent(true)
	logger := Ctx("test", "silent")
	logger.Info().Msg("should not appear")

	assert.Empty(t, buf.String())

	SetSilent(false)
	logger.Info().Msg("should appear")
	assert.Contains(t, buf.String(), "should appear")
}

func TestSetKeys(t *testing.T) {
	// Restore default keys after test
	defer SetKeys("pkg", "component", "event", "result")

	SetKeys("p", "c", "e", "r")

	var buf bytes.Buffer
	InitWithWriter("info", &buf)

	logger := Ctx("mypkg", "mycomp")
	logger.Info().Msg("test keys")

	var result map[string]any
	err := json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, "mypkg", result["p"])
	assert.Equal(t, "mycomp", result["c"])
	assert.NotContains(t, result, "pkg")
	assert.NotContains(t, result, "component")
}

// Ctx reads the key names configured by SetKeys; run them concurrently so
// `go test -race` catches unsynchronized access.
func TestCtxConcurrentWithSetKeys(t *testing.T) {
	defer SetKeys("pkg", "component", "event", "result")

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				SetKeys("pkg", "component", "event", "result")
			}
		}()
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = Ctx("test", "concurrent")
			}
		}()
	}
	wg.Wait()
}

func TestProxyWriterRedirection(t *testing.T) {
	// This tests that loggers created BEFORE InitWithWriter still redirect correctly.
	defer SetSilent(false)

	// Create logger before initialization
	logger := Ctx("before", "init")

	var buf bytes.Buffer
	InitWithWriter("info", &buf)

	logger.Info().Msg("redirected log")

	assert.Contains(t, buf.String(), "redirected log")
}

func TestInitWithNoWriters(t *testing.T) {
	defer SetSilent(false)
	InitWithWriters("info")
	// Should not panic and output should be Discard
}

func TestInitWithInvalidLevel(t *testing.T) {
	defer SetSilent(false)
	InitWithWriters("invalid", io.Discard)
	// Should default to InfoLevel
}

// A custom writer that reenters logg (e.g. calls Ctx) inside Write must not
// deadlock against a concurrent SetSilent/SetKeys: if proxyWriter.Write held
// mu.RLock during the call, the queued write lock would block Ctx's RLock.
type reentrantWriter struct {
	entered chan struct{}
	once    sync.Once
}

func (w *reentrantWriter) Write(b []byte) (int, error) {
	w.once.Do(func() { close(w.entered) })
	// Give the blocked write lock time to queue before Ctx takes RLock.
	time.Sleep(20 * time.Millisecond)
	_ = Ctx("inner", "writer")
	return len(b), nil
}

func TestReentrantWriterDoesNotDeadlock(t *testing.T) {
	defer SetSilent(false)
	defer Init("info")

	w := &reentrantWriter{entered: make(chan struct{})}
	InitWithWriter("info", w)

	go func() {
		<-w.entered
		SetSilent(true)
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		logger := Ctx("test", "reentrant")
		logger.Info().Msg("x")
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("deadlock: reentrant writer blocked on the log lock")
	}
}

func TestInitConsole(t *testing.T) {
	// Restoring state
	defer Init("info")

	// InitConsole is hardcoded to os.Stderr, so we just verify it doesn't panic
	// and sets the global level correctly.
	InitConsole("debug")

	assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
}
