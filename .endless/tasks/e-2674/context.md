Filed from gomion's planning of RunCLIArgs (the RunCLI signature change), 2026-11-24, at Mike's request.

How it works today: testutil.BufferedWriter (buffered_writer.go) implements cliutil.Writer for tests, capturing output in shared stdout and stderr buffers. It has mutating SetQuiet(bool) and SetVerbosity(int), the latter panicking outside 1-3. Its V2(), V3() and Loud() Writers are cached and share the buffers.

What prompted it: go-cliutil is adding Quiet(), Verbosity(), WithQuiet(bool) Writer and WithVerbosity(Verbosity) Writer to the cliutil.Writer interface (Mike, 2026-11-24), so BufferedWriter must implement them to keep satisfying it. gomion's in-process CLI test will hand RunCLI a BufferedWriter and RunCLI will call WithQuiet/WithVerbosity on it; the test reads the original BufferedWriter, so the returned Writers must share its buffers. Whether SetVerbosity keeps its int parameter or is removed is open; changing it to cliutil.Verbosity breaks existing callers.

**Verbosity 0 (2026-11-25, Mike).** cliutil.Verbosity 0 (NoVerbosity) stays a real level: nothing is printed by Printf. BufferedWriter.SetVerbosity panics outside 1-3 today; its WithVerbosity accepts 0..3, matching cliutil.NewWriter after E-2673.

**Order (2026-11-26).** This task goes before E-2673, because go-cliutil's own test/ module uses BufferedWriter as a cliutil.Writer (buffered_writer_test.go, and EchoKey in terminal_test.go). Once the interface declares the four methods, those tests compile only against a go-testutil tag whose BufferedWriter already has them. BufferedWriter can gain them now, as go-testutil's required go-cliutil v0.4.0 already has cliutil.Verbosity, and this task tags a go-testutil release after landing for E-2673 to require.
