---
rfc: 0006
title: Server bundle
status: accepted
created: 2026-07-17
updated: 2026-07-18
supersedes: null
superseded-by: null
schema: schemas/pawn-bundle.schema.json
---

## Summary

This RFC defines the bundle format used to install and update a SA-MP or
open.mp server. A bundle contains the server files, manifest, checksums,
configuration references, and runtime requirements.

## Motivation

A server installation combines a binary, gamemode, filterscripts, native
extensions, and configuration. There is no standard manifest that says which
files belong together or verifies them before installation. `pawnserver` needs
that contract for installation and updates.

## Current behavior

Neither SA-MP nor open.mp publishes a bundle manifest. Server archives rely on
file-placement conventions and a hand-edited `config.json` or `server.cfg`.
The open.mp `config.json` controls the runtime; it does not describe or verify
an installation. This RFC references that configuration instead of copying it
into another format.

## Proposal

`schemas/pawn-bundle.schema.json` defines a bundle manifest
(`pawn-bundle.json`, conventionally at a bundle's root):

- `schemaVersion` (integer).
- `name`, `version` (semver): bundle identity.
- `runtimeProfile` (string): RFC 0001 profile ID (`samp-037` or `openmp`)
  the bundle targets.
- `server` (object): `binary` (`{ path, checksum, platform }` array, one
  per supported platform, per the shared baseline's supported-platform
  list), `version` (upstream server version string).
- `entryPoints` (object): `gamemode` (path to the compiled `.amx` shipped
  in this bundle), `filterscripts` (array of paths), matching the naming
  open.mp's own `config.json` `pawn.main_scripts`/`pawn.side_scripts` use,
  so a bundle's entry points and its embedded `config.json` cannot silently
  disagree (`pawnserver` validation, not a schema-level cross-check, since
  JSON Schema cannot easily assert cross-document consistency here: see
  Open questions).
- `plugins`/`components` (array): `{ name, path, checksum, platform }`,
  matching open.mp's `legacy_plugins`/component distinction.
- `configuration` (object): `path` (relative path to the bundle's
  `config.json`), `schema` (fixed reference to
  `https://schemas.pawnkit.dev/openmp-config/v1/schema.json` for open.mp
  targets; `null`/omitted for `samp-037` targets using `server.cfg`, which
  this RFC does not schema-ize since it is a flat key=value format, not
  JSON: see `language/implementation-defined.md`).
- `services` (array, optional): auxiliary processes the bundle expects
  running (e.g. a MySQL instance for a gamemode dependency): `{ name, kind,
  required }`.
- `migrations` (array, optional): ordered persistence migration steps:
  `{ id, description, appliesToVersion }`; this RFC defines only the
  manifest-level record that migrations exist and their identity/order, not
  a migration execution engine (out of scope, `pawnserver`'s job).
- `persistence` (object, optional): `{ paths: [string] }`: directories/
  files that must survive an update (player data, SQLite databases).
- `health` (object, optional): `{ checkCommand` or `checkPort`, `timeout`
  in seconds`}` describing how to determine the bundle is serving.
- `checksum` (string): checksum of the bundle manifest's own canonical
  content, for update-integrity verification (the manifest checksumming
  itself, distinct from per-file checksums above).

## PawnKit extensions

The bundle is a PawnKit format. Its reference to open.mp `config.json` uses the
upstream format as-is; RFC 0007 documents that file separately.

## Compatibility impact

- [x] Additive (no existing consumer needs to change to keep working).

This is the first version of the bundle format.

## Alternatives considered

- **Embed the full `config.json` inline in the bundle manifest** instead of
  referencing it by path: rejected; it would force every bundle manifest
  change to also touch runtime configuration, and would duplicate a schema
  already owned separately (RFC 0007 / `openmp-config.schema.json`).
- **Model `server.cfg` (SA-MP legacy config) as a JSON Schema by treating it
  as key-value pairs**: rejected for this RFC; `server.cfg` is a flat
  whitespace-delimited text format, not JSON, and shoehorning it into JSON
  Schema would validate a *transcoded* representation, not the actual file
  bytes a `samp-037` bundle ships. Left as an open question below rather
  than forcing a schema that doesn't describe the real artifact.

## Security considerations

- Bundle installation handles untrusted archives. `pawnserver` MUST apply
  traversal protection and size and recursion limits when unpacking a bundle.
  This schema requires `checksum` fields so an
  installer can verify before extracting further.
- `persistence.paths` and `server.binary[].path` are relative paths this
  schema cannot itself sandbox; implementations MUST reject paths that
  escape the bundle root.
- Bundles containing native plugins/components should be treated as
  requiring the same process-isolation caution as `pawn-plugin-host`
  applies elsewhere in the ecosystem (baseline 6.10, "keep native plugins
  out of the main process when isolation is possible"): noted here as
  guidance for `pawnserver`, not enforced by the schema.

## Migration plan

Not applicable: this is the first version.

## Reference implementation status

`pawnserver` is the reference implementation. It validates paths, checksums,
platform selection, archives, configuration consistency, and transactional
installation against this contract.

## Conformance tests

`examples/pawn-bundle/valid.json` is validated by `tools/validate`.
Archive, checksum, path, platform, and installation cases are owned by
`pawnserver`.

## Open questions

None for version 1. Resolved decisions:

- Cross-document consistency (bundle
  `entryPoints.gamemode` matching the referenced `config.json`'s
  `pawn.main_scripts`) cannot be expressed in JSON Schema because Draft
  2020-12 cannot compare separate documents. `pawnserver` enforces it.
- `server.cfg`-based (`samp-037`) bundles do not get a JSON Schema because it
  would describe a converted object, not the text file the server reads.
- SHA-256 is mandatory in version 1. Every checksum uses the
  `sha256:<lowercase hex>` form.
