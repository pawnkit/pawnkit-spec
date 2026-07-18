---
rfc: 0000
title: <Short, specific title>
status: draft
created: YYYY-MM-DD
updated: YYYY-MM-DD
supersedes: null
superseded-by: null
schema: null
---

<!--
Copy this file to rfcs/NNNN-short-title.md and use the next free number.
Keep every section. A draft may mark implementation and conformance work as
open. GOVERNANCE.md defines the RFC lifecycle and review periods.
-->

## Summary

What does this RFC define, and why is it shared across repositories?

## Motivation

What problem does it solve? Name the tools or users currently blocked.

## Current behavior

Describe what happens without this RFC. Cite compiler source, existing formats,
or other primary documentation where possible.

If the RFC formalizes an existing format, document that format before adding
PawnKit behavior.

## Proposal

Define the fields or behavior, related schemas, and examples.

Use MUST, SHOULD, and MAY only when the distinction affects interoperability.

## PawnKit extensions

List additions to an existing external format separately. Extensions MUST be
optional and additive unless this RFC introduces a breaking major version with
a migration plan.

Write "None" when there are no extensions.

## Compatibility impact

- [ ] Additive
- [ ] Breaking

Explain what an existing consumer must change. For a new format, write "Not
applicable; this is the first version."

## Alternatives considered

Record the alternatives and why they were rejected. Include doing nothing when
it is a realistic option.

## Security considerations

Cover untrusted input, resource limits, path traversal, secret redaction, and
network access where relevant. If the shared engineering rules are sufficient,
say why.

## Migration plan

For a breaking change, describe manual steps, a `pawnmigrate` rule, or a
compatibility reader. Support the current and previous major version during the
transition, as required by `docs/compatibility.md`.

For a new format, write "Not applicable; this is the first version."

## Reference implementation status

Name the implementing repository and current status. This may remain open in a
draft but is required before acceptance.

## Conformance tests

Name the conformance files or implementation test suite. This may remain open
in a draft but is required before acceptance.

## Open questions

List questions that can be settled with a source reference or repeatable test.
Do not present an unverified answer as fact.
