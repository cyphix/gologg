package logg

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

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

func TestInitConsole(t *testing.T) {
	// Restoring state
	defer Init("info")

	// InitConsole is hardcoded to os.Stderr, so we just verify it doesn't panic
	// and sets the global level correctly.
	InitConsole("debug")

	assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
}
