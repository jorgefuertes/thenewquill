package log_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/jorgefuertes/thenewquill/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testWriter es un writer personalizado para capturar la salida del log
type testWriter struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

func (w *testWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

func (w *testWriter) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf.Reset()
}

func (w *testWriter) lines() []string {
	s := w.String()
	if s == "" {
		return []string{}
	}
	return strings.Split(strings.TrimSuffix(s, "\n"), "\n")
}

// setupTest configura el entorno de test y devuelve el writer y una función de limpieza
func setupTest(level log.LogLevel) (*testWriter, func()) {
	w := &testWriter{}
	log.SetOutput(w)
	log.SetLevel(level)

	return w, func() {
		log.SetOutput(nil)
		log.SetLevel(log.WarningLevel)
	}
}

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    log.LogLevel
		expected string
	}{
		{log.DebugLevel, "DEBUG"},
		{log.InfoLevel, "INFO"},
		{log.WarningLevel, "WARN"},
		{log.ErrorLevel, "ERROR"},
		{log.FatalLevel, "FATAL"},
		{log.NoLevel, "DEBUG"}, // NoLevel devuelve DEBUG por defecto
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.String())
		})
	}
}

func TestSetLevel(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	t.Run("debug level shows all messages", func(t *testing.T) {
		w.Reset()
		log.SetLevel(log.DebugLevel)

		log.Debug("debug message")
		log.Info("info message")
		log.Warning("warning message")
		log.Error("error message")

		lines := w.lines()
		require.Len(t, lines, 4)
	})

	t.Run("info level hides debug", func(t *testing.T) {
		w.Reset()
		log.SetLevel(log.InfoLevel)

		log.Debug("debug message")
		log.Info("info message")

		lines := w.lines()
		require.Len(t, lines, 1)
		assert.Contains(t, lines[0], "info message")
	})

	t.Run("warning level hides debug and info", func(t *testing.T) {
		w.Reset()
		log.SetLevel(log.WarningLevel)

		log.Debug("debug message")
		log.Info("info message")
		log.Warning("warning message")

		lines := w.lines()
		require.Len(t, lines, 1)
		assert.Contains(t, lines[0], "warning message")
	})

	t.Run("error level hides debug info and warning", func(t *testing.T) {
		w.Reset()
		log.SetLevel(log.ErrorLevel)

		log.Debug("debug message")
		log.Info("info message")
		log.Warning("warning message")
		log.Error("error message")

		lines := w.lines()
		require.Len(t, lines, 1)
		assert.Contains(t, lines[0], "error message")
	})
}

func TestSetOutput(t *testing.T) {
	t.Run("nil output defaults to stdout", func(t *testing.T) {
		log.SetOutput(nil)
		// No podemos verificar stdout directamente, pero no debe hacer panic
		assert.NotPanics(t, func() {
			log.SetLevel(log.ErrorLevel) // evitar output real
			log.Debug("test")
		})
	})

	t.Run("custom output captures messages", func(t *testing.T) {
		w, cleanup := setupTest(log.DebugLevel)
		defer cleanup()

		log.Debug("test message")

		assert.Contains(t, w.String(), "test message")
	})
}

func TestDebug(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	log.Debug("debug %s %d", "test", 42)

	output := w.String()
	assert.Contains(t, output, "[DEBUG]")
	assert.Contains(t, output, "debug test 42")
	assert.True(t, strings.HasSuffix(output, "\n"))
}

func TestInfo(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	log.Info("info %s", "message")

	output := w.String()
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "info message")
}

func TestWarning(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	log.Warning("warning %d", 123)

	output := w.String()
	assert.Contains(t, output, "[WARN]")
	assert.Contains(t, output, "warning 123")
}

func TestError(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	log.Error("error occurred: %s", "something bad")

	output := w.String()
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "error occurred: something bad")
}

func TestWithoutLevel(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	log.WithoutLevel("plain message")

	output := w.String()
	// NoLevel no debe tener prefijo de nivel
	assert.NotContains(t, output, "[DEBUG]")
	assert.NotContains(t, output, "[INFO]")
	assert.NotContains(t, output, "[WARN]")
	assert.NotContains(t, output, "[ERROR]")
	assert.Contains(t, output, "plain message")
	assert.True(t, strings.HasSuffix(output, "\n"))
}

func TestWithoutFormat(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	log.WithoutFormat(log.InfoLevel, "raw line")

	output := w.String()
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "raw line")
}

func TestMultilineOutput(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	log.Info("line1\nline2\nline3")

	lines := w.lines()
	require.Len(t, lines, 3)

	// Cada línea debe tener el prefijo de nivel
	for _, line := range lines {
		assert.Contains(t, line, "[INFO]")
	}

	assert.Contains(t, lines[0], "line1")
	assert.Contains(t, lines[1], "line2")
	assert.Contains(t, lines[2], "line3")
}

func TestFormatting(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	t.Run("no args uses format as message", func(t *testing.T) {
		w.Reset()
		log.Info("simple message without args")

		assert.Contains(t, w.String(), "simple message without args")
	})

	t.Run("with args formats correctly", func(t *testing.T) {
		w.Reset()
		log.Info("value: %d, string: %s, float: %.2f", 42, "test", 3.14)

		assert.Contains(t, w.String(), "value: 42, string: test, float: 3.14")
	})

	t.Run("special characters", func(t *testing.T) {
		w.Reset()
		log.Info("special: %% tab:\there")

		output := w.String()
		assert.Contains(t, output, "special: %")
		assert.Contains(t, output, "tab:\there")
	})
}

func TestOutputFormat(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	// El formato esperado es: "[LEVEL] mensaje\n"
	// Cuando no es terminal, el color() devuelve "[LEVEL]"

	t.Run("debug format", func(t *testing.T) {
		w.Reset()
		log.Debug("test")
		assert.Equal(t, "[DEBUG] test\n", w.String())
	})

	t.Run("info format", func(t *testing.T) {
		w.Reset()
		log.Info("test")
		assert.Equal(t, "[INFO] test\n", w.String())
	})

	t.Run("warning format", func(t *testing.T) {
		w.Reset()
		log.Warning("test")
		assert.Equal(t, "[WARN] test\n", w.String())
	})

	t.Run("error format", func(t *testing.T) {
		w.Reset()
		log.Error("test")
		assert.Equal(t, "[ERROR] test\n", w.String())
	})

	t.Run("without level format", func(t *testing.T) {
		w.Reset()
		log.WithoutLevel("test")
		assert.Equal(t, "test\n", w.String())
	})
}

func TestLevelFiltering(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	tests := []struct {
		name          string
		level         log.LogLevel
		expectedCount int
	}{
		{"DebugLevel shows 5", log.DebugLevel, 5},
		{"InfoLevel shows 4", log.InfoLevel, 4},
		{"WarningLevel shows 3", log.WarningLevel, 3},
		{"ErrorLevel shows 2", log.ErrorLevel, 2},
		{"FatalLevel shows 1", log.FatalLevel, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w.Reset()
			log.SetLevel(tt.level)

			log.Debug("debug")
			log.Info("info")
			log.Warning("warning")
			log.Error("error")
			log.WithoutLevel("nolevel") // NoLevel siempre se muestra si level <= NoLevel

			lines := w.lines()
			assert.Len(t, lines, tt.expectedCount)
		})
	}
}

func TestColorFallback(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	// Cuando no es terminal, color() debe devolver "[LEVEL]" sin códigos ANSI
	log.Debug("test")

	output := w.String()
	// No debe contener secuencias de escape ANSI
	assert.NotContains(t, output, "\x1b[")
	assert.Contains(t, output, "[DEBUG]")
}

func TestEmptyMessage(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	log.Info("")

	lines := w.lines()
	require.Len(t, lines, 1)
	assert.Equal(t, "[INFO] ", lines[0])
}

func TestLogLevelOrder(t *testing.T) {
	// Verificar que los niveles están en orden correcto
	assert.True(t, log.DebugLevel < log.InfoLevel)
	assert.True(t, log.InfoLevel < log.WarningLevel)
	assert.True(t, log.WarningLevel < log.ErrorLevel)
	assert.True(t, log.ErrorLevel < log.FatalLevel)
	assert.True(t, log.FatalLevel < log.NoLevel)
}

func TestConcurrentWrites(t *testing.T) {
	w, cleanup := setupTest(log.DebugLevel)
	defer cleanup()

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(n int) {
			log.Info("message %d", n)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	lines := w.lines()
	assert.Len(t, lines, 10)
}
