package main

import (
	"fmt"
	"log/slog"

	"github.com/mikeschinkel/go-testutil"
)

// Example demonstrating basic usage of go-testutil for testing helpers
func main() {
	fmt.Println("go-testutil Basic Usage Example")
	fmt.Printf("================================%s", "\n\n")

	// Example 1: BufferedWriter for capturing output
	fmt.Println("1. BufferedWriter Example")
	fmt.Println("   Captures stdout/stderr for testing CLI applications")

	bw := testutil.NewBufferedWriter()

	// Write to stdout and stderr
	bw.Printf("This goes to stdout\n")
	bw.Errorf("This goes to stderr\n")

	// Retrieve what was written
	fmt.Printf("   Captured stdout: %q\n", bw.GetStdout())
	fmt.Printf("   Captured stderr: %q\n", bw.GetStderr())

	// Example 2: BufferedWriter with verbosity
	fmt.Println("\n2. BufferedWriter with Verbosity")
	fmt.Println("   Controls output based on verbosity level")

	bwVerbose := testutil.NewBufferedWriter()
	bwVerbose.SetVerbosity(2)
	bwVerbose.Printf("Verbose output\n")
	fmt.Printf("   Captured: %q\n", bwVerbose.GetStdout())

	// Example 3: BufferedLogHandler for capturing log output
	fmt.Println("\n3. BufferedLogHandler Example")
	fmt.Println("   Captures slog output for testing logging behavior")

	handler := testutil.NewBufferedLogHandler()
	logger := slog.New(handler)

	// Log some messages
	logger.Info("Application started", "version", "1.0.0")
	logger.Warn("Configuration missing", "file", "config.yaml")
	logger.Error("Failed to connect", "error", "connection refused")

	// Retrieve captured log entries
	entries, err := handler.GetLogEntriesByLevel(slog.LevelInfo)
	if err != nil {
		fmt.Printf("   Error reading entries: %v\n", err)
	} else if len(entries) > 0 {
		fmt.Printf("   Captured %d INFO log entries\n", len(entries))
		fmt.Printf("     Example: [%s] %s\n", entries[0].Level, entries[0].Message)
	}

	// Example 4: Get all log output as string
	fmt.Println("\n4. Log Output as String")
	fmt.Println("   Retrieving all logged output:")

	allLogs := handler.String()
	if len(allLogs) > 0 {
		lines := len(allLogs) / 50 // Rough estimate
		if lines > 3 {
			lines = 3
		}
		fmt.Printf("   Log buffer contains ~%d lines of JSON output\n", lines)
	}

	// Example 5: BufferedLogger (alternative logger interface)
	fmt.Println("\n5. BufferedLogger Example")
	fmt.Println("   Simple buffered logger for basic testing needs")

	buffLogger := testutil.GetBufferedLogger()
	buffLogger.Info("First log message")
	buffLogger.Info("Second log message")

	logEntries, err := testutil.GetBufferedLogHandler().GetLogEntries()
	if err != nil {
		fmt.Printf("   Error reading entries: %v\n", err)
	} else {
		fmt.Printf("   Logger captured %d messages\n", len(logEntries))
		for i, entry := range logEntries {
			fmt.Printf("   - Entry %d:\n", i+1)
			for key, value := range entry {
				fmt.Printf("     - %s: %v\n", key, value)
			}
		}
	}

	// Example 6: Quiet mode
	fmt.Println("\n6. Quiet Mode Example")
	fmt.Println("   Suppressing output with quiet flag")

	bwQuiet := testutil.NewBufferedWriter()
	bwQuiet.SetQuiet(true)
	bwQuiet.Printf("This is suppressed\n")

	output := bwQuiet.GetStdout()
	fmt.Printf("   Buffer is empty (output suppressed): %v\n", len(output) == 0)

	fmt.Println("\nUsage Notes:")
	fmt.Println("- Use BufferedWriter in tests to verify CLI output")
	fmt.Println("- Use BufferedLogHandler to test logging behavior")
	fmt.Println("- Use verbosity levels to control test output detail")
	fmt.Println("- All buffers are safe for concurrent access")
}
