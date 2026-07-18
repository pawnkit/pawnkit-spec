---
rfc: 0001
title: Target language and runtime profiles
status: experimental
created: 2026-07-17
updated: 2026-07-18
supersedes: null
superseded-by: null
schema: schemas/pawn-profile.schema.json
---

## Summary

This RFC defines machine-readable profiles for Pawn dialects, runtimes, and
lint policies. PawnKit tools use the same profile IDs instead of interpreting
"targeting open.mp" independently.

## Motivation

"Pawn v3" does not identify a compiler lineage, runtime, or policy. Without a
shared profile document:

- `pawn-parser` and `pawn-analysis` may disagree about accepted syntax.
- `pawnlint` has no stable way to say that a rule applies to `strict` but not
  `legacy`.
- `pawn-project` cannot validate that a manifest's declared target is one
  PawnKit actually supports.
- `pawnserver` cannot know which AMX/host API behavior a bundle assumes.

## Current behavior

Target selection is usually implicit and tool-specific.

- sampctl has a `preset` field with values including `samp` and `openmp`. It
  selects sampctl defaults, not a complete language contract.
- pawn-lang/compiler is a maintained fork of Pawn 3.2.3664. Its compatibility
  mode restores selected behavior from the older compiler.
- open.mp uses the same compiler lineage. The pinned compiler probes in
  [`conformance/compiler`](../conformance/compiler/README.md) record behavior
  shared by current pawn-lang and open.mp builds.

Neither project publishes a versioned schema of accepted syntax, directives,
and diagnostics. Consumers otherwise have to infer those facts from compiler
source and release notes.

## Proposal

### Profile kinds

A profile document has a `kind`, one of:

- `source` describes syntax, preprocessing, tags, operators, and directives.
  `pawn-parser` and `pawn-analysis` consume it.
- `runtime` describes AMX and host behavior, including native availability,
  calling conventions, and numeric behavior. `goamx` and `pawnserver` consume
  it.
- `policy` changes diagnostic severity without changing validity. A policy
  profile MUST declare a `basedOn` profile ID.

A project may use `openmp` syntax with a `samp-037`-compatible runtime shim
during migration. One version string cannot describe that combination. See
`docs/compatibility.md` for the profile status table.

### Profile identity

Every profile document has:

- `id`: a stable lowercase kebab-case string such as `openmp`. IDs MUST be
  unique within `profiles/`.
- `schemaVersion`: the major version of `pawn-profile.schema.json`.
- `version`: the profile document's own semantic version.
- `kind`: `source`, `runtime`, or `policy`.
- `displayName` and `description`.
- `compiler` (object, `source`/`runtime` kinds only): `lineage` (e.g.
  `"pawn-3.2.3664"`), `vendor` (`"pawn-lang"` or `"openmultiplayer"`),
  `minVersion`/`maxVersion` (informational, not enforced by the schema
  beyond their string format), and `notes`.
- `dialectFlags`: named scalar values fixed by the profile, such as
  `automaticIncludeGuards` and `tagMismatchWarningLevel`. Profiles may add
  flags as compiler behavior is confirmed, so the schema does not fix the
  property names.
- `basedOn` (string, `policy` kind only): the profile ID this overlays.

### The five shipped profiles

| ID | Kind | Summary |
|---|---|---|
| `legacy` | source | Original Pawn 3.2.3664 syntax/semantics as commonly targeted by pre-2015 SA-MP gamemodes, without pawn-lang/open.mp compiler-specific fixes. Existing scripts are assumed *valid but not recommended* under this profile, not modernized. |
| `samp-037` | source + runtime | The dialect accepted by the pawn-lang/compiler-based toolchains bundled with SA-MP 0.3.7 server releases, plus the SA-MP 0.3.7 AMX runtime/native surface. |
| `openmp` | source + runtime | The dialect accepted by the open.mp-bundled compiler (3.10.11 lineage) plus the open.mp AMX runtime/native surface, including `omp-stdlib` additions. |
| `recommended` | policy, `basedOn: openmp` | Modern-practice defaults: warns on legacy patterns (see `language/implementation-defined.md` for the valid/unsafe/legacy/modern/style distinction from principle 3.7) without rejecting valid legacy code. |
| `strict` | policy, `basedOn: openmp` | Superset of `recommended` that upgrades most `recommended` warnings to errors and adds stricter tag/tag-mismatch enforcement. Intended for new projects, not migration targets for legacy codebases. |

### Consumption contract

- `pawn-project` manifests (RFC 0002) reference a profile by `id` in the
  `profiles`/`build` sections; an unknown ID is a manifest validation error.
- `pawnlint` rule metadata may declare "applies under profile X" using the
  same IDs.
- `pawn-analysis`/`pawn-parser` use `dialectFlags` to decide which
  productions/behaviors to accept without needing this repository's prose,
  once the flag is confirmed (not an open question).

## PawnKit extensions

The profile format and the `recommended` and `strict` policies are specific to
PawnKit. sampctl's `preset` field only selects its own defaults.

`pawn-project` SHOULD map `preset: samp` to `samp-037` and `preset: openmp` to
`openmp` when a manifest does not select a profile explicitly. RFC 0002 defines
that mapping.

## Compatibility impact

- [x] Additive

This is the first PawnKit profile format.

## Alternatives considered

- A single version string cannot express separate source and runtime targets or
  name a policy overlay.
- Keeping profile facts only in `pawn-parser` would make them unavailable to
  other implementations and editors.

## Security considerations

Profile documents are static data. They contain no executable content, network
references, or filesystem paths.

## Migration plan

Not applicable: this is the first version.

## Reference implementation status

`pawn-project` loads and validates profile IDs. The profile documents and
schema are available in this repository. Full use of every `dialectFlags`
entry by `pawn-parser` and `pawn-analysis` is still incomplete.

`pawn-parser` also has `ProfileLossless`, `ProfileAnalysis`, and
`ProfileTokensOnly`. Those select parser output, not a language dialect.

## Conformance tests

`tools/validate` checks every profile document against the schema. The fixtures
under `conformance/compiler` cover the language facts documented here. New
flags need a source reference or a compiler probe.

## Open questions

- Is the pawn-lang compiler's bug-fix list stable enough to represent as
  `dialectFlags`?
- Should a `library` profile exist for code meant to run unmodified across
  multiple runtime profiles (e.g. a published include library targeting
  both `samp-037` and `openmp`)?
