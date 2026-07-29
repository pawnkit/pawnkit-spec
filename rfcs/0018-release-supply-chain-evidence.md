---
rfc: 0018
title: Release supply-chain evidence
status: accepted
created: 2026-07-29
updated: 2026-07-29
supersedes: null
superseded-by: null
schema: schemas/pawn-release-set-v3.schema.json
---

## Summary

This RFC adds SBOM and build-attestation records to release-set artifacts.

## Motivation

PawnKit release workflows publish checksums, SBOMs, and GitHub artifact
attestations. The tested release set records only the archive checksum, so a
consumer cannot locate or verify the other evidence from the set.

## Current behavior

Release-set v1 records an immutable archive URL, size, and SHA-256 checksum.
Several tool releases also publish a CycloneDX SBOM and a signed GitHub
artifact attestation, but their location and workflow identity are not part of
the shared document.

## Proposal

Release-set v3 artifacts contain:

- `sbom`: an immutable release asset URL, byte size, and SHA-256 checksum;
- `provenance`: the GitHub repository, workflow identity, and artifact subject
  digest covered by the attestation.

The SBOM URL must belong to the same repository and tag as the artifact.
`provenance.subject` must equal the artifact checksum. Consumers verify the
attestation signature and signer identity with GitHub's artifact-attestation
verification, not by trusting this record alone.

Every v3 artifact contains both records. Earlier schema versions remain
unchanged.

## PawnKit extensions

None.

## Compatibility impact

- [x] Additive
- [ ] Breaking

Version 3 requires both fields. Versions 1 and 2 are unchanged.

## Alternatives considered

Embedding SBOMs would make the release set large and duplicate release assets.
Recording only an attestations page would not bind the evidence to an artifact
digest or workflow identity.

## Security considerations

Consumers must verify archive and SBOM hashes before use. Attestation
verification must validate the signature, transparency-log entry, repository,
and workflow identity. URLs and downloaded files remain untrusted and subject
to the release-set size limits.

## Migration plan

Generators must add evidence before publishing a v3 set. Version 1 and version
2 consumers continue to use their existing formats.

## Reference implementation status

`pawn-actions` verifies archive and SBOM hashes, then uses GitHub's attestation
verifier to check the signature, signer workflow, source tag, source commit,
and transparency-log record.

## Conformance tests

`pawnkit-spec/examples/pawn-release-set-v3` covers the schema. The
`pawn-actions/releaseset` suite checks mismatched subjects, release URLs, and
SBOM checksums.

## Open questions

None.
