---
rfc: 0004
title: Diagnostics and edits
status: experimental
created: 2026-07-17
updated: 2026-07-24
supersedes: null
superseded-by: null
schema: schemas/pawn-diagnostic-v2.schema.json
---

## Summary

This RFC defines the diagnostic and text-edit format used by PawnKit CLIs,
lint baselines, and editor integrations.

## Motivation

`pawnlint`, `pawnfmt`, `pawn-analysis`, `pawntest`, and `pawnlsp` all
produce findings about source code. Without one shared shape:

- Baselines (suppress known findings) written for one tool would not
  compose with another's output.
- `pawnlsp` would need per-tool adapters instead of one translation to the
  Language Server Protocol.
- SARIF exporters would need a different source shape for each tool.

## Current behavior

The published v1 schema and the `pawnkit-core` v1 wire format disagree. The
schema uses `range`, `relatedLocations`, and `fixes`; released core code uses
`primary`, `related`, `safeFixes`, and `reviewFixes`.

Version 1 remains published at its original URL. Version 2 matches the released
core data model and becomes the current format.

## Proposal

`schemas/pawn-diagnostic-v2.schema.json` defines a `Diagnostic` object:

- `schemaVersion` (integer): `2`.
- `code` (string, stable, namespaced e.g. `pawnlint:no-goto`): never
  silently repurposed (shared baseline 6.9).
- `source` (string): the tool/subsystem that produced it (e.g. `pawnlint`,
  `pawn-analysis`, `pawnfmt`).
- `severity` (enum): `error`, `warning`, `info`, `hint`.
- `message` (string): human-readable, imperative-neutral description.
- `primary` (object): primary URI and half-open byte span. A derived
  line/character range may also be present.
- `related` (array, optional): locations with a short message.
- `notes` (array of strings, optional): supplementary non-fix text.
- `help` (string, optional): actionable guidance text.
- `docsUrl` (string, optional): link to rule or behavior docs.
- `tags` (array of strings, optional): e.g. `deprecated`, `unnecessary`,
  matching LSP `DiagnosticTag` naming where overlapping, for cheap
  translation.
- `safeFixes` (array, optional): fixes that may be applied without review.
- `reviewFixes` (array, optional): fixes that need user review.
- `suppressed` (object, optional): suppression kind and reason.

A fix contains a message, a fixed kind, and a workspace edit. Workspace edits
group byte-span replacements by URI and may include a document version.
Consumers MUST reject overlapping edits within one document and SHOULD offer a
preview before applying a multi-file change.

Line and character values are zero-based. Their encoding is negotiated by the
transport. Byte spans are authoritative and use zero-based UTF-8 byte offsets.

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

- [x] Breaking (migration plan required).

The v1 URL and checked-in schema remain unchanged. Producers move to v2.
Consumers must accept v1 and v2 during the migration window.

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

1. Publish the v2 schema without modifying the v1 schema.
2. Update `pawnkit-core` to emit v2 and read both versions.
3. Release producers before consumers.
4. Keep v1 decoding for at least the current and previous supported core
   release.

The old schema and the released core v1 format used the same version number for
different shapes. A v1 decoder must inspect the fields: `range` identifies the
published schema, while `primary` identifies the released core format. Both map
to the v2 data model. V1 paths become URI strings, and each old fix becomes a
single-document workspace edit.

## Reference implementation status

`pawnkit-core` v0.2.0 is the reference implementation. `pawn-project` v0.3.0
and `pawnkit-cli` v1.3.0 exchange v2 diagnostics through build results.
`pawn-analysis`, `pawn-parser`, and `pawnlint` produce core diagnostic types;
`pawnlsp` translates them to LSP.

## Conformance tests

Valid and invalid v2 examples are checked by `tools/validate`. `pawnkit-core`
freezes both v1 shapes and v2. The Actions release-set smoke verifies v2 build
results against the small corpus projects.

## Open questions

- Exact list of standard `tags` values beyond `deprecated`/`unnecessary`
  (the two given as examples in the shared engineering baseline) is not
  yet enumerated; recorded as an open item rather than inventing a closed
  list prematurely.
