---
rfc: 0010
title: Support and lifecycle
status: draft
created: 2026-07-23
updated: 2026-07-23
supersedes: null
superseded-by: null
schema: null
---

## Summary

This RFC defines PawnKit maturity labels and the evidence required for a
support claim. It also proposes a small repository-owned support record.

## Motivation

PawnKit repositories currently use labels such as preview and experimental
without shared meanings. CI matrices imply platform support, while README
files and release notes may describe something different.

Users need to know which tools, platforms, profiles, and compiler versions
have current test coverage. Maintainers need a clear point at which a release
becomes deprecated or unsupported.

## Current behavior

Each repository chooses its own wording. The organisation audit uses stable,
preview, experimental, infrastructure, and deprecated as working labels, but
these labels are not a public contract.

The tested release set records targets and profiles. It does not state the
maturity of each component or the compiler versions covered by its evidence.

## Proposal

PawnKit uses five maturity labels:

- `stable`: public contracts are accepted and compatibility is maintained;
- `preview`: intended for users, but some public contracts are not accepted;
- `experimental`: available for testing without compatibility guarantees;
- `infrastructure`: supports PawnKit delivery and is not an end-user product;
- `deprecated`: still supported for a stated period, with a named replacement
  or removal plan.

A stable component MUST:

- use accepted contracts for its public cross-repository behavior;
- pass its required checks on every supported platform;
- state supported profiles and compiler versions;
- publish a migration path before removing supported behavior;
- follow the security and compatibility policies in this repository.

A preview component MUST pass its advertised platform checks. It MUST identify
draft or experimental contracts that can still change. A release note MUST
call out incompatible changes.

An experimental component MUST say which guarantees are missing. It MUST NOT
be required by the primary editor, terminal, or CI workflow.

Infrastructure components MUST state which published artifacts or services
they maintain. They do not inherit stable product status from those artifacts.

A deprecated component MUST name its replacement when one exists and set an
end-of-support date or release milestone. Security fixes continue until that
point. Removal after the stated period is a breaking change.

Each repository SHOULD own `.pawnkit/support.json`. Version 1 contains:

- `schemaVersion`: `1`;
- `repository`: the `pawnkit/<name>` repository;
- `maturity`: one of the five labels;
- `platforms`: tested operating-system and architecture targets;
- `profiles`: tested Pawn profile names;
- `compilers`: tested compiler names and exact versions;
- `contracts`: consumed PawnKit contracts and supported major versions;
- `limitations`: short confirmed limits;
- `deprecated`: an optional replacement and end-of-support date or milestone.

An empty array means that a dimension does not apply. It does not mean all
values are supported.

The repository CI MUST validate its support record and reject claims not
covered by configured jobs. The tested release set copies the maturity label
and relevant support dimensions from the tagged component. The website reads
owner records and release-set evidence instead of maintaining a separate
support table.

Support applies only to released versions. A branch build may be tested, but
it is not a supported release until its tag and artifacts exist.

## PawnKit extensions

None.

## Compatibility impact

- [x] Additive
- [ ] Breaking

Existing repositories can add support records without changing their public
APIs. Current maturity labels remain informal until a record is published.

## Alternatives considered

README-only support statements are easy to write but difficult to compare or
validate.

Using CI matrices alone omits profiles, compiler versions, maturity, and
deprecation plans.

Putting every support claim in `pawnkit.dev` would make the website a second
source of truth.

## Security considerations

Support records contain public metadata only. Validators must still bound file
size, reject duplicate values, and avoid network access during normal local
validation.

Security support ends only at the recorded date or milestone. A deprecated
component must not silently stop receiving fixes before then.

## Migration plan

Repositories add `.pawnkit/support.json` before their next maturity change.
Until then, existing labels are descriptive and cannot be used as release-set
support evidence.

## Reference implementation status

Open. `pawnkit-spec` owns the schema and validator. Repository CI templates in
`pawn-actions` will check support records.

## Conformance tests

Open. The schema needs valid and invalid examples. `pawn-actions` needs tests
that compare support claims with configured release targets.

## Open questions

- Should compiler support allow inclusive version ranges after exact-version
  evidence exists?
- Should a repository record support for source builds separately from release
  artifacts?
