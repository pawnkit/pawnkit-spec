---
rfc: 0005
title: API metadata (natives, callbacks, constants, tags)
status: draft
created: 2026-07-17
updated: 2026-07-17
supersedes: null
superseded-by: null
schema: schemas/pawn-api.schema.json
---

## Summary

This RFC defines a machine-readable index of SA-MP and open.mp APIs. It covers
signatures, constants, tags, availability, deprecations, and provenance.
`pawn-api` produces the index; analysis and documentation tools consume it.

## Motivation

Every tool that needs to know "does `SetPlayerPos` exist under `samp-037`,
what are its parameter types, is it deprecated in favor of something else"
would otherwise parse `omp-stdlib` includes or maintain its own table. This RFC
defines the interchange format produced by `pawn-api`, so consumers do not need
its extraction logic.

## Current behavior

`omp-stdlib` (`openmultiplayer/omp-stdlib`) declares APIs in `.inc` files and
documents them in comments. No upstream JSON or YAML index has been recorded.
The open.mp documentation site
(`open.mp/docs/scripting/language/reference/Contents`) presents function
references as rendered documentation pages, not a downloadable schema.
PawnKit therefore maintains its own index, generated from upstream sources by
`pawn-api`.

## Proposal

`schemas/pawn-api.schema.json` defines an API metadata document as a
`schemaVersion` plus an `entries` array. Each entry:

- `id` (string, stable, namespaced, e.g. `native:SetPlayerPos`,
  `callback:OnPlayerConnect`, `const:MAX_PLAYERS`, `tag:Float`): unique
  across the document; `tools/validate` checks this.
- `kind` (enum): `native`, `callback`, `function`, `constant`, `tag`,
  `define`.
- `name` (string): the Pawn-visible identifier.
- `signature` (object, natives/callbacks/functions only): `parameters`
  (array of `{ name, tag, default? }`), `returnTag`.
- `value` (string or number, constants/defines only): literal value where
  statically known; omitted when the value is host-computed.
- `tags` (array of strings, optional): associated Pawn tags for a
  constant/native return, distinct from the `tag` API-metadata `kind`
  above (naming collision noted; see Open questions).
- `availability` (array of objects): one entry per RFC 0001 profile this
  API element is available under: `{ "profile": "openmp", "since": "...",
  "until": null }`. Absence of a profile in this array means "not available
  under that profile," which is itself meaningful (e.g. an open.mp-only
  native absent from `samp-037`).
- `deprecated` (object, optional): `{ "since": "...", "replacement":
  "native:NewName", "reason": "..." }`.
- `source` (object): provenance, required per the shared baseline's
  "derived API data must record its source and licensing constraints":
  `repository` (e.g. `openmultiplayer/omp-stdlib`), `path` (file path
  within that repository), `commit` (the commit the fact was extracted
  from), `license` (SPDX identifier of the source repository).
- `documentationUrl` (string, optional): link to rendered open.mp docs for
  this entry, when known.

## PawnKit extensions

Not applicable in the RFC 0002 sense: there is no pre-existing external
API-metadata format being formalized. Every field here is a new PawnKit
contract populated from upstream Pawn source and documentation by `pawn-api`.

## Compatibility impact

- [x] Additive (no existing consumer needs to change to keep working).

This is the first version of the API metadata format.

## Alternatives considered

- Storing full documentation inline would make the interchange file noisy and
  expensive to diff. `documentationUrl` links to long-form documentation.
- One file per API entry would complicate interchange. `pawn-api` may still use
  that layout internally before generating the combined document.

## Security considerations

- `source.repository`/`source.path` are provenance metadata, not fetched by
  schema validation; no network access is implied by validating a
  `pawn-api`-shaped document.
- Consumers must not treat `value` for constants as safe to interpolate
  into generated code without the same escaping care as any other derived
  data.

## Migration plan

Not applicable: this is the first version.

## Reference implementation status

`pawn-api` is the reference producer. It validates source entries and generates
the versioned interchange document consumed by other tools.

## Conformance tests

The schema example is validated by `tools/validate`. `pawn-api` tests generated
shape conformance and keeps source provenance for reviewed entries. Extraction
correctness remains the producer's responsibility.

## Open questions

- Confirm whether `omp-stdlib` or open.mp's documentation pipeline already
  publishes any machine-readable API table (JSON/YAML) that this RFC should
  reference or align field names with, rather than defining names in
  isolation? None is currently recorded, but the search was not exhaustive.
- Should `tags` (Pawn language tags associated with an entry) and `kind:
  "tag"` (an entry that itself documents a tag like `Float:`) share the
  `id` namespace prefix `tag:`? Current draft uses `tag:` for both the
  entry-kind and the associated-tags field's *values*, which is a naming
  collision that still needs a decision before acceptance.
- `omp-stdlib` is MPL-2.0. `pawn-api` records that licence in generated
  provenance for extracted entries.
