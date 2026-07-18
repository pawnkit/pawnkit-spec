// Command validate checks pawnkit-spec's schemas, profiles, examples,
// conformance documents, and RFC front matter without network access.
// Directories are identified by base name, so argument order doesn't matter.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var validStatuses = map[string]bool{
	"draft":        true,
	"experimental": true,
	"accepted":     true,
	"deprecated":   true,
	"superseded":   true,
}

type failure struct {
	file   string
	reason string
}

func (f failure) String() string {
	return fmt.Sprintf("%s: %s", f.file, f.reason)
}

type validator struct {
	schemasDir       string
	profilesDir      string
	examplesDir      string
	conformanceDir   string
	rfcsDir          string
	compiled         map[string]*jsonschema.Schema // name -> compiled
	failures         []failure
	documentsChecked int
}

func main() {
	start := time.Now()
	v := &validator{compiled: map[string]*jsonschema.Schema{}}

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: validate <dir>...")
		os.Exit(2)
	}

	for _, arg := range os.Args[1:] {
		abs, err := filepath.Abs(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid path %q: %v\n", arg, err)
			os.Exit(2)
		}
		switch filepath.Base(abs) {
		case "schemas":
			v.schemasDir = abs
		case "profiles":
			v.profilesDir = abs
		case "examples":
			v.examplesDir = abs
		case "conformance":
			v.conformanceDir = abs
		case "rfcs":
			v.rfcsDir = abs
		default:
			fmt.Fprintf(os.Stderr, "unrecognized directory (expected schemas/profiles/examples/conformance/rfcs): %s\n", abs)
			os.Exit(2)
		}
	}

	if v.schemasDir != "" {
		v.loadSchemas()
	}
	if v.conformanceDir != "" {
		v.checkConformanceSchema()
	}
	if v.profilesDir != "" {
		v.checkProfiles()
	}
	if v.examplesDir != "" {
		v.checkExamples()
	}
	if v.rfcsDir != "" {
		v.checkRFCs()
	}

	elapsed := time.Since(start)
	if len(v.failures) > 0 {
		sort.Slice(v.failures, func(i, j int) bool { return v.failures[i].file < v.failures[j].file })
		fmt.Fprintln(os.Stderr, "FAIL:")
		for _, f := range v.failures {
			fmt.Fprintln(os.Stderr, "  "+f.String())
		}
		fmt.Fprintf(os.Stderr, "\n%d failure(s) across %d document(s) checked in %s\n", len(v.failures), v.documentsChecked, elapsed)
		os.Exit(1)
	}

	fmt.Printf("ok: validated %d documents in %s\n", v.documentsChecked, elapsed)
}

func (v *validator) loadSchemas() {
	entries, err := os.ReadDir(v.schemasDir)
	if err != nil {
		v.fail(v.schemasDir, fmt.Sprintf("cannot read schemas dir: %v", err))
		return
	}

	seenIDs := map[string]string{} // $id -> file
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".schema.json") {
			continue
		}
		path := filepath.Join(v.schemasDir, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			v.fail(path, fmt.Sprintf("cannot read: %v", err))
			continue
		}
		var doc map[string]any
		if err := json.Unmarshal(raw, &doc); err != nil {
			v.fail(path, fmt.Sprintf("invalid JSON: %v", err))
			continue
		}
		id, _ := doc["$id"].(string)
		if id == "" {
			v.fail(path, "missing $id")
			continue
		}
		if !strings.HasPrefix(id, "https://schemas.pawnkit.dev/") {
			v.fail(path, fmt.Sprintf("$id %q does not use the https://schemas.pawnkit.dev/ convention", id))
		}
		if prev, ok := seenIDs[id]; ok {
			v.fail(path, fmt.Sprintf("duplicate $id %q also used by %s", id, prev))
		}
		seenIDs[id] = e.Name()

		schemaVal, _ := doc["$schema"].(string)
		if !strings.Contains(schemaVal, "2020-12") {
			v.fail(path, fmt.Sprintf("$schema %q is not JSON Schema 2020-12", schemaVal))
		}

		if err := c.AddResource(id, bytes.NewReader(raw)); err != nil {
			v.fail(path, fmt.Sprintf("cannot register schema: %v", err))
			continue
		}
		v.documentsChecked++
	}

	for id, file := range seenIDs {
		sch, err := c.Compile(id)
		if err != nil {
			v.fail(file, fmt.Sprintf("does not compile as JSON Schema 2020-12: %v", err))
			continue
		}
		// matched against examples/<name>/ by directory name
		name := strings.TrimSuffix(file, ".schema.json")
		v.compiled[name] = sch
	}
}

func mustReDecode(raw []byte) any {
	// jsonschema needs numbers as json.Number, not float64
	var v any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		panic(err)
	}
	return v
}

func (v *validator) checkConformanceSchema() {
	path := filepath.Join(v.conformanceDir, "expected-results.schema.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		v.fail(path, fmt.Sprintf("cannot read: %v", err))
		return
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		v.fail(path, fmt.Sprintf("invalid JSON: %v", err))
		return
	}
	id, _ := doc["$id"].(string)
	if id == "" {
		v.fail(path, "missing $id")
		return
	}
	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020
	if err := c.AddResource(id, bytes.NewReader(raw)); err != nil {
		v.fail(path, fmt.Sprintf("cannot register schema: %v", err))
		return
	}
	sch, err := c.Compile(id)
	if err != nil {
		v.fail(path, fmt.Sprintf("does not compile: %v", err))
		return
	}
	v.compiled["conformance"] = sch
	v.documentsChecked++

	manifestPath := filepath.Join(v.conformanceDir, "manifest.md")
	if _, err := os.Stat(manifestPath); err != nil {
		v.fail(manifestPath, "conformance/manifest.md is required and missing")
	}
}

func (v *validator) checkProfiles() {
	sch := v.compiled["pawn-profile"]
	if sch == nil {
		v.fail(v.profilesDir, "schemas/pawn-profile.schema.json was not loaded/compiled; cannot validate profiles")
		return
	}
	entries, err := os.ReadDir(v.profilesDir)
	if err != nil {
		v.fail(v.profilesDir, fmt.Sprintf("cannot read profiles dir: %v", err))
		return
	}
	seenIDs := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(v.profilesDir, e.Name())
		v.validateAgainst(path, sch)

		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var doc map[string]any
		if err := json.Unmarshal(raw, &doc); err != nil {
			continue
		}
		id, _ := doc["id"].(string)
		if id == "" {
			continue
		}
		if prev, ok := seenIDs[id]; ok {
			v.fail(path, fmt.Sprintf("duplicate profile id %q also used by %s", id, prev))
		}
		seenIDs[id] = e.Name()
	}
}

func (v *validator) checkExamples() {
	entries, err := os.ReadDir(v.examplesDir)
	if err != nil {
		v.fail(v.examplesDir, fmt.Sprintf("cannot read examples dir: %v", err))
		return
	}

	seenDiagnosticCodes := map[string]string{}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		var sch *jsonschema.Schema
		if name == "conformance" {
			sch = v.compiled["conformance"]
		} else {
			sch = v.compiled[name]
		}
		if sch == nil {
			v.fail(filepath.Join(v.examplesDir, name), fmt.Sprintf("no compiled schema found for examples/%s (expected schemas/%s.schema.json)", name, name))
			continue
		}
		dir := filepath.Join(v.examplesDir, name)
		files, err := os.ReadDir(dir)
		if err != nil {
			v.fail(dir, fmt.Sprintf("cannot read: %v", err))
			continue
		}
		found := false
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
				continue
			}
			found = true
			path := filepath.Join(dir, f.Name())
			v.validateAgainst(path, sch)

			if name == "pawn-diagnostic" {
				raw, err := os.ReadFile(path)
				if err == nil {
					var doc map[string]any
					if json.Unmarshal(raw, &doc) == nil {
						if code, ok := doc["code"].(string); ok {
							if prev, ok := seenDiagnosticCodes[code]; ok {
								v.fail(path, fmt.Sprintf("duplicate diagnostic code %q also used by %s", code, prev))
							}
							seenDiagnosticCodes[code] = f.Name()
						}
					}
				}
			}
		}
		if !found {
			v.fail(dir, "no example .json files found")
		}
	}
}

func (v *validator) validateAgainst(path string, sch *jsonschema.Schema) {
	raw, err := os.ReadFile(path)
	if err != nil {
		v.fail(path, fmt.Sprintf("cannot read: %v", err))
		return
	}
	if len(raw) > 1<<20 {
		v.fail(path, "exceeds 1 MiB size limit (docs/performance.md)")
		return
	}
	inst := mustReDecode(raw)
	if err := sch.Validate(inst); err != nil {
		v.fail(path, fmt.Sprintf("schema validation failed: %v", err))
		return
	}
	v.documentsChecked++
}

var frontMatterRe = regexp.MustCompile(`(?s)\A---\n(.*?)\n---\n`)
var dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
var rfcNumRe = regexp.MustCompile(`^\d{4}-`)

func (v *validator) checkRFCs() {
	entries, err := os.ReadDir(v.rfcsDir)
	if err != nil {
		v.fail(v.rfcsDir, fmt.Sprintf("cannot read rfcs dir: %v", err))
		return
	}
	seenNumbers := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(v.rfcsDir, e.Name())
		v.documentsChecked++
		raw, err := os.ReadFile(path)
		if err != nil {
			v.fail(path, fmt.Sprintf("cannot read: %v", err))
			continue
		}
		isTemplate := e.Name() == "0000-template.md"

		m := frontMatterRe.FindSubmatch(raw)
		if m == nil {
			v.fail(path, "missing YAML front matter (--- ... ---) block")
			continue
		}
		fm, err := parseFrontMatter(string(m[1]))
		if err != nil {
			v.fail(path, fmt.Sprintf("cannot parse front matter: %v", err))
			continue
		}

		required := []string{"rfc", "title", "status", "created", "updated", "supersedes", "superseded-by", "schema"}
		for _, key := range required {
			if _, ok := fm[key]; !ok {
				v.fail(path, fmt.Sprintf("front matter missing required key %q", key))
			}
		}

		if !rfcNumRe.MatchString(e.Name()) {
			v.fail(path, "filename does not start with a 4-digit RFC number")
		} else if rfcVal, ok := fm["rfc"]; ok {
			wantNum := e.Name()[:4]
			if rfcVal != wantNum && rfcVal != strings.TrimLeft(wantNum, "0") {
				v.fail(path, fmt.Sprintf("front matter rfc=%q does not match filename number %q", rfcVal, wantNum))
			}
			if prev, ok := seenNumbers[wantNum]; ok {
				v.fail(path, fmt.Sprintf("duplicate RFC number %q also used by %s", wantNum, prev))
			}
			seenNumbers[wantNum] = e.Name()
		}

		if status, ok := fm["status"]; ok && !isTemplate {
			if !validStatuses[status] {
				v.fail(path, fmt.Sprintf("status %q is not one of draft/experimental/accepted/deprecated/superseded", status))
			}
		}

		if !isTemplate {
			for _, key := range []string{"created", "updated"} {
				if val, ok := fm[key]; ok && !dateRe.MatchString(val) {
					v.fail(path, fmt.Sprintf("%s=%q is not YYYY-MM-DD", key, val))
				}
			}
		}

		requiredSections := []string{
			"## Summary", "## Motivation", "## Compatibility impact",
			"## Alternatives considered", "## Security considerations",
			"## Migration plan", "## Reference implementation status",
			"## Conformance tests", "## Open questions",
		}
		if isTemplate {
			requiredSections = append(requiredSections, "## Current behavior", "## Proposal")
		}
		body := string(raw)
		for _, section := range requiredSections {
			if !strings.Contains(body, section) {
				v.fail(path, fmt.Sprintf("missing required section %q", section))
			}
		}
	}
}

// minimal "key: value" parser; not a full YAML implementation
func parseFrontMatter(block string) (map[string]string, error) {
	out := map[string]string{}
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimRight(line, " \t")
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx < 0 {
			return nil, fmt.Errorf("malformed line %q", line)
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"`)
		out[key] = val
	}
	return out, nil
}

func (v *validator) fail(file, reason string) {
	v.failures = append(v.failures, failure{file: file, reason: reason})
}
