---
rfc: 0020
title: Server runtime artifact index
status: experimental
created: 2026-07-30
updated: 2026-07-30
supersedes: null
superseded-by: null
schema: schemas/pawn-runtime-index.schema.json
---

## Summary

This RFC defines an immutable list of server runtime downloads. Each entry
pins the upstream release, host target, archive layout, byte size, and hashes
needed to install a runtime without consulting a mutable release page.

## Motivation

PawnKit can build a project with a verified compiler, but `pawn run` still
needs a server installed by the user or sampctl. The CLI, Actions, and editor
need one source for the runtime archive and its expected contents.

## Current behavior

`pawnserver` verifies PawnKit server bundles after they have been assembled.
It does not acquire an upstream SA-MP or open.mp server. The open.mp
installation guide directs users to GitHub release archives, while the
project manifest records only a runtime version.

## Proposal

`schemas/pawn-runtime-index.schema.json` defines version 1.

An index contains:

- `schemaVersion`: `1`;
- `id`: an immutable lowercase identifier;
- `generatedAt`: when the entries were reviewed;
- `artifacts`: one or more runtime artifacts.

Each artifact contains:

- `vendor`, exact `version`, and PawnKit `profile`;
- `target`, describing the host operating system and architecture;
- `source`, naming the upstream repository, release tag, and commit;
- `archive`, with an immutable HTTPS URL, format, byte size, and SHA-256 hash;
- `root`, the runtime directory inside the archive;
- `executable`, with its archive path, architecture, and SHA-256 hash.

The tuple `(vendor, version, target)` MUST be unique. Consumers MUST verify
the outer index checksum supplied by a tested PawnKit release, then verify the
archive and executable before installation. A cache hit remains valid only
while the executable hash matches.

Index documents are immutable. Corrections get a new `id`. Checked-in indexes
live under `runtime-indexes/` and are published at an immutable URL containing
that ID.

## PawnKit extensions

This is a PawnKit format. It does not change upstream server archives,
sampctl manifests, or the bundle format from RFC 0006.

## Compatibility impact

- [x] Additive
- [ ] Breaking

Existing projects may continue to use a local server. Managed runtime
acquisition is an additional resolution step.

## Alternatives considered

Using the latest GitHub release would make an unchanged project download
different code over time. Keeping download tables in the CLI, Actions, and
editor would duplicate policy and hashes.

Embedding the runtime archive in a project lockfile was rejected because the
archive is platform-specific and shared by many projects.

## Security considerations

Indexes and archives are untrusted input. Implementations must bound their
size, require HTTPS, reject unsafe paths and links, limit extraction, and use
recoverable cache writes. The archive and selected executable must both match
their hashes before any bundled code runs.

Native components and plugins are executable code. Hash verification provides
integrity, not isolation.

## Migration plan

Not applicable; this is the first version.

## Reference implementation status

The schema and first reviewed open.mp index are implemented in
`pawnkit-spec`. Runtime selection and installation remain to be implemented
in `pawn-project`, `pawnserver`, and `pawnkit-cli`.

## Conformance tests

The offline validator checks the schema, every checked-in index, duplicate
coordinates, and an unsafe-root example. Runtime installers must test archive
traversal, links, size limits, target selection, checksum failures, executable
permissions, and recoverable updates.

## Open questions

SA-MP runtime archives are not included until PawnKit can name an authorised,
immutable distribution source and verify its redistribution terms.
