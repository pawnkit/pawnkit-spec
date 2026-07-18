# Governance

`pawnkit-spec` uses a lightweight RFC process. A contributor should be able to
read an RFC and its schema, then implement the contract without reverse
engineering another PawnKit repository.

## RFC lifecycle and states

Every RFC has exactly one status at a time, recorded in its front matter:

| Status | Meaning |
|---|---|
| `draft` | Open for discussion and may change incompatibly. |
| `experimental` | Ready for prototypes, but not yet frozen. Changes are recorded and versioned. |
| `accepted` | Stable at its current major version. Breaking changes require a replacement RFC and migration path. |
| `deprecated` | Supported for existing consumers but discouraged for new work. The RFC names a replacement when one exists. |
| `superseded` | Replaced by another RFC. Both files link to each other. |

An RFC reaches `accepted` only when both conditions are met:

1. At least one reference implementation (or a documented plan with an
   owning repository and milestone) exists for the behavior it describes.
2. Conformance tests exist in `conformance/` or in the owning
   implementation repository's test suite, and are referenced from the
   RFC's "Conformance tests" section.

The existence of an implementation repository does not make an RFC accepted.
Use the status in the RFC front matter.

## Decision process

1. Open a pull request using `rfcs/0000-template.md`. A draft may leave its
   implementation and conformance sections incomplete, but must address every
   other section.
2. Contributors from any PawnKit repository may review it. Drafts have no
   fixed quorum.
3. Breaking changes stay open for public discussion for at least 14 calendar
   days after the pull request is marked ready for review. This applies to
   accepted schema versions and normative language behavior already shipped in
   a profile. Maintainers follow the same rule. Typos, examples, and other
   non-normative edits are exempt.
4. A maintainer may merge after discussion settles. Acceptance also requires
   the implementation and conformance evidence listed above.
5. Draft and experimental RFCs may be amended in place. An accepted RFC may
   receive clarifications, but behavioral changes require a superseding RFC.

## Compatibility and migration requirements

Every RFC and schema change records the following in its compatibility and
migration sections:

- Whether the change is additive or breaking.
- For a breaking change, manual migration steps, a `pawnmigrate` rule, or a
  reader that accepts both versions for a stated period.
- The minimum schema major version consumers must support during the
  transition (current + one prior major, per the shared engineering
  baseline).

## Critical compatibility rule for the project manifest

RFC 0002 and RFC 0003 MUST remain compatible with sampctl's `pawn.json`,
`pawn.yaml`, and `pawn.lock`. They must not reuse those filenames for a
different contract. PawnKit-only fields are additive, documented in the RFC,
and marked `x-pawnkit-extension: true` in the schema.

Changing this rule requires an RFC and the full breaking-change discussion
period.

## Roles

- Editors merge pull requests, maintain RFC numbering and status metadata, and
  enforce discussion periods. They do not decide technical questions alone.
- Anyone contributing to `rfcs`, `language`, `profiles`, or `schemas` is a
  contributor.
- Repository maintainers hold merge access and apply these rules.

## Numbering

RFCs are numbered sequentially starting at `0001` (`0000` is reserved for the
template). Numbers are never reused, even if an RFC is withdrawn before
merge. A withdrawn RFC's number is marked `superseded` or `withdrawn` in the
RFC index rather than recycled.
