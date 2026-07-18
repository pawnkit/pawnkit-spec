# Security policy

Security issues in this repository usually concern a contract that permits unsafe behavior, such as escaping archive paths, unbounded recursion, or validation that requires executing project code. The local validation tool is also in scope.

Report issues through GitHub's private security advisory form for `pawnkit-spec`. Include the affected files, likely impact, and a minimal document when possible. Do not open a public issue before a fix or mitigation is available.

The validator reads local repository files, performs no network requests, and does not execute Pawn code. Every remote schema identifier must have an offline copy under `schemas`.

Security fixes target the latest main branch and most recent tag. A compatible schema fix keeps the current major version; a breaking fix follows the RFC migration process.
