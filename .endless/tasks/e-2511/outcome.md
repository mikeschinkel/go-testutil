## Release proposal for go-testutil

go-testutil has one versioned module, the root `github.com/mikeschinkel/go-testutil`. The `test/` and `examples/basic_usage/` modules are unversioned helpers with a `replace` to the root and are never tagged.

**Last tag:** v0.2.1 (2025-12-01, commit 6ce0e95).

**Changes since v0.2.1, outside .endless/:**

- Breaking: six `BufferedWriter` accessors are renamed to say what they capture. `GetStdout` → `GetPrintOutput`, `GetStderr` → `GetErrorOutput`, `GetStdoutLines` → `GetPrintLines`, `GetStderrLines` → `GetErrorLines`, `CountStdoutLines` → `CountPrintLines`, `CountStderrLines` → `CountErrorLines`. `ContainsStdout`, `ContainsStderr`, `ContainsOutput`, `GetAllOutput` and `Reset` keep their names. Callers of the old names fail to compile.
- Dependencies: go-cliutil v0.2.0 → v0.4.0 and go-dt v0.2.6 → v0.6.1 (go-dt/appinfo v0.3.0 and go-dt/dtx v0.2.2 indirect). Consumers pick up the newer siblings transitively.
- Toolchain: `go 1.27.0` in every go.mod. A consumer on an older toolchain upgrades through GOTOOLCHAIN or by hand.
- Build and CI: golangci-lint pinned at v2.15.0; both workflows on Go 1.27 with GOEXPERIMENT=jsonv2 dropped; release notes exclude commits that touch only .endless/; /bin/ gitignored.

**Proposed version: v0.3.0.** The accessor rename is a breaking API change, and under v0.x semver a breaking change bumps the minor. go-cliutil/test and go-cfgstore/test already require go-testutil v0.3.0 behind a local replace, so this is the version those modules expect.

**Draft release description** (the `description` input of the Release workflow):

> Renames BufferedWriter's accessors to say what they capture: GetPrintOutput, GetErrorOutput, GetPrintLines, GetErrorLines, CountPrintLines and CountErrorLines replace GetStdout, GetStderr, GetStdoutLines, GetStderrLines, CountStdoutLines and CountStderrLines; callers of the old names must update. Requires go-cliutil v0.4.0, go-dt v0.6.1 and Go 1.27.

**Order of operations for the owner:** land E-2511, push main, run the Release workflow with version `v0.3.0` and the description above. gomion can then drop its replace for go-testutil and require v0.3.0.
