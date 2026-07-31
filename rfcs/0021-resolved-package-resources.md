---
rfc: 0021
title: Resolved package resources
status: draft
created: 2026-07-31
updated: 2026-07-31
supersedes: null
superseded-by: null
schema: schemas/pawn-lock.schema.json
---

## Summary

This RFC adds resolved plugin, component, filterscript, and include assets to
the sampctl-compatible lockfile. PawnKit needs the final URL, checksum, and
extraction map before it can install package resources without sampctl.

## Motivation

Package manifests select release assets with patterns. The current lockfile
records the package commit but not the asset chosen for a host. A later release
can replace an asset at the same URL, and different clients can resolve a
pattern differently.

`pawn restore` can reproduce source checkouts today. It cannot safely restore
the native files declared by `plugin://`, `component://`, `filterscript://`, or
`includes://` dependencies.

## Current behavior

sampctl 1.14.1 writes version 1 lockfiles from
`src/pkg/package/lockfile`. Its dependency entries preserve a scheme and source
commit. Resource selection still comes from `pawn.json`.

The sampctl reader uses `encoding/json.Unmarshal`, so it ignores unknown
fields. Its writer marshals the known Go structure and therefore removes
unknown fields when it rewrites the file.

PawnKit's temporary legacy lock shape contains `platformArtifacts`, but writers
must now use the sampctl version 1 shape described by RFC 0003.

## Proposal

Version 1 lockfiles MAY contain a top-level `pawnkit` object:

```json
{
  "pawnkit": {
    "schema_version": 1,
    "resources": [
      {
        "package": "plugin://samp-incognito/samp-streamer-plugin",
        "resource": "streamer",
        "target": "linux-amd64",
        "url": "https://example.invalid/streamer.zip",
        "size": 1234,
        "checksum": "sha256:...",
        "archive": "zip",
        "files": [
          {
            "source": "plugins/streamer.so",
            "destination": "plugins/streamer.so",
            "size": 456,
            "checksum": "sha256:..."
          }
        ]
      }
    ]
  }
}
```

`package` MUST match a dependency key in the same lockfile. `resource` is the
manifest resource name. `target` uses PawnKit host target names.

`url` MUST be a final HTTPS URL without credentials. `size` and `checksum`
cover the downloaded file. `archive` is `zip`, `tar.gz`, or `file`.

Each file records its archive source, project-relative destination, extracted
size, and checksum. A `file` resource has one file whose source is the
downloaded filename.

Coordinates formed by package, resource, and target MUST be unique. Installers
MUST select an exact target and MUST NOT fall back to another architecture.

## PawnKit extensions

The optional top-level `pawnkit` object is a PawnKit extension to sampctl's
version 1 lockfile. Its schema definition will use
`x-pawnkit-extension: true`.

sampctl 1.14.1 can read a lockfile containing this object, but removes the
object when it saves the file. PawnKit MUST report missing resolved resources
after such a rewrite and regenerate them before installing assets.

## Compatibility impact

- [x] Additive
- [ ] Breaking

Existing sampctl lockfiles remain valid. sampctl can read the extension, though
it does not preserve it when writing. PawnKit readers that do not install
resources may ignore it.

## Alternatives considered

A second resource lockfile would survive sampctl rewrites, but splits one
dependency resolution across two files.

Re-resolving assets from `pawn.json` keeps the current format but is not
reproducible. Recording only the final URL still permits undetected replacement
and unsafe extraction.

Changing sampctl's core version 1 fields would break compatibility and is not
required.

## Security considerations

Readers must bound the lockfile, downloads, entry count, individual files, and
total extraction size. They must verify the archive before extraction and each
installed file afterward.

Absolute paths, traversal, links, device files, duplicate destinations, and
case-colliding destinations must be rejected. Redirects must remain on HTTPS
and must not introduce credentials.

Installers must stage writes and replace destinations recoverably. Package
scripts are not executed during resolution or installation.

## Migration plan

Existing lockfiles remain usable for source restoration. `pawn restore` will
report scheme dependencies without resolved resources and direct the user to a
future resource-resolution command.

The temporary legacy `platformArtifacts` reader remains available on the
schedule in RFC 0003. `pawnmigrate` may copy complete legacy artifact records
into the extension when every required field is present.

## Reference implementation status

Open. `pawn-project` will own parsing, validation, resolution, and safe
installation. `pawnkit-cli` will only invoke that capability.

Implementation must not begin until this RFC is accepted.

## Conformance tests

Open. `pawnkit-spec` needs valid Linux and Windows examples plus failures for
duplicate coordinates, insecure URLs, missing checksums, traversal, links,
case collisions, size limits, and target fallback.

`pawn-project` needs archive and single-file installation tests on Linux and
Windows. A sampctl 1.14.1 compatibility test must confirm that the extended
lockfile loads and document that a save removes the extension.

## Open questions

- Should PawnKit preserve the extension when invoking a sampctl operation that
  rewrites the lockfile, or only warn and require regeneration?
- Which command resolves resources without changing dependency commits?
- Should a future sampctl release preserve the namespaced object?
