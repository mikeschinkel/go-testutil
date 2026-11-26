package testutil

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/mikeschinkel/go-cliutil"
)

// BufferedWriter implements cliutil.Writer and captures all output in buffers
// for testing. Every Writer derived from it through V2, V3, Loud, WithQuiet or
// WithVerbosity shares its buffers and its mutex, so a test reads everything
// written through any of them from the original BufferedWriter.
type BufferedWriter struct {
	stdBuf     *bytes.Buffer
	errBuf     *bytes.Buffer
	mu         *sync.RWMutex
	quiet      bool
	loud       bool
	verbosity  cliutil.Verbosity
	useLevel   int
	loudWriter cliutil.Writer
	v2Writer   cliutil.Writer
	v3Writer   cliutil.Writer
}

// Verify BufferedWriter implements cliutil.Writer interface
var _ cliutil.Writer = (*BufferedWriter)(nil)

// NewBufferedWriter creates a new BufferedWriter with default settings
func NewBufferedWriter() *BufferedWriter {
	return &BufferedWriter{
		stdBuf:    &bytes.Buffer{},
		errBuf:    &bytes.Buffer{},
		mu:        &sync.RWMutex{},
		quiet:     false,
		verbosity: cliutil.HighVerbosity, // Default to max verbosity for testing
		useLevel:  1,                     // Default level
	}
}

// Printf writes formatted output to stdout buffer
func (w *BufferedWriter) Printf(format string, args ...any) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.quiet && !w.loud {
		return
	}
	if int(w.verbosity) < w.useLevel {
		return
	}

	formatted := fmt.Sprintf(format, args...)
	w.stdBuf.WriteString(formatted)
}

// Errorf writes formatted error output to doterr buffer
func (w *BufferedWriter) Errorf(format string, args ...any) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Process error arguments to flatten newlines (same as cliWriter)
	processedArgs := make([]any, len(args))
	for i, arg := range args {
		if err, ok := arg.(error); ok {
			processedArgs[i] = strings.ReplaceAll(err.Error(), "\n", "; ")
		} else {
			processedArgs[i] = arg
		}
	}

	formatted := fmt.Sprintf(format, processedArgs...)
	w.errBuf.WriteString(formatted)
}

// derive returns a new BufferedWriter over w's buffers and mutex with w's
// settings and empty V2/V3/Loud caches. The caller holds w.mu.
func (w *BufferedWriter) derive() *BufferedWriter {
	return &BufferedWriter{
		stdBuf:    w.stdBuf, // Share the same buffers
		errBuf:    w.errBuf,
		mu:        w.mu, // and the same mutex
		quiet:     w.quiet,
		loud:      w.loud,
		verbosity: w.verbosity,
		useLevel:  w.useLevel,
	}
}

// Loud returns a Writer that prints whatever its quiet setting
func (w *BufferedWriter) Loud() cliutil.Writer {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.loudWriter != nil {
		return w.loudWriter
	}

	lw := w.derive()
	lw.loud = true
	w.loudWriter = lw
	return w.loudWriter
}

// V2 returns a Writer for verbosity level 2
func (w *BufferedWriter) V2() cliutil.Writer {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.v2Writer != nil {
		return w.v2Writer
	}

	v2 := w.derive()
	v2.useLevel = 2
	w.v2Writer = v2
	return w.v2Writer
}

// V3 returns a Writer for verbosity level 3
func (w *BufferedWriter) V3() cliutil.Writer {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.v3Writer != nil {
		return w.v3Writer
	}

	v3 := w.derive()
	v3.useLevel = 3
	w.v3Writer = v3
	return w.v3Writer
}

// Quiet reports whether Printf output is suppressed
func (w *BufferedWriter) Quiet() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.quiet
}

// Verbosity returns the verbosity level Printf output is filtered by
func (w *BufferedWriter) Verbosity() cliutil.Verbosity {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.verbosity
}

// WithQuiet returns a Writer over w's buffers with the quiet setting changed,
// leaving w unchanged
func (w *BufferedWriter) WithQuiet(quiet bool) cliutil.Writer {
	w.mu.RLock()
	defer w.mu.RUnlock()

	qw := w.derive()
	qw.quiet = quiet
	return qw
}

// WithVerbosity returns a Writer over w's buffers with the verbosity changed,
// leaving w unchanged. It accepts 0 (cliutil.NoVerbosity, where Printf writes
// nothing) to 3, and panics outside that range.
func (w *BufferedWriter) WithVerbosity(verbosity cliutil.Verbosity) cliutil.Writer {
	if verbosity < cliutil.NoVerbosity || verbosity > cliutil.HighVerbosity {
		panic(fmt.Sprintf("Invalid verbosity for BufferedWriter; must be between 0-3; got %d", verbosity))
	}
	w.mu.RLock()
	defer w.mu.RUnlock()

	vw := w.derive()
	vw.verbosity = verbosity
	return vw
}

// Testing helper methods

// GetPrintOutput returns the current Printf buffer contents
func (w *BufferedWriter) GetPrintOutput() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.stdBuf.String()
}

// GetErrorOutput returns the current Errorf buffer contents
func (w *BufferedWriter) GetErrorOutput() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.errBuf.String()
}

// GetAllOutput returns both Printf and Errorf outout combined
func (w *BufferedWriter) GetAllOutput() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.stdBuf.String() + w.errBuf.String()
}

// ContainsStdout returns true if Printf buffer contains the specified substring
func (w *BufferedWriter) ContainsStdout(s string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return strings.Contains(w.stdBuf.String(), s)
}

// ContainsStderr returns true if Errorf buffer contains the specified substring
func (w *BufferedWriter) ContainsStderr(s string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return strings.Contains(w.errBuf.String(), s)
}

// ContainsOutput returns true if either Printf or Errorf buffer contains the specified substring
func (w *BufferedWriter) ContainsOutput(s string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return strings.Contains(w.stdBuf.String(), s) || strings.Contains(w.errBuf.String(), s)
}

// Reset clears both Printf and Errorf buffers
func (w *BufferedWriter) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.stdBuf.Reset()
	w.errBuf.Reset()
}

// GetPrintLines returns stdout content split into lines (excluding empty lines)
func (w *BufferedWriter) GetPrintLines() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	content := w.stdBuf.String()
	if content == "" {
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(content), "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}

// GetErrorLines returns doterr content split into lines (excluding empty lines)
func (w *BufferedWriter) GetErrorLines() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	content := w.errBuf.String()
	if content == "" {
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(content), "\n")
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}

// CountPrintLines returns the number of non-empty lines in stdout
func (w *BufferedWriter) CountPrintLines() int {
	return len(w.GetPrintLines())
}

// CountErrorLines returns the number of non-empty lines in doterr
func (w *BufferedWriter) CountErrorLines() int {
	return len(w.GetErrorLines())
}

func (w *BufferedWriter) Writer() io.Writer {
	return w.stdBuf
}

func (w *BufferedWriter) ErrWriter() io.Writer {
	return w.errBuf
}
