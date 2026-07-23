---
rfc: 0016
title: Schema versioning and migration
status: draft
created: 2026-07-23
updated: 2026-07-23
supersedes: null
superseded-by: null
schema: null
---

## Summary

This RFC defines how PawnKit identifies, publishes, changes, and retires JSON
schemas.

## Motivation

Schemas are shared by command-line tools, editor clients, Actions, and the
website. A changed file at an existing URL can break old releases without
changing their binaries.

## Current behavior

Schemas use draft 2020-12 and versioned `$id` URLs under
`https://schemas.pawnkit.dev`. Documents usually carry an integer
`schemaVersion`. The spec validator checks examples offline, and the website
publishes checked-in schemas. Migration and support rules are not yet shared.

## Proposal

Every published schema has:

- a stable name;
- an integer major in its `$id`, using
  `https://schemas.pawnkit.dev/<name>/v<major>/schema.json`;
- a matching document `schemaVersion`;
- a SHA-256 hash recorded by tested release sets.

A published major URL is immutable. Fixing wording outside validation rules is
allowed only through a new checked-in spec release; the served schema bytes
must still match the hash recorded for a tested release.

The following changes are additive within a draft major:

- adding an optional property;
- widening an enum or accepted value range;
- adding a new definition that existing documents do not use.

The following require a new major after a schema is accepted:

- adding a required property;
- removing or renaming a property;
- narrowing accepted values;
- changing a field's meaning or units;
- changing unknown-property handling.

Producers MUST emit one declared major. Readers MUST reject unknown majors
clearly and SHOULD read the current and previous supported major during a
migration. The tested release set records the required majors and hashes.

Each new major needs valid examples, failure examples, a migration note, and a
compatibility reader or `pawnmigrate` rule where projects store the document.

Schema URLs are published only from `pawnkit-spec`. CI must compare the served
bytes with the checked-in file and fail on a missing URL, unexpected redirect,
or hash mismatch.

## PawnKit extensions

None.

## Compatibility impact

- [x] Additive
- [ ] Breaking

This policy does not change an existing schema. Accepted schemas gain immutable
major URLs and a support window.

## Alternatives considered

Using package versions in schema URLs was rejected because package and contract
releases move independently. Serving only an unversioned latest schema was
rejected because old tools could receive incompatible validation rules.

## Security considerations

Schema fetching must use HTTPS, bounded responses, and pinned hashes in release
automation. Validation must not resolve arbitrary network references. Schema
errors must not echo secrets from input documents.

## Migration plan

Before acceptance, CI will inventory published schema URLs and hashes. Any
accepted schema that changed incompatibly in place must publish a new major and
retain the previous one.

## Reference implementation status

`pawnkit-spec/tools/validate` checks local schemas and examples. Remote
immutability and migration checks remain open.

## Conformance tests

The offline validator, valid and invalid examples, and a CI URL/hash check form
the conformance suite.

## Open questions

- How long should the previous accepted major remain supported?
- Should unversioned convenience URLs redirect or be omitted?
