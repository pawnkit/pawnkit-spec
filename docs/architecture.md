# Architecture

`pawnkit-spec` owns contracts, not implementations. It publishes schemas, profiles, RFCs, and language notes as data that other repositories consume.

```text
pawnkit-spec
    |
    +-- pawn-project   manifests, lockfiles, tool settings
    +-- pawn-api       API metadata
    +-- pawnserver     bundles
    +-- Pawn tools     profiles and diagnostics
```

These arrows mean conformance, not Go imports. This repository has no public Go library.

## Repository areas

| Area | Content |
|---|---|
| `rfcs` | Decisions and migration plans |
| `language` | Confirmed Pawn dialect behavior and explicit unknowns |
| `profiles` | Machine-readable source and runtime profiles |
| `schemas` | JSON Schema 2020-12 contracts |
| `examples` | Documents checked against those schemas |
| `conformance` | A common format for implementation results |

## Ownership

The schema and meaning of a shared format belong here. Code that reads, writes, or acts on that format belongs in its implementation repository.

| Concern | Owner |
|---|---|
| Pawn syntax | `pawn-parser` |
| Preprocessing and semantics | `pawn-analysis` |
| Project discovery | `pawn-project` |
| SA-MP and open.mp API facts | `pawn-api` |
| Lint rules | `pawnlint` |
| Bundle installation | `pawnserver` |

## Design rules

### Extend sampctl manifests

PawnKit settings are optional additions to the existing sampctl manifest. RFC 0002 does not introduce a competing project format.

### Keep source and runtime profiles distinct

A source profile describes accepted language behavior. A runtime profile describes AMX and host behavior. One profile document may contain both sections, but tools must not treat them as the same kind of fact.

### Share one diagnostic shape

CLI reports, baselines, and editor integrations use the diagnostic schema in this repository rather than defining private formats.

### Record uncertainty

If compiler behavior has not been confirmed against a primary source or conformance test, the language notes mark it as unknown. They do not turn a plausible guess into a contract.

### Validate offline

`tools/validate` resolves schema references from this checkout and makes no network requests. Publishing the same files at `schemas.pawnkit.dev` does not change local validation.

## Adding a contract

Add profiles under `profiles` and validate them against the profile schema. A new kind of shared field or file format needs an RFC. Tool-private session or cache formats should stay with the tool unless another implementation needs to exchange them.
