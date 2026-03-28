# gz

gz is a tool for single commit, chained git branches.

## Overview

For years, I have followed a git workflow that consists of feature branches with a single commit where the vast majority of my commits are done via `git commit -a --amend --no-edit`. These feature branches can be chained where a sub-feature branch has a single commit for its local changes, then the single commit for the parent branch, then commits from `main`. One of the tricky parts is when I commit changes to the parent branch. I typically visit all the child branches, reset them to the new `HEAD` of the parent branch, and `cherry-pick` the local changes back into the branch. This is all a little tedious to do by hand. That's what `gz` is for, it provides automations for this workflow.

## Notes

How does one resolve a branch once it is merged? Imagine:

```
main <- feature-1 <-  feature-2
```

`feature-1` is merged into `main` but it is merged via a squash merge. We wanted an automated way to detect the merge and end up
in this state.

```
main (containing feature-1 changes in squashed commit) <- feature-2
```

This kind of works if `feature-1` is a single commit because `git` detects the single commit as a cherry-pick. But if `feature-1` is multiple commits, a rebase will create a merge conflict. One way around this is to have an explicit command to re-parent ... or to have a command to assert that a branch is in the parent, so blindly reset the local branch.

## Current State Mar 28, 2026

### Makefile and build

| Command | What it does |
|--------|----------------|
| `make` / `make all` | Builds `bin/gz` from `cmd/gz` (and regenerates `gz.pb.go` from `gz.proto` when needed, using vendored `bin/protoc` and `bin/protoc-gen-go` if missing). |
| `make clean` | Removes `bin/` and generated `gz.pb.go`. |
| `make test` | Runs `go test -v=true ./...`. |

You can also run the CLI with `go run ./cmd/gz …` or install/build with `go build -o gz ./cmd/gz`.

### `gz` (main binary)

Global flag on almost every command: `-r` / `--root` — repository root directory (default `.`). On `gz rebase` only, `-r` / `--root` is the root-update strategy (see that row), not the repo path.

| Command | Aliases | Arguments / flags | Purpose |
|--------|---------|-------------------|---------|
| `gz` | — | (none; prints help) | Entrypoint. |
| `gz alias` | — | `<branch> <alias>…` (at least one alias) | Record aliases for a branch. |
| `gz checkout` | `co` | `<branch>` | Check out a branch. |
| `gz create` | `+`, `push` | `<name>` (cobra allows a second positional argument; only the first is used) | Create a stacked feature branch. Flags: `-f` / `--from` parent branch; `-a` / `--alias` (repeatable). Prints branch JSON on success. |
| `gz commit` | `save` | — | Commit current changes into the branch via `git commit -a` (with `--amend` when the branch already has commits and `--append` is not set). Flags: `-a` / `--append` append a new commit instead of amending; `-m` / `--message`; `-e` / `--edit` omits `--no-edit` on amend so git can open the editor for the message. Prints branch JSON on success. |
| `gz rebase` | — | — | Rebase the current branch. Flag: `-r` / `--root` strategy: `nothing` (default), `fetch-and-rebase`, or `rebase`. |
| `gz reset` | — | — | Reset current branch to the parent’s `HEAD`. Prints branch JSON on success. |
| `gz store` | — | — | Debug/admin; running without a subcommand shows help. |
| `gz store get` | — | — | Print the current branch record from the internal DB as JSON (no output if missing). |
| `gz store edit` | — | — | Edit the current branch in the DB via editor; prints JSON after. |
| `gz completion` | — | — | Cobra shell completion: subcommands `bash`, `fish`, `powershell`, `zsh`. |

### Other programs in the repo

| Location | How to run | Notes |
|----------|------------|--------|
| `etc/example-squash` | `go run ./etc/example-squash [-path dir]` | Small demo that builds a sample git repo (default path `repo`); not part of the `gz` CLI. |
