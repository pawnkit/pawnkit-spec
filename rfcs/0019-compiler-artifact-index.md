---
rfc: 0019
title: Compiler artifact index
status: draft
created: 2026-07-30
updated: 2026-07-30
supersedes: null
superseded-by: null
schema: schemas/pawn-compiler-index.schema.json
---

## Summary

This RFC defines a signed-off list of Pawn compiler downloads. Each entry pins
the compiler, host target, archive layout, byte size, and SHA-256 hashes needed
to install it without guessing from a release page.

## Motivation

`pawn-project` can verify and cache a compiler when a caller supplies an
artifact URL and checksum. The CLI, Actions, and editor do not have a shared
source for those values, so they either require a local `pawncc` or maintain
their own download tables.

PawnKit needs one reviewed index before it can acquire compilers consistently
across the terminal, CI, and editor.

## Current behavior

Projects may select a compiler vendor and version through RFC 0002. sampctl
lockfiles record a compiler version and preset but do not pin the compiler
archive URL or checksum.

The compiler probe workflow downloads Pawn 3.10.10 from a fixed GitHub release
URL and checks a hard-coded SHA-256 value. `pawn-project/toolchain` accepts the
same facts through `ResolveOptions`, but no public PawnKit contract supplies
them.

## Proposal

`schemas/pawn-compiler-index.schema.json` defines version 1.

An index contains:

- `schemaVersion`: `1`;
- `id`: an immutable lowercase identifier;
- `generatedAt`: when the entries were reviewed;
- `artifacts`: one or more compiler artifacts.

Each artifact contains:

- `vendor` and exact `version`;
- the PawnKit `profiles` for which it is a valid compiler;
- `target`, describing the host operating system and architecture;
- `source`, naming the upstream repository, release tag, and commit;
- `archive`, with an immutable HTTPS URL, format, byte size, and SHA-256 hash;
- `executable`, with its path inside the archive, architecture, and SHA-256
  hash.

The tuple `(vendor, version, target)` MUST be unique. Paths use forward slashes,
must be relative, and must not contain `.` or `..` segments.

Consumers MUST select an exact vendor, version, and target. They MUST verify the
index checksum supplied by their tested PawnKit release before reading it,
verify the archive before extraction, and verify the compiler executable after
extraction. A cache hit is valid only while the executable hash still matches.

Index documents are immutable. A correction gets a new `id`. The checked-in
documents live under `compiler-indexes/`, and the website publishes them at an
immutable URL containing that ID. A tested release set names the index ID and
SHA-256 hash used by its compatibility run.

## PawnKit extensions

This is a PawnKit format. It does not change compiler releases, sampctl
manifests, or sampctl lockfiles.

## Compatibility impact

- [x] Additive
- [ ] Breaking

Existing projects may continue to use a configured compiler or `pawncc` from
`PATH`. Managed compiler acquisition is an additional resolution step.

## Alternatives considered

Reading the latest GitHub release makes builds change without a project or
tool update. Deriving asset names from a version cannot authenticate the
download and breaks when upstream packaging changes.

Keeping separate tables in Actions and VS Code would repeat the current
version-selection problem.

Adding archive data to `pawn.lock` would make sampctl compatibility harder and
would copy platform-wide release facts into every project.

## Security considerations

The index and archives are untrusted input. Implementations must bound their
size, require HTTPS, reject redirects to non-HTTPS URLs, reject duplicate
coordinates, and validate paths before extraction. Archive and executable
hashes must be checked before the compiler runs.

An index checksum authenticates the selected document; the hashes inside it do
not replace that outer check. Download failures must leave the previous cache
entry intact.

## Migration plan

Not applicable; this is the first version.

## Reference implementation status

Planned for `pawn-project/toolchain`. `pawnkit-cli`, `pawn-actions`, and
`vscode-pawn` will consume that implementation rather than maintaining their
own compiler tables.

## Conformance tests

The offline `pawnkit-spec` validator checks valid and invalid index examples.
Implementation tests still need to cover duplicate coordinates, target
selection, size and hash failures, archive traversal, interrupted updates, and
offline cache use.

## Open questions

- Should compiler indexes be referenced by release-set v3 through an additive
  field, or should that reference wait for the next release-set major?
- Which upstream open.mp package should be the source of its compiler artifact:
  the compiler release or the complete server bundle?
