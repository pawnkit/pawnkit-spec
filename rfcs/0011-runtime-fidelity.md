---
rfc: 0011
title: Runtime fidelity
status: draft
created: 2026-07-23
updated: 2026-07-23
supersedes: null
superseded-by: null
schema: null
---

## Summary

This RFC defines runtime-fidelity tiers for PawnKit test, debug, and plugin
results. A result names the environment it actually used.

## Motivation

PawnKit can execute AMX bytecode, model open.mp behavior, call legacy native
plugins, or launch a real server. These environments catch different bugs.
Calling all of them a runtime test hides those differences.

Pawntest needs a stable way to report its backend. The debugger, plugin host,
and compatibility reports need to use the same terms.

## Current behavior

Pawntest uses goamx and a modelled set of platform natives by default. Pawn
providers can supply more natives. Legacy plugins run through
`pawn-plugin-host` workers.

Goamx tests AMX behavior without claiming to implement an open.mp server.
Pawnserver assembles real server bundles, but it does not share a result
contract with pawntest.

Reports do not currently name a common fidelity tier.

## Proposal

PawnKit defines four runtime-fidelity tiers:

- `pure-amx`: executes AMX instructions and host-supplied natives without a
  model of SA-MP or open.mp behavior;
- `platform-simulation`: uses deterministic models for selected platform
  natives, callbacks, and state;
- `native-plugin-integration`: adds a legacy native plugin through an isolated
  worker process;
- `real-server-integration`: runs against a named SA-MP or open.mp server
  build.

The tiers are descriptive, not a ranking. A higher tier can introduce
nondeterminism or unsafe native code while still matching production more
closely.

Every machine-readable runtime result MUST include:

- `runtimeTier`: one of the four names;
- `engine`: the runtime or server name and exact version;
- `target`: the Pawn profile used to compile the program;
- `capabilities`: the modelled or connected subsystems used by the test;
- `limitations`: confirmed differences that affect the result.

A result MAY list providers, adapters, plugins, or server components. It MUST
NOT include credentials, private paths, or environment values.

Pure AMX results cover VM behavior only. They MUST NOT claim platform-native or
callback compatibility.

Platform simulation results MUST list the simulated capability groups. Missing
natives must fail unless the test explicitly allows or mocks them. Pawn
providers remain part of this tier because they execute inside the controlled
AMX environment.

Native plugin integration MUST name the plugin architecture and worker
protocol version. Process isolation limits crashes and hangs; it is not a
security sandbox.

Real server integration MUST name the server artifact, version, configuration
profile, and adapter protocol. It MUST distinguish server startup from a
completed test session.

Pawntest plain output SHOULD print the tier once before the summary when it is
not `platform-simulation`. JSON, TAP, JUnit, coverage, and compatibility
reports MUST retain the tier in a format-appropriate field.

Pawndebug reports the tier of its backend. It does not promote a simulated
session to real-server fidelity. Compatibility reports compare results only
within the same tier unless they explicitly describe a differential.

## PawnKit extensions

None.

## Compatibility impact

- [x] Additive
- [ ] Breaking

Machine-readable report formats gain runtime metadata. Readers must ignore
unknown additive fields during the transition.

## Alternatives considered

Backend names such as `goamx` describe an implementation, not the behavior a
test covered.

Using only simulated and real groups would hide pure VM tests and native
plugin risks.

Treating native plugins as real-server tests would imply server behavior that
the worker does not provide.

## Security considerations

AMX programs, native plugins, server bundles, adapter messages, and reports
are untrusted input. Implementations must retain instruction, memory, process,
output, and timeout limits appropriate to their tier.

Native plugin workers are not sandboxes. Real-server tests must bind network
listeners deliberately, avoid public interfaces by default, and redact
credentials and private paths.

## Migration plan

Writers add runtime metadata before consumers require it. During one report
major, missing metadata is interpreted as the writer's documented legacy
default and marked as inferred.

## Reference implementation status

Open. `goamx`, `pawntest`, `pawndebug`, `pawn-plugin-host`, and `pawnserver`
own their tier-specific behavior. `pawntest` owns test-result reporting.

## Conformance tests

Open. Pawntest needs one report fixture per implemented tier. Goamx, plugin
host, and real-server adapters need differential and failure-path fixtures.

## Open questions

- Which report formats need a major version before runtime metadata becomes
  required?
- Should capability names be free-form until the first simulation inventory
  is complete?
