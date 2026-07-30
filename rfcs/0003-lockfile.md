---
rfc: 0003
title: Lockfile (pawn.lock)
status: draft
created: 2026-07-17
updated: 2026-07-30
supersedes: null
superseded-by: null
schema: schemas/pawn-lock.schema.json
---

## Summary

PawnKit uses sampctl's version 1 `pawn.lock` format. The lockfile records exact
dependency revisions and may also record runtime files and build output.

## Motivation

Local builds and CI need the same dependency graph. Existing projects already
have a lockfile format, so PawnKit should read it instead of creating another
contract with the same filename.

## Current behavior

sampctl 1.14.1 reads and writes version 1 lockfiles in
`src/pkg/package/lockfile`. The top-level fields are:

- `version`
- `generated`
- `sampctl_version`
- `dependencies`
- optional `runtime`
- optional `build`

Dependencies are stored in an object keyed by their resolved package identity.
Each value records its constraint, resolved revision, commit, source identity,
scheme, integrity value, and reverse dependency edges.

An earlier PawnKit draft used `schemaVersion` and `packages`. That shape is not
sampctl-compatible and must not be written by PawnKit tools.

## Proposal

PawnKit readers MUST accept sampctl version 1 lockfiles. PawnKit writers MUST
emit that shape without renaming fields or changing their meaning.

`dependencies` is an object whose values contain:

- `constraint`, `resolved`, `commit`, `user`, and `repo`
- optional `integrity`, `site`, `path`, `branch`, and `local`
- optional `transitive` and `required_by`
- optional `scheme`: `plugin`, `component`, `includes`, or `filterscript`

`integrity` accepts sampctl's `sha256:<hex>` and `commit:<sha>` forms.
Implementations MUST verify a present integrity value before using installed
content.

`runtime` records the runtime version, platform, type, and installed files.
Each file records its path, size, hash, and mode.

`build` may record the compiler version and preset, entry file, output file,
and output hash.

Object keys in `dependencies` are identities, not filesystem paths. Readers
must use the dependency fields and scheme to choose an install target.

### Earlier PawnKit draft

Readers MAY accept the earlier `schemaVersion`/`packages` shape until
2027-07-30. This is a read-only migration path. Writers MUST use the sampctl
shape.

## PawnKit extensions

There are no PawnKit-only fields in version 1. A future extension must be
additive, namespaced, ignored safely by sampctl, and marked
`x-pawnkit-extension: true` in the schema.

## Compatibility impact

- [x] Additive
- [ ] Breaking

The proposal adopts the format already written by sampctl. It also keeps a
temporary reader path for the earlier PawnKit draft.

## Alternatives considered

- Keeping the earlier PawnKit shape would make the same filename mean two
  different things.
- A separate PawnKit lockfile would leave sampctl projects with competing
  dependency graphs.
- Reconstructing locked state from `pawn.json` would lose exact revisions and
  integrity data.

## Security considerations

- Treat dependency identities, local paths, runtime files, and build paths as
  untrusted input.
- Reject absolute paths, traversal, and archive entries outside their target.
- Verify `integrity` before using installed dependencies.
- Do not execute package scripts while reading or installing dependencies.
- Bound downloads, extracted file counts, individual files, and total output.
- A missing integrity value is compatible with existing lockfiles but should
  produce a clear audit finding.

## Migration plan

1. Update the schema and `pawn-project` reader to accept both shapes.
2. Make PawnKit writers emit sampctl version 1.
3. Add `pawnmigrate` support for converting the earlier PawnKit draft.
4. Remove the legacy reader after 2027-07-30 in a release that documents the
   change.

## Reference implementation status

sampctl 1.14.1 is the source implementation for the version 1 shape.
`pawn-project` currently implements only the earlier PawnKit draft and must be
updated before native dependency installation uses the lockfile.

## Conformance tests

`examples/pawn-lock/sampctl-v1.json` covers the sampctl shape.
`examples/pawn-lock/valid.json` keeps the temporary legacy shape covered.

`pawn-project` must test both readers, reverse-edge conversion, schemes,
integrity values, local dependencies, malformed paths, and unsupported
versions.

## Open questions

- Should a future version require integrity for every remote dependency?
- Should PawnKit contribute a forward dependency list to sampctl, or continue
  deriving it from `required_by`?
