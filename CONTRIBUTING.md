# Contributing

PawnKit is maintained by volunteers, so reviews may take a little time.

Corrections, examples, and implementation evidence are welcome. Standards work
can start with a small compiler test or a clearly described compatibility gap.

This repository contains Markdown and JSON plus the validation tool under `tools/validate`. It requires Go 1.26 or later.

Run the same offline check as CI:

```sh
cd tools/validate
go run . ../../schemas ../../profiles ../../examples ../../conformance ../../rfcs
```

## RFCs

Copy `rfcs/0000-template.md` to the next numbered file and keep every section. A draft may leave implementation and conformance work open. Breaking changes need a migration plan and the discussion period defined in [GOVERNANCE.md](GOVERNANCE.md).

## Schemas

Schemas live under `schemas` and use a stable `$id` containing their major version. Add or update an example with every shape change. Additive changes keep the major version; breaking changes require the RFC process.

Parser code, project resolution, API data, and generated website pages belong in their implementation repositories. [docs/architecture.md](docs/architecture.md) lists the owners.

There are no generated files today. If a generator is added, document its exact command and make CI reject stale output.

## Releases

Validate the repository, update [CHANGELOG.md](CHANGELOG.md), and tag the source bundle with semantic versioning. Individual schemas keep their own major versions independently of the repository tag.

Be direct when reviewing language behavior. Cite compiler source or a repeatable test. If neither settles the question, record it as unknown in `language/implementation-defined.md`.
