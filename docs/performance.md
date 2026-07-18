# Performance

Only the validation tool executes code in this repository. Its job is to give contributors quick local feedback.

Run it from `tools/validate`:

```sh
go run . ../../schemas ../../profiles ../../examples ../../conformance ../../rfcs
```

The complete check should finish in under two seconds on a warm single-core run. Individual fixtures are limited to 1 MiB unless the validator records a reviewed exception.

The tool prints its document count and elapsed time. A separate benchmark suite is unnecessary at the current size; add one if generated bundles or a much larger conformance set make the command slow enough to investigate.
