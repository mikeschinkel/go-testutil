# Lessons Learned

A capture log of corrections, newest last. Entries are appended by
`endless lesson write`, which commits this file on the main checkout as it
writes — so recording a lesson never waits on a task landing, and never
dirties a worktree.

## Format

```
### [YYYY-MM-DD] One-line summary
<the lesson — what went wrong, why, and the rule that replaces it>
- **Project**: <project>
```

---

## Entries

### [2026-10-28] verify.toml is the only suite format; verify.sh is deprecated
In E-2411 (go-testutil) the spawn prompt said to write .endless/tasks/e-<id>/verify.sh sourcing _harness.sh. go-testutil had no harness, so I copied endless's _harness.sh and _guard.sh in. Mike dropped that: verify.toml is the only way moving forward and verify.sh is deprecated, so write a task's suite as a verify.toml manifest (a sh/TAP check if a script is needed), even when a prompt names verify.sh.
- **Project**: go-testutil
