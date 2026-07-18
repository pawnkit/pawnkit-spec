---
rfc: 0007
title: Tool configuration discovery
status: draft
created: 2026-07-17
updated: 2026-07-17
supersedes: null
superseded-by: null
schema: schemas/pawnlint.schema.json, schemas/pawnfmt.schema.json, schemas/pawntest.schema.json
---

## Summary

This RFC defines where PawnKit tools find project settings and how settings
from a manifest and a standalone file are merged. It also defines the schemas
for `pawnlint`, `pawnfmt`, and `pawntest` configuration.

## Motivation

`pawnlint`, `pawnfmt`, and `pawntest` each need per-project configuration.
Without a shared discovery rule, each tool would invent its own precedence
between "settings embedded in `pawn.json`" and "a dedicated config file,"
producing inconsistent results for projects that use more than one tool.

## Current behavior

`pawn-project` implements the shared discovery and merge rules. sampctl has no
equivalent linter, formatter, or test configuration contract, so these rules
are specific to PawnKit.

## Proposal

### Discovery order

For a given tool `<name>` (e.g. `pawnlint`) and a project root determined by
`pawn-project`'s manifest discovery:

1. **Standalone file, closest to the file being processed, wins for
   file-scoped overrides.** Tools MAY support directory-scoped standalone
   files (e.g. a `pawnlint.json` in a subdirectory overriding the root one)
   using the same merge rule as step 3 below, applied per directory level
   from root to leaf.
2. **Project-root standalone file** (e.g. `<projectRoot>/pawnlint.json`,
   `<projectRoot>/.pawnfmt.json`, `<projectRoot>/pawntest.json`) is read if
   present.
3. **Manifest-embedded settings** (`pawnkit.tool.<name>` in `pawn.json`/
   `pawn.yaml`, RFC 0002) are read if present.
4. **Tool built-in defaults** apply for anything not set by 1-3.

### Merge rule

When both a standalone file and manifest-embedded settings exist for the
same tool, they are **merged**, not one replacing the other outright:

- Scalar and enum fields: the standalone file's value wins if both set the
  same field; the manifest's value is used only if the standalone file
  omits that field.
- Array fields (e.g. rule lists, ignore globs): concatenated, standalone
  file's entries first, then manifest entries, with de-duplication of
  identical entries: **unless** the field is documented by that tool's own
  schema as "replace" (e.g. a `rules` object keyed by rule ID naturally
  merges key-by-key; a bare ignore-glob array concatenates). Each of
  `pawnlint.schema.json`/`pawnfmt.schema.json`/`pawntest.schema.json`
  documents, per field, whether it merges or replaces, using an
  `x-pawnkit-merge` annotation (`"replace"` or `"concat"`) in the schema so
  tools can discover the rule without interpreting this prose.
- This "standalone wins on conflict, both contribute where non-conflicting"
  rule exists so a contributor can override one setting locally
  (standalone file, often gitignored or personal) without losing team-wide
  defaults committed in the manifest.

### The three standalone schemas

- `schemas/pawnlint.schema.json`: `schemaVersion`, `extends` (array of
  named or path-based rule-set references), `rules` (object keyed by rule
  ID, each `{ severity: "error"|"warning"|"off", options?: object }`),
  `ignore` (array of glob strings, `x-pawnkit-merge: concat`), `profile`
  (RFC 0001 profile ID override for which dialect rules assume).
- `schemas/pawnfmt.schema.json`: `schemaVersion`, `indent` (`{ style:
  "space"|"tab", size: integer }`), `lineWidth` (integer), `braceStyle`
  (enum), `ignore` (array of glob strings, `x-pawnkit-merge: concat`).
  Kept intentionally small: `pawnfmt` must stay lossless-syntax-driven and
  fast per principle 3.5, not grow semantic configuration.
- `schemas/pawntest.schema.json`: `schemaVersion`, `entries` (array of
  test-suite entry `.pwn` files or globs), `runtimeProfile` (RFC 0001
  profile ID), `timeoutSeconds` (integer, per-test default), `reporters`
  (array of enum `"human"|"json"|"junit"`).

## PawnKit extensions

Not applicable in the RFC 0002 sense: this is a new format with no external
precedent.

## Compatibility impact

- [x] Additive (no existing consumer needs to change to keep working).

This is the first version of the tool configuration contract.

## Alternatives considered

- Letting the standalone file replace all manifest settings would discard team
  defaults when a contributor changes one local option.
- Letting the manifest always win would make standalone configuration mostly
  useless.
- One schema for all tools would mix unrelated lint, formatting, and test
  settings. Callers already know which tool owns the file, so separate schemas
  are simpler.

## Security considerations

- Standalone config files and manifest-embedded settings are project-
  authored, not project-*script*, content: reading them must not require
  executing the project (shared baseline 6.10). None of the three schemas
  defines an executable field (no shell-out hooks); if a future revision
  adds one, it MUST document sandboxing/opt-in requirements explicitly.
- `ignore`/glob fields are patterns, not paths to execute; no traversal
  risk beyond normal file-matching, but implementations should still bound
  match-expansion cost (shared baseline "expansion... limits").

## Migration plan

Not applicable: this is the first version.

## Reference implementation status

`pawn-project` implements shared configuration discovery and merging.
`pawnlint`, `pawnfmt`, and `pawntest` consume their own validated settings.

## Conformance tests

Schema examples are validated by `tools/validate`. Merge precedence and
tool-specific settings are tested in their owning repositories.

## Open questions

- Should directory-scoped standalone files (step 1) be required of every
  tool, or optional per-tool (e.g. `pawnfmt` might reasonably support this
  for formatting-style-per-directory, while `pawntest` entry discovery may
  not need it)? Current draft makes it a MAY; revisit once `pawnlint`
  exists and has a real use case pushing on this.
- Exact glob dialect (`.gitignore`-style vs. full glob with `**`) is not
  pinned down yet; recorded here rather than assumed, since it affects both
  schema pattern validation and every consuming tool's matching library
  choice.
