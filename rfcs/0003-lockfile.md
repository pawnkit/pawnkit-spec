---
rfc: 0003
title: Lockfile (pawn.lock)
status: draft
created: 2026-07-17
updated: 2026-07-18
supersedes: null
superseded-by: null
schema: schemas/pawn-lock.schema.json
---

## Summary

This RFC defines the PawnKit contract for `pawn.lock`. The file records exact
dependency sources, revisions, checksums, graph edges, platform artifacts, and
the compiler and runtime selected during resolution.

## Motivation

`pawn-project` needs a reviewable record of dependency resolution so local and
CI builds use the same graph.

## Current behavior

sampctl added `pawn.lock` in 1.13.0. The version 1 format remains current in
1.14.0 and includes dependency resolution, integrity, runtime, and build data.
PawnKit MUST read that format without breaking sampctl projects. Its own fields
need a namespace or a versioned migration rather than an incompatible file with
the same name.

## Proposal

The current PawnKit schema predates sampctl's lockfile and is not compatible
with the version 1 format in sampctl 1.14.0. The field list below records the
earlier draft.
It MUST NOT move to `experimental` until this RFC adopts the upstream shape or
defines namespaced extensions that sampctl safely ignores.

`schemas/pawn-lock.schema.json` defines a `pawn.lock` document with:

- `schemaVersion` (integer): major version of this schema.
- `generatedAt` (RFC 3339 timestamp): informational, not used for cache
  invalidation logic by the schema itself.
- `manifestChecksum` (string, e.g. `sha256:<hex>`): checksum of the
  `pawn.json`/`pawn.yaml` this lockfile was resolved from, so tools can
  detect drift.
- `compiler` (object): `vendor`, `version`, `checksum` of the compiler used
  at resolution/build time. Ties to RFC 0001's profile `compiler` shape.
- `runtimeProfile` (string): RFC 0001 profile ID describing the runtime
  this lockfile assumes.
- `packages` (array): one entry per resolved dependency:
  - `name` (string): the `user/repo` identifier from RFC 0002.
  - `resolved` (string): the exact ref sampctl-style strings can point at:
    tag, branch, or commit as resolved (e.g. `user/repo#abc123...`).
  - `version` (string, optional): semantic version if the upstream tag is
    semver-shaped.
  - `commit` (string): resolved commit SHA. Required: a lockfile entry
    without an exact commit is not reproducible.
  - `source` (object): `type` (`git`, `archive`, `local`), `url`.
  - `checksum` (string, e.g. `sha256:<hex>`): of the fetched artifact,
    required when `source.type` is `archive`; recommended otherwise.
  - `kind` (string enum): `dependency`, `dev-dependency`, `plugin`,
    `component`, `includes`, `filterscript`, mirroring RFC 0002's dependency
    scheme prefixes so a lockfile entry can be traced back to its manifest
    declaration.
  - `platformArtifacts` (array, optional): per-platform binary info for
    `plugin`/`component` kinds: `platform` (e.g. `linux-x86_64`), `url` or
    `path`, `checksum`.
  - `dependencies` (array of strings, optional): names of other `packages`
    entries this one depends on, forming the resolved dependency graph as
    edges between `name` values already present in this document.

### Example shape

```json
{
  "schemaVersion": 1,
  "generatedAt": "2026-07-17T00:00:00Z",
  "manifestChecksum": "sha256:1111111111111111111111111111111111111111111111111111111111111c",
  "compiler": { "vendor": "openmultiplayer", "version": "3.10.11" },
  "runtimeProfile": "openmp",
  "packages": [
    {
      "name": "pawn-lang/YSI-Includes",
      "resolved": "pawn-lang/YSI-Includes@5.x",
      "commit": "2222222222222222222222222222222222222a",
      "source": { "type": "git", "url": "https://github.com/pawn-lang/YSI-Includes" },
      "kind": "dependency"
    }
  ]
}
```

## PawnKit extensions

The extension design is unresolved. sampctl owns the existing top-level
version 1 fields: `version`, `generated`, `sampctl_version`, `dependencies`,
`runtime`, and `build`. PawnKit-specific data must not reuse those names with
different meanings.

## Compatibility impact

- [ ] Additive
- [x] Breaking

The current schema rejects sampctl 1.14.0 lockfiles and assigns different
meanings to the same filename. That conflict must be fixed before adoption.

## Alternatives considered

- Reusing `pawn.json` resources would mix declared intent with resolved state.
- A foreign format such as Cargo-style TOML would add another parser without
  improving compatibility with sampctl's JSON lockfile.

## Security considerations

- `checksum` fields exist specifically so implementations can verify
  downloaded artifacts (shared baseline 6.10, "verify download checksums
  when available"). The schema marks `checksum` required for `archive`
  sources precisely to make skipping verification a validation error, not
  merely a missed best practice.
- `platformArtifacts[].path`/`url` are untrusted strings; implementations
  MUST apply the same path-traversal and archive-extraction protections as
  for manifest paths (see RFC 0002 security considerations).
- A lockfile is otherwise inert data; validating it requires no network
  access.

## Migration plan

Not applicable: this is the first version.

## Reference implementation status

`pawn-project` implements the earlier PawnKit draft. It does not yet satisfy
the sampctl 1.14.0 compatibility requirement.

## Conformance tests

Examples under `examples/pawn-lock` are validated by `tools/validate`.
`pawn-project` tests loading, graph validation, and deterministic output.

## Open questions

- Should PawnKit adopt sampctl's version 1 shape directly, contribute required
  fields upstream, or store PawnKit-only metadata under a namespaced object?
- Should `pawn.lock` be required to be deterministic/reproducible byte-for-
  byte (ignoring `generatedAt`), or only field-for-field (allowing key
  reordering)? This RFC does not mandate key ordering; JSON Schema cannot
  express byte-level determinism, so this is a `pawn-project`
  implementation concern to resolve, recorded here so it isn't silently
  assumed.
- Should the lockfile record the resolved `dialectFlags` snapshot from RFC
  0001 inline, or only the `runtimeProfile` ID (relying on the profile
  document at build time)? This RFC currently only stores the ID
  (simpler, but not self-contained if a profile document changes
  retroactively). Open for discussion before this RFC can move to
  `experimental`.
