Filed from gomion's planning of RunCLIArgs (the RunCLI signature change), 2026-11-24, at Mike's request.

How it works today: testutil.BufferedWriter (buffered_writer.go) implements cliutil.Writer for tests, capturing output in shared stdout and stderr buffers. It has mutating SetQuiet(bool) and SetVerbosity(int), the latter panicking outside 1-3. Its V2(), V3() and Loud() Writers are cached and share the buffers.

What prompted it: go-cliutil is adding Quiet(), Verbosity(), WithQuiet(bool) Writer and WithVerbosity(Verbosity) Writer to the cliutil.Writer interface (Mike, 2026-11-24), so BufferedWriter must implement them to keep satisfying it. gomion's in-process CLI test will hand RunCLI a BufferedWriter and RunCLI will call WithQuiet/WithVerbosity on it; the test reads the original BufferedWriter, so the returned Writers must share its buffers. Whether SetVerbosity keeps its int parameter or is removed is open; changing it to cliutil.Verbosity breaks existing callers.
