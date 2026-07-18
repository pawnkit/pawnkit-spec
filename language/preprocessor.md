# Preprocessor profile: SA-MP 0.3.7 and open.mp

This document covers the `samp-037` and `openmp` source profiles from RFC
0001.

Preprocessing is a distinct language stage. `pawn-analysis` MUST retain the
original and expanded token streams, with source mapping between them. This
document defines the required behavior, not the internal implementation.

## Directive lines

A directive starts with `#` as the first non-whitespace character on a logical
line. Line splicing happens before directive processing.

Directives are processed from top to bottom. Conditional directives decide
which lines take part in later processing.

## Includes

```pawn
#include <a_samp>
#include "myfile.inc"
```

Angle brackets search configured include roots. Quoted paths search relative to
the including file first, then the include roots. RFC 0002 defines project
resolution.

Current pawn-lang and open.mp compilers do not generate guards for `.inc` files
by default. `-Z+` or `#pragma compat` restores the historical behavior.

`#tryinclude` follows the same search rules but does not fail when the file is
missing.

## Macros

Object-like macro:

```pawn
#define MAX_HEALTH 100
```

Function-like macro:

```pawn
#define CLAMP(%1,%2,%3) ((%1) < (%2) ? (%2) : ((%1) > (%3) ? (%3) : (%1)))
```

Pawn macros use positional parameters such as `%1`, rather than C-style named
parameters. A trailing `\` continues a macro onto the next physical line.

`#undef` removes a macro. Redefinition without `#undef` emits warning 201 and
uses the new definition.

The reference compilers do not bound direct recursive expansion. A macro such
as `#define SELF SELF` can hang compilation. PawnKit tools MUST apply the limits
described below instead of copying this behavior.

## Conditional compilation

```pawn
#if defined SOME_FLAG
    // ...
#elseif OTHER_FLAG > 5
    // ...
#else
    // ...
#endif
```

`defined NAME` tests whether a macro or symbol exists.

Pawn uses `#elseif`. Current compilers do not recognize `#elif` as a
conditional directive.

The preprocessor must scan inactive regions far enough to match nested
conditionals. `pawn-analysis` MUST preserve inactive branch tokens so editors
can inspect and edit either branch.

## Pragmas

Common pragmas include:

- `#pragma pack`, with its unpack counterpart, controls the default string
  representation.
- `#pragma compat` enables compatibility behavior from the older compiler.
- `#pragma dynamic <cells>` sets runtime stack and heap space.
- `#pragma tabsize <n>` sets the compiler's tab width for diagnostics.
- `#pragma unused <name>` suppresses the compiler warning for one identifier.

The list is not exhaustive. Tools should recognize `#pragma unused` as an
intentional suppression instead of reporting the same finding again.

## Compile-time assertions

```pawn
#assert MAX_PLAYERS <= 1000
```

`#assert` fails compilation when its constant expression is false. It is not a
runtime `assert` native.

## File and line remapping

`#file` and `#line` change the location reported by the compiler for following
diagnostics. They do not change parsing or program semantics.

`pawn-analysis` should retain this mapping for compiler-compatible diagnostic
locations.

## Expansion source maps

A diagnostic produced inside expanded code MUST retain:

- The expanded position where the issue was found.
- The original invocation position.

Tools should also retain the macro definition position when available.

## Limits

`pawn-parser` and `pawn-analysis` MUST bound include depth, expansion depth,
expanded output size, and processing time for untrusted input. This
specification does not set the numeric limits. Each implementation must
document and test its chosen bounds.
