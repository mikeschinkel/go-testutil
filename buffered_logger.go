package testutil

import (
	"log/slog"
)

// Current design with package level logger cannot support parallel tests
// TODO: Redesign this to support parallel tests

var bufferedLogger *slog.Logger
var bufferedLogHandler *BufferedLogHandler

func GetBufferedLogger() *slog.Logger {
	if bufferedLogger != nil {
		goto end
	}
	ResetBufferedLogger()
end:
	return bufferedLogger
}

func GetBufferedLogHandler() *BufferedLogHandler {
	if bufferedLogHandler == nil {
		GetBufferedLogger()
	}
	return bufferedLogHandler
}

func ResetBufferedLogger() {
	bufferedLogHandler = NewBufferedLogHandler()
	bufferedLogger = slog.New(bufferedLogHandler)
}
