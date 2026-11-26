package test

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/mikeschinkel/go-cliutil"
	"github.com/mikeschinkel/go-testutil"
)

func TestBufferedWriter_Basic(t *testing.T) {
	writer := testutil.NewBufferedWriter()

	// Test basic Printf
	writer.Printf("Hello %s", "World")
	if !writer.ContainsStdout("Hello World") {
		t.Errorf("Expected stdout to contain 'Hello World', got: %q", writer.GetPrintOutput())
	}

	// Test basic Errorf
	writer.Errorf("Error: %s", "something failed")
	if !writer.ContainsStderr("Error: something failed") {
		t.Errorf("Expected doterr to contain 'Error: something failed', got: %q", writer.GetErrorOutput())
	}
}

func TestBufferedWriter_Interface(t *testing.T) {
	var w cliutil.Writer = testutil.NewBufferedWriter()

	// Test that all interface methods are available
	w.Printf("test")
	w.Errorf("error test")
	loudWriter := w.Loud()
	v2Writer := w.V2()
	v3Writer := w.V3()

	if loudWriter == nil {
		t.Error("Loud() should return a non-nil Writer")
	}
	if v2Writer == nil {
		t.Error("V2() should return a non-nil Writer")
	}
	if v3Writer == nil {
		t.Error("V3() should return a non-nil Writer")
	}
}

func TestBufferedWriter_VerbosityLevels(t *testing.T) {
	writer := testutil.NewBufferedWriter()
	w := writer.WithVerbosity(cliutil.MediumVerbosity)

	// Level 1 should work (verbosity 2 >= useLevel 1)
	w.Printf("Level 1 message")
	if !writer.ContainsStdout("Level 1 message") {
		t.Error("Level 1 message should be captured")
	}

	writer.Reset()

	// Level 2 should work (verbosity 2 >= useLevel 2)
	v2 := w.V2()
	v2.Printf("Level 2 message")
	if !writer.ContainsStdout("Level 2 message") {
		t.Error("Level 2 message should be captured")
	}

	writer.Reset()

	// Level 3 should NOT work (verbosity 2 < useLevel 3)
	v3 := w.V3()
	v3.Printf("Level 3 message")
	if writer.ContainsStdout("Level 3 message") {
		t.Error("Level 3 message should NOT be captured with verbosity 2")
	}
}

func TestBufferedWriter_QuietMode(t *testing.T) {
	writer := testutil.NewBufferedWriter()
	w := writer.WithQuiet(true)

	// Printf should be suppressed in quiet mode
	w.Printf("Should not appear")
	if writer.ContainsStdout("Should not appear") {
		t.Error("Printf should be suppressed in quiet mode")
	}

	// Errorf should still work in quiet mode
	w.Errorf("Error should appear")
	if !writer.ContainsStderr("Error should appear") {
		t.Error("Errorf should work even in quiet mode")
	}

	// Loud() should bypass quiet mode
	loud := w.Loud()
	loud.Printf("Loud message")
	if !writer.ContainsStdout("Loud message") {
		t.Error("Loud() should bypass quiet mode")
	}
}

func TestBufferedWriter_ErrorFormatting(t *testing.T) {
	writer := testutil.NewBufferedWriter()

	// Test error with newlines gets flattened
	err := errors.New("line 1\nline 2\nline 3")
	writer.Errorf("Error occurred: %v", err)

	stderr := writer.GetErrorOutput()
	if !strings.Contains(stderr, "line 1; line 2; line 3") {
		t.Errorf("Expected error newlines to be replaced with semicolons, got: %q", stderr)
	}
}

func TestBufferedWriter_HelperMethods(t *testing.T) {
	writer := testutil.NewBufferedWriter()

	// Test line counting
	writer.Printf("Line 1\n")
	writer.Printf("Line 2\n")
	writer.Printf("\n") // Empty line should be ignored
	writer.Printf("Line 3\n")

	lines := writer.GetPrintLines()
	if len(lines) != 3 {
		t.Errorf("Expected 3 non-empty lines, got %d: %v", len(lines), lines)
	}

	if writer.CountPrintLines() != 3 {
		t.Errorf("Expected CountPrintLines() to return 3, got %d", writer.CountPrintLines())
	}

	// Test reset
	writer.Reset()
	if writer.GetPrintOutput() != "" {
		t.Error("Expected stdout to be empty after Reset()")
	}
	if writer.GetErrorOutput() != "" {
		t.Error("Expected doterr to be empty after Reset()")
	}
}

func TestBufferedWriter_SharedBuffers(t *testing.T) {
	writer := testutil.NewBufferedWriter()

	// Test that V2 and V3 share the same buffers
	writer.Printf("Main message\n")
	writer.V2().Printf("V2 message\n")
	writer.V3().Printf("V3 message\n")

	stdout := writer.GetPrintOutput()
	if !strings.Contains(stdout, "Main message") {
		t.Error("Main message should be in stdout")
	}
	if !strings.Contains(stdout, "V2 message") {
		t.Error("V2 message should be in shared stdout buffer")
	}
	if !strings.Contains(stdout, "V3 message") {
		t.Error("V3 message should be in shared stdout buffer")
	}
}

func TestBufferedWriter_ConcurrentAccess(t *testing.T) {
	writer := testutil.NewBufferedWriter()

	// Test concurrent writes don't cause data races
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			writer.Printf("goroutine1-%d\n", i)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			writer.Errorf("goroutine2-%d\n", i)
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// Just verify we don't crash and have some content
	if writer.CountPrintLines() == 0 {
		t.Error("Expected some stdout lines from concurrent writes")
	}
	if writer.CountErrorLines() == 0 {
		t.Error("Expected some doterr lines from concurrent writes")
	}
}

func TestBufferedWriter_Getters(t *testing.T) {
	writer := testutil.NewBufferedWriter()
	if writer.Quiet() {
		t.Error("NewBufferedWriter() should not be quiet")
	}
	if writer.Verbosity() != cliutil.HighVerbosity {
		t.Errorf("NewBufferedWriter() verbosity: want %d, got %d", cliutil.HighVerbosity, writer.Verbosity())
	}

	qw := writer.WithQuiet(true).(*testutil.BufferedWriter)
	if !qw.Quiet() {
		t.Error("WithQuiet(true) should return a quiet Writer")
	}
	if qw.Verbosity() != writer.Verbosity() {
		t.Errorf("WithQuiet() should keep verbosity %d, got %d", writer.Verbosity(), qw.Verbosity())
	}

	vw := writer.WithVerbosity(cliutil.LowVerbosity).(*testutil.BufferedWriter)
	if vw.Verbosity() != cliutil.LowVerbosity {
		t.Errorf("WithVerbosity(1) verbosity: want 1, got %d", vw.Verbosity())
	}
	if vw.Quiet() != writer.Quiet() {
		t.Error("WithVerbosity() should keep the quiet setting")
	}
}

func TestBufferedWriter_WithLeavesReceiverUnchanged(t *testing.T) {
	writer := testutil.NewBufferedWriter()
	_ = writer.WithQuiet(true)
	_ = writer.WithVerbosity(cliutil.NoVerbosity)

	if writer.Quiet() {
		t.Error("WithQuiet() should not change the receiver's quiet setting")
	}
	if writer.Verbosity() != cliutil.HighVerbosity {
		t.Errorf("WithVerbosity() should not change the receiver's verbosity; got %d", writer.Verbosity())
	}

	writer.Printf("receiver prints\n")
	if !writer.ContainsStdout("receiver prints") {
		t.Error("Receiver should still print after WithQuiet(true)")
	}
}

func TestBufferedWriter_WithSharesBuffers(t *testing.T) {
	writer := testutil.NewBufferedWriter()

	writer.WithQuiet(false).Printf("from WithQuiet\n")
	writer.WithVerbosity(cliutil.LowVerbosity).Printf("from WithVerbosity\n")
	writer.WithQuiet(true).Errorf("error from quiet\n")

	if !writer.ContainsStdout("from WithQuiet") {
		t.Error("WithQuiet Writer output should appear in the original's buffer")
	}
	if !writer.ContainsStdout("from WithVerbosity") {
		t.Error("WithVerbosity Writer output should appear in the original's buffer")
	}
	if !writer.ContainsStderr("error from quiet") {
		t.Error("WithQuiet(true) Writer should still record Errorf in the original's buffer")
	}

	writer.WithQuiet(true).Printf("suppressed\n")
	writer.WithVerbosity(cliutil.LowVerbosity).V2().Printf("above verbosity\n")
	if writer.ContainsStdout("suppressed") || writer.ContainsStdout("above verbosity") {
		t.Errorf("With* Writers should apply their own setting; got %q", writer.GetPrintOutput())
	}
}

func TestBufferedWriter_WithKeepsUseLevelAndFreshCaches(t *testing.T) {
	writer := testutil.NewBufferedWriter()
	v3 := writer.V3().(*testutil.BufferedWriter)

	// A V3 Writer at verbosity 2 keeps use level 3, so prints nothing
	v3.WithVerbosity(cliutil.MediumVerbosity).Printf("level 3 at verbosity 2\n")
	if writer.ContainsStdout("level 3 at verbosity 2") {
		t.Error("WithVerbosity() should keep the receiver's use level")
	}

	// A cached V2 of the receiver must not be reused by the derived Writer
	_ = writer.V2()
	qw := writer.WithQuiet(true)
	qw.V2().Printf("quiet v2\n")
	if writer.ContainsStdout("quiet v2") {
		t.Error("WithQuiet() should return a Writer with empty V2/V3/Loud caches")
	}
}

func TestBufferedWriter_NoVerbosity(t *testing.T) {
	writer := testutil.NewBufferedWriter()
	w := writer.WithVerbosity(cliutil.NoVerbosity)

	w.Printf("should not appear\n")
	w.Errorf("error appears\n")

	if writer.GetPrintOutput() != "" {
		t.Errorf("Verbosity 0 should record nothing from Printf; got %q", writer.GetPrintOutput())
	}
	if !writer.ContainsStderr("error appears") {
		t.Error("Verbosity 0 should still record Errorf")
	}
}

func TestBufferedWriter_WithVerbosityPanicsOutOfRange(t *testing.T) {
	for _, v := range []cliutil.Verbosity{-1, 4} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("WithVerbosity(%d) should panic", v)
				}
			}()
			testutil.NewBufferedWriter().WithVerbosity(v)
		}()
	}
}

func TestBufferedWriter_QuietLoudAndLevels(t *testing.T) {
	writer := testutil.NewBufferedWriter()
	qw := writer.WithQuiet(true)

	qw.V2().Printf("quiet v2\n")
	qw.V3().Printf("quiet v3\n")
	if writer.GetPrintOutput() != "" {
		t.Errorf("V2 and V3 of a quiet Writer should record nothing; got %q", writer.GetPrintOutput())
	}

	qw.Loud().Printf("loud after quiet\n")
	if !writer.ContainsStdout("loud after quiet") {
		t.Error("Loud() should record Printf output regardless of WithQuiet(true)")
	}

	// Loud of the receiver, then quieted, still prints
	loud := writer.Loud().(*testutil.BufferedWriter)
	loud.WithQuiet(true).Printf("loud then quiet\n")
	if !writer.ContainsStdout("loud then quiet") {
		t.Error("A Loud Writer should print whatever its quiet setting")
	}
}

func TestBufferedWriter_SharedMutexRace(t *testing.T) {
	writer := testutil.NewBufferedWriter()
	writers := []cliutil.Writer{
		writer,
		writer.V2(),
		writer.V3(),
		writer.Loud(),
		writer.WithQuiet(false),
		writer.WithVerbosity(cliutil.HighVerbosity),
	}

	var wg sync.WaitGroup
	for i, w := range writers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				w.Printf("w%d-%d\n", i, j)
				w.Errorf("w%d-%d\n", i, j)
			}
		}()
	}
	wg.Wait()

	if got, want := writer.CountPrintLines(), len(writers)*100; got != want {
		t.Errorf("Expected %d stdout lines, got %d", want, got)
	}
	if got, want := writer.CountErrorLines(), len(writers)*100; got != want {
		t.Errorf("Expected %d stderr lines, got %d", want, got)
	}
}
