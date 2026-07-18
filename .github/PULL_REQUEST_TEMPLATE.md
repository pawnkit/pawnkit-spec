## Summary

<!-- What does this PR change, and why? Link the RFC/issue if applicable. -->

## Type of change

- [ ] New or amended RFC
- [ ] Language profile documentation
- [ ] Machine-readable profile (`profiles/`)
- [ ] JSON Schema (`schemas/`) + example (`examples/`)
- [ ] Conformance format (`conformance/`)
- [ ] Validation tooling (`tools/validate`)
- [ ] Docs/governance/process only

## Compatibility

- [ ] Additive
- [ ] Breaking; the RFC includes compatibility and migration details

## Checklist

- [ ] `cd tools/validate && go run . ../../schemas ../../profiles ../../examples ../../conformance ../../rfcs` passes locally.
- [ ] New/changed schemas have at least one example under `examples/` that validates.
- [ ] New/changed RFCs use `rfcs/0000-template.md` and have well-formed front matter.
- [ ] Language claims cite compiler source or a repeatable probe; unresolved behavior is marked as unknown.
