Planned from the E-2628 session, 2026-11-26.

**Setters.** Options were: (A) keep SetQuiet and SetVerbosity(int) as they are; (B) keep them with SetVerbosity taking cliutil.Verbosity and accepting 0; (C) remove both in favour of WithQuiet/WithVerbosity. A and B leave the cached V2/V3/Loud Writers with stale settings after a set, and two ways to change settings. Mike chose C, 2026-11-26.

**Order.** This task goes before E-2673: go-cliutil's own tests use BufferedWriter as a cliutil.Writer, so once E-2673's interface declares the four methods they compile only against a go-testutil tag that already has them. go-testutil requires go-cliutil v0.4.0, which already defines cliutil.Verbosity, so nothing here waits on E-2673.

**Found while planning.** The Writers V2(), V3() and Loud() return have their own mutex over shared buffers, a data race under concurrent writes; fixed here since the new With* Writers share the same buffers.
