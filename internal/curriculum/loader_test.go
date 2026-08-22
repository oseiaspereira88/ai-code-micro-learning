package curriculum

import (
	"slices"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadValidPackBuildsQueryableCatalog(t *testing.T) {
	catalog, diags, err := Load("testdata/valid", DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %+v", diags)
	}
	if catalog == nil {
		t.Fatal("expected a catalog")
	}

	item, ok := catalog.Get("fixture.challenge-one")
	if !ok || item.Kind != KindChallenge {
		t.Fatalf("Get(challenge) = %+v, %v", item, ok)
	}

	challenges := catalog.List(KindChallenge, "slices")
	if len(challenges) != 1 || challenges[0].ID != "fixture.challenge-one" {
		t.Fatalf("List(challenge, slices) = %+v", challenges)
	}
	if len(catalog.List(KindChallenge, "no-such-theme")) != 0 {
		t.Fatal("theme filter should exclude unrelated challenges")
	}
}

func TestChallengeHidesReservedVariants(t *testing.T) {
	catalog, _, err := Load("testdata/valid", DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ch, ok := catalog.Challenge("fixture.challenge-one")
	if !ok {
		t.Fatal("expected fixture.challenge-one to exist")
	}
	if ch.Variants != nil {
		t.Fatalf("Variants should be hidden from the public view, got %v", ch.Variants)
	}
}

func TestLoadIsDeterministicRegardlessOfManifestOrder(t *testing.T) {
	c1, _, err := Load("testdata/valid", DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c2, _, err := Load("testdata/valid", DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a, _ := c1.Get("fixture.challenge-one")
	b, _ := c2.Get("fixture.challenge-one")
	if a != b {
		t.Fatalf("Load is not deterministic: %+v vs %+v", a, b)
	}
}

func TestLoadRejectsDuplicateID(t *testing.T) {
	assertBlockingDiagnostic(t, "testdata/duplicate_id", DiagDuplicateID)
}

func TestLoadRejectsMissingReference(t *testing.T) {
	catalog, diags, err := Load("testdata/missing_reference", DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if catalog == nil {
		t.Fatal("missing reference is non-blocking; catalog should still be built")
	}
	if !hasDiagnostic(diags, DiagMissingReference) {
		t.Fatalf("expected DiagMissingReference, got %+v", diags)
	}
}

func TestLoadDetectsPrerequisiteCycle(t *testing.T) {
	assertBlockingDiagnostic(t, "testdata/prerequisite_cycle", DiagPrerequisiteCycle)
}

func TestLoadRejectsCriteriaWithoutEvidence(t *testing.T) {
	assertBlockingDiagnostic(t, "testdata/criteria_without_evidence", DiagCriteriaWithoutEvidence)
}

func TestLoadRejectsIncompatibleSchemaVersion(t *testing.T) {
	assertBlockingDiagnostic(t, "testdata/incompatible_schema_version", DiagIncompatibleSchemaVersion)
}

func TestLoadRejectsPathEscapingPackDirectory(t *testing.T) {
	assertBlockingDiagnostic(t, "testdata/path_escape", DiagPathEscapesPack)
}

func TestLoadEnforcesFileSizeLimit(t *testing.T) {
	_, _, err := Load("testdata/valid", Limits{MaxFileBytes: 1, MaxNodeDepth: 32, MaxAliases: 64})
	if err == nil {
		t.Fatal("expected an error when the manifest itself exceeds the size limit")
	}
}

func TestCheckYAMLLimitsRejectsExcessiveDepth(t *testing.T) {
	node := nestedSequence(10)
	if err := checkYAMLLimits(node, Limits{MaxFileBytes: 1 << 20, MaxNodeDepth: 5, MaxAliases: 64}); err == nil {
		t.Fatal("expected a depth-limit error")
	}
	if err := checkYAMLLimits(node, Limits{MaxFileBytes: 1 << 20, MaxNodeDepth: 20, MaxAliases: 64}); err != nil {
		t.Fatalf("unexpected error within depth limit: %v", err)
	}
}

// nestedSequence builds a YAML sequence nested depth levels deep:
// [[[...]]]. Depth is what checkYAMLLimits bounds.
func nestedSequence(depth int) *yaml.Node {
	leaf := &yaml.Node{Kind: yaml.ScalarNode, Value: "leaf"}
	node := leaf
	for range depth {
		node = &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{node}}
	}
	return node
}

func assertBlockingDiagnostic(t *testing.T, dir string, code DiagnosticCode) {
	t.Helper()
	catalog, diags, err := Load(dir, DefaultLimits)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if catalog != nil {
		t.Fatal("expected no catalog when a blocking diagnostic is present")
	}
	if !hasDiagnostic(diags, code) {
		t.Fatalf("expected %s, got %+v", code, diags)
	}
}

func hasDiagnostic(diags []Diagnostic, code DiagnosticCode) bool {
	return slices.ContainsFunc(diags, func(d Diagnostic) bool { return d.Code == code })
}

// FuzzMalformedPackNeverPanics feeds arbitrary bytes through the exact path
// an adversarial pack file would take: YAML parse, security limits, decode
// into Pack, then Validate. It never touches the filesystem, so it runs at
// fuzzing speed. The only assertion is "no panic" — malformed input is
// expected to fail cleanly with a decode error or a diagnostic.
func FuzzMalformedPackNeverPanics(f *testing.F) {
	f.Add([]byte("schema_version: 1\nid: seed\nversion: 1.0.0\n"))
	f.Add([]byte("a: &x [*x]\n"))               // self-referential alias
	f.Add([]byte("a: &x\n  b: *x\n"))           // aliased mapping
	f.Add([]byte(""))                           // empty document
	f.Add([]byte("- - - - - - - - - - - - []")) // deep nesting, no aliases

	f.Fuzz(func(t *testing.T, data []byte) {
		var doc yaml.Node
		if err := yaml.Unmarshal(data, &doc); err != nil {
			return // malformed YAML is an expected, clean rejection
		}
		if err := checkYAMLLimits(&doc, DefaultLimits); err != nil {
			return // limits rejected it before any Decode ran
		}
		if len(doc.Content) == 0 {
			return // empty document
		}
		var pack Pack
		if doc.Decode(&pack) != nil {
			return // not a mapping shaped like a Pack
		}
		_ = Validate([]Pack{pack})
	})
}
