# Compiler probes

These probes answer compatibility questions that are awkward to infer from
documentation alone. Their source files live in
[`pawn-corpus`](https://github.com/pawnkit/pawn-corpus/tree/v0.1.6) so other
tools can use the same evidence.

The current results were recorded against:

- `pawn-lang/compiler` at `134ad7a836c581546665340aedb59efd4636e269`
  (Pawn 3.10.10).
- `openmultiplayer/compiler` at
  `29f1a8c7fd2e01929807cd8c50353fbc93bbd651` (Pawn 3.10.11).

Both compilers produced the same results in this set. The probes confirm that
binary literals, digit separators, tag unions, `#elseif`, packed `!"..."`
strings, and nested block shadowing are accepted. Shadowing emits warning 219.
Tag mismatches emit warning 213 and name the expected and actual tags.

They also confirm that `$...h` hexadecimal literals, local `@` names, and
`#elif` are rejected. Repeating an include fails by default and succeeds in
compatibility mode. Macro redefinition emits warning 201. A directly
self-referential macro does not terminate within two seconds.

## Run the probes

Check out `pawn-corpus` v0.1.6 beside this repository, then pass one or more
`pawncc` paths:

```sh
git clone --branch v0.1.6 https://github.com/pawnkit/pawn-corpus.git ../pawn-corpus
./conformance/compiler/run.sh /path/to/pawncc
```

Set `PAWN_CORPUS` if the checkout is elsewhere.

The runner expects GNU `timeout` and gives each compilation two seconds. A
timeout is expected only for the self-referential macro probe.

The behavior is also visible in compiler source. Numeric parsing, escapes,
macro substitution, directives, and compatibility include guards are handled
in [`sc2.c`](https://github.com/pawn-lang/compiler/blob/134ad7a836c581546665340aedb59efd4636e269/source/compiler/sc2.c).
UTF-8 detection is handled in
[`sci18n.c`](https://github.com/pawn-lang/compiler/blob/134ad7a836c581546665340aedb59efd4636e269/source/compiler/sci18n.c).
