# pawnkit-spec

`pawnkit-spec` contains the contracts shared by PawnKit tools: RFCs, language profiles, JSON Schemas, and conformance data for SA-MP 0.3.7 and open.mp projects.

## Status

RFC 0006 is accepted. The remaining RFCs are draft or experimental, so most
version 1 schemas are not frozen yet. [GOVERNANCE.md](GOVERNANCE.md) explains
the states and acceptance rules.

## Find a contract

| Need | Directory |
|---|---|
| Pawn language behavior | [`language`](language) |
| SA-MP and open.mp target profiles | [`profiles`](profiles) |
| Shared file formats | [`schemas`](schemas) |
| Format examples | [`examples`](examples) |
| Proposed or accepted decisions | [`rfcs`](rfcs) |
| Implementation results | [`conformance`](conformance) |

Implementations live in their own repositories. For example, `pawn-parser` owns parsing and `pawn-project` owns project discovery. This repository defines the contracts they follow.

## Schema stability

Every schema declares a stable `$id` of the form:

```text
https://schemas.pawnkit.dev/<name>/v<major>/schema.json
```

`https://schemas.pawnkit.dev/` does not serve the files yet. For now, treat each `$id` as an identifier and load schemas from a pinned commit or tag in this repository.

Formats with a `schemaVersion` field use the same major version as their schema URL. The transition policy is in [docs/compatibility.md](docs/compatibility.md).

## Validating locally

Run the same schema, example, profile, conformance, and RFC checks used by CI:

```sh
cd tools/validate
go run . ../../schemas ../../profiles ../../examples ../../conformance ../../rfcs
```

The Go module exists only inside `tools/validate`; `pawnkit-spec` is not a public Go library.

## Contributing

Compiler evidence, compatibility notes, and small schema fixes are welcome.
See [CONTRIBUTING.md](CONTRIBUTING.md) before changing a shared contract.

## License

MIT, see [LICENSE](LICENSE).
