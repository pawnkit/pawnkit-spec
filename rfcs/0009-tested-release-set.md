---
rfc: 0009
title: Tested release set
status: experimental
created: 2026-07-23
updated: 2026-07-23
supersedes: null
superseded-by: null
schema: schemas/pawn-release-set.schema.json
---

## Summary

This RFC defines one machine-readable set of PawnKit versions, release
artifacts, schemas, targets, test evidence, and known limits.

## Motivation

PawnKit Actions, the CLI, VS Code extension, server tooling, and website select
versions independently. A release can therefore pass its own tests while an
official consumer installs a combination that was never tested together.

The release set gives those consumers one pinned input. It also gives users a
short answer to "which versions work together?"

## Current behavior

Repositories pin their direct dependencies and managed tools in local files.
Release workflows publish checksums, but no shared document connects those
checksums to a cross-tool test run.

There is no current PawnKit format for a tested set.

## Proposal

`schemas/pawn-release-set.schema.json` defines version 1.

A release set MUST contain:

- `schemaVersion`: `1`;
- `id`: an immutable identifier;
- `generatedAt`: the UTC completion time;
- `source`: the repository and commit that produced the set;
- `targets`: the operating-system and architecture pairs tested;
- `profiles`: the Pawn profiles tested;
- `components`: exact repository versions and commits;
- `schemas`: exact schema URLs and SHA-256 hashes;
- `evidence`: the workflow, commit, projects, targets, and completion time;
- `knownLimits`: confirmed limits that affect the named set.

Downloadable component artifacts MUST use immutable HTTPS release URLs and
record their target, byte size, and SHA-256 hash. Components without release
artifacts MAY omit `artifacts`.

Every version and commit MUST exist publicly before the set is published. The
set generator MUST reject local replacements, PawnKit pseudo-versions, missing
artifacts, checksum mismatches, duplicate component names, and untested
targets.

Official PawnKit consumers SHOULD select their default versions from the latest
compatible tested set. A user MAY override a version, but the consumer must not
describe that combination as tested.

Published sets are immutable. A corrected set gets a new `id`.
Checked-in sets live under `release-sets/<id>.json`. The website MAY publish
the newest accepted set at `/release-sets/latest.json`.

## PawnKit extensions

None.

## Compatibility impact

- [x] Additive
- [ ] Breaking

Not applicable; this is the first version.

## Alternatives considered

Keeping versions in each consumer leaves compatibility implicit and makes
updates easy to miss.

Using a Git tag alone does not record artifacts, hashes, schemas, targets, or
test evidence.

Reading the newest release from GitHub would produce an unreviewed combination
and would make builds change without a repository update.

## Security considerations

Consumers must treat the document and downloaded artifacts as untrusted.
Implementations must bound document and download size, require HTTPS, reject
duplicate names and targets, verify byte size and SHA-256 before extraction,
and keep existing archive traversal checks.

The set contains public release metadata only. Workflows must not add tokens,
private URLs, logs, or environment values.

## Migration plan

Not applicable; this is the first version.

## Reference implementation status

In progress in `pawn-actions`. Checked-in sets live in `pawnkit-spec`.

## Conformance tests

The offline `pawnkit-spec` validator checks schema examples. Acceptance also
requires `pawn-actions` tests for missing releases, stale hashes, duplicate
components, unsupported targets, and unavailable artifacts.

## Open questions

None.
