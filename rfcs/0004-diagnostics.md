---
rfc: 0004
title: Diagnostics and edits
status: experimental
created: 2026-07-17
updated: 2026-07-17
supersedes: null
superseded-by: null
schema: schemas/pawn-diagnostic.schema.json
---

## Summary

This RFC defines the diagnostic and text-edit format used by PawnKit CLIs,
lint baselines, and editor integrations.

## Motivation

`pawnlint`, `pawnfmt`, `pawn-analysis`, `pawntest`, and `pawn-language-server` all
produce findings about source code. Without one shared shape:

- Baselines (suppress known findings) written for one tool would not
  compose with another's output.
- `pawn-language-server` would need per-tool adapters instead of one translation to the
  Language Server Protocol.
- SARIF exporters would need a different source shape for each tool.

## Current behavior

PawnKit tools now exchange the shared diagnostic types from `pawnkit-core`.
They may keep richer internal findings, but serialized diagnostics, CLI JSON,
and baselines use this RFC's shape.

## Proposal

`schemas/pawn-diagnostic.schema.json` defines a `Diagnostic` object:

- `code` (string, stable, namespaced e.g. `pawnlint:no-goto`): never
  silently repurposed (shared baseline 6.9).
- `source` (string): the tool/subsystem that produced it (e.g. `pawnlint`,
  `pawn-analysis`, `pawnfmt`).
- `severity` (enum): `error`, `warning`, `info`, `hint`.
- `message` (string): human-readable, imperative-neutral description.
- `range` (object): primary file + range: `file` (string, project-relative
  path), `start`/`end` (`{ "line": int, "column": int, "offset": int }`,
  1-based line/column, 0-based byte offset, matching common LSP + tooling
  convention: offsets are UTF-8 byte offsets per `language/lexical.md`'s
  source-encoding rules).
- `relatedLocations` (array, optional): same range shape plus a `message`
  each, for secondary locations ("previous declaration here").
- `notes` (array of strings, optional): supplementary non-fix text.
- `help` (string, optional): actionable guidance text.
- `documentationUrl` (string, optional): link to rule/behavior docs.
- `tags` (array of strings, optional): e.g. `deprecated`, `unnecessary`,
  matching LSP `DiagnosticTag` naming where overlapping, for cheap
  translation.
- `fixes` (array of `Edit`, optional): **safe** fixes, appliable without
  review.
- `unsafeFixes` (array of `Edit`, optional): fixes requiring review (e.g.
  behavior-changing).
- `suppression` (object, optional): `suppressed` (bool), `reason` (string),
  `mechanism` (string, e.g. `inline-comment`, `baseline-file`).

An `Edit` object:

- `file` (string).
- `range` (same shape as above).
- `newText` (string).
- `version` (integer or string, optional): document version this edit
  applies to, for version-aware application and staleness detection (shared
  baseline 6.4, "version-aware text edits with overlap validation").

Multi-file edit sets are represented as an array of `Edit` under a
diagnostic's `fixes`/`unsafeFixes`; a tool applying them MUST validate that
edits do not overlap within the same file and MUST support preview
(showing the diff) before transactional application, per the shared
baseline: this RFC defines the data shape, not the application algorithm,
which belongs to the applying tool (e.g. `pawnkit-cli`, `pawn-language-server`).

### Exit codes

Tools reporting these diagnostics as CLI findings SHOULD use the shared
exit-code classes from the shared engineering baseline (section 6.3):
`0` success, `1` findings present, `2` invalid invocation, `3` toolchain
failure, `4` internal failure. This RFC does not re-specify those classes;
it only requires that a `Diagnostic` with `severity: error` corresponds to
at least one finding causing a non-zero exit under class `1`.

## PawnKit extensions

Not applicable in the RFC 0002 sense: there is no pre-existing external
diagnostic format this formalizes (SARIF is a possible *export* target, not
a source format PawnKit tools already emit). The `tags` field is
deliberately named to match LSP `DiagnosticTag` values where they overlap,
to ease that translation, but this is a design choice, not compatibility
with a pre-existing PawnKit format.

## Compatibility impact

- [x] Additive (no existing consumer needs to change to keep working).

This is the first version of the shared diagnostic format.

## Alternatives considered

- SARIF is too verbose for the internal format. PawnKit exports SARIF when a
  consumer needs it.
- Per-tool diagnostic shapes would multiply the formats every consumer must
  understand.

## Security considerations

- `message`/`notes`/`help` are free text and MUST NOT be assumed safe to
  render as HTML without escaping by consumers (e.g. `pawnkit.dev`); this
  schema does not itself sanitize content.
- Diagnostic bundles containing file paths/messages from untrusted project
  content MUST have secrets redacted by producers before publishing logs or
  bundles, per the shared baseline; this schema has no dedicated secret
  field and diagnostic producers are responsible for not leaking one into
  `message`.

## Migration plan

Not applicable: this is the first version.

## Reference implementation status

`pawn-analysis`, `pawn-parser`, and `pawnlint` produce shared diagnostic
types through `pawnkit-core`. `pawn-language-server` translates them to LSP.

## Conformance tests

The schema example is validated by `tools/validate`. `pawnkit-core` freezes
wire-format fixtures, while producing tools test their own diagnostic codes.

## Open questions

- Should `range.start`/`end` use UTF-16 code units instead of byte offsets
  for the LSP-facing translation layer (LSP conventionally uses UTF-16)?
  This RFC specifies UTF-8 byte offsets as the canonical interchange
  representation and leaves UTF-16 conversion to the language-server adapter.
  `pawn-parser` and `pawn-analysis` operate on UTF-8 source bytes.
- Exact list of standard `tags` values beyond `deprecated`/`unnecessary`
  (the two given as examples in the shared engineering baseline) is not
  yet enumerated; recorded as an open item rather than inventing a closed
  list prematurely.
