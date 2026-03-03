package filter_test

import (
	"testing"

	"github.com/Arkestone/mcp/pkg/filter"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"JWT Authentication", []string{"jwt", "authentication"}},
		{"mcp-instructions", []string{"mcp", "instructions"}},
		{"**/*.go", []string{"go"}},
		{"go", []string{"go"}}, // 2-char kept; false positives prevented by minPrefixLen=3 in Score
		{"k8s", []string{"k8s"}},
		{"", nil},
		// Stopwords are filtered out.
		{"how to use jwt auth", []string{"jwt", "auth"}},
		{"add and get with the set", nil},
	}
	for _, tt := range tests {
		got := filter.Tokenize(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("Tokenize(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range tt.want {
			if got[i] != tt.want[i] {
				t.Errorf("Tokenize(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestScore_EmptyQuery(t *testing.T) {
	s := filter.Score("", "any name", "any desc", []string{"tag"})
	if s != 1 {
		t.Errorf("Score with empty query = %d, want 1", s)
	}
}

// A query that reduces entirely to stopwords is treated as empty (pass-through).
func TestScore_AllStopwordsQuery(t *testing.T) {
	s := filter.Score("how to use the", "jwt auth", "token auth", []string{"auth"})
	if s != 1 {
		t.Errorf("all-stopword query should return 1 (pass-through), got %d", s)
	}
}

func TestScore_NoMatch(t *testing.T) {
	s := filter.Score("python", "Go development", "Build Go projects", []string{"golang", "build"})
	if s != 0 {
		t.Errorf("Score = %d, want 0", s)
	}
}

func TestScore_Precision(t *testing.T) {
	// "go" should NOT match "google", "cargo", "goroutine"
	s := filter.Score("go", "Google Cloud", "use cargo for rust", []string{"goroutine"})
	if s != 0 {
		t.Errorf("Score(go) matched false positives: %d", s)
	}
}

func TestScore_ExactMatchHigher(t *testing.T) {
	sExact := filter.Score("jwt", "jwt", "", nil)
	sPrefix := filter.Score("jw", "jwt", "", nil)
	if sExact <= sPrefix {
		t.Errorf("exact score %d should be > prefix score %d", sExact, sPrefix)
	}
}

func TestScore_NameHigherThanDesc(t *testing.T) {
	sName := filter.Score("auth", "auth", "", nil)
	sDesc := filter.Score("auth", "something", "auth description", nil)
	if sName <= sDesc {
		t.Errorf("name score %d should be > desc score %d", sName, sDesc)
	}
}

func TestScore_CoverageBonus(t *testing.T) {
	// Matching 2 terms should score higher than matching 1 term.
	s1 := filter.Score("jwt", "jwt auth", "token management", []string{"security"})
	s2 := filter.Score("jwt auth", "jwt auth", "token management", []string{"security"})
	if s2 <= s1 {
		t.Errorf("2-term match (%d) should be > 1-term match (%d)", s2, s1)
	}
}

// Stemming: inflected forms should match their common stem.
func TestScore_Stemming(t *testing.T) {
	tests := []struct {
		query string
		name  string
		desc  string
		want  string // "match" or "no-match"
	}{
		// Plural / singular normalisation
		{"instruction", "instructions", "", "match"},
		{"instructions", "instruction", "", "match"},
		// -ing / -ed / -ion
		{"testing", "test", "run test cases", "match"},
		{"deployment", "deploy", "deploy your service", "match"},
		{"authentication", "auth", "jwt auth middleware", "match"},
	}
	for _, tt := range tests {
		s := filter.Score(tt.query, tt.name, tt.desc, nil)
		if tt.want == "match" && s == 0 {
			t.Errorf("Score(%q, %q, %q) = 0, expected non-zero (stem match)", tt.query, tt.name, tt.desc)
		}
		if tt.want == "no-match" && s != 0 {
			t.Errorf("Score(%q, %q, %q) = %d, expected 0", tt.query, tt.name, tt.desc, s)
		}
	}
}

// Co-occurrence bonus: matching 2+ terms in the same field scores higher
// than matching the same terms spread across different fields.
func TestScore_CooccurrenceBonus(t *testing.T) {
	// Both "jwt" and "auth" in name → co-occurrence bonus.
	sBoth := filter.Score("jwt auth", "jwt auth service", "", nil)
	// "jwt" in name, "auth" only in description → no co-occurrence bonus.
	sSplit := filter.Score("jwt auth", "jwt service", "requires auth token", nil)
	if sBoth <= sSplit {
		t.Errorf("co-occurrence in same field (%d) should be > split across fields (%d)", sBoth, sSplit)
	}
}

func TestSortByScore_Order(t *testing.T) {
	type item struct {
		name string
	}
	items := []item{{"low"}, {"high"}, {"medium"}}
	scores := map[string]int{"low": 1, "high": 10, "medium": 5}
	sorted := filter.SortByScore(items, func(it item) int { return scores[it.name] })

	want := []string{"high", "medium", "low"}
	for i, s := range want {
		if sorted[i].name != s {
			t.Errorf("sorted[%d] = %q, want %q", i, sorted[i].name, s)
		}
	}
}

func TestSortByScore_Stable(t *testing.T) {
	// Equal-score items must keep their original relative order.
	type item struct{ name string }
	items := []item{{"a"}, {"b"}, {"c"}} // all same score
	sorted := filter.SortByScore(items, func(item) int { return 5 })
	for i, want := range []string{"a", "b", "c"} {
		if sorted[i].name != want {
			t.Errorf("stable order violated: sorted[%d] = %q, want %q", i, sorted[i].name, want)
		}
	}
}

func TestSortByScore_ExcludesZero(t *testing.T) {
	type item struct{ name string }
	items := []item{{"match"}, {"nomatch"}}
	sorted := filter.SortByScore(items, func(it item) int {
		if it.name == "match" {
			return 5
		}
		return 0
	})
	if len(sorted) != 1 || sorted[0].name != "match" {
		t.Errorf("expected only 'match', got %v", sorted)
	}
}


// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkTokenize(b *testing.B) {
s := "JWT Authentication middleware for secure API access with OAuth2 and OpenID Connect support"
b.ResetTimer()
for range b.N {
_ = filter.Tokenize(s)
}
}

func BenchmarkScore_MultiTerm(b *testing.B) {
query := "jwt authentication security"
name := "JWT Auth Middleware"
desc := "Secure your API endpoints with JWT token authentication and role-based access control"
tags := []string{"jwt", "security", "authentication", "middleware"}
b.ResetTimer()
for range b.N {
_ = filter.Score(query, name, desc, tags)
}
}

func BenchmarkScore_NoMatch(b *testing.B) {
query := "kubernetes deployment"
name := "JWT Auth Middleware"
desc := "Secure your API with tokens"
tags := []string{"jwt", "security"}
b.ResetTimer()
for range b.N {
_ = filter.Score(query, name, desc, tags)
}
}

func BenchmarkSortByScore_100Items(b *testing.B) {
type item struct{ name, desc string }
items := make([]item, 100)
for i := range items {
items[i] = item{name: "skill", desc: "description of a skill that does something useful"}
}
b.ResetTimer()
for range b.N {
_ = filter.SortByScore(items, func(it item) int {
return filter.Score("skill description", it.name, it.desc, nil)
})
}
}

// ---------------------------------------------------------------------------
// Phrase adjacency bonus tests
// ---------------------------------------------------------------------------

func TestScore_PhraseAdjacentInName(t *testing.T) {
// Query "code review" → item with name "code-review" should score higher
// than item with name "code-style-review" (non-adjacent tokens).
adjacent := filter.Score("code review", "code-review", "", nil)
nonAdj := filter.Score("code review", "code-style-review", "", nil)
if adjacent <= nonAdj {
t.Errorf("adjacent phrase should score higher: adjacent=%d, non-adjacent=%d", adjacent, nonAdj)
}
}

func TestScore_PhraseAdjacentStemmed(t *testing.T) {
// "unit testing" phrase should match "unit-test" via stemming
// (stem("testing")="test" == stem("test")="test").
score := filter.Score("unit testing", "unit-test", "", nil)
if score == 0 {
t.Error("stemmed phrase should still score > 0")
}
// And it should score higher than non-adjacent tokens.
other := filter.Score("unit testing", "testing-unit-helpers", "", nil)
if score <= other {
t.Errorf("adjacent stemmed phrase should outscore non-adjacent: phrase=%d other=%d", score, other)
}
}

func TestScore_PhraseNotAdjacent(t *testing.T) {
// "code review" with reversed order "review-code" should NOT get phrase bonus.
forward := filter.Score("code review", "code-review", "", nil)
reversed := filter.Score("code review", "review-code", "", nil)
// Both score > 0 (both terms hit name), but forward phrase gets extra bonus.
if forward <= reversed {
t.Errorf("forward phrase order should score higher: forward=%d reversed=%d", forward, reversed)
}
}

// ---------------------------------------------------------------------------
// Short-tag exact match tests
// ---------------------------------------------------------------------------

func TestScore_ShortTagExact(t *testing.T) {
// Tag "go" is too short for Tokenize, but the short-tag bonus should fire.
withTag := filter.Score("go", "my-lib", "", []string{"go"})
withoutTag := filter.Score("go", "my-lib", "", []string{"python"})
if withTag <= 0 {
t.Errorf("short tag 'go' should score > 0, got %d", withTag)
}
if withTag <= withoutTag {
t.Errorf("matching short tag should outscore non-matching: with=%d without=%d", withTag, withoutTag)
}
}

func TestScore_ShortTagCaseInsensitive(t *testing.T) {
lower := filter.Score("go", "my-lib", "", []string{"Go"})
if lower == 0 {
t.Error("short tag match should be case-insensitive")
}
}

func TestScore_ShortTagQueryCI(t *testing.T) {
// Query "CI" (uppercased) should match tag "ci".
score := filter.Score("CI", "pipeline", "", []string{"ci"})
if score == 0 {
t.Errorf("query 'CI' should match tag 'ci', got %d", score)
}
}

// ---------------------------------------------------------------------------
// Name precision bonus tests
// ---------------------------------------------------------------------------

func TestScore_NamePrecision_ConciseWins(t *testing.T) {
// Query "jwt auth": name "jwt-auth" (2 tokens, 100% covered)
// vs name "jwt-auth-middleware-guide" (4 tokens, 50% covered).
concise := filter.Score("jwt auth", "jwt-auth", "", nil)
verbose := filter.Score("jwt auth", "jwt-auth-middleware-guide", "", nil)
if concise <= verbose {
t.Errorf("concise name should outscore verbose: concise=%d verbose=%d", concise, verbose)
}
}

func TestScore_NamePrecision_LowCoverageNoBonus(t *testing.T) {
// "go testing" → name "go-testing-framework-complete-guide" has low coverage.
// Score should still be > 0 (terms match) but no precision bonus.
score := filter.Score("go testing", "go-testing-framework-complete-guide", "", nil)
if score == 0 {
t.Error("low-precision name should still score > 0 when terms match")
}
}

func TestScore_NamePrecision_SingleToken(t *testing.T) {
// Single-token name that matches query → 100% coverage → precision bonus applies.
single := filter.Score("auth", "auth", "", nil)
multi := filter.Score("auth", "auth-helper-utils-extended", "", nil)
if single <= multi {
t.Errorf("single precise name should outscore long name: single=%d multi=%d", single, multi)
}
}
