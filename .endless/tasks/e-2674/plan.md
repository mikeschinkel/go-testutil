**New methods** on *BufferedWriter (buffered_writer.go):

```go
func (w *BufferedWriter) Quiet() bool
func (w *BufferedWriter) Verbosity() cliutil.Verbosity
func (w *BufferedWriter) WithQuiet(quiet bool) cliutil.Writer
func (w *BufferedWriter) WithVerbosity(v cliutil.Verbosity) cliutil.Writer
```

WithQuiet and WithVerbosity return a new *BufferedWriter over the receiver's buffers with one setting changed, the same use level, and empty V2/V3/Loud caches; the receiver is unchanged. A test keeps reading the original BufferedWriter, which sees every line because the buffers are shared. The verbosity field becomes cliutil.Verbosity. WithVerbosity accepts 0 to 3 (cliutil.NoVerbosity: Printf writes nothing; Mike, 2026-11-25) and panics outside it.

**Setters removed** (Mike, 2026-11-26): SetQuiet and SetVerbosity go. go-testutil's tests (buffered_writer_test.go, fuzz_test.go, corpus_test.go) and examples/basic_usage move to WithQuiet/WithVerbosity. go-cliutil's test/buffered_writer_test.go calls both setters too; E-2673 changes those calls when it bumps go-cliutil's test/go.mod to this task's tag.

**One mutex.** Today the Writers V2(), V3() and Loud() return each have their own sync.RWMutex but write into the parent's buffers, so concurrent writes through two of them race. Every Writer over the same buffers shares one *sync.RWMutex, including those from WithQuiet and WithVerbosity.

**Loud.** BufferedWriter gains a loud field that Loud() sets; a loud Writer prints whatever its quiet setting, matching cliWriter after E-2673.

**Release.** go-testutil's remote already has v0.3.0 (not fetched locally). Removing the setters breaks callers, so after landing this task tags v0.4.0, which E-2673 requires.

**Tests** in test/buffered_writer_test.go: the getters; WithQuiet and WithVerbosity leave the receiver unchanged and their output appears in the original's buffers under the new setting; V2/V3 of a quiet Writer print nothing; a Loud Writer prints after WithQuiet(true); verbosity 0 prints nothing from Printf and still records Errorf; WithVerbosity(4) panics; `var _ cliutil.Writer = (*BufferedWriter)(nil)` still compiles; concurrent writes through a parent and its V2 and a With* Writer pass `go test -race`.
