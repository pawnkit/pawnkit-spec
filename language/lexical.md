# Lexical profile: SA-MP 0.3.7 and open.mp

This document covers the `samp-037` and `openmp` source profiles from RFC
0001. A statement applies to both profiles unless it names one of them.

## Source encoding

Pawn compilers inspect the whole file for a UTF-8 byte-order mark. A marked
file must contain valid UTF-8 or compilation fails with error 077. Without the
mark, invalid UTF-8 falls back to the configured code page.

`pawn-parser` and `pawn-analysis` MUST use UTF-8 byte offsets as their canonical
position representation. See RFC 0004.

Tools SHOULD preserve the source encoding when editing an existing file. New
files SHOULD use UTF-8 without a byte-order mark.

## Newlines

Both LF (`\n`) and CRLF (`\r\n`) are accepted. A trailing `\` joins the next
physical line to the current logical line.

Tools count CRLF as one line ending. The compiler's exact column counting for
CRLF has not been verified.

## Whitespace and comments

Spaces, tabs, and newlines separate tokens where needed. They otherwise have no
semantic meaning.

- `//` starts a line comment.
- `/*` and `*/` delimit a block comment.
- Block comments do not nest.

Comments are trivia. A lossless parser retains them, but the compiler assigns
them no meaning.

## Identifiers

An identifier starts with an ASCII letter or `_`. Later characters may also be
digits. Identifiers are case-sensitive.

At file scope, `@` marks a symbol as public. It is not a general escape for
reserved words and is rejected on local declarations.

A trailing `:` denotes a tag and is not part of the identifier token. `Float:`
therefore consists of `Float` followed by `:`. The parser distinguishes a tag
from a label or ternary delimiter.

## Reserved words

Common reserved words include:

```text
if else while do for switch case default break continue return goto
new static stock public const native forward state enum struct
sizeof tagof defined
```

This list is informative. `pawn-parser` SHOULD take its normative list from the
compiler lexer table. An omission here does not make a word available as an
identifier.

## Numeric literals

- Decimal integer: `123`
- Hexadecimal integer: `0x1A2B`
- Binary integer: `0b1010`
- Floating point: `1.5`, `1.0e10`

A leading zero does not introduce an octal integer in these profiles.

Underscores may separate digits, as in `1_000_000` and `0b1010_0011`. Current
compilers reject the legacy `$1A2Bh` hexadecimal form. Binary literals predate
the maintained Git history and are supported by both profiles.

## Character and string literals

A single-quoted literal such as `'a'` produces the character's cell value.
Named escapes are `\a`, `\b`, `\e`, `\f`, `\n`, `\r`, `\t`, and `\v`. Quotes,
the escape character, and `%` may also be escaped. `\xNN` reads hexadecimal;
`\NNN` reads decimal, not octal. A semicolon may terminate either numeric form.

A double-quoted literal such as `"abc"` produces a packed or unpacked string.
Strings are unpacked by default unless a source marker or `#pragma pack`
changes the representation. See `syntax.md` and `semantics.md`.

The `!` prefix inverts the current packing mode. With the default unpacked
mode, `!"abc"` is packed. Under `#pragma pack`, the same prefix makes it
unpacked.

Strings do not normally cross a logical line. An unterminated string is a
lexical error unless line splicing continues it.

## Path case

Identifiers are case-sensitive. Include-path case follows the host filesystem.
This makes a project that works on Windows liable to fail on a case-sensitive
host.

`pawn-project` SHOULD warn when an include path does not match the on-disk case.

## Position conventions

Unless a tool states otherwise, examples use one-based lines and columns with
zero-based UTF-8 byte offsets. RFC 0004 defines diagnostic ranges.
