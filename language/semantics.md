# Semantic profile: SA-MP 0.3.7 and open.mp

This document covers the `samp-037` and `openmp` source profiles from RFC
0001. `pawn-analysis` owns the implementation of these rules.

## Cells

Pawn stores values in fixed-width cells. SA-MP and open.mp use 32-bit cells;
these profiles do not cover the 16-bit or 64-bit variants from the general Pawn
specification.

Integers, `Float:` values, booleans, characters, and unpacked array elements
each occupy one cell. A `Float:` value stores an IEEE 754 single-precision bit
pattern in that cell.

Arrays are contiguous. Multi-dimensional arrays use compiler-computed strides,
not pointers to independent arrays.

## Tags

Tags are weak compile-time types. They have no runtime representation. At the
AMX level, tagged and untagged cells with the same bits are identical.

A tag mismatch is normally a compiler warning, not an error. Assigning an
untagged value to `Float:`, for example, still compiles even though interpreting
the bits is usually incorrect. Policy profiles may raise the severity without
making the source invalid.

Warning 213 names both the expected and actual tags. The pinned pawn-lang and
open.mp compilers use the same wording.

Tag unions such as `{TagA, TagB}:value` allow more than one compatible tag.

The `_` tag, and an absent tag, are compatible with other tags without a
mismatch warning. This supports tag-agnostic functions and natives.

## Operator overloading

Code may define an operator for specific tags:

```pawn
Fixed:operator+(Fixed:a, Fixed:b)
{
    return Fixed:(_:a + _:b);
}
```

The compiler selects an overload by operand tags. If none matches, it uses the
built-in operator for compatible operands. There is no runtime dispatch.

`Tag:value` and `_:value` change the compile-time tag without changing bits.
They do not perform numeric conversion. Natives such as `float()` and
`floatround()` perform actual conversion.

## Automatons and states

A function may have several state-qualified implementations. The current value
of an implicit or explicit automaton selects the implementation at runtime.

State-qualified functions are uncommon in SA-MP and open.mp projects, but
`pawn-analysis` MUST support them because neither source profile disables the
feature.

## Scope and lifetime

At file scope, `new` declares a global. Inside a function, it declares a local.

File-level `static` limits visibility to the current preprocessed compilation
unit. Function-level `static` preserves a local value between calls.

Nested blocks introduce local scopes. A nested local may shadow an outer local;
the reference compilers accept it with warning 219.

## Function resolution

Pawn has no general function overloading by parameter type. Tag-specific
operator overloads are the exception. Reusing a normal function name with a
different signature is a redeclaration error.

The host finds `public` functions by their exact AMX symbol name. Renaming a
public function can therefore change runtime behavior. A rename tool must treat
that edit as review-required under RFC 0004.

## Preprocessing and scope

Semantic analysis uses the expanded token stream and retains a map to the
original source, as required by `preprocessor.md`.

A name introduced through a macro resolves at the expansion site. Pawn macros
are textual substitutions, not hygienic macros with definition-site scope.

## Scope of this document

This profile does not define native availability, callback behavior, or AMX
instruction execution. `pawn-api` owns API facts under RFC 0005. `goamx` owns
AMX runtime behavior.
