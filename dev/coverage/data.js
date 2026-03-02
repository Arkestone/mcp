window.BENCHMARK_DATA = {
  "lastUpdate": 1772458144354,
  "repoUrl": "https://github.com/Arkestone/mcp",
  "entries": {
    "Test Coverage": [
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
          "id": "b9ba161cef6e9579ed927cf2c825ec91e7c336d3",
          "message": "fix(workflows): fix 3 runtime failures found in verification\n\ncontainer-structure.yml:\n- Fix 'field metadataTests not found in type v2.StructureTest' — rename\n  metadataTests (plural list) to metadataTest (singular object) per v2 schema\n- Add '/bin/sh shouldExist: false' to explicitly verify distroless has no shell\n\nbenchmark-trend.yml / coverage-trend.yml / binary-size.yml:\n- Add 'Bootstrap gh-pages if missing' step before benchmark-action to handle\n  first-run case where gh-pages branch doesn't exist yet\n- Add concurrency group to binary-size.yml (gh-pages-binary-sizes,\n  cancel-in-progress: false) matching the pattern of the other two\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T14:25:17+01:00",
          "tree_id": "acd8c17bbf8296898678cbe416560ae991bf20c8",
          "url": "https://github.com/Arkestone/mcp/commit/b9ba161cef6e9579ed927cf2c825ec91e7c336d3"
        },
        "date": 1772457952897,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 84.2,
            "unit": "%"
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
          "id": "bce961c8fd5d91a98578b7738b070d4c6aacf78b",
          "message": "fix(ci): align coverage threshold to actual baseline (84% not 85%)\n\nCoverage has been stable at 84.2% since project inception. The 85%\nthreshold was aspirational but blocks CI. Lowering to 84% to reflect\nreality — coverage improvements should be tracked via coverage-trend.yml.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T14:28:33+01:00",
          "tree_id": "250c5db32cd9c3d5f001280e83b452ad2b4ae6f7",
          "url": "https://github.com/Arkestone/mcp/commit/bce961c8fd5d91a98578b7738b070d4c6aacf78b"
        },
        "date": 1772458143600,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 84.2,
            "unit": "%"
          }
        ]
      }
    ]
  }
}