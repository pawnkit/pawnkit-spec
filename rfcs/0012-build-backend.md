---
rfc: 0012
title: Build backend protocol
status: draft
created: 2026-07-23
updated: 2026-07-23
supersedes: null
superseded-by: null
schema: null
---

## Summary

This RFC defines a local process protocol for restoring, building, and running
Pawn projects. `pawn-project` resolves configuration; a backend executes the
resolved request.

## Motivation

The CLI can negotiate optional build and test commands, but its current
adapter passes only a project directory. Each executable must then rediscover
the profile, entry point, include paths, defines, compiler, and output.

PawnKit needs to support sampctl projects and direct compiler use without
putting either implementation inside the CLI.

## Current behavior

`pawn check` accepts `--build-tool` and `--test-tool`. It probes:

```text
<tool> capabilities --output json
```

It then runs the selected command with `--project` and `--output json`.
Protocol version 1 returns a status and message, but it does not carry the
resolved project model or structured diagnostics.

There are no standalone `pawn restore`, `pawn build`, or `pawn run` commands.

## Proposal

PawnKit defines build backend protocol version 1. Backends are local
executables. The CLI does not evaluate a command through a shell.

Capability negotiation runs:

```text
<backend> capabilities --output json
```

The response contains:

- `protocolVersion`: `1`;
- `name` and `version`: backend identity;
- `operations`: any of `restore`, `build`, and `run`;
- `profiles`: supported Pawn profile names, or an empty list for any resolved
  profile;
- `features`: optional capabilities such as dependency restore, compiler
  acquisition, or server launch.

Execution runs:

```text
<backend> execute --input REQUEST --output RESULT
```

`REQUEST` and `RESULT` are files created with owner-only permissions. `-`
selects standard input or output. File arguments avoid command-line length
limits on projects with many includes or defines.

A request contains:

- `schemaVersion`: `1`;
- `operation`: `restore`, `build`, or `run`;
- `projectRoot`: the absolute canonical project root;
- `profile`: the selected Pawn profile;
- `target`: the selected runtime target;
- `entry`: the absolute source entry point when required;
- `output`: the absolute AMX or bundle output when required;
- `includePaths`: absolute paths in canonical search order;
- `defines`: resolved names and values;
- `compiler`: an optional exact path, version, and checksum;
- `arguments`: build or runtime arguments from the selected manifest profile.

Restore requests omit entry and output. Build requests require both. Run
requests require an existing output or identify a build request that must run
first; the first version SHOULD use an existing output to keep execution
explicit.

The backend MUST treat the request as resolved input. It MUST NOT search for a
different project manifest or silently select another profile. It MAY validate
that files still match the request.

A result contains:

- `schemaVersion`: `1`;
- `status`: `passed`, `failed`, or `cancelled`;
- `backend`: name and version;
- `artifacts`: paths, media types, sizes, and optional SHA-256 hashes;
- `diagnostics`: PawnKit diagnostic version 1 values;
- `process`: optional exit code and bounded standard output and error;
- `runtime`: runtime-fidelity metadata for a completed run.

Paths in results must remain inside the project output directory unless the
request explicitly names another destination.

Cancellation closes standard input and sends the platform's normal process
termination signal. A backend SHOULD stop child processes and remove partial
artifacts. The caller MAY force termination after a bounded grace period.

Sampctl is the first recommended adapter because it already owns dependency
restoration and existing project commands. A direct compiler adapter covers
projects whose dependencies and compiler are already resolved. Project-defined
commands and container backends remain experimental until they can satisfy the
same request and security rules.

`pawn-project` owns request construction from manifests and lockfiles.
`pawnkit-cli` selects and invokes a backend. Actions and editors call the CLI
or submit the same request; they do not rebuild it.

## PawnKit extensions

The protocol maps sampctl-compatible project fields into resolved requests. It
does not change `pawn.json`, `pawn.yaml`, or `pawn.lock`.

## Compatibility impact

- [x] Additive
- [ ] Breaking

The existing `pawn check --build-tool` adapter can remain during one CLI minor
release. New backends use `execute`; old adapters are not described as tested
build backends.

## Alternatives considered

Making sampctl the only backend would block direct compiler and controlled CI
use cases.

Invoking the compiler only would duplicate sampctl dependency and runtime
behavior.

Letting every backend reload the manifest recreates the project disagreement
this protocol is meant to remove.

Passing every field as a command-line flag is fragile on Windows and exposes
arguments in process listings.

## Security considerations

Manifests, requests, compiler binaries, source files, dependencies, artifacts,
and process output are untrusted. Callers must bound request and output sizes,
verify managed compiler downloads, reject traversal, avoid shells, and limit
execution time.

Backends must not print credentials or environment values. Run backends bind
network listeners only when requested and should default to loopback.

## Migration plan

1. `pawnkit-spec` publishes request, result, and capability schemas.
2. `pawn-project` constructs deterministic requests.
3. Sampctl and direct compiler adapters implement protocol version 1.
4. `pawnkit-cli` adds restore, build, and run commands.
5. The old project-directory adapter remains for one CLI minor release.

## Reference implementation status

Open. `pawn-project` owns request construction. `pawnkit-cli` owns invocation.
The first adapters belong with the CLI unless a reusable implementation
emerges.

## Conformance tests

Open. Tests need valid and invalid messages, Windows path cases, cancellation,
partial writes, bounded output, a clean sampctl project, and the small SA-MP
and open.mp corpus projects.

## Open questions

- Which sampctl versions can produce deterministic restore and build results?
- Should run requests name one AMX artifact or a complete server bundle?
