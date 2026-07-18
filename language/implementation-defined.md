# Implementation-defined and unknown behavior

This file defines the vocabulary tools use for compatibility advice and the
evidence required to settle compiler behavior.

Unknown behavior stays unknown until there is evidence. Do not turn a likely
answer into a requirement.

## Classification

| Class | Meaning | Example | Expected tool behavior |
|---|---|---|---|
| Invalid | The active compiler profile rejects it. | An unbalanced `#if` and `#endif`. | Report a parser or compiler error. |
| Valid but unsafe | It compiles but has a known correctness risk. | Reinterpreting a `Float:` cell as an integer through a tag-mismatch assignment. | `pawnlint` warns under `recommended` and errors under `strict`. |
| Valid legacy | It compiles and was common before open.mp and modern includes. | Manually indexed player-data arrays. | `legacy` and `samp-037` accept it by default. Other profiles may suggest a migration. |
| Recommended modern | It is the preferred form for open.mp and `omp-stdlib`. | Using callback hooks supplied by modern includes. | `recommended` and `strict` prefer it. Legacy profiles do not require it. |
| Style preference | It has no compatibility or correctness effect. | Brace placement and spacing. | `pawnfmt` owns it. `pawnlint` does not treat it as language law. |

A construct has one class within a profile. Its class may differ between
profiles. RFC 0001 uses profiles for this reason.

## Migration

A migration from unsafe or legacy code to a modern form must be incremental,
previewable, and reversible. This specification does not require a particular
interface.

RFCs and schemas that define fixes must keep enough information for review and
reversal. RFC 0004 does this with edits that contain a range and replacement
text. A tool must not hide the change behind an in-place mutation.

## Compiler evidence

The language behavior previously listed as open is covered by the pinned probes
in [`conformance/compiler`](../conformance/compiler/README.md). New uncertainty
should stay explicit until source or a repeatable test settles it.

### RFC questions

Each RFC has its own "Open questions" section for unresolved contract design.
This file does not repeat those entries.

## Resolving a question

Use one of these sources:

1. A minimal reproduction run against a named `pawn-lang/compiler` or
   `openmultiplayer/compiler` version.
2. The relevant compiler source at a recorded commit.

Record the version or commit. Then move the confirmed behavior into the
relevant language document or add a conformance fixture.

Do not silently remove an open question. Resolve it with evidence or leave it
listed.
