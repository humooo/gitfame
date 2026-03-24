package internal

import "testing"

func TestBuildLanguageFilter_UnknownLanguagesDoNotRestrict(t *testing.T) {
	mapping := []MappingEntity{{Name: "Go"}, {Name: "Markdown"}}

	filter, unknown := buildLanguageFilter(mapping, []string{"unknownlang"})
	if len(filter) != 0 {
		t.Fatalf("expected empty filter for unknown-only languages, got %v", filter)
	}
	if len(unknown) != 1 || unknown[0] != "unknownlang" {
		t.Fatalf("unexpected unknown languages: %v", unknown)
	}

	if !IsAcceptableLanguage("", filter) {
		t.Fatalf("empty filter must not restrict files")
	}
}

func TestBuildLanguageFilter_KnownAndUnknown(t *testing.T) {
	mapping := []MappingEntity{{Name: "Go"}, {Name: "Markdown"}}

	filter, unknown := buildLanguageFilter(mapping, []string{"go", "unknownlang"})
	if len(filter) != 1 {
		t.Fatalf("expected single known language in filter, got %v", filter)
	}
	if _, ok := filter["go"]; !ok {
		t.Fatalf("expected go to be present in filter")
	}
	if len(unknown) != 1 || unknown[0] != "unknownlang" {
		t.Fatalf("unexpected unknown languages: %v", unknown)
	}

	if !IsAcceptableLanguage("Go", filter) {
		t.Fatalf("expected Go to match language filter")
	}
	if IsAcceptableLanguage("Markdown", filter) {
		t.Fatalf("did not expect Markdown to match go-only filter")
	}
}

func TestMatchesPatterns_InvalidPatternReturnsError(t *testing.T) {
	_, err := MatchesPatterns("foo.go", []string{"["})
	if err == nil {
		t.Fatalf("expected invalid glob pattern error")
	}
}
