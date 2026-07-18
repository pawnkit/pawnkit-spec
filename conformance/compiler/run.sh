#!/usr/bin/env bash

set -u

case_dir=$(CDPATH= cd -- "$(dirname -- "$0")/cases" && pwd)
work_dir=$(mktemp -d)
trap 'rm -rf "$work_dir"' EXIT

if [ "$#" -gt 0 ]; then
  compilers=("$@")
elif [ -n "${PAWNCC:-}" ]; then
  compilers=("$PAWNCC")
else
  echo "usage: PAWNCC=/path/to/pawncc $0" >&2
  echo "   or: $0 /path/to/pawncc [...]" >&2
  exit 2
fi

failures=0

run_case() {
  compiler=$1
  name=$2
  expected=$3
  shift 3

  output="$work_dir/${name}.txt"
  artifact="$work_dir/${name}.amx"
  status=0
  (
    cd "$case_dir" || exit 1
    LD_LIBRARY_PATH="$(dirname -- "$compiler")${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}" \
      timeout 2s "$compiler" "$@" "$name.pwn" -o"$artifact"
  ) >"$output" 2>&1 || status=$?

  passed=false
  case "$expected" in
    pass) [ "$status" -eq 0 ] && passed=true ;;
    fail) [ "$status" -ne 0 ] && [ "$status" -ne 124 ] && passed=true ;;
    timeout) [ "$status" -eq 124 ] && passed=true ;;
    warning:*)
      code=${expected#warning:}
      [ "$status" -eq 0 ] && grep -Eq "warning 0*$code:" "$output" && passed=true
      ;;
  esac

  if $passed; then
    printf 'ok   %-22s %s\n' "$name" "$expected"
  else
    printf 'FAIL %-22s expected %s, exit %s\n' "$name" "$expected" "$status"
    sed 's/^/     /' "$output"
    failures=$((failures + 1))
  fi
}

for compiler in "${compilers[@]}"; do
  compiler=$(realpath "$compiler")
  echo "== $compiler =="
  run_case "$compiler" at-global pass
  run_case "$compiler" at-identifier fail
  run_case "$compiler" binary-literal pass
  run_case "$compiler" block-shadowing warning:219
  run_case "$compiler" digit-separator pass
  run_case "$compiler" hex-dollar-suffix fail
  run_case "$compiler" include-twice fail
  run_case "$compiler" include-twice pass -Z+
  run_case "$compiler" macro-recursion timeout
  run_case "$compiler" macro-redefinition warning:201
  run_case "$compiler" preproc-elif fail
  run_case "$compiler" preproc-elseif pass
  run_case "$compiler" string-prefix pass
  run_case "$compiler" tag-mismatch warning:213
  run_case "$compiler" tag-union pass
done

exit "$failures"
