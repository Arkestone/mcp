window.BENCHMARK_DATA = {
  "lastUpdate": 1772484372935,
  "repoUrl": "https://github.com/Arkestone/mcp",
  "entries": {
    "Go Benchmarks": [
      {
        "commit": {
          "author": {
            "email": "5368160+Aadryn@users.noreply.github.com",
            "name": "aadryn",
            "username": "Aadryn"
          },
          "committer": {
            "email": "5368160+Aadryn@users.noreply.github.com",
            "name": "aadryn",
            "username": "Aadryn"
          },
          "distinct": true,
          "id": "1228642ab474af47f08f390c33da139e4ae2311c",
          "message": "test: add FilterByFilePath, FilterByQuery, FetchDirRecursive tests and filter benchmarks\n\n- FilterByFilePath: 7 tests (instructions loader) → 81.2% → 95.5% coverage\n- FilterByQuery: 5 tests + lifecycle (prompts loader) → 59.2% → 71.4% coverage\n- FetchDirRecursive: 5 tests with recursiveServer() helper → pkg/github 71.9% → 98.4%\n- Filter benchmarks: Tokenize, Score (multi-term / no-match), SortByScore(100 items)\n- Fix TestFilterByQuery_SortsByScore: use 'hammering utility' to avoid stem collision with 'build'\n- CI: macOS added to test matrix (stable only); coverage/Codecov gated on ubuntu-latest\n- nightly.yml: fuzz 30s/target, pprof memory profiling, HTML coverage archive 30d\n- outdated.yml: go list -u weekly check + mod tidy drift detection\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:32:48+01:00",
          "tree_id": "6a7a2afcee9d39000f9a58bc4840983fbea199ca",
          "url": "https://github.com/Arkestone/mcp/commit/1228642ab474af47f08f390c33da139e4ae2311c"
        },
        "date": 1772483634010,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1417,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "857030 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1417,
            "unit": "ns/op",
            "extra": "857030 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "857030 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "857030 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1423,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "744032 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1423,
            "unit": "ns/op",
            "extra": "744032 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "744032 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "744032 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1421,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "719538 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1421,
            "unit": "ns/op",
            "extra": "719538 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "719538 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "719538 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 7282,
            "unit": "ns/op\t    1280 B/op\t      47 allocs/op",
            "extra": "165026 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 7282,
            "unit": "ns/op",
            "extra": "165026 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1280,
            "unit": "B/op",
            "extra": "165026 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 47,
            "unit": "allocs/op",
            "extra": "165026 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 7146,
            "unit": "ns/op\t    1280 B/op\t      47 allocs/op",
            "extra": "163124 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 7146,
            "unit": "ns/op",
            "extra": "163124 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1280,
            "unit": "B/op",
            "extra": "163124 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 47,
            "unit": "allocs/op",
            "extra": "163124 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 7177,
            "unit": "ns/op\t    1280 B/op\t      47 allocs/op",
            "extra": "164401 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 7177,
            "unit": "ns/op",
            "extra": "164401 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1280,
            "unit": "B/op",
            "extra": "164401 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 47,
            "unit": "allocs/op",
            "extra": "164401 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3597,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "316644 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3597,
            "unit": "ns/op",
            "extra": "316644 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "316644 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "316644 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3598,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "317431 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3598,
            "unit": "ns/op",
            "extra": "317431 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "317431 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "317431 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3573,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "321456 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3573,
            "unit": "ns/op",
            "extra": "321456 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "321456 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "321456 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 159395,
            "unit": "ns/op\t   52440 B/op\t    2105 allocs/op",
            "extra": "6745 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 159395,
            "unit": "ns/op",
            "extra": "6745 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 52440,
            "unit": "B/op",
            "extra": "6745 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2105,
            "unit": "allocs/op",
            "extra": "6745 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 159735,
            "unit": "ns/op\t   52440 B/op\t    2105 allocs/op",
            "extra": "6722 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 159735,
            "unit": "ns/op",
            "extra": "6722 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 52440,
            "unit": "B/op",
            "extra": "6722 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2105,
            "unit": "allocs/op",
            "extra": "6722 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 159928,
            "unit": "ns/op\t   52440 B/op\t    2105 allocs/op",
            "extra": "6675 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 159928,
            "unit": "ns/op",
            "extra": "6675 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 52440,
            "unit": "B/op",
            "extra": "6675 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2105,
            "unit": "allocs/op",
            "extra": "6675 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "5368160+Aadryn@users.noreply.github.com",
            "name": "aadryn",
            "username": "Aadryn"
          },
          "committer": {
            "email": "5368160+Aadryn@users.noreply.github.com",
            "name": "aadryn",
            "username": "Aadryn"
          },
          "distinct": true,
          "id": "42a2d6c49f369a75c52d2fd77ad0c1e0746d1715",
          "message": "test: add syncRepo/syncAllRepos tests with httptest server\n\n- TestSyncRepo_DownloadsPromptFiles: verifies files cached from fake GitHub API\n- TestSyncRepo_APIError: verifies error propagation on 404\n- TestSyncAllRepos_NoRepos: verifies no-op on empty config\n- prompts loader coverage: 71.4% → 86.7%\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:34:22+01:00",
          "tree_id": "9687d21d8ddd693852bc6a8a2f67ec08537e7666",
          "url": "https://github.com/Arkestone/mcp/commit/42a2d6c49f369a75c52d2fd77ad0c1e0746d1715"
        },
        "date": 1772483707144,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1455,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "820965 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1455,
            "unit": "ns/op",
            "extra": "820965 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "820965 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "820965 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1462,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "734722 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1462,
            "unit": "ns/op",
            "extra": "734722 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "734722 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "734722 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1460,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "750614 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1460,
            "unit": "ns/op",
            "extra": "750614 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "750614 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "750614 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 7247,
            "unit": "ns/op\t    1280 B/op\t      47 allocs/op",
            "extra": "157504 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 7247,
            "unit": "ns/op",
            "extra": "157504 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1280,
            "unit": "B/op",
            "extra": "157504 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 47,
            "unit": "allocs/op",
            "extra": "157504 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 7271,
            "unit": "ns/op\t    1280 B/op\t      47 allocs/op",
            "extra": "150404 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 7271,
            "unit": "ns/op",
            "extra": "150404 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1280,
            "unit": "B/op",
            "extra": "150404 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 47,
            "unit": "allocs/op",
            "extra": "150404 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 7382,
            "unit": "ns/op\t    1280 B/op\t      47 allocs/op",
            "extra": "161438 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 7382,
            "unit": "ns/op",
            "extra": "161438 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1280,
            "unit": "B/op",
            "extra": "161438 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 47,
            "unit": "allocs/op",
            "extra": "161438 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3672,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "309993 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3672,
            "unit": "ns/op",
            "extra": "309993 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "309993 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "309993 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3776,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "311050 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3776,
            "unit": "ns/op",
            "extra": "311050 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "311050 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "311050 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3658,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "313453 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3658,
            "unit": "ns/op",
            "extra": "313453 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "313453 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "313453 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 165409,
            "unit": "ns/op\t   52440 B/op\t    2105 allocs/op",
            "extra": "6624 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 165409,
            "unit": "ns/op",
            "extra": "6624 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 52440,
            "unit": "B/op",
            "extra": "6624 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2105,
            "unit": "allocs/op",
            "extra": "6624 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 163819,
            "unit": "ns/op\t   52440 B/op\t    2105 allocs/op",
            "extra": "6787 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 163819,
            "unit": "ns/op",
            "extra": "6787 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 52440,
            "unit": "B/op",
            "extra": "6787 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2105,
            "unit": "allocs/op",
            "extra": "6787 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 164328,
            "unit": "ns/op\t   52440 B/op\t    2105 allocs/op",
            "extra": "6768 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 164328,
            "unit": "ns/op",
            "extra": "6768 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 52440,
            "unit": "B/op",
            "extra": "6768 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2105,
            "unit": "allocs/op",
            "extra": "6768 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "5368160+Aadryn@users.noreply.github.com",
            "name": "aadryn",
            "username": "Aadryn"
          },
          "committer": {
            "email": "5368160+Aadryn@users.noreply.github.com",
            "name": "aadryn",
            "username": "Aadryn"
          },
          "distinct": true,
          "id": "b5a1ddfce9f50c2de367aeffe98da6348717e12d",
          "message": "feat(filter): phrase adjacency bonus + short-tag exact match\n\nPhrase adjacency bonus (detectPhrases):\n- +25% when consecutive query tokens appear adjacent in name field\n  e.g. query 'code review' → name 'code-review' earns +25%\n  while 'code-style-review' (non-adjacent) does not\n- +15% per adjacent pair in tags, +10% in description\n- Uses stem matching so 'unit testing' matches 'unit-test' adjacency\n- Order-sensitive: 'code review' does not bonus 'review-code'\n\nShort-tag exact match bonus:\n- Raw tag string matched case-insensitively against raw query words\n- Enables short tags ('go', 'ci', 'api', 'k8s') that Tokenize drops (≤1-char)\n  to score properly — gives +8 pts per matching tag\n- 'go' query now correctly ranks items tagged 'Go' above untagged\n\nNew constants: phraseBonusName=25, phraseBonusTag=15, phraseBonusDesc=10, shortTagBonus=8\nNew function: detectPhrases(queryTokens, targets []string) int\nTests added: 6 new filter tests (3 phrase, 3 short-tag)\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:45:26+01:00",
          "tree_id": "8a264c42ae938a09aedcaa8fb10cda0909ec9c7a",
          "url": "https://github.com/Arkestone/mcp/commit/b5a1ddfce9f50c2de367aeffe98da6348717e12d"
        },
        "date": 1772484372031,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1372,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "814149 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1372,
            "unit": "ns/op",
            "extra": "814149 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "814149 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "814149 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1365,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "746672 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1365,
            "unit": "ns/op",
            "extra": "746672 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "746672 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "746672 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1427,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "747459 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1427,
            "unit": "ns/op",
            "extra": "747459 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "747459 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "747459 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 9239,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "129696 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 9239,
            "unit": "ns/op",
            "extra": "129696 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "129696 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "129696 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 9200,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "128433 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 9200,
            "unit": "ns/op",
            "extra": "128433 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "128433 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "128433 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 9268,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "127801 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 9268,
            "unit": "ns/op",
            "extra": "127801 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "127801 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "127801 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3396,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "338660 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3396,
            "unit": "ns/op",
            "extra": "338660 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "338660 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "338660 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3378,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "336542 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3378,
            "unit": "ns/op",
            "extra": "336542 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "336542 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "336542 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3374,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "331318 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3374,
            "unit": "ns/op",
            "extra": "331318 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "331318 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "331318 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 209055,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "5280 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 209055,
            "unit": "ns/op",
            "extra": "5280 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "5280 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "5280 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 206179,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "5232 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 206179,
            "unit": "ns/op",
            "extra": "5232 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "5232 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "5232 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 209517,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "5311 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 209517,
            "unit": "ns/op",
            "extra": "5311 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "5311 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "5311 times\n4 procs"
          }
        ]
      }
    ]
  }
}