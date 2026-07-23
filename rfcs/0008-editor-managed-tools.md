---
rfc: 0008
title: Editor-managed tool state
status: draft
created: 2026-07-23
updated: 2026-07-23
supersedes: null
superseded-by: null
schema: null
---

## Summary

This RFC defines how an editor tells `pawnlsp` about include roots installed
with managed tools. It keeps project settings in `pawn-project` while allowing
an editor to report files that exist outside the workspace.

## Motivation

`vscode-pawn` installs tools in extension storage. Some tools include Pawn
sources, such as the `pawntest` include. The language server needs the storage
path, but it must not treat editor settings as a second project manifest.

The current client sends a generic `includePaths` array during initialization
and through `workspace/didChangeConfiguration`. That name does not distinguish
managed tool files from project include paths.

## Current behavior

`vscode-pawn` sends this implementation-specific initialization value:

```json
{
  "includePaths": ["/editor-storage/tools/pawntest/v1.1.2/include"]
}
```

It sends the same value under `settings.pawn.includePaths` after installing
`pawntest`. `pawnlsp` passes the roots to `pawn-project` as managed include
roots. Other project state comes from `pawn.json` or `pawn.yaml`.

## Proposal

The client SHOULD send this object as `initializationOptions`:

```json
{
  "pawnkit": {
    "protocolVersion": 1,
    "managedIncludeRoots": [
      "/editor-storage/tools/pawntest/v1.1.2/include"
    ]
  }
}
```

`protocolVersion` MUST be `1`.

`managedIncludeRoots` contains absolute filesystem paths owned by tools that
the editor installed. It MUST NOT contain project include paths, compiler
settings, profiles, defines, or API data.

The server MUST clean and de-duplicate the roots, reject relative paths, and
pass accepted roots to `pawn-project` through `ManagedIncludeRoots`. Project
roots and dependencies retain their normal priority.

After managed tools change, the client SHOULD send:

```text
pawnkit/didChangeManagedTools
```

with:

```json
{
  "protocolVersion": 1,
  "managedIncludeRoots": [
    "/editor-storage/tools/pawntest/v1.1.2/include"
  ]
}
```

The server MUST reload affected project state after accepting the notification.
An empty array removes all editor-managed roots.

Clients MAY omit the object when they do not manage tools. Servers MUST then
use the project model without additional roots.

## PawnKit extensions

This is a PawnKit extension to LSP initialization options and notifications.
It does not change the LSP base protocol.

## Compatibility impact

- [x] Additive
- [ ] Breaking

During the version 1 transition, `pawnlsp` accepts the existing top-level
`includePaths` initialization value and `settings.pawn.includePaths`
configuration value as deprecated aliases. New clients use only the versioned
fields.

## Alternatives considered

Restarting the server after every tool installation would avoid a notification,
but it would discard useful editor state and interrupt active requests.

Making the extension edit `pawn.json` would mix machine-local storage paths
into a shared project file.

Letting `pawnlsp` search VS Code storage would couple an editor-independent
server to one client.

Keeping the generic `includePaths` name would leave project paths and managed
tool state indistinguishable.

## Security considerations

The server accepts only absolute paths and applies the existing
`pawn-project` include resolver. Implementations should bound the number of
roots and must not scan unrelated parent directories. The notification carries
paths, not file contents or credentials, and performs no network access.

Workspace trust and download verification remain the editor's responsibility.

## Migration plan

1. `pawnlsp` adds the versioned fields while retaining the deprecated aliases.
2. `vscode-pawn` sends the versioned initialization object and notification.
3. Tests confirm that initialization, installation, reload, and document
   lifecycle events use the same project state.
4. The aliases remain supported through the next `pawnlsp` minor release. A
   later removal requires a breaking release and an RFC update.

## Reference implementation status

Planned for `pawnlsp` and `vscode-pawn`.

## Conformance tests

Planned in the `pawnlsp/lsp` protocol tests and the `vscode-pawn` client tests.
They will cover empty roots, invalid relative roots, updates, reloads, and
compatibility aliases.

## Open questions

- Should a later protocol version report managed tool identity and version as
  well as its include root?
