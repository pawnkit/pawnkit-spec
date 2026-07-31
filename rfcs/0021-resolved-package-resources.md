---
rfc: 0021
title: Resolved package resources
status: experimental
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

`package` MUST match a dependency key in the same lockfile. This may be a
resource-scheme key such as `plugin://owner/repo` or sampctl's usual
`github.com/owner/repo` form. Packages can declare release resources without
using a resource scheme. `resource` is the manifest resource name. `target`
uses PawnKit host target names.

`url` MUST be a final HTTPS URL without credentials. `size` and `checksum`
cover the downloaded file. `archive` is `zip`, `tar.gz`, or `file`.

Each file records its archive source, project-relative destination, extracted
size, and checksum. A `file` resource has one file whose source is the
downloaded filename.

Coordinates formed by package, resource, and target MUST be unique. Installers
MUST select an exact target and MUST NOT fall back to another architecture.

When a resource pattern matches several release assets, resolvers MUST select
the first match in the release provider's asset order. This matches sampctl and
supports existing manifests whose patterns also match debug archives. The lock
records the selected URL and checksum, so restoration does not repeat this
selection.

Some existing manifests name a plugin's old directory after a release moves
the binary. If an exact plugin path selects no file, a resolver MAY select the
sole archive file with the same basename. It MUST reject the fallback when
more than one file has that basename. The locked source path records the file
actually selected.

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

`pawn-project` owns parsing, validation, resolution, and safe installation.
The first prototype will consume complete lock records before it adds network
resolution. `pawnkit-cli` only invokes those capabilities.

## Conformance tests

The lock schema includes a valid resolved-resource example and rejects
insecure URLs, missing checksums, unsafe paths, unsupported archives, and
unbounded resource lists.

Coordinate uniqueness, archive links, case collisions, aggregate size limits,
and exact target selection require implementation tests because JSON Schema
cannot inspect archive contents or compare compound keys.

`pawn-project` needs archive and single-file installation tests on Linux and
Windows. A sampctl 1.14.1 compatibility test must confirm that the extended
lockfile loads and document that a save removes the extension.

## Open questions

`pawn restore` only consumes locked data. A later `pawn install` command may
resolve missing resources and update the lockfile without changing dependency
commits.

PawnKit will not copy resource records across a sampctl lockfile rewrite. They
may be stale if dependency data changed. It reports the missing records and
requires resolution again.

A future sampctl release may preserve the namespaced object, but PawnKit does
not depend on that behavior.
