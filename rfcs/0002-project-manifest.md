---
rfc: 0002
title: Project manifest (pawn.json / pawn.yaml)
status: experimental
created: 2026-07-17
updated: 2026-07-17
supersedes: null
superseded-by: null
schema: schemas/pawn-project.schema.json
---

## Summary

This RFC provides a versioned schema for sampctl's `pawn.json` and `pawn.yaml`
manifests. It also reserves a `pawnkit` section for profiles, tool settings,
and tests. Existing sampctl manifests remain valid.

## Motivation

PawnKit tools need the same answers: where the project starts, how it builds,
and what it depends on. A shared schema keeps each tool from growing its own
slightly different manifest reader.

## Current behavior

sampctl's package definition (`pawn.json`, or equivalently `pawn.yaml`) is
documented in `Southclaws/sampctl` `docs/package-definition-reference.md`
and the project wiki. As observed there:

### Core fields

| Field | Type | Required | Purpose |
|---|---|---|---|
| `entry` | string | recommended | The `.pwn` file to compile. |
| `output` | string | recommended | Destination path for the compiled `.amx`. |
| `user` | string | recommended | GitHub username/org that owns the package. |
| `repo` | string | recommended | GitHub repository name. |
| `dependencies` | array of strings | no | Packages needed to build and run. |
| `dev_dependencies` | array of strings | no | Packages needed only for development/testing. |
| `preset` | string | no | Selects runtime/compiler defaults; observed values `samp`, `openmp`. |
| `local` | bool | no, default `false` | Build/run inside the project folder instead of a temporary runtime folder. |
| `include_path` | string | no | Subfolder containing additional `.inc` sources. |
| `resources` | object | no | Extra platform-specific files/archives (documented as primarily relevant to plugin packages). |
| `extract_ignore_patterns` | array of strings | no | Patterns excluded during archive extraction. |
| `contributors` | array | no | Package metadata. |
| `website` | string | no | Package metadata (project URL). |
| `experimental.build_file` | bool | no, default `true` | Whether sampctl generates `sampctl_build_file.inc` with version/platform/git metadata. |

### Dependency string format

Dependencies are `user/repo` strings with optional version pinning (source:
`Southclaws/sampctl` `docs/dependencies.md`):

- `user/repo:1.2.3`: pins a released tag.
- `user/repo@branch-name`: targets a branch.
- `user/repo#<sha1>`: locks to an exact commit.
- Prefixed schemes select install target: `plugin://` (binaries into
  `./plugins/`), `component://` (open.mp components into `./components/`),
  `includes://` (adds include search paths), `filterscript://`
  (filterscript resources).

### `build` / `builds`

A single `build` object or a `builds` array (each entry requires `name`).
Observed shape:

```json
{
  "build": {
    "args": ["-d3", "-;+"],
    "compiler": {
      "site": "github.com",
      "user": "sampctl",
      "repo": "compilers",
      "version": "3.10.10"
    }
  }
}
```

`builds` entries additionally support a `constants` object for preprocessor
`-D`-style definitions.

### `runtime` / `runtimes`

A single `runtime` object or a `runtimes` array (each entry requires `name`,
and typically `port`). This is effectively `server.cfg` re-expressed as
JSON/YAML, plus a few sampctl-only control fields. Fields observed on the
sampctl wiki's Runtime Configuration Reference (non-exhaustive field list;
see the wiki for the authoritative source):

`version`, `endpoint`, `mode`, `extra` (object, e.g. plugin-specific
settings), `gamemodes` (array), `rcon_password`, `announce`, `maxplayers`,
`port`, `lanmode`, `query`, `rcon`, `logqueries`, `stream_rate`,
`stream_distance`, `sleep`, `maxnpc`, `onfoot_rate`, `incar_rate`,
`weapon_rate`, `chatlogging`, `timestamp`, `bind`, `password`, `hostname`,
`language`, `mapname`, `weburl`, `gamemodetext`, `filterscripts` (array),
`plugins` (array), `nosign`, `logtimeformat`, `messageholelimit`,
`messageslimit`, `ackslimit`, `playertimeout`, `minconnectiontime`,
`lagcompmode`, `connseedtime`, `db_logging`, `db_log_queries`,
`conncookies`, `cookielogging`. Field names map 1:1 to `server.cfg` /
environment variables (`SAMP_<UPPER_SNAKE_NAME>`).

sampctl synchronizes project state after manifest edits via `sampctl
ensure`; there is no committed lockfile format documented for sampctl
itself (see RFC 0003, which is therefore a PawnKit addition, not a
formalization of an existing sampctl file).

## Proposal

`schemas/pawn-project.schema.json` accepts every field enumerated above,
using the same names and semantics, as **optional** (matching sampctl's own
"recommended," not "required," posture for most fields: a manifest missing
`user`/`repo` is still valid sampctl input for a purely local project). The
schema also defines PawnKit extensions (see
below), a way to reference RFC 0001 profiles, RFC 0007 tool configuration,
and `pawntest` test declarations.

### Field mapping from `preset` to profiles

`pawn-project` SHOULD map `preset: "samp"` to the `samp-037` profile ID and
`preset: "openmp"` to the `openmp` profile ID (RFC 0001) when no explicit
`pawnkit.profile` extension field is present, so existing sampctl manifests
get a sensible profile without edits.

### Include/source roots

sampctl's documented shape only names a single `entry` file and an optional
`include_path` string. Real multi-directory projects often need multiple
include roots. The schema accepts `include_path` as sampctl defines it
(string) **and** an additive `pawnkit.includePaths` (array of strings) for
projects that need more than one; a `pawn-project` implementation MUST
treat `include_path` as shorthand for a single-element `includePaths` list
when both could apply, and MUST NOT require `pawnkit.includePaths` to be
present.

## PawnKit extensions

All extension fields live under an optional top-level `pawnkit` object so
they can never collide with a sampctl-defined field name, and a sampctl-only
implementation can ignore the entire object safely. Extensions in this
initial version:

- `pawnkit.schemaVersion` (integer, required if `pawnkit` object present):
  the major version of `pawn-project.schema.json` used by the extensions.
- `pawnkit.profile` (string): an explicit RFC 0001 profile ID, overriding
  the `preset`-based default mapping above.
- `pawnkit.includePaths` (array of strings): additional include roots
  beyond `include_path`/the project root.
- `pawnkit.tests` (object): `pawntest` test entry points and options; see
  RFC 0007 for how this interacts with a standalone `pawntest.json`.
- `pawnkit.tool` (object): free-form per-tool settings keyed by tool name
  (e.g. `pawnkit.tool.pawnlint`), formalized in RFC 0007's precedence rules.
- `pawnkit.lockfile` (string, default `"pawn.lock"`): path to the RFC 0003
  lockfile, in case a project relocates it.

No extension field is required. A manifest with no `pawnkit` object at all
is a valid sampctl-shaped manifest.

## Compatibility impact

- [x] Additive (no existing consumer needs to change to keep working).

A manifest that satisfies sampctl's documented shape today validates
against `pawn-project.schema.json` without modification. This is the first
version of the PawnKit-side schema; there is no prior PawnKit manifest
schema to be compatible with.

## Alternatives considered

- A separate `pawnkit.json` would force existing sampctl projects to migrate
  to renamed fields.
- Top-level extension fields could collide with future sampctl fields. All
  PawnKit additions therefore live under `pawnkit`.

## Security considerations

- `dependencies`/`dev_dependencies` strings reference remote repositories
  and, via `plugin://`/`component://`, native binaries. The schema itself
  only validates shape; it does not fetch anything. Implementations (not
  this repository) MUST verify checksums where available and MUST NOT
  execute project scripts merely to read the manifest (shared engineering
  baseline 6.10).
- `include_path`/`pawnkit.includePaths` are relative path strings; the
  schema constrains them to strings but cannot prevent path traversal by
  itself: `pawn-project` MUST reject `..`-escaping paths per the shared
  baseline's archive/path traversal protection requirement. Noted here so
  the requirement travels with the schema.

## Migration plan

Not applicable: this is the first version, and it is designed to accept
existing sampctl manifests unmodified.

## Reference implementation status

`pawn-project` implements JSON and YAML loading, validation, profile mapping,
include paths, and the optional `pawnkit` extensions. sampctl is not required
to understand the extension object.

## Conformance tests

Examples under `examples/pawn-project` are validated by `tools/validate`.
`pawn-project` carries manifest fixtures and schema-conformance tests.

## Open questions

- Is `resources` (platform-specific files/archives) fully specified anywhere
  beyond "applies to plugin packages," or does its shape vary by package in
  practice? Available sampctl documentation describes it only at a high level.
  The schema keeps `resources` open pending review of sampctl's `rook` and
  `types` packages.
- Does sampctl enforce any required fields at all (vs. failing later at
  build/run time with a clearer error)? This RFC assumes "no required
  fields" based on the reference table being labeled informally
  ("required"/"recommended"/"optional" is this RFC's own paraphrase of the
  source table's structure, not a verbatim requirement level from sampctl).
- Does `pawn.yaml` support every `pawn.json` field with the same names? This
  has not been checked against sampctl's YAML loader.
