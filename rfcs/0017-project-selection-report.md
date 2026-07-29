---
rfc: 0017
title: Resolved project selection report
status: experimental
created: 2026-07-29
updated: 2026-07-29
supersedes: null
superseded-by: null
schema: schemas/pawn-project-selection.schema.json
---

## Summary

This RFC defines a small JSON report for the profile, build, runtime, and paths
selected from a Pawn project manifest.

## Motivation

Editors and command-line frontends need to show the same project selection used
by build and analysis tools. Reading `pawn.json` or `pawn.yaml` in each client
would duplicate profile fallback, path resolution, and sampctl compatibility
rules owned by `pawn-project`.

## Current behavior

`pawn-project` resolves the active profile, build, runtime, entry point, and
project root. Consumers can access that state through its Go API, but there is
no process-level report for editors or other languages.

`vscode-pawn` can run `pawn doctor`, but that report contains findings rather
than the resolved project selection.

## Proposal

The PawnKit CLI provides:

```text
pawn project --project DIR --output human|json
```

JSON output follows `pawn-project-selection.schema.json` and contains:

- `schemaVersion`: report major, initially `1`;
- `root`: absolute project root;
- `manifest`: absolute manifest path;
- `entry`: absolute entry source path, or an empty string when unset;
- `profile`: resolved RFC 0001 profile ID, or an empty string when unset;
- `build`: selected build name, or an empty string;
- `runtime`: selected runtime name, or an empty string.

Resolution uses RFC 0002 and the same `pawn-project` operation used by builds.
The command does not modify the manifest.

Consumers MUST reject unknown schema majors. Paths are local filesystem paths
and MUST NOT be included in telemetry or public reports without redaction.

## PawnKit extensions

None.

## Compatibility impact

- [x] Additive
- [ ] Breaking

This is a new command and report. Existing consumers do not change.

## Alternatives considered

Adding selection fields to `pawn doctor` would mix project identity with health
findings. Parsing manifests in the editor would duplicate `pawn-project`.

A manifest-writing command is not included. Preserving comments and formatting
across JSON and YAML needs a separate design.

## Security considerations

The command only reads local project files. It must apply the existing project
path checks and must not fetch dependencies or execute project code. Absolute
paths may identify a user or workspace and should be treated as private.

## Migration plan

Not applicable; this is the first version.

## Reference implementation status

`pawn-project` owns resolution. `pawnkit-cli` is the report producer, and
`vscode-pawn` is the first consumer.

## Conformance tests

`pawnkit-spec/examples/pawn-project-selection` validates the report shape.
`pawnkit-cli/pkg/cli` tests JSON output from JSON and YAML manifests.
`vscode-pawn` tests report decoding and status labels.

## Open questions

- Should a later version expose all available named builds and runtimes?
- Should manifest updates use a previewable patch protocol?
