# Compatibility

## Supported profiles

| Profile | File | Kind | Status |
|---|---|---|---|
| `legacy` | `profiles/legacy.json` | source | experimental |
| `samp-037` | `profiles/samp-037.json` | source + runtime | experimental |
| `openmp` | `profiles/openmp.json` | source + runtime | experimental |
| `recommended` | `profiles/recommended.json` | source (policy overlay) | draft |
| `strict` | `profiles/strict.json` | source (policy overlay) | draft |

The first three profiles describe compiler or runtime behavior. `recommended`
and `strict` only change diagnostic policy; they do not change which Pawn
programs are valid. Both currently build on `openmp`. RFC 0001 explains the
source/runtime split.

## Version policy

Repository releases use semantic version tags (`vX.Y.Z`) for the complete
documentation and schema bundle.

A schema's major version appears in its `$id`, for example
`.../v1/schema.json`. Formats with a `schemaVersion` field use the same major
version. Additive fields and relaxed constraints do not require a new major
version.

RFCs follow the lifecycle in `GOVERNANCE.md`: `draft`, `experimental`, then
`accepted`. They may later become `deprecated` or `superseded`. A schema is not
stable while its RFC is still a draft or experiment.

Profiles have two versions. `schemaVersion` selects the profile schema;
`version` tracks revisions to that particular profile.

## Backward compatibility during transitions

Tools should read the current schema major version and at least one earlier
stable version.

A breaking change needs a migration path. Depending on the format, that may be
manual instructions, a `pawnmigrate` rule, or a reader that supports both
versions for a stated period. The RFC records which approach applies.

Existing names keep their meaning. Reusing a diagnostic code, profile ID, or
schema version for a different purpose is a breaking change.

## Initial release note

Most formats in this repository are still on their first published version.
Their RFCs record any known compatibility work. RFC 0003 is the exception: its
early PawnKit lockfile draft now overlaps with the `pawn.lock` introduced by
sampctl 1.13.0 and must be reconciled before acceptance.

## Sampctl compatibility commitment

`pawn-project` must continue to read sampctl manifests. PawnKit extensions stay
under the optional `pawnkit` key so sampctl projects do not need a second
manifest.

sampctl also owns the version 1 `pawn.lock` shape introduced in 1.13.0. RFC
0003 is currently a draft because the older PawnKit schema used the same
filename for a different structure. PawnKit must adopt the upstream shape or
use extensions that do not break sampctl before that RFC can advance.
