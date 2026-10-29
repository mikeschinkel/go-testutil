#!/usr/bin/env bash
#
# E-2411's suite, run by `endless task verify E-2411` through verify.toml.
# Emits TAP on stdout; diagnostics go to stderr.
#
# It proves the tree this suite sits in: every module builds, vets, passes and
# is tidy from published tags alone, requires go-cliutil v0.4.0 and go-dt
# v0.6.1 with no cross-repo replace, lints clean with the pinned golangci-lint,
# the dependents build and pass against it, and wip/E-2411 holds the snapshot
# of main's uncommitted work.

set -u

: "${ENDLESS_VERIFY_DIR:?run this through: endless task verify E-2411}"

WT="$(cd "${ENDLESS_VERIFY_DIR}/../../.." && pwd)"
MAIN="$(cd "$(git -C "${WT}" rev-parse --git-common-dir)/.." && pwd)"
PKGS="$(dirname "${MAIN}")"
SCRATCH="$(mktemp -d)"
trap 'rm -rf "${SCRATCH}"' EXIT

# Every module must build from published tags alone: no workspace, and no
# go.mod/go.sum rewrites hiding a missing entry.
export GOWORK=off GOFLAGS=-mod=readonly

N=0
FAILED=0
ok() { N=$((N + 1)); echo "ok ${N} - $1"; }
not_ok() {
    N=$((N + 1)); FAILED=$((FAILED + 1))
    echo "not ok ${N} - $1"
    [[ -n "${2:-}" ]] && printf '%s\n' "$2" | tail -20 | sed 's/^/# /'
}
# check <description> <dir> <command...>: passes when the command exits 0 in <dir>.
check() {
    local desc="$1" dir="$2" out
    shift 2
    if out="$(cd "${dir}" && "$@" 2>&1)"; then ok "${desc}"; else not_ok "${desc}" "${out}"; fi
}
# same <description> <expected> <actual>
same() {
    if [[ "$3" == "$2" ]]; then ok "$1"; else not_ok "$1" "expected: $2"$'\n'"actual:   $3"; fi
}
replaces() {
    (cd "${WT}/$1" && go mod edit -json | jq -r '.Replace // [] | .[] | "\(.Old.Path) => \(.New.Path)"')
}
version_of() {
    (cd "${WT}/$1" && go list -m -f '{{.Version}}' "$2" 2>&1)
}

MODULES=(. test examples/basic_usage)

# Every module builds, vets and passes its tests with GOWORK=off.
check "root: go build" "${WT}" go build ./...
check "root: go vet" "${WT}" go vet ./...
check "test/: go test -race" "${WT}/test" go test -race ./...
check "test/: go vet" "${WT}/test" go vet ./...
check "examples/basic_usage: go build" "${WT}/examples/basic_usage" go build -o /dev/null ./...
check "examples/basic_usage: go vet" "${WT}/examples/basic_usage" go vet ./...

# Every module is tidy and gofmt-clean.
for m in "${MODULES[@]}"; do
    check "${m}: go mod tidy -diff is empty" "${WT}/${m}" go mod tidy -diff
done
same "gofmt -l reports nothing" "" "$(cd "${WT}" && gofmt -l . 2>&1)"

# go-cliutil v0.4.0 and go-dt v0.6.1, with no replace for another repo's module.
for m in . test; do
    same "${m}: requires go-cliutil v0.4.0" "v0.4.0" "$(version_of "${m}" github.com/mikeschinkel/go-cliutil)"
done
for m in "${MODULES[@]}"; do
    same "${m}: requires go-dt v0.6.1" "v0.6.1" "$(version_of "${m}" github.com/mikeschinkel/go-dt)"
done
same "root: no replace directives" "" "$(replaces .)"
same "test/: replaces only go-testutil, with .." \
    "github.com/mikeschinkel/go-testutil => .." "$(replaces test)"
same "examples/basic_usage: replaces only go-testutil, with ../../" \
    "github.com/mikeschinkel/go-testutil => ../../" "$(replaces examples/basic_usage)"

# golangci-lint v2.15.0, pinned the same in the Makefile and CI.
LINTER="github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0"
check "Makefile pins golangci-lint v2.15.0" "${WT}" grep -qF "LINTER = \"${LINTER}\"" Makefile
check "test.yml pins golangci-lint v2.15.0" "${WT}" grep -qF "version: v2.15.0" .github/workflows/test.yml
check "root: golangci-lint reports no issues" "${WT}" go run "${LINTER}" run ./... --timeout=5m

# Dependents build and pass against this branch. Each reaches go-testutil
# through a replace pointing at the main checkout; a scratch go.work points it
# at this worktree instead.
for dep in go-cfgstore/test go-cliutil/test; do
    printf 'go 1.25.3\nuse (\n\t%s\n\t%s\n)\n' "${WT}" "${PKGS}/${dep}" >"${SCRATCH}/dep.work"
    check "${dep}: go vet against this branch" "${PKGS}/${dep}" \
        env GOWORK="${SCRATCH}/dep.work" GOFLAGS= go vet ./...
    check "${dep}: go test against this branch" "${PKGS}/${dep}" \
        env GOWORK="${SCRATCH}/dep.work" GOFLAGS= go test ./...
done

# wip/E-2411 holds the snapshot of main's uncommitted work.
same "wip/E-2411's tip is the snapshot commit" "WIP snapshot of main working tree" \
    "$(git -C "${WT}" log -1 --format=%s wip/E-2411 2>&1)"
same "the snapshot touches only the six go.mod/go.sum files" \
    "$(printf '%s\n' examples/basic_usage/go.mod examples/basic_usage/go.sum go.mod go.sum test/go.mod test/go.sum)" \
    "$(git -C "${WT}" diff --name-only wip/E-2411^ wip/E-2411 2>&1)"

echo "1..${N}"
(( FAILED == 0 ))
