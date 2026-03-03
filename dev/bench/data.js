window.BENCHMARK_DATA = {
  "lastUpdate": 1772527300487,
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
          "id": "baf7a690ed99fd5d7fa48adfa087239304b922d0",
          "message": "feat(filter): name precision ratio bonus\n\nName precision bonus (+20% when ≥75% of name tokens covered by query):\n- Rewards items with concise names closely matching the query\n- 'jwt-auth' (2/2 = 100% coverage) beats 'jwt-auth-middleware-guide' (2/4 = 50%)\n  for query 'jwt auth'\n- Coverage computed via reverse pass: for each name token, check if any\n  query token matches (exact, stem, or prefix)\n- Only fires when name is non-empty and coverage >= 75% threshold\n\nNew constants: namePrecisionThreshold=75, namePrecisionBonus=20\nTests added: 3 new precision tests (concise wins, low coverage no bonus, single token)\npkg/filter coverage: 94.8% → 95.4%\n\nFull scoring pipeline now: tokenize → match → coverage bonus → co-occur bonus\n→ phrase adjacency → short-tag → name precision\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:46:45+01:00",
          "tree_id": "b5b5e629e3d984cf904d2d568be7f770c1950954",
          "url": "https://github.com/Arkestone/mcp/commit/baf7a690ed99fd5d7fa48adfa087239304b922d0"
        },
        "date": 1772484452168,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1437,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "821294 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1437,
            "unit": "ns/op",
            "extra": "821294 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "821294 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "821294 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1439,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "736543 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1439,
            "unit": "ns/op",
            "extra": "736543 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "736543 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "736543 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1438,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "758631 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1438,
            "unit": "ns/op",
            "extra": "758631 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "758631 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "758631 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10583,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "109533 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10583,
            "unit": "ns/op",
            "extra": "109533 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "109533 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "109533 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10867,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "108931 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10867,
            "unit": "ns/op",
            "extra": "108931 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "108931 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "108931 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10648,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "110204 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10648,
            "unit": "ns/op",
            "extra": "110204 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "110204 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "110204 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3646,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "300812 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3646,
            "unit": "ns/op",
            "extra": "300812 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "300812 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "300812 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3640,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "316779 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3640,
            "unit": "ns/op",
            "extra": "316779 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "316779 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "316779 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3674,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "321946 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3674,
            "unit": "ns/op",
            "extra": "321946 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "321946 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "321946 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 234645,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4719 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 234645,
            "unit": "ns/op",
            "extra": "4719 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4719 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4719 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 234683,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4737 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 234683,
            "unit": "ns/op",
            "extra": "4737 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4737 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4737 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 233649,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4870 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 233649,
            "unit": "ns/op",
            "extra": "4870 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4870 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4870 times\n4 procs"
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
          "id": "df2e7de1b55fa9349fa5201c4f5eb8b8abba1534",
          "message": "feat: frontmatter files: glob pattern for prompts and skills\n\nBoth mcp-prompts and mcp-skills now support a files: field in frontmatter\nthat restricts which file paths the item applies to — identical semantics\nto instructions' existing applyTo: field.\n\nFrontmatter example:\n  ---\n  description: TypeScript code reviewer\n  tags: [typescript, code-review]\n  files: \"**/*.ts\"\n  ---\n\n  or as a list:\n  files:\n    - \"**/*.ts\"\n    - \"**/*.tsx\"\n\nChanges:\n- Prompt.Files []string: parsed from frontmatter files: (string or list)\n- Skill.Files []string: parsed from frontmatter files: (string or list)\n- FilterByFilePath() added to both loader and scanner packages\n  - Items with no Files: always included (global scope)\n  - Items with Files: included only when at least one pattern matches\n  - Empty filePath: all items returned unchanged (backward compatible)\n- file_path parameter added to list-prompts, get-context, optimize-prompts,\n  list-skills, get-context (skills), optimize-skills tool handlers\n- Files exposed in list output (files field in ListEntry)\n- toStringSlice() helper added to scanner (same as instructions loader)\n\nTests: 5 new FilterByFilePath tests per package (unit + integration with\n  real frontmatter parsing); 21 parseFrontmatter call sites updated to\n  unpack new 5th return value\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T22:16:08+01:00",
          "tree_id": "698c4c4529570c98c416ee91f717ac8edebf24fe",
          "url": "https://github.com/Arkestone/mcp/commit/df2e7de1b55fa9349fa5201c4f5eb8b8abba1534"
        },
        "date": 1772486225898,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1374,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "854737 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1374,
            "unit": "ns/op",
            "extra": "854737 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "854737 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "854737 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1373,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "813100 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1373,
            "unit": "ns/op",
            "extra": "813100 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "813100 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "813100 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1400,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "806210 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1400,
            "unit": "ns/op",
            "extra": "806210 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "806210 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "806210 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10057,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "119416 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10057,
            "unit": "ns/op",
            "extra": "119416 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "119416 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "119416 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10004,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "118400 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10004,
            "unit": "ns/op",
            "extra": "118400 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "118400 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "118400 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10185,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "112606 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10185,
            "unit": "ns/op",
            "extra": "112606 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "112606 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "112606 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3446,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "326767 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3446,
            "unit": "ns/op",
            "extra": "326767 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "326767 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "326767 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3430,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "342988 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3430,
            "unit": "ns/op",
            "extra": "342988 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "342988 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "342988 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3422,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "338640 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3422,
            "unit": "ns/op",
            "extra": "338640 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "338640 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "338640 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 224734,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "5018 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 224734,
            "unit": "ns/op",
            "extra": "5018 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "5018 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "5018 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 224585,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "5179 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 224585,
            "unit": "ns/op",
            "extra": "5179 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "5179 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "5179 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 226395,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "5040 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 226395,
            "unit": "ns/op",
            "extra": "5040 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "5040 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "5040 times\n4 procs"
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
          "id": "a7f8955ad3bea6a0c3c54fb92cb5be83d4b96b21",
          "message": "chore: update changelogs for v1.1.0\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T23:33:46+01:00",
          "tree_id": "dc626441666a1d8f9759a8fd79b118ae84fce21c",
          "url": "https://github.com/Arkestone/mcp/commit/a7f8955ad3bea6a0c3c54fb92cb5be83d4b96b21"
        },
        "date": 1772490872695,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1478,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "811693 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1478,
            "unit": "ns/op",
            "extra": "811693 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "811693 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "811693 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1483,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "746631 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1483,
            "unit": "ns/op",
            "extra": "746631 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "746631 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "746631 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1488,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "732890 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1488,
            "unit": "ns/op",
            "extra": "732890 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "732890 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "732890 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10656,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "108320 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10656,
            "unit": "ns/op",
            "extra": "108320 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "108320 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "108320 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10772,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "110679 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10772,
            "unit": "ns/op",
            "extra": "110679 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "110679 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "110679 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10759,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "103339 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10759,
            "unit": "ns/op",
            "extra": "103339 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "103339 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "103339 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3663,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "310615 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3663,
            "unit": "ns/op",
            "extra": "310615 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "310615 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "310615 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3736,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "313639 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3736,
            "unit": "ns/op",
            "extra": "313639 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "313639 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "313639 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3867,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "315913 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3867,
            "unit": "ns/op",
            "extra": "315913 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "315913 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "315913 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 235952,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4731 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 235952,
            "unit": "ns/op",
            "extra": "4731 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4731 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4731 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 235622,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4652 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 235622,
            "unit": "ns/op",
            "extra": "4652 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4652 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4652 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 235380,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4880 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 235380,
            "unit": "ns/op",
            "extra": "4880 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4880 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4880 times\n4 procs"
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
          "id": "7735fd0df42f854d4ecf27eb15eab711617efbc2",
          "message": "fix: distinguish rate-limit 403 from private-repo 403 in GitHub client\n\nWhen calling the GitHub API without a token, public repos that exceed\nthe 60 req/hr unauthenticated rate limit also return HTTP 403. The\nprevious error message was misleading: 'is this a private repo?'.\n\nNow httpError reads the response body (JSON) and the X-RateLimit-Remaining\nheader to detect rate-limit errors and returns a clear message:\n  'HTTP 403 for <path>: rate limit exceeded — set GITHUB_TOKEN to increase\n   the rate limit (5000 req/hr vs 60 req/hr)'\n\nAuth/private-repo 403s (body does not mention 'rate limit') still produce\nthe original 'no GITHUB_TOKEN set — is this a private repo?' hint.\n\nAdded tests: TestFetchDir_HTTP403RateLimitViaBody,\nTestFetchDir_HTTP403RateLimitViaHeader, TestFetchFile_HTTP403RateLimitViaBody\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T23:55:49+01:00",
          "tree_id": "45344908c493139cd195484b94f544736ba01f41",
          "url": "https://github.com/Arkestone/mcp/commit/7735fd0df42f854d4ecf27eb15eab711617efbc2"
        },
        "date": 1772492206572,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1446,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "796056 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1446,
            "unit": "ns/op",
            "extra": "796056 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "796056 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "796056 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1454,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "762200 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1454,
            "unit": "ns/op",
            "extra": "762200 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "762200 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "762200 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1453,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "762051 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1453,
            "unit": "ns/op",
            "extra": "762051 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "762051 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "762051 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10797,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "107551 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10797,
            "unit": "ns/op",
            "extra": "107551 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "107551 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "107551 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10752,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "109978 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10752,
            "unit": "ns/op",
            "extra": "109978 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "109978 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "109978 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10692,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "108416 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10692,
            "unit": "ns/op",
            "extra": "108416 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "108416 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "108416 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3708,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "275550 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3708,
            "unit": "ns/op",
            "extra": "275550 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "275550 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "275550 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3705,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "312801 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3705,
            "unit": "ns/op",
            "extra": "312801 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "312801 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "312801 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3724,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "313318 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3724,
            "unit": "ns/op",
            "extra": "313318 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "313318 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "313318 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 233624,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4783 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 233624,
            "unit": "ns/op",
            "extra": "4783 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4783 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4783 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 232228,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 232228,
            "unit": "ns/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 236356,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4872 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 236356,
            "unit": "ns/op",
            "extra": "4872 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4872 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4872 times\n4 procs"
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
          "id": "2bacf53a966b719833be9fa0cfcee133eca39829",
          "message": "fix: async initial sync + better rate-limit error messages\n\npkg/syncer: make initial sync non-blocking\n- syncer.Start() now runs the first sync call inside the background\n  goroutine rather than synchronously in the caller goroutine.\n- Effect: HTTP server binds and responds to /healthz within ~300ms even\n  when the GitHub repo sync takes minutes (e.g. 175 files from\n  github/awesome-copilot).\n- Updated TestStartCallsFnImmediately to use a channel-based wait\n  instead of an immediate atomic read (the fn still runs 'immediately'\n  but asynchronously).\n\npkg/github: distinguish rate-limit 403 from private-repo 403\n- httpError now reads the response body (JSON) and the\n  X-RateLimit-Remaining header to detect rate-limit errors.\n- Rate-limited public repos now get a clear message:\n  'rate limit exceeded — set GITHUB_TOKEN to increase the rate limit'\n  instead of the misleading 'is this a private repo?' hint.\n- Added 3 new tests for rate-limit detection via body and header.\n\nTested in real conditions against github/awesome-copilot (178 instructions):\n- TypeScript file (src/auth/login.ts): 77/178 matched\n- Go file (pkg/server/handler.go): 46/178 matched\n- Terraform (infra/main.tf): 48/178 matched (includes azure-verified-modules-terraform)\n- Kubernetes YAML (k8s/deployment.yaml): 51/178 matched\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-03T00:08:37+01:00",
          "tree_id": "68f5cf182a220a1b76fe74d10e6d222efe6dea12",
          "url": "https://github.com/Arkestone/mcp/commit/2bacf53a966b719833be9fa0cfcee133eca39829"
        },
        "date": 1772492958446,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1468,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "837145 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1468,
            "unit": "ns/op",
            "extra": "837145 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "837145 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "837145 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1446,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "715782 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1446,
            "unit": "ns/op",
            "extra": "715782 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "715782 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "715782 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1454,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "722714 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1454,
            "unit": "ns/op",
            "extra": "722714 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "722714 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "722714 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10647,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "110574 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10647,
            "unit": "ns/op",
            "extra": "110574 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "110574 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "110574 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10918,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "111463 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10918,
            "unit": "ns/op",
            "extra": "111463 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "111463 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "111463 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 11637,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "95328 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 11637,
            "unit": "ns/op",
            "extra": "95328 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "95328 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "95328 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3719,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "305140 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3719,
            "unit": "ns/op",
            "extra": "305140 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "305140 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "305140 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3663,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "315099 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3663,
            "unit": "ns/op",
            "extra": "315099 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "315099 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "315099 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3666,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "317832 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3666,
            "unit": "ns/op",
            "extra": "317832 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "317832 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "317832 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 232615,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4813 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 232615,
            "unit": "ns/op",
            "extra": "4813 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4813 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4813 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 231579,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4909 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 231579,
            "unit": "ns/op",
            "extra": "4909 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4909 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4909 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 231642,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4868 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 231642,
            "unit": "ns/op",
            "extra": "4868 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4868 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4868 times\n4 procs"
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
          "id": "1ef541574f4281e2b9ada1b71413ea9aa3a80d55",
          "message": "chore: release v1.2.0\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-03T09:23:24+01:00",
          "tree_id": "fd20a68f28959d50094f812ea954576c17226f14",
          "url": "https://github.com/Arkestone/mcp/commit/1ef541574f4281e2b9ada1b71413ea9aa3a80d55"
        },
        "date": 1772526256613,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1436,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "800937 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1436,
            "unit": "ns/op",
            "extra": "800937 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "800937 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "800937 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1456,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "732219 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1456,
            "unit": "ns/op",
            "extra": "732219 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "732219 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "732219 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1447,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "736840 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1447,
            "unit": "ns/op",
            "extra": "736840 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "736840 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "736840 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10619,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "110847 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10619,
            "unit": "ns/op",
            "extra": "110847 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "110847 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "110847 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10857,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "110491 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10857,
            "unit": "ns/op",
            "extra": "110491 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "110491 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "110491 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10623,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "110011 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10623,
            "unit": "ns/op",
            "extra": "110011 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "110011 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "110011 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3689,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "312402 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3689,
            "unit": "ns/op",
            "extra": "312402 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "312402 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "312402 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3667,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "303519 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3667,
            "unit": "ns/op",
            "extra": "303519 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "303519 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "303519 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3705,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "310286 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3705,
            "unit": "ns/op",
            "extra": "310286 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "310286 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "310286 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 231332,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 231332,
            "unit": "ns/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 232143,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4809 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 232143,
            "unit": "ns/op",
            "extra": "4809 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4809 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4809 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 232451,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4698 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 232451,
            "unit": "ns/op",
            "extra": "4698 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4698 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4698 times\n4 procs"
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
          "id": "3017a7ef7dc2c4b82cfcc99c5f3bb7bfd94fec29",
          "message": "docs: MANDATORY RULE — never add Co-authored-by trailers to commits",
          "timestamp": "2026-03-03T09:30:08+01:00",
          "tree_id": "d5def8ea0cac6222720aa5b7c087742b52de2628",
          "url": "https://github.com/Arkestone/mcp/commit/3017a7ef7dc2c4b82cfcc99c5f3bb7bfd94fec29"
        },
        "date": 1772526656564,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1434,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "809800 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1434,
            "unit": "ns/op",
            "extra": "809800 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "809800 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "809800 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1439,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "785835 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1439,
            "unit": "ns/op",
            "extra": "785835 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "785835 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "785835 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1440,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "727555 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1440,
            "unit": "ns/op",
            "extra": "727555 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "727555 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "727555 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10728,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "110542 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10728,
            "unit": "ns/op",
            "extra": "110542 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "110542 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "110542 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10923,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "109213 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10923,
            "unit": "ns/op",
            "extra": "109213 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "109213 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "109213 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 11101,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "110284 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 11101,
            "unit": "ns/op",
            "extra": "110284 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "110284 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "110284 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3648,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "316156 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3648,
            "unit": "ns/op",
            "extra": "316156 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "316156 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "316156 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3678,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "309297 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3678,
            "unit": "ns/op",
            "extra": "309297 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "309297 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "309297 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3713,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "316610 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3713,
            "unit": "ns/op",
            "extra": "316610 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "316610 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "316610 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 238512,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 238512,
            "unit": "ns/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4780 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 242716,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4816 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 242716,
            "unit": "ns/op",
            "extra": "4816 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4816 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4816 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 240689,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4752 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 240689,
            "unit": "ns/op",
            "extra": "4752 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4752 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4752 times\n4 procs"
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
          "id": "a0932949970a410592e367c4d50955fcfd813179",
          "message": "docs: remove all Co-authored-by references from instructions",
          "timestamp": "2026-03-03T09:33:40+01:00",
          "tree_id": "b98b7018dc2d5167839def446d03482483e6f507",
          "url": "https://github.com/Arkestone/mcp/commit/a0932949970a410592e367c4d50955fcfd813179"
        },
        "date": 1772526860583,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1360,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "835756 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1360,
            "unit": "ns/op",
            "extra": "835756 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "835756 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "835756 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1424,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "827082 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1424,
            "unit": "ns/op",
            "extra": "827082 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "827082 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "827082 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1366,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "814504 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1366,
            "unit": "ns/op",
            "extra": "814504 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "814504 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "814504 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 9977,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "120184 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 9977,
            "unit": "ns/op",
            "extra": "120184 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "120184 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "120184 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10072,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "118621 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10072,
            "unit": "ns/op",
            "extra": "118621 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "118621 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "118621 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 9956,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "119456 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 9956,
            "unit": "ns/op",
            "extra": "119456 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "119456 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "119456 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3412,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "337615 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3412,
            "unit": "ns/op",
            "extra": "337615 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "337615 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "337615 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3461,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "344678 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3461,
            "unit": "ns/op",
            "extra": "344678 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "344678 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "344678 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3405,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "341994 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3405,
            "unit": "ns/op",
            "extra": "341994 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "341994 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "341994 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 226763,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "5023 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 226763,
            "unit": "ns/op",
            "extra": "5023 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "5023 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "5023 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 224922,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4929 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 224922,
            "unit": "ns/op",
            "extra": "4929 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4929 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4929 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 228434,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "5041 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 228434,
            "unit": "ns/op",
            "extra": "5041 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "5041 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "5041 times\n4 procs"
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
          "id": "5b3d1e5017daa9956701632baf5420f6413a367a",
          "message": "docs: clarify mandatory rule — forbid Copilot co-author trailers",
          "timestamp": "2026-03-03T09:34:15+01:00",
          "tree_id": "1c08bd80896847f48f2363976e68c4e8e5b43972",
          "url": "https://github.com/Arkestone/mcp/commit/5b3d1e5017daa9956701632baf5420f6413a367a"
        },
        "date": 1772526899680,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1440,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "813670 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1440,
            "unit": "ns/op",
            "extra": "813670 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "813670 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "813670 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1456,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "755300 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1456,
            "unit": "ns/op",
            "extra": "755300 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "755300 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "755300 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1447,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "735964 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1447,
            "unit": "ns/op",
            "extra": "735964 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "735964 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "735964 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10652,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "111597 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10652,
            "unit": "ns/op",
            "extra": "111597 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "111597 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "111597 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10727,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "111861 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10727,
            "unit": "ns/op",
            "extra": "111861 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "111861 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "111861 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10949,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "110740 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10949,
            "unit": "ns/op",
            "extra": "110740 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "110740 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "110740 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3655,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "315280 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3655,
            "unit": "ns/op",
            "extra": "315280 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "315280 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "315280 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3692,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "311247 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3692,
            "unit": "ns/op",
            "extra": "311247 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "311247 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "311247 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3737,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "315825 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3737,
            "unit": "ns/op",
            "extra": "315825 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "315825 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "315825 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 232935,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4743 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 232935,
            "unit": "ns/op",
            "extra": "4743 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4743 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4743 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 232723,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 232723,
            "unit": "ns/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4759 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 233338,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4801 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 233338,
            "unit": "ns/op",
            "extra": "4801 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4801 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4801 times\n4 procs"
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
          "id": "3de03a9f965e55a1d84062e8ec655946791c00c3",
          "message": "docs: enforce mandatory rule: no Copilot authorship trailers",
          "timestamp": "2026-03-03T09:34:15+01:00",
          "tree_id": "1c08bd80896847f48f2363976e68c4e8e5b43972",
          "url": "https://github.com/Arkestone/mcp/commit/3de03a9f965e55a1d84062e8ec655946791c00c3"
        },
        "date": 1772527128499,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1445,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "818232 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1445,
            "unit": "ns/op",
            "extra": "818232 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "818232 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "818232 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1457,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "741568 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1457,
            "unit": "ns/op",
            "extra": "741568 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "741568 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "741568 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1454,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "742568 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1454,
            "unit": "ns/op",
            "extra": "742568 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "742568 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "742568 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10944,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "109747 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10944,
            "unit": "ns/op",
            "extra": "109747 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "109747 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "109747 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10919,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "107904 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10919,
            "unit": "ns/op",
            "extra": "107904 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "107904 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "107904 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10907,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "106641 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10907,
            "unit": "ns/op",
            "extra": "106641 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "106641 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "106641 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3743,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "316832 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3743,
            "unit": "ns/op",
            "extra": "316832 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "316832 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "316832 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3743,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "295398 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3743,
            "unit": "ns/op",
            "extra": "295398 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "295398 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "295398 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3749,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "314802 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3749,
            "unit": "ns/op",
            "extra": "314802 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "314802 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "314802 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 237269,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4927 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 237269,
            "unit": "ns/op",
            "extra": "4927 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4927 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4927 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 236919,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4719 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 236919,
            "unit": "ns/op",
            "extra": "4719 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4719 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4719 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 238684,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4820 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 238684,
            "unit": "ns/op",
            "extra": "4820 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4820 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4820 times\n4 procs"
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
          "id": "32e5478c2b2691b6969d5ec6cbc53c3c176b6ca4",
          "message": "feat: initial release",
          "timestamp": "2026-03-03T09:40:44+01:00",
          "tree_id": "1c08bd80896847f48f2363976e68c4e8e5b43972",
          "url": "https://github.com/Arkestone/mcp/commit/32e5478c2b2691b6969d5ec6cbc53c3c176b6ca4"
        },
        "date": 1772527300028,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkTokenize",
            "value": 1442,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "834913 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1442,
            "unit": "ns/op",
            "extra": "834913 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "834913 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "834913 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1439,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "766950 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1439,
            "unit": "ns/op",
            "extra": "766950 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "766950 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "766950 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize",
            "value": 1440,
            "unit": "ns/op\t     728 B/op\t      21 allocs/op",
            "extra": "759088 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - ns/op",
            "value": 1440,
            "unit": "ns/op",
            "extra": "759088 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - B/op",
            "value": 728,
            "unit": "B/op",
            "extra": "759088 times\n4 procs"
          },
          {
            "name": "BenchmarkTokenize - allocs/op",
            "value": 21,
            "unit": "allocs/op",
            "extra": "759088 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10624,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "111774 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10624,
            "unit": "ns/op",
            "extra": "111774 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "111774 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "111774 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10655,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "111805 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10655,
            "unit": "ns/op",
            "extra": "111805 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "111805 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "111805 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm",
            "value": 10851,
            "unit": "ns/op\t    1328 B/op\t      48 allocs/op",
            "extra": "112004 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - ns/op",
            "value": 10851,
            "unit": "ns/op",
            "extra": "112004 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - B/op",
            "value": 1328,
            "unit": "B/op",
            "extra": "112004 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_MultiTerm - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "112004 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3684,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "316116 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3684,
            "unit": "ns/op",
            "extra": "316116 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "316116 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "316116 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3672,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "316687 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3672,
            "unit": "ns/op",
            "extra": "316687 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "316687 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "316687 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch",
            "value": 3649,
            "unit": "ns/op\t     552 B/op\t      28 allocs/op",
            "extra": "315921 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - ns/op",
            "value": 3649,
            "unit": "ns/op",
            "extra": "315921 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - B/op",
            "value": 552,
            "unit": "B/op",
            "extra": "315921 times\n4 procs"
          },
          {
            "name": "BenchmarkScore_NoMatch - allocs/op",
            "value": 28,
            "unit": "allocs/op",
            "extra": "315921 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 231165,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4771 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 231165,
            "unit": "ns/op",
            "extra": "4771 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4771 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4771 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 230278,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4818 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 230278,
            "unit": "ns/op",
            "extra": "4818 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4818 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4818 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items",
            "value": 230747,
            "unit": "ns/op\t   55640 B/op\t    2205 allocs/op",
            "extra": "4836 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - ns/op",
            "value": 230747,
            "unit": "ns/op",
            "extra": "4836 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - B/op",
            "value": 55640,
            "unit": "B/op",
            "extra": "4836 times\n4 procs"
          },
          {
            "name": "BenchmarkSortByScore_100Items - allocs/op",
            "value": 2205,
            "unit": "allocs/op",
            "extra": "4836 times\n4 procs"
          }
        ]
      }
    ]
  }
}