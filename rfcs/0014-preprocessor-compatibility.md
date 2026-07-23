---
rfc: 0014
title: Preprocessor compatibility
status: draft
created: 2026-07-23
updated: 2026-07-23
supersedes: null
superseded-by: null
schema: null
---

## Summary

This RFC defines the shared boundary for Pawn preprocessing. `pawn-parser`
recognises directive structure, while `pawn-analysis` evaluates it for a
resolved project profile.

## Motivation

Includes, conditional branches, and macro expansion affect every later result.
If linting, navigation, migration, and documentation preprocess files
differently, they disagree about the same project.

## Current behavior

`pawn-parser` retains directives and original source structure.
`pawn-analysis/preprocess` resolves includes, evaluates conditions, expands
macros, tracks active regions, and maps expanded tokens back to source.
`pawn-project` supplies roots, profiles, defines, and include paths.

The shared corpus covers common directives and several community include
patterns. The contract and supported compiler divergences are not yet recorded
in one place.

## Proposal

A preprocessing request consists of:

- a source URI and bytes;
- the resolved target profile and defines;
- an include resolver supplied by `pawn-project`;
- expansion, include-depth, output, and diagnostic limits;
- a revision identifying the project inputs.

The result contains:

- original and expanded tokens;
- the include graph;
- active and inactive source regions;
- macro definitions and expansions;
- source mappings for expanded output;
- diagnostics with original-source locations;
- limit and cancellation status.

`pawn-parser` MUST preserve directive syntax without choosing active branches.
`pawn-analysis` MUST own evaluation and macro semantics. Consumers MUST use the
analysis result rather than expanding source themselves.

Quoted includes search relative to the including file before resolved include
roots. Angle includes use resolved include roots. Platform path separators are
normalised by `pawn-project`; analysis does not guess new roots.

Expansion MUST be bounded. A limit hit returns partial data and a diagnostic,
not a panic or unbounded allocation. Cancellation MUST stop include loading and
expansion promptly.

Clean and incremental analysis MUST produce the same active regions, expanded
tokens, mappings, and diagnostics for the same request revision.

Compiler-specific behavior belongs to a named profile. An intentional
divergence must have a compiler-backed corpus fixture and compatibility note.

## PawnKit extensions

None. This RFC describes PawnKit's model of existing compiler behavior.

## Compatibility impact

- [x] Additive
- [ ] Breaking

The RFC formalises the current ownership boundary. A consumer with its own
preprocessor must migrate to `pawn-analysis`.

## Alternatives considered

Letting each tool preprocess only the syntax it needs was rejected because it
caused different include graphs and diagnostics. Moving evaluation into the
parser was rejected because active branches depend on project inputs.

## Security considerations

Include paths and source are untrusted. Resolvers must reject traversal outside
their permitted roots, and expansion must bound recursion, output, diagnostics,
and loaded files. Network access is outside this protocol.

## Migration plan

Consumers move preprocessing to `pawn-analysis` and pass project inputs from
`pawn-project`. Existing result fields remain during the pre-1.0 migration.

## Reference implementation status

`pawn-analysis/preprocess` is the reference implementation. Ownership and core
limits are implemented; the compatibility table remains open.

## Conformance tests

`pawn-corpus/preprocessor` and the compiler differential suite are the
conformance source. The RFC needs a checked-in index mapping each required
behavior to its fixture.

## Open questions

- Which compiler revisions define each named compatibility profile?
- Which warning-only compiler quirks should PawnKit reproduce?
