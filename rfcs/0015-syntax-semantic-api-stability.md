---
rfc: 0015
title: Syntax and semantic API stability
status: experimental
created: 2026-07-23
updated: 2026-07-25
supersedes: null
superseded-by: null
schema: null
---

## Summary

This RFC defines stability rules for the public Go APIs of `pawn-parser` and
`pawn-analysis`, including syntax kinds, symbols, references, source maps, and
semantic diagnostics.

## Motivation

Linting, LSP features, migrations, and documentation depend on these packages.
Unlabelled changes can silently alter diagnostics or force several repositories
to update at once.

## Current behavior

Both modules are versioned and tested by downstream tools. They expose public
Go packages, but they do not label individual surfaces as stable, preview, or
internal. Compatibility is inferred from release notes and compilation.

## Proposal

Public surfaces use the lifecycle terms from RFC 0010:

- **stable** surfaces follow semantic versioning and retain the previous major
  during a documented migration window;
- **preview** surfaces may change in a minor pre-1.0 release and must be named
  in release notes;
- **internal** surfaces stay under Go `internal` packages and carry no
  compatibility promise.

The stable parser surface covers token kinds, lossless source access, node
ranges, parse diagnostics, and documented tree traversal. Numeric syntax-kind
values are not wire identifiers unless a separate schema says so.

The stable analysis surface covers request options, symbols, references,
source mappings, diagnostics, and cancellation. Symbol identity must remain
deterministic for the same project revision. Diagnostic codes are governed by
RFC 0004.

Adding an optional field or a new named kind is additive. Removing or
reinterpreting a field, changing range units, changing identity rules, or
reusing a diagnostic code is breaking.

Consumers SHOULD test the smallest public surface they use. Cross-repository
tests MUST use tagged modules and MUST NOT use local replacements.

## PawnKit extensions

None.

## Compatibility impact

- [x] Additive
- [ ] Breaking

This RFC labels existing APIs. It does not change their current data.

## Alternatives considered

Freezing every exported identifier was rejected because some packages are
still experimental. Treating all pre-1.0 releases as unstable was rejected
because downstream tools need a usable compatibility promise.

## Security considerations

API stability does not make parsed input trusted. Stable request APIs must
retain cancellation and resource limits. Diagnostics must not expose source
outside the requested project.

## Migration plan

Repositories document stable and preview packages before this RFC is accepted.
A breaking stable change needs release notes, downstream updates, and a
supported compatibility path.

## Reference implementation status

`pawn-parser` labels its stable public packages and `pawn-analysis` labels its
preview public packages. Both modules compile a small supported-surface test
from an external test package.

## Conformance tests

Each module has an external-package test for its supported entry points and
types. The tagged downstream build matrix and tested release set record the
versions used.

## Open questions

- Should syntax-kind names receive a language-neutral interchange schema?
