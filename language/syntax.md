# Syntax profile: SA-MP 0.3.7 and open.mp

This document covers grammar accepted by the `samp-037` and `openmp` source
profiles from RFC 0001. `pawn-parser` owns the executable grammar.

## Top-level declarations

A source file may contain:

- Preprocessor directives.
- Function declarations and definitions.
- Variable and constant declarations.
- Enums.
- Tagged declarations.
- State-qualified declarations.

Function modifiers include `stock`, `public`, `native`, `forward`, and
`static`. Variable modifiers include `new`, `static`, `const`, and `stock`.

Tags use the `TagName:` prefix. This dialect has no separate `tag` or
`typedef`-style declaration.

## Functions

```pawn
forward OnPlayerConnect(playerid);
public OnPlayerConnect(playerid)
{
    return 1;
}

stock Float:GetDistance(Float:x1, Float:y1, Float:x2, Float:y2)
{
    return floatsqroot(floatpower(x2 - x1, 2.0) + floatpower(y2 - y1, 2.0));
}

native SetPlayerPos(playerid, Float:x, Float:y, Float:z);
```

`forward` declares a signature without a body. `public` exposes a function to
the host by name. `native` declares a function implemented by the host or a
plugin. `stock` permits an unused function without an unused-symbol warning.

An untagged parameter or return value has the weak, generally compatible tag
described in `semantics.md`.

Parameters may have defaults:

```pawn
SetTimer(name[], interval, repeat = 0)
```

Pawn has no general variadic syntax. Variadic natives such as `printf` are host
features built on the calling convention.

## Variables and arrays

```pawn
new playerName[MAX_PLAYER_NAME];
new Float:health = 100.0;
new const DAYS[] = {"Mon", "Tue", "Wed"};
new matrix[3][4];
```

`new` declares a local or global variable. An initializer may infer an omitted
array size. Multi-dimensional arrays use one `[size]` suffix per dimension.

A string assigned to an unpacked array becomes a null-terminated cell array.
See `lexical.md` and `semantics.md` for string representation.

## Tags

```pawn
new Float:x;
new bool:flag = true;
new PlayerState:currentState;
```

A tag appears immediately before a declaration name, parameter, or function
return position.

Tag unions use the form `{Tag1, Tag2}:name`. Both profiles accept them.

`Float:` and `bool:` are common tags. Libraries define others, including
`Text:`, `PlayerText:`, and `Menu:`. Tags have no runtime representation beyond
their cell value.

## Enums

```pawn
enum E_PLAYER_DATA
{
    E_PLAYER_SCORE,
    E_PLAYER_LEVEL,
    Float:E_PLAYER_HEALTH,
}

new playerData[MAX_PLAYERS][E_PLAYER_DATA];
```

Members start at zero and increase by one unless assigned explicitly. After an
explicit value, the next member continues from that value.

Members may have tags. This supports the common enum-indexed array pattern used
as a struct substitute.

A named enum may also introduce a tag for its members.

## Operators

Pawn supports:

- Arithmetic: `+ - * / %`
- Bitwise: `& | ^ ~ << >> >>>`
- Logical: `&& || !`
- Comparison: `== != < > <= >=`
- Assignment: `= += -= *= /= %= &= |= ^= <<= >>= >>>=`
- Increment and decrement: `++ --`
- Ternary: `?:`

`>>>` is a logical right shift. `>>` is an arithmetic right shift.

`sizeof` and `tagof` are compile-time operators despite their call-like syntax.

Tagged operator declarations use function syntax:

```pawn
Fixed:operator+(Fixed:a, Fixed:b) = ...
```

`semantics.md` defines overload resolution.

## Control flow

Pawn supports `if`, `else`, `while`, `do`, `for`, `switch`, `case`, `default`,
`break`, `continue`, `return`, and `goto`.

A `case` may contain comma-separated values and `..` ranges:

```pawn
case 1, 3..5:
```

A label uses `identifier:` in statement position. Parser context distinguishes
it from a tag prefix.

## States

```pawn
state STATE_ALIVE:
GivePlayerWeapon(playerid, weaponid, ammo)
{
    // Active while the automaton is in STATE_ALIVE.
}
```

The current automaton state selects a state-qualified function implementation.
See `semantics.md`.

## Preprocessor directives

Directive lines form a separate grammar layer and begin with `#` on a logical
line. `preprocessor.md` defines `#include`, macros, conditionals, assertions,
pragmas, and location remapping.
