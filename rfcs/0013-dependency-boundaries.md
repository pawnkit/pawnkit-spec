---
rfc: 0013
title: Repository dependency boundaries
status: draft
created: 2026-07-23
updated: 2026-07-23
supersedes: null
superseded-by: null
schema: null
---

## Summary

This RFC defines allowed dependency directions between PawnKit repositories.
Shared owners do not import their consumers, and orchestration stays at the
edge of the graph.

## Motivation

PawnKit tools share project loading, syntax, analysis, API facts, diagnostics,
and runtime code. Unrestricted imports can create cycles or move a fact into
the first consumer that needs it.

Release testing must also reject local replacements and unpublished module
versions before a downstream tag is created.

## Current behavior

The current Go graph is acyclic. `pawnkit-core`, `pawn-parser`,
`pawn-project`, `pawn-analysis`, `pawn-api`, and `goamx` act as shared owners.
Tools consume those modules, while the CLI, LSP, and editor compose tools near
the edge.

This direction is documented in repository architecture notes and workspace
rules, but CI does not validate it as one organisation-wide contract.

## Proposal

PawnKit uses these dependency layers:

1. `foundation`: `pawnkit-core`;
2. `syntax-project-runtime`: `pawn-parser`, `pawn-project`, and `goamx`;
3. `facts-analysis`: `pawn-api` and `pawn-analysis`;
4. `tools`: formatter, linter, test, migration, documentation, debug, and
   plugin-host libraries;
5. `adapters`: `pawnlsp`, `pawnkit-cli`, `pawn-actions`, `vscode-pawn`, and
   `pawnserver`;
6. `publication`: `pawnkit.dev`.

A repository MAY depend on its own layer or a lower layer only when the import
does not create a cycle. The following ownership rules are stricter than the
layer order:

- `pawnkit-core` has no PawnKit runtime dependency;
- `pawn-parser` may depend on core but not project, analysis, or API data;
- `pawn-project` may depend on core but not parsing or analysis;
- `goamx` does not depend on project or language tooling;
- `pawn-analysis` may depend on parser and core, but project state is supplied
  as input;
- `pawn-api` may use parser and core for extraction, but it does not depend on
  analysis or consumers;
- tools consume canonical owners and do not import other tools merely to copy
  private behavior;
- adapters may compose tools but must not become an alternative owner;
- publication reads released owner data and contains no canonical copy.

An interface belongs with the consumer that needs it unless it is itself a
versioned cross-repository contract. Implementations remain with the owner of
the underlying behavior.

`pawnkit-spec` is a contract repository, not a Go runtime dependency.
Implementations consume tagged schemas, examples, or generated bindings.
`pawn-corpus` is test evidence and MUST NOT become a production dependency.

Test-only dependencies MAY point from a lower layer to a higher-level
integration harness when they do not enter a published module's runtime
graph. Imported fixtures retain their licence, source commit, and expected
result metadata.

Generated data MAY be committed in a consumer only when runtime constraints
prevent loading the owner artifact. It MUST record the owner version and
generation command, and CI MUST reject stale output. Hand-edited forks are not
allowed.

Every released Go module MUST use public tags. A release MUST fail when a
PawnKit dependency uses:

- a local `replace` directive;
- an unpublished or missing tag;
- a pseudo-version;
- a commit different from the named tested release set;
- a dependency edge that violates this RFC;
- a cycle in the organisation graph.

Equivalent checks apply to npm packages, downloaded tools, Actions, and
bundled data. Development workspaces may use local replacements, but release
checks reject them.

The tested release set SHOULD record the direct PawnKit dependency graph and
the commit that produced it. Validators compare the graph with tagged module
metadata before publication.

## PawnKit extensions

None.

## Compatibility impact

- [x] Additive
- [ ] Breaking

The current public module graph already follows these directions. A repository
with a new reversed edge must move the behavior to its owner or add a reviewed
contract before release.

## Alternatives considered

Relying on Go's cycle detection covers packages in one build, not repository
ownership, Node packages, Actions, tools, or generated data.

Allowing any acyclic graph would still permit a shared library to depend on an
editor or CLI.

Putting all shared behavior in `pawnkit-core` would make core large and couple
unrelated domains.

## Security considerations

Release checks treat module metadata, archives, generated data, and downloaded
manifests as untrusted. They must bound input, require immutable HTTPS sources,
and verify hashes where the package system does not.

Dependency reports must redact credentials from repository URLs and local
configuration.

## Migration plan

CI first reports violations without blocking development branches. A
repository must remove release-blocking edges before its next tagged release.
Moving a public API between repositories follows the normal RFC and
deprecation process.

## Reference implementation status

Open. `pawn-actions` and the release-set validator own organisation graph
checks. Individual repositories keep their native dependency checks.

## Conformance tests

Open. Fixtures need valid layered graphs plus cycles, reversed edges, local
replacements, pseudo-versions, unpublished tags, and stale generated data.

## Open questions

- Should same-layer tool imports require an explicit allow-list?
- How should generated bindings record their owner when one output combines
  several tagged sources?
