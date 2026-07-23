#!/usr/bin/env bash

set -u

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_dir=$(CDPATH= cd -- "$script_dir/../.." && pwd)
corpus_dir=${PAWN_CORPUS:-"$repo_dir/../pawn-corpus"}
work_dir=$(mktemp -d)
trap 'rm -rf "$work_dir"' EXIT

if [ ! -d "$corpus_dir" ]; then
  echo "pawn-corpus not found at $corpus_dir" >&2
  echo "set PAWN_CORPUS to a pawn-corpus v0.1.5 checkout" >&2
  exit 2
fi

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
  source=$4
  source_dir=$(dirname -- "$source")
  source_name=$(basename -- "$source")
  compiler_dir=$(dirname -- "$compiler")
  library_dir=$compiler_dir
  if [ -d "$compiler_dir/../lib" ]; then
    library_dir=$(CDPATH= cd -- "$compiler_dir/../lib" && pwd)
  fi
  shift 4

  output="$work_dir/${name}.txt"
  artifact="$work_dir/${name}.amx"
  status=0
  (
    cd "$corpus_dir/$source_dir" || exit 1
    LD_LIBRARY_PATH="$library_dir${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}" \
      timeout 2s "$compiler" "$@" "$source_name" -o"$artifact"
  ) >"$output" 2>&1 || status=$?

  passed=false
  case "$expected" in
    pass) [ "$status" -eq 0 ] && passed=true ;;
    fail) [ "$status" -eq 1 ] && passed=true ;;
    timeout) [ "$status" -eq 124 ] && passed=true ;;
    error-at:*)
      location=${expected#error-at:}
      location_file=${location%:*}
      location_line=${location##*:}
      [ "$status" -eq 1 ] &&
        grep -Fq "$location_file($location_line)" "$output" &&
        passed=true
      ;;
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
  run_case "$compiler" at-global pass syntax/valid/compiler/at_global.pwn
  run_case "$compiler" at-identifier fail semantics/compiler_at_identifier.pwn
  run_case "$compiler" binary-literal pass lexer/compiler_binary_literal.pwn
  run_case "$compiler" block-shadowing warning:219 syntax/valid/compiler/block_shadowing.pwn
  run_case "$compiler" digit-separator pass lexer/compiler_digit_separator.pwn
  run_case "$compiler" hex-dollar-suffix fail lexer/compiler_hex_dollar_suffix.pwn
  run_case "$compiler" include-twice fail preprocessor/compiler_include_twice/main.pwn
  run_case "$compiler" include-twice pass preprocessor/compiler_include_twice/main.pwn -Z+
  run_case "$compiler" include-order pass preprocessor/compiler_include_order/main.pwn
  run_case "$compiler" include-location error-at:broken.inc:2 preprocessor/compiler_source_location/main.pwn
  run_case "$compiler" macro-recursion timeout preprocessor/compiler_macro_self_recursion.pwn
  run_case "$compiler" macro-redefinition warning:201 preprocessor/compiler_macro_redefinition.pwn
  run_case "$compiler" active-regions pass preprocessor/active_regions_nested.pwn
  run_case "$compiler" profile-openmp pass preprocessor/profile_openmp_define.pwn __OPEN_MP__=1
  run_case "$compiler" preproc-elif fail preprocessor/compiler_elif.pwn
  run_case "$compiler" preproc-elseif pass preprocessor/compiler_elseif.pwn
  run_case "$compiler" string-prefix pass lexer/compiler_string_prefix.pwn
  run_case "$compiler" tag-mismatch warning:213 semantics/compiler_tag_mismatch.pwn
  run_case "$compiler" tag-union pass semantics/compiler_tag_union.pwn
done

exit "$failures"
