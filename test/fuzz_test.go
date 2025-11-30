package test

import (
	"log/slog"
	"testing"

	"github.com/mikeschinkel/go-testutil"
)

// FuzzBufferedWriter tests BufferedWriter with random string inputs
func FuzzBufferedWriter(f *testing.F) {
	// Seed corpus with various output types
	seeds := []string{
		"",
		"a",
		"Hello, World!",
		"Line 1\nLine 2\n",
		"Tab\tseparated\tvalues",
		"Unicode: 你好世界",
		"Mixed\r\nLine\rEndings\n",
		"Very long output",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Create buffered writer and ensure writing doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("BufferedWriter panicked with input %q: %v", input, r)
			}
		}()

		bw := testutil.NewBufferedWriter()

		// Test writing via Printf/Errorf
		bw.Printf("%s", input)
		bw.Errorf("%s", input)

		// Verify buffers are accessible
		_ = bw.GetStdout()
		_ = bw.GetStderr()
	})
}

// FuzzBufferedLogHandler tests BufferedLogHandler with random log messages
func FuzzBufferedLogHandler(f *testing.F) {
	// Seed corpus with various log message types
	seeds := []string{
		"",
		"simple message",
		"error: something went wrong",
		"Multi\nLine\nMessage",
		"Unicode: 你好世界",
		"Special chars: !@#$%^&*()",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, message string) {
		// Create handler and ensure logging doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("BufferedLogHandler panicked with message %q: %v", message, r)
			}
		}()

		handler := testutil.NewBufferedLogHandler()
		logger := slog.New(handler)

		// Test various log levels
		logger.Info(message)
		logger.Warn(message)
		logger.Error(message)
		logger.Debug(message)

		// Verify output is accessible
		_ = handler.String()
		_, _ = handler.GetLogEntries()
	})
}

// FuzzBufferedWriterVerbosity tests BufferedWriter with various verbosity levels
func FuzzBufferedWriterVerbosity(f *testing.F) {
	// Seed corpus with valid verbosity levels
	seeds := []int{1, 2, 3}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, verbosity int) {
		// Only test valid verbosity levels (1-3)
		if verbosity < 1 || verbosity > 3 {
			return
		}

		// Ensure doesn't panic with various configurations
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("BufferedWriter panicked with verbosity=%d: %v", verbosity, r)
			}
		}()

		bw := testutil.NewBufferedWriter()
		bw.SetVerbosity(verbosity)
		bw.SetQuiet(true)
		bw.Printf("test message")
		_ = bw.GetStdout()
	})
}
