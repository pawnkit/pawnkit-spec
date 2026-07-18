# Conformance reports

A conformance report summarizes which PawnKit contracts an implementation has tested. It does not replace that repository's unit, corpus, fuzz, or integration tests.

Reports validate against `expected-results.schema.json`.

## Required fields

- `subject` identifies the repository and version under test.
- `specVersion` pins the PawnKit specification version.
- `checks` contains individual assertions tied to an RFC and, when relevant, a schema.
- `summary` contains the matching pass, fail, and skip counts.

Use `skip` with an explanation when a capability is not implemented. Do not omit the check and make the report look complete.

## Check IDs

The repository defining a check owns its ID. Use `<owner>/<short-name>`, for example `pawn-project/accepts-sampctl-manifest`. Keep a published ID stable so reports remain comparable.

## Producing a report

1. Pin the specification version.
2. Record one check for each useful assertion, not one broad check for an entire schema.
3. Use an RFC-only check for behavior that JSON Schema cannot express.
4. Compute the summary from the checks before publishing.

The [valid example](../examples/conformance/valid.json) shows the full shape. Implementation repositories may publish reports with their CI artifacts or documentation and link them from the relevant RFC.
