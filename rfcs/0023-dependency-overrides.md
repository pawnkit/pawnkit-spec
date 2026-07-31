---
rfc: 0023
title: Root dependency overrides
status: draft
created: 2026-07-31
updated: 2026-07-31
supersedes: null
superseded-by: null
schema: schemas/pawn-project.schema.json
---

## Summary

This RFC lets a root project replace an obsolete transitive package source or
constraint. Overrides are explicit in the project manifest and apply during
RFC 0022 dependency resolution.

## Motivation

Old Pawn packages often keep working after their manifests become stale. A
transitive dependency may point to a deleted repository, an old fork, or a
constraint that no longer matches the root project. Today the user must fork
that package or accept an install failure.

ScavengeSurvive provides a repeatable case. `Southclaws/samp-ini` requests
`Misiur/YSI-Includes`, while the project and its current libraries use
`pawn-lang/YSI-Includes@5.x`. Both repositories install to `YSI-Includes`, so
installing both would overwrite files.

## Current behavior

RFC 0022 resolves each package identity and rejects incompatible constraints,
cycles, and install-path collisions. The sampctl manifest shape has no field
for replacing a transitive source.

## Proposal

The optional `pawnkit.dependencyOverrides` object maps one dependency identity
to a replacement dependency reference:

```json
{
  "pawnkit": {
    "schemaVersion": 1,
    "dependencyOverrides": {
      "Misiur/YSI-Includes": "pawn-lang/YSI-Includes@5.x"
    }
  }
}
```

Keys identify a provider, package, and optional resource scheme. They do not
include a version constraint. Values use the dependency-reference format from
RFC 0002.

Only the root manifest may define overrides. The resolver MUST ignore or
reject overrides found in dependency manifests. It applies the map before
resolving a matching transitive request. Direct root dependencies are already
under the project owner's control and are not rewritten.

The replacement becomes the package identity stored in `pawn.lock`. The root
manifest remains the reviewable record of the substitution. Normal lock reuse
MUST consider the active override map; changing or removing an override makes
affected lock entries stale.

An override MUST NOT change between ordinary and resource dependency schemes.
Resolvers MUST reject duplicate normalized keys, replacement cycles, and
ambiguous matches. Resolution limits from RFC 0022 still apply after
substitution.

## PawnKit extensions

Add optional `pawnkit.dependencyOverrides` to the RFC 0002 extension object.
It maps dependency identities to dependency references.

## Compatibility impact

- [x] Additive
- [ ] Breaking

Existing manifests and consumers remain valid. Consumers that do not support
the extension may ignore the `pawnkit` object, but they will not reproduce an
overridden dependency graph.

## Alternatives considered

Forking every stale package shifts maintenance to users and hides the reason
for the fork.

Allowing repositories with the same final path to overwrite each other makes
the result order-dependent and unsafe.

Treating forks as aliases is incorrect because forks can contain unrelated
code and history.

## Security considerations

Overrides use the existing credential-free dependency syntax and provider
limits. Implementations must reject cycles and scheme changes before network
access. Diagnostics must not include provider credentials or unbounded remote
output.

## Migration plan

No migration is required. Projects add an override only when a reviewed
transitive replacement is necessary. Removing one may require `pawn install
--update` to refresh the graph.

## Reference implementation status

Planned for `pawn-project` and `pawnkit-cli` after review of this draft.

## Conformance tests

Planned schema examples will cover valid source and constraint replacements,
credentials, scheme changes, and malformed identities. `pawn-project` tests
will cover substitution, cycles, direct dependency precedence, lock drift, and
deterministic output. ScavengeSurvive will provide the public integration test.

## Open questions

- Should a later version support target- or profile-specific overrides?
