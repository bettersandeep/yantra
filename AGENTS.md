# AGENTS.md

Rules for coding agents and human contributors in this repo.
Qodo reads this file as review context.

## What this repo is

yantra is a software factory.
It gives failing tests from a public test suite to coding agents as tickets, checks each change, sends it through Qodo review, and decides what merges.
This repo holds the factory.
The first product that the factory builds is in [yantra-sql](https://github.com/bettersandeep/yantra-sql).

## Setup

- Go 1.23 or later, and git.
- Build: `go build ./cmd/yantra`

## Verify

Run both commands before each commit.
CI runs the same two commands.

```
go vet ./...
go test -race ./...
```

To run them without a local Go install:

```
docker run --rm --network none -v "$PWD":/src:ro -w /src -e GOCACHE=/tmp/gocache -e GOFLAGS=-buildvcs=false golang:1.23 sh -c 'go vet ./... && go test -race ./...'
```

## Rules

### 1. A check that cannot run is a failure

- Scope: `internal/referee/**`, `cmd/yantra/**`, and each future gate or check.
- Rule: when a check cannot read its input, cannot parse it, or does not recognise it, the check returns an error and the command exits with code 2. It never reports "clean". Do not skip an input that you do not understand.
- In practice, for each output of `git`, `gh` or another program, and for each file that a check reads:
  - JSON that parses is not JSON that is correct. Reject `null`, the wrong type, and an object without the fields that the check uses.
  - A list that is as long as the limit that was asked for is incomplete. Treat it as an error.
  - A value from a fixed set, such as a git file status, is an error when it is not in the set. Do not send it to a default branch.
  - An empty string is not a value. Check it after the trim.
- Reason: pull request #2. The first version of `referee diff` passed its tests and passed CI, and Qodo found three paths where it reported "clean" without a look at the file.
- Second reason: pull requests #7 and #8. This rule was in this file, and the agent broke it again four times in two tickets. Qodo found each one. The list above is those four cases. One was a test file with a non-ASCII name. A gate that fails open is worse than no gate, because people trust it.

### 2. Resolve caller input before it reaches another program

- Scope: `cmd/yantra/**` and each place that runs `git` or another program.
- Rule: never pass a caller-supplied string as an argument to `git`. Resolve a revision to a commit ID first with `git rev-parse --verify --end-of-options`. Read paths with `-z`.
- Reason: pull request #2. `--base=--output=/tmp/x` reached `git diff` as an option, and git wrote the file.

### 3. Do not weaken a test to make it pass

- Scope: all `*_test.go` files.
- Rule: do not remove an assertion, delete a test, add a skip, or add a `nolint` comment to get a green build. If a test is wrong, say so in the pull request and change it in a separate commit with the reason.
- Reason: this is the cheat that the factory exists to catch. The repo that catches it must not do it.

### 4. A fix arrives with a test that fails on the old code

- Scope: each commit that fixes a bug or a review finding.
- Rule: add a test case that fails before the fix and passes after it. Keep the fix in a separate commit from the feature.
- Reason: a finding is closed by proof, not by a reply. The separate commit lets a reviewer, and Qodo's re-review, see the fix alone.

### 5. Standard library only

- Scope: `go.mod`.
- Rule: do not add a dependency. If the standard library cannot do the job, open an issue first and get approval from a human.
- Reason: this program decides what merges. Each dependency is code in the gate that nobody here reviewed.

## Pull requests

- One issue per pull request. The description starts with `Closes #<issue>`.
- A human merges. An agent does not merge its own pull request.
