window.BENCHMARK_DATA = {
  "lastUpdate": 1772457950556,
  "repoUrl": "https://github.com/Arkestone/mcp",
  "entries": {
    "Binary Sizes": [
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
        "date": 1772457949114,
        "tool": "customSmallerIsBetter",
        "benches": [
          {
            "name": "mcp-instructions",
            "value": 8.274,
            "unit": "MB"
          },
          {
            "name": "mcp-skills",
            "value": 8.61,
            "unit": "MB"
          },
          {
            "name": "mcp-adr",
            "value": 8.602,
            "unit": "MB"
          },
          {
            "name": "mcp-memory",
            "value": 8.285,
            "unit": "MB"
          },
          {
            "name": "mcp-prompts",
            "value": 8.61,
            "unit": "MB"
          },
          {
            "name": "mcp-graph",
            "value": 7.914,
            "unit": "MB"
          }
        ]
      }
    ]
  }
}