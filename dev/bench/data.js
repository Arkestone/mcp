window.BENCHMARK_DATA = {
  "lastUpdate": 1772483634677,
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
      }
    ]
  }
}