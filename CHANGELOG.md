# Changelog

All notable changes to this repository are documented here. Format loosely
follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Schema
major versions are tracked independently in each schema's `$id`; this log
calls out schema version changes explicitly.

## [0.1.13] - 2026-07-23

### Added

- Added draft RFC 0010 for support and lifecycle terms.
- Added draft RFC 0011 for runtime-fidelity reporting.

## [0.1.12] - 2026-07-23

### Changed

- Promoted the command-line tool candidate after its three-platform smoke test.

## [0.1.11] - 2026-07-23

### Added

- Added a bootstrap release-set candidate for the primary command-line tools.

## [0.1.10] - 2026-07-23

### Added

- Added the first tested release set for CLI project discovery.
- Added offline validation for checked-in release sets.

## [0.1.9] - 2026-07-23

### Added

- Added a release-set candidate with published CLI artifacts and schema hashes.

## [0.1.8] - 2026-07-23

### Changed

- Accepted RFC 0009 after the `pawn-actions` v1.1.0 implementation passed CI.
- Updated the README to list both accepted RFCs.

## [0.1.7] - 2026-07-23

### Changed

- Moved RFC 0009 to experimental for the `pawn-actions` prototype.

## [0.1.6] - 2026-07-23

### Added

- Draft RFC 0009 and schema for tested release sets.

## [0.1.5] - 2026-07-23

### Added

- Added the canonical generated project-manifest example.

## [0.1.4] - 2026-07-23

### Changed

- Recorded the RFC 0008 reference implementations and tests.

## [0.1.3] - 2026-07-23

### Added

- Draft RFC 0008 for editor-managed tool state.

## [0.1.2] - 2026-07-23

### Added

- Compiler checks for include order, source locations, active regions, and
  profile defines.

## [0.1.1] - 2026-07-23

### Changed

- Moved compiler probe sources to `pawn-corpus` v0.1.5.
- Added a pinned compiler differential to CI.

## [0.1.0] - 2026-07-20

### Added

- Initial governance process and RFC template (`GOVERNANCE.md`,
  `rfcs/0000-template.md`).
- Seven initial RFCs: target profiles (0001), project manifest (0002), lockfile
  (0003), diagnostics (0004), API metadata (0005), server bundle (0006), and
  tool configuration discovery (0007).
- Language profile documentation for the SA-MP 0.3.7 / open.mp Pawn dialect:
  `language/lexical.md`, `language/syntax.md`, `language/preprocessor.md`,
  `language/semantics.md`, `language/implementation-defined.md`.
- Five machine-readable target profiles (`profiles/legacy.json`,
  `profiles/samp-037.json`, `profiles/openmp.json`,
  `profiles/recommended.json`, `profiles/strict.json`), all schema version 1.
- Ten JSON Schema 2020-12 documents under `schemas/`, all `v1`:
  `pawn-project`, `pawn-lock`, `pawn-api`, `pawn-diagnostic`, `pawn-bundle`,
  `pawnlint`, `pawnfmt`, `pawntest`, `openmp-config`, `pawn-profile`.
- One or more validating examples per schema under `examples/`.
- Conformance-result format (`conformance/expected-results.schema.json`) and
  reporting manual (`conformance/manifest.md`).
- Offline validation tool (`tools/validate`) wired into
  `.github/workflows/ci.yml`.
- Compiler probes pinned to pawn-lang 3.10.10 and open.mp 3.10.11.

### Changed

- Updated RFC 0003 for sampctl's version 1 `pawn.lock`, introduced in 1.13.0
  and checked against 1.14.0. The RFC is draft until the schemas are
  reconciled.

### Notes

- This is the first commit series for `pawnkit-spec`. Most formats have no
  earlier PawnKit version. The lockfile is an exception because it must remain
  compatible with sampctl's existing file.
