// Package filter provides keyword relevance scoring for MCP item filtering.
//
// Design goals:
//   - Precise: word-boundary tokenization, not substring matching
//     ("go" does NOT match "google" or "cargo")
//   - Relevant: weighted scoring — name > tags > description; stem + prefix matching
//   - Reproducible: stable sort (score desc, then original order for ties)
//   - Efficient: stateless pure functions, no allocations beyond the sort slice
//
// Scoring pipeline:
//  1. Tokenize + stopword-filter query and all fields
//  2. Match each query token per field via exact → stem → prefix
//  3. Coverage bonus: +10% per additional matched query term (beyond first)
//  4. Co-occurrence bonus: +30% per field where 2+ query terms match
//  5. Phrase adjacency bonus: +25%/+15%/+10% per adjacent pair in name/tags/desc
//  6. Short-tag exact bonus: raw tag string match for short tags ("go", "ci"…)
package filter

import (
	"sort"
	"strings"
	"unicode"
)

// Field matching weights.
const (
	weightNameExact  = 10
	weightNamePrefix = 5
	weightTagExact   = 8
	weightTagPrefix  = 3
	weightDescExact  = 3
	weightDescPrefix = 1

	// minPrefixLen is the minimum token length for prefix/stem-prefix matching.
	minPrefixLen = 3

	// cooccurBonus is the extra percentage when 2+ query terms match within the same field.
	cooccurBonus = 30

	// phraseBonusName/Tag/Desc is applied per adjacent query-token pair detected
	// in that field (e.g. query "code review" → name "code-review" earns +25%).
	phraseBonusName = 25
	phraseBonusTag  = 15
	phraseBonusDesc = 10

	// shortTagBonus is added (flat points) when a raw tag string exactly matches
	// a query token, enabling short tags like "go", "ci", "api" to score even
	// though Tokenize drops ≤1-char tokens.
	shortTagBonus = weightTagExact

	// namePrecisionThreshold is the minimum fraction of name tokens that must be
	// covered by query tokens to earn the precision bonus.
	namePrecisionThreshold = 75 // 75% — as integer hundredths

	// namePrecisionBonus is applied (+%) when name token coverage >= threshold.
	// Rewards concise names that closely match the query over verbose long names.
	namePrecisionBonus = 20
)

// stopwords are high-frequency tokens that carry little semantic signal.
// Filtering them prevents "how", "use", "add", "the" from inflating scores
// and diluting the relevance of real query terms.
var stopwords = map[string]bool{
	// 2-char function words (prepositions, articles, conjunctions)
	"to": true, "of": true, "in": true, "on": true, "at": true,
	"is": true, "it": true, "if": true, "do": true, "be": true,
	"as": true, "an": true, "or": true, "by": true,
	// 3+ char stopwords
	"the": true, "and": true, "for": true, "with": true,
	"use": true, "add": true, "get": true, "set": true,
	"from": true, "into": true, "via": true, "how": true,
	"that": true, "this": true, "are": true, "was": true,
	"not": true, "but": true, "all": true, "one": true,
	"can": true, "you": true, "your": true, "any": true,
	"has": true, "its": true, "also": true, "used": true,
	"will": true, "when": true, "than": true, "each": true,
}

// stemSuffixes are stripped (longest-first) to normalise inflected forms.
// A root must keep at least 4 characters after stripping.
var stemSuffixes = []string{
	"ations", "ation", "tions", "tion", "ments", "ment",
	"ings", "ated", "ates", "ing", "ions", "ion",
	"ers", "ors", "eds", "ed", "es", "er", "or",
}

// stem strips common English suffixes for stem-level matching.
// Examples: "authentication"→"authentic", "instructions"→"instruct", "testing"→"test".
func stem(s string) string {
	for _, sfx := range stemSuffixes {
		if strings.HasSuffix(s, sfx) {
			if root := s[:len(s)-len(sfx)]; len(root) >= 4 {
				return root
			}
		}
	}
	// Strip trailing 's' only for words long enough to avoid over-stemming.
	if len(s) >= 5 && strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}
	return s
}

// Score computes a relevance score for name/description/tags against query.
//
// The scoring pipeline:
//  1. Tokenize + stopword-filter query and all fields
//  2. Match each query token per field via exact → stem → prefix
//  3. Coverage bonus: +10% per additional matched query term (beyond first)
//  4. Co-occurrence bonus: +30% per field where 2+ query terms match together
//  5. Phrase adjacency bonus: +25%/+15%/+10% when adjacent query terms
//     appear adjacent in name/tags/desc ("code review" → "code-review" +25%)
//  6. Short-tag bonus: raw tag == raw query word gives +8 (enables "go", "ci")
//  7. Name precision bonus: +20% when ≥75% of name tokens are covered by query
//     (rewards "jwt-auth" over "jwt-auth-middleware-guide" for query "jwt auth")
//
// Returns 0 if no query term matches (item excluded).
// Returns 1 if query is empty or reduces to all stopwords (pass-through).
func Score(query, name, description string, tags []string) int {
	if query == "" {
		return 1
	}
	queryTokens := Tokenize(query)
	if len(queryTokens) == 0 {
		return 1
	}

	nameTokens := Tokenize(name)
	descTokens := Tokenize(description)
	tagTokenSets := make([][]string, len(tags))
	for i, tag := range tags {
		tagTokenSets[i] = Tokenize(tag)
	}

	totalScore := 0
	termsMatched := 0
	nameMatched := 0
	tagMatched := 0
	descMatched := 0

	for _, qt := range queryTokens {
		nameHit := matchTokens(qt, nameTokens, weightNameExact, weightNamePrefix)
		tagHit := 0
		for _, tagTokens := range tagTokenSets {
			tagHit += matchTokens(qt, tagTokens, weightTagExact, weightTagPrefix)
		}
		descHit := matchTokens(qt, descTokens, weightDescExact, weightDescPrefix)

		termScore := nameHit + tagHit + descHit
		if termScore > 0 {
			termsMatched++
			totalScore += termScore
			if nameHit > 0 {
				nameMatched++
			}
			if tagHit > 0 {
				tagMatched++
			}
			if descHit > 0 {
				descMatched++
			}
		}
	}

	if totalScore == 0 {
		return 0
	}

	// Coverage bonus: +10% per additional matched query term (beyond first).
	if termsMatched > 1 {
		totalScore = totalScore * (100 + 10*(termsMatched-1)) / 100
	}

	// Co-occurrence bonus: +30% for each field where 2+ query terms matched.
	// Rewards items where multiple concepts of the query appear together.
	cooccur := 0
	if nameMatched >= 2 {
		cooccur++
	}
	if tagMatched >= 2 {
		cooccur++
	}
	if descMatched >= 2 {
		cooccur++
	}
	if cooccur > 0 {
		totalScore = totalScore * (100 + cooccurBonus*cooccur) / 100
	}

	// Phrase adjacency bonus: rewards items where consecutive query tokens
	// appear adjacent in the same field (phrase-level matching).
	// "code review" query → name "code-review" gets +25%, while
	// "code-style-review" (non-adjacent) does not.
	if np := detectPhrases(queryTokens, nameTokens); np > 0 {
		totalScore = totalScore * (100 + phraseBonusName*np) / 100
	}
	tagPhrases := 0
	for _, tagTokens := range tagTokenSets {
		tagPhrases += detectPhrases(queryTokens, tagTokens)
	}
	if tagPhrases > 0 {
		totalScore = totalScore * (100 + phraseBonusTag*tagPhrases) / 100
	}
	if dp := detectPhrases(queryTokens, descTokens); dp > 0 {
		totalScore = totalScore * (100 + phraseBonusDesc*dp) / 100
	}

	// Short-tag exact bonus: raw tag strings can be short ("go", "ci", "api")
	// and are dropped by Tokenize's ≥2-char filter. Match them against raw
	// query tokens before tokenization for full short-tag coverage.
	for _, rawQuery := range strings.Fields(strings.ToLower(query)) {
		for _, tag := range tags {
			if strings.ToLower(tag) == rawQuery {
				totalScore += shortTagBonus
				break
			}
		}
	}

	// Name precision bonus: if ≥75% of name tokens are covered by the query,
	// reward the item for having a concise, on-topic name.
	// Example: query "jwt auth" → name "jwt-auth" (2/2 = 100% ≥ 75% → +20%)
	//          query "jwt auth" → name "jwt-auth-middleware-guide" (2/4 = 50% → no bonus)
	if len(nameTokens) > 0 {
		covered := 0
		for _, nt := range nameTokens {
			ntStem := stem(nt)
			for _, qt := range queryTokens {
				if nt == qt || stem(qt) == ntStem || strings.HasPrefix(nt, qt) {
					covered++
					break
				}
			}
		}
		if covered*100/len(nameTokens) >= namePrecisionThreshold {
			totalScore = totalScore * (100 + namePrecisionBonus) / 100
		}
	}

	return totalScore
}

// detectPhrases counts how many consecutive query-token pairs appear adjacent
// (at neighbouring positions, using stem matching) in the target token list.
// Each query pair is counted at most once, even if it appears multiple times.
//
// Example: query tokens ["code","review"], targets ["code","review","guide"]
// → 1 (the pair code+review is adjacent at positions 0,1).
func detectPhrases(queryTokens, targets []string) int {
	if len(queryTokens) < 2 || len(targets) < 2 {
		return 0
	}
	n := 0
	for i := 0; i < len(queryTokens)-1; i++ {
		q1 := stem(queryTokens[i])
		q2 := stem(queryTokens[i+1])
		for j := 0; j < len(targets)-1; j++ {
			if stem(targets[j]) == q1 && stem(targets[j+1]) == q2 {
				n++
				break // count each query pair once
			}
		}
	}
	return n
}

// matchTokens returns the best score for a query token against a target token set.
// Priority: exact match > stem match > prefix/stem-prefix match.
func matchTokens(qt string, targets []string, exactW, prefixW int) int {
	qStem := stem(qt)
	best := 0
	for _, t := range targets {
		// 1. Exact word match — highest confidence, return immediately.
		if t == qt {
			return exactW
		}
		// 2. Stem match — normalised forms count as exact.
		//    "instructions"↔"instruction", "testing"↔"test", "deployment"↔"deploy"
		tStem := stem(t)
		if qStem == tStem && exactW > best {
			best = exactW
			continue
		}
		// 3. Prefix / stem-prefix match — partial match at word boundary.
		//    Forward: query is a prefix of target ("auth" → "authentication")
		//    Reverse: query's stem is broader, target's stem is a prefix of it
		//             ("authentication" query → "auth" target via stem)
		if len(qt) >= minPrefixLen && strings.HasPrefix(t, qt) && prefixW > best {
			best = prefixW
			continue
		}
		if len(qStem) >= minPrefixLen && len(tStem) >= minPrefixLen {
			if (strings.HasPrefix(tStem, qStem) || strings.HasPrefix(qStem, tStem)) && prefixW > best {
				best = prefixW
			}
		}
	}
	return best
}

// Tokenize splits s into lowercase alphanumeric tokens on non-alphanumeric
// boundaries. Single-character tokens and stopwords are discarded.
//
// Examples:
//
//	"JWT Authentication"      → ["jwt", "authentication"]
//	"how to use mcp-skills"   → ["mcp", "skills"]    (stopwords filtered)
//	"mcp-instructions"        → ["mcp", "instructions"]
//	"**/*.go"                 → ["go"]
func Tokenize(s string) []string {
	s = strings.ToLower(s)
	var tokens []string
	var cur strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cur.WriteRune(r)
		} else if cur.Len() > 0 {
			if cur.Len() > 1 {
				if tok := cur.String(); !stopwords[tok] {
					tokens = append(tokens, tok)
				}
			}
			cur.Reset()
		}
	}
	if cur.Len() > 1 {
		if tok := cur.String(); !stopwords[tok] {
			tokens = append(tokens, tok)
		}
	}
	return tokens
}

// SortByScore filters and sorts items by relevance to query.
//
//   - Items scoring 0 are excluded when query is non-empty.
//   - Items are sorted by score descending.
//   - Equal-score items keep their original relative order (stable, reproducible).
//   - If query is empty, all items are returned in original order unchanged.
func SortByScore[T any](items []T, scoreFn func(T) int) []T {
	type ranked struct {
		item  T
		score int
		idx   int
	}

	candidates := make([]ranked, 0, len(items))
	for i, item := range items {
		if s := scoreFn(item); s > 0 {
			candidates = append(candidates, ranked{item, s, i})
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score != candidates[j].score {
			return candidates[i].score > candidates[j].score
		}
		return candidates[i].idx < candidates[j].idx // stable tie-break by insertion order
	})

	out := make([]T, len(candidates))
	for i, r := range candidates {
		out[i] = r.item
	}
	return out
}
