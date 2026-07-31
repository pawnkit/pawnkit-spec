---
rfc: 0022
title: Dependency resolution and lock generation
status: experimental
created: 2026-07-31
updated: 2026-07-31
supersedes: null
superseded-by: null
schema: schemas/pawn-lock.schema.json
---

## Summary

This RFC defines how PawnKit resolves project dependencies and writes a
sampctl-compatible `pawn.lock`. It covers exact revisions, transitive
dependencies, conflicts, and repeatable installs.

## Motivation

`pawn install` can restore an existing lockfile, but a new project still needs
sampctl to create that file. PawnKit cannot replace sampctl as the main
workflow until it can resolve `pawn.json` and `pawn.yaml` itself.

The resolver belongs in `pawn-project`. The CLI, Actions, language server, and
other consumers must not choose revisions independently.

## Current behavior

sampctl 1.14.1 accepts tag (`:`), branch (`@`), and commit (`#`) references. It
records each selected commit in a version 1 lockfile and discovers transitive
dependencies from the selected package manifest.

PawnKit implements the same manifest references and lock shape through RFCs
0002 and 0003. It can create, reconcile, update, and restore this graph.

## Proposal

`pawn-project` will expose one resolver for supported project manifests. Given
a manifest, an optional existing lock, and a resolution mode, it produces the
complete version 1 dependency object.

Normal install mode MUST reuse an existing entry when its source identity and
constraint still match the manifest. It MUST resolve entries that are missing
or stale. Updating a valid locked revision requires an explicit update mode.

The resolver applies references as follows:

- `#commit` resolves that commit.
- `@branch` resolves the current branch tip when resolution is required.
- `:tag` resolves that exact tag.
- an unqualified dependency resolves the provider's default branch.

Every remote result records the full commit ID and
`integrity:commit:<full-commit>`. `resolved` records the tag, branch, abbreviated
commit, or `HEAD` form required by the sampctl version 1 contract.

After selecting a revision, the resolver reads that revision's package
manifest and adds its dependencies. Development dependencies are included only
for the root project. Transitive entries set `transitive` and record every
direct parent in `required_by`.

Package identity includes its scheme and source repository. If two graph paths
request the same identity with incompatible constraints, resolution MUST fail
and report both paths. It MUST NOT silently choose one constraint.

Resolution order and output order MUST be deterministic. Network completion
order must not affect the lockfile.

The lockfile writer MUST replace the dependency object recoverably while
preserving supported runtime, build, and PawnKit extension data. Resource
records whose package or commit changed MUST be removed and resolved again.
For a new lock, `sampctl_version` records the compatible format baseline
`1.14.1`; it does not imply that sampctl created the file.

## PawnKit extensions

None. This RFC writes the sampctl version 1 fields defined by RFC 0003 and the
optional resource extension defined by RFC 0021.

## Compatibility impact

- [x] Additive
- [ ] Breaking

Existing sampctl lockfiles remain valid. Normal install mode preserves valid
locked revisions, so opening a project with PawnKit does not update dependency
versions.

## Alternatives considered

Calling sampctl would keep two dependency engines in the main workflow and
would not complete the replacement goal.

Resolving independently in each command would let build, lint, CI, and the
editor disagree about the dependency graph.

Always selecting the newest revision would make normal installs change over
time and weaken lockfile review.

## Security considerations

Repository URLs and manifests are untrusted. The resolver must use bounded
process output and manifest reads, reject credentials and unsupported URL
schemes, limit graph depth and package count, and detect cycles.

Package scripts are not executed. Revisions are verified after checkout.
Lockfile replacement must be recoverable and must reject links and non-regular
targets.

Credentials may be supplied to providers but must not appear in URLs,
diagnostics, lockfiles, or command output.

## Migration plan

No file migration is required. Projects with a sampctl lock keep using it.
Projects without a lock can create the same version 1 shape through
`pawn install`.

## Reference implementation status

`pawn-project` contains the graph resolver, drift check, and lock writer.
`pawnkit-cli` provides GitHub revision lookup. `pawn install` creates missing
locks, reconciles changed direct dependencies, and refreshes the full graph
when passed `--update`.

## Conformance tests

`pawn-project/dependency/resolution_test.go` covers constraints, transitive and
development dependencies, cycles, conflicts, matching lock reuse, explicit
updates, and deterministic output. `pawn-project/lockfile/dependencies_write_test.go`
covers lock creation and preservation. `pawnkit-cli` tests GitHub responses,
offline lock reuse, stale-lock reconciliation, updates, and clean
manifest-only installation.

## Open questions

- Should update mode later accept selected package names?
- Which non-GitHub providers should the first implementation support?
