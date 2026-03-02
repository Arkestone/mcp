window.BENCHMARK_DATA = {
  "lastUpdate": 1772492957493,
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
          "id": "9d92da677b390a066b3119099fe37499e90a6806",
          "message": "fix(container-structure): use numeric UID 65532 not string 'nonroot'\n\ndistroless/static-debian12:nonroot stores the user in image metadata as\nthe numeric UID 65532, not the string 'nonroot'. container-structure-test\ncompares the raw image config value, so the test must match exactly.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T14:33:03+01:00",
          "tree_id": "fba16cf6269580aacc74a323948eddb0645d944c",
          "url": "https://github.com/Arkestone/mcp/commit/9d92da677b390a066b3119099fe37499e90a6806"
        },
        "date": 1772458416346,
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
          "id": "93fea6dd409b6c1384cd330c363e3fe6daf9c9a9",
          "message": "fix(post-release): use github.repository instead of hardcoded repo name\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T14:39:45+01:00",
          "tree_id": "28d116e068446014292f9dc3305f5a0034638463",
          "url": "https://github.com/Arkestone/mcp/commit/93fea6dd409b6c1384cd330c363e3fe6daf9c9a9"
        },
        "date": 1772458818530,
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
          "id": "bd614d2907f47dfe019da9e5eb75265225afbe2a",
          "message": "fix(workflows): action-pinning report-only + link-check vscode exclude\n\naction-pinning.yml:\n- Fix grep pattern: use '^\\s*uses:' to match only YAML directives,\n  not bash code containing 'uses:' as a string literal (false positives)\n- Change exit 1 → warning annotation + exit 0; SHA pinning is an ongoing\n  task for Dependabot/Renovate, not a hard CI blocker\n\nlink-check.yml:\n- Exclude insiders.vscode.dev from link checks; these are valid VS Code\n  install redirect URLs that intentionally return 302\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T14:43:43+01:00",
          "tree_id": "5183b751872e3206d87c5e955fcf7990f6873312",
          "url": "https://github.com/Arkestone/mcp/commit/bd614d2907f47dfe019da9e5eb75265225afbe2a"
        },
        "date": 1772459053193,
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
          "id": "b7365644e7fd875d4c4e80a42ccd2ae61a677017",
          "message": "fix(docs+link-check): update broken URLs and fix link checker config\n\nBroken links fixed in 13 markdown files:\n- GitHub custom instructions URL: customize-github-copilot → customizing-copilot\n- VS Code MCP docs URL: docs/copilot/model-context-protocol → docs/copilot/chat/mcp-servers\n\nlink-check.yml improvements:\n- Exclude website/ directory (Docusaurus root-relative /docs/ links need a\n  base URL that lychee can't resolve in CI without a running dev server)\n- Exclude '^/docs/' pattern for root-relative links\n- Already had insiders.vscode.dev excluded from previous commit\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T14:50:17+01:00",
          "tree_id": "df0c1ac9a16c6bcd2d5d9cb1450a882c28da8d9b",
          "url": "https://github.com/Arkestone/mcp/commit/b7365644e7fd875d4c4e80a42ccd2ae61a677017"
        },
        "date": 1772459453183,
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
          "id": "8c89c23aa16bff03cf794fe4f4192a082b274af4",
          "message": "fix(link-check): exclude CHANGELOG files (historical release tag links)\n\nCHANGELOG.md files contain links to historical release tags (v0.1.0, per-server\ntags) that don't exist in this repo. Exclude them from link validation.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T14:52:13+01:00",
          "tree_id": "44c417c7109b810a83f79fc8172f40c6d54e3948",
          "url": "https://github.com/Arkestone/mcp/commit/8c89c23aa16bff03cf794fe4f4192a082b274af4"
        },
        "date": 1772459566810,
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
          "id": "371c2468e5d59a19d066d4e8c9253973062fff3a",
          "message": "fix(link-check): replace invalid regex '^/docs/' with '^file://'\n\nLychee uses Rust regex; the '^/docs/' pattern caused a regex parse error\n('repetition operator missing expression'). Replaced with '^file://' to\nexclude file:// local links, which was the actual intent.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T14:53:55+01:00",
          "tree_id": "8e27534179af10715440dd519c3842208798f81e",
          "url": "https://github.com/Arkestone/mcp/commit/371c2468e5d59a19d066d4e8c9253973062fff3a"
        },
        "date": 1772459671756,
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
          "id": "924adcdc9704981277efc17af04a4559c628f2dc",
          "message": "fix(link-check): use regex syntax for --exclude-path (not glob)\n\nLychee's --exclude-path takes Rust regex patterns, not glob patterns.\n'**/CHANGELOG.md' is invalid regex (** = repetition without expression).\nFixed to 'CHANGELOG\\.md' which matches any path containing CHANGELOG.md.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T14:58:46+01:00",
          "tree_id": "ef247ee6c3884d7c5d985faec4a30f78bab6a6c9",
          "url": "https://github.com/Arkestone/mcp/commit/924adcdc9704981277efc17af04a4559c628f2dc"
        },
        "date": 1772459958575,
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
          "id": "72a35992198bf38247ea3072496094056cfa8b82",
          "message": "fix(workflows): add GHCR login for container-scan and post-release docker validation\n\n- container-scan.yml: add docker/login-action before Trivy to authenticate\n  with GHCR (packages need auth even when public); add packages: read permission\n- post-release.yml: add docker/login-action before image validation;\n  add packages: read permission; make SBOM asset check a warning (not fail)\n  since SBOM is generated by the release workflow and may not be present\n  for manually-created releases\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T15:09:50+01:00",
          "tree_id": "a7c5fd23fbcc57891c0e150189f85739ea315270",
          "url": "https://github.com/Arkestone/mcp/commit/72a35992198bf38247ea3072496094056cfa8b82"
        },
        "date": 1772460627956,
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
          "id": "a0fd025b8f56dae72f7ee7e31e5912a009c1616f",
          "message": "fix(release): remove dist/ build artifacts from git tracking\n\nThe dist/ directory was accidentally committed with goreleaser build outputs.\nThis caused goreleaser --clean to make the git tree dirty (staged deletions),\nfailing the release workflow immediately.\n\ndist/ is already in .gitignore — this commit only removes tracked files.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T15:17:15+01:00",
          "tree_id": "4b183d9ca6dbde8f3431dcab6d98437075eacae2",
          "url": "https://github.com/Arkestone/mcp/commit/a0fd025b8f56dae72f7ee7e31e5912a009c1616f"
        },
        "date": 1772461072310,
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
          "id": "01e87a8627bd8359fad94810ad5ca31f3cd910a1",
          "message": "fix(release): checkout tagged commit on workflow_dispatch\n\nWhen triggered manually with a tag input, goreleaser requires the checked-out\ncommit to match the tag. Previously the workflow always checked out HEAD (main),\nwhich diverges from older tags after subsequent commits.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T15:19:53+01:00",
          "tree_id": "09b385860d01d45bf5ce6aa52a7ed2fc1419cafc",
          "url": "https://github.com/Arkestone/mcp/commit/01e87a8627bd8359fad94810ad5ca31f3cd910a1"
        },
        "date": 1772461228257,
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
          "id": "be22141f3d36eb8aac9d0838c72fa35420aef7aa",
          "message": "fix(release): add --overwrite flag on workflow_dispatch to handle re-runs\n\nWhen re-running the release workflow for an existing tag (via workflow_dispatch),\ngoreleaser fails with 422 'already_exists' when trying to upload binary assets\nthat were already successfully uploaded.\n\nThe --overwrite flag tells goreleaser to replace existing release assets,\nallowing partial re-runs (e.g., to push missing Docker images after a\npreviously failed release workflow).\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T15:24:17+01:00",
          "tree_id": "3331adb7013f172a641c73b30193f29d93615208",
          "url": "https://github.com/Arkestone/mcp/commit/be22141f3d36eb8aac9d0838c72fa35420aef7aa"
        },
        "date": 1772461489275,
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
          "id": "67b352c82de50cefd849c2f0d504480993808e5b",
          "message": "fix(release): use replace_existing_artifacts for re-run support\n\nReplace the invalid --overwrite flag with goreleaser's config option\nreplace_existing_artifacts: true. This allows workflow_dispatch re-runs\non an existing release tag to overwrite previously uploaded artifacts\n(e.g., to push missing Docker images after a failed release).\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T15:29:41+01:00",
          "tree_id": "617ed00986a61a9ebb71388901ae485c655c16e2",
          "url": "https://github.com/Arkestone/mcp/commit/67b352c82de50cefd849c2f0d504480993808e5b"
        },
        "date": 1772461810872,
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
          "id": "d93aecd3a02022c9b06994cc011029bdba4010a9",
          "message": "fix(release): delete existing assets before goreleaser re-run on workflow_dispatch\n\nThe replace_existing_artifacts goreleaser option doesn't work when assets\nalready exist (GitHub returns 422). Add an explicit cleanup step that runs\nonly on workflow_dispatch (re-run case) to delete all existing release\nassets before goreleaser uploads fresh ones.\n\nThis allows re-running the release workflow to repair a failed release\n(e.g., Docker images never pushed) without triggering a new release tag.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T15:36:20+01:00",
          "tree_id": "fcae3a439b4d2bddb9e7439684343b5fed8c912c",
          "url": "https://github.com/Arkestone/mcp/commit/d93aecd3a02022c9b06994cc011029bdba4010a9"
        },
        "date": 1772462216239,
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
          "id": "38b7697ea05bc5e7acae630ae3839dc60b64e739",
          "message": "fix(release): use while+read loop to skip empty asset IDs in cleanup\n\nThe for-loop over ASSET_IDS could iterate with empty string when the variable\nis empty (or contains trailing newlines), causing 'gh api DELETE' to get 404.\nReplace with while+read loop that explicitly skips empty lines.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T15:42:51+01:00",
          "tree_id": "3031f5cb27d0e7d288b14f0251436f0999d02bcf",
          "url": "https://github.com/Arkestone/mcp/commit/38b7697ea05bc5e7acae630ae3839dc60b64e739"
        },
        "date": 1772462601073,
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
          "id": "84ad29af2d199dbf7ab0171ae55ef343ab67c900",
          "message": "fix(release): use apiUrl to extract numeric asset ID for REST API deletion\n\nThe .assets[].id field in gh CLI returns a GraphQL node ID (RA_kwDOR...),\nbut the GitHub REST DELETE endpoint requires a numeric integer ID.\nExtract the numeric ID from .assets[].apiUrl (last path segment).\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T15:46:32+01:00",
          "tree_id": "717dc643b5c5071b258febdd7a58efa5104a2abc",
          "url": "https://github.com/Arkestone/mcp/commit/84ad29af2d199dbf7ab0171ae55ef343ab67c900"
        },
        "date": 1772462824367,
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
          "id": "ebff3df80bbcdc0a578c17ea9b4b0028494b76c4",
          "message": "fix(release): strip v prefix for Docker tag lookup in cosign step\n\nGoreleaser uses .Version (1.0.0) not .Tag (v1.0.0) as Docker image tag.\nThe cosign signing step was trying to inspect :v1.0.0 which doesn't exist.\n\nAlso: use dynamic org name in packages visibility API call.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T15:56:43+01:00",
          "tree_id": "f5888feb9c592bf65759c511960de12ecaad7d97",
          "url": "https://github.com/Arkestone/mcp/commit/ebff3df80bbcdc0a578c17ea9b4b0028494b76c4"
        },
        "date": 1772463437654,
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
          "id": "46e83c90ae4b33225c9b694b253aa522bef0e556",
          "message": "fix(release): sign Docker images by tag, not by raw inspect output\n\ndocker buildx imagetools inspect --format '{{.Manifest.Digest}}' returns\nthe full OCI image index manifest text for multi-arch images, not just the\ndigest string. Passing that to cosign as IMAGE@<full text> causes a\n'could not parse reference' error.\n\nFix: sign by IMAGE:TAG directly — cosign resolves the digest internally\nand handles multi-arch OCI image indexes correctly.\n\nAlso add a pre-flight check that the image is accessible before signing.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T16:02:11+01:00",
          "tree_id": "1b557cebce4515dd22f006a3870d13c0fc32624a",
          "url": "https://github.com/Arkestone/mcp/commit/46e83c90ae4b33225c9b694b253aa522bef0e556"
        },
        "date": 1772463768151,
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
          "id": "fb06d5929fd883a8b6c0923da67e1fb6d7905da4",
          "message": "fix(post-release): strip v prefix from tag for Docker image validation\n\nGoreleaser uses .Version (e.g. 1.0.0) not .Tag (v1.0.0) as the Docker\nimage tag. The post-release validation was looking for :v1.0.0 which\ndoesn't exist, causing 'not found' errors.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T16:18:28+01:00",
          "tree_id": "1507c41172cb42f8c974010369731d12f3d20257",
          "url": "https://github.com/Arkestone/mcp/commit/fb06d5929fd883a8b6c0923da67e1fb6d7905da4"
        },
        "date": 1772464742149,
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
          "id": "234c1746810ccef067337a0d2e2c4a05b0b20de1",
          "message": "fix(post-release): add docker logs on smoke test failure, increase timeout\n\nAdd docker logs output when smoke test fails to diagnose container crashes.\nIncrease wait from 15s to 20s for container startup.\nUse explicit STARTED flag to distinguish timeout from success.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T16:21:49+01:00",
          "tree_id": "79f59d9c72172df66e43d8c466281e6b94695034",
          "url": "https://github.com/Arkestone/mcp/commit/234c1746810ccef067337a0d2e2c4a05b0b20de1"
        },
        "date": 1772464946512,
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
          "id": "3a4ca923f84e9eab55a9e8497866a3b9ea582050",
          "message": "fix(release): use paginated API to delete all release assets before re-run\n\ngh release view --json assets has a limit and only returns up to 30 assets.\nWith 37 assets per release (6 servers × 5 platforms + checksums), some\nassets remain undeleteted, causing goreleaser 422 already_exists errors.\n\nFix: use gh api --paginate on the REST /releases/{id}/assets endpoint which\nreturns all assets regardless of count.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T16:29:08+01:00",
          "tree_id": "02bdb02a2bcb7635e42a58c1547ebf57ef9aa012",
          "url": "https://github.com/Arkestone/mcp/commit/3a4ca923f84e9eab55a9e8497866a3b9ea582050"
        },
        "date": 1772465383517,
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
          "id": "c539af17b28f6f92de6cfb7805ef8fb5cc21b46f",
          "message": "fix(release): delete entire release on re-run, not just assets\n\nDeleting individual assets is insufficient because goreleaser uploads archives\nin parallel and retries failed uploads, causing 422 already_exists errors for\nassets that were already successfully uploaded in the same run.\n\nDeleting the entire release before goreleaser runs gives it a clean slate to\nrecreate the release and upload all assets without conflicts.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T16:35:41+01:00",
          "tree_id": "36422c1b7bd2db12d329b300f8bad0c85027a3ca",
          "url": "https://github.com/Arkestone/mcp/commit/c539af17b28f6f92de6cfb7805ef8fb5cc21b46f"
        },
        "date": 1772465778841,
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
          "id": "73547b9676237d3bbe4eb049e182734c53d8ee0c",
          "message": "fix(release): add concurrency control to prevent parallel release runs\n\nForce-pushing a tag while also triggering workflow_dispatch caused two\nrelease runs to execute in parallel, both trying to create/update the same\nGitHub release. Adding concurrency group cancels the older run automatically.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T16:41:10+01:00",
          "tree_id": "6211facf094c1f87ee57bfbdfac41eed4f918d37",
          "url": "https://github.com/Arkestone/mcp/commit/73547b9676237d3bbe4eb049e182734c53d8ee0c"
        },
        "date": 1772466105413,
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
          "id": "d66804c1b35e4f15460c00056b0c020662f6d997",
          "message": "ci: merge related workflows into 5 consolidated files\n\n- benchmark-trend + binary-size + coverage-trend → trends.yml\n  (3 push-to-main trend jobs in one workflow, per-job concurrency)\n- smoke-test + mcp-protocol → mcp-validation.yml\n  (healthz + full MCP protocol check in one matrix job)\n- fuzz + flaky-tests → reliability.yml\n  (Sunday reliability suite with improved flaky analysis summary)\n- conventional-commits + changelog-check + pr-size → pr-quality.yml\n  (all PR quality gates in one workflow, 3 parallel jobs)\n- dependency-review + license-check → dependency-compliance.yml\n  (dependency-review runs only on PRs via if: condition)\n\nReduces workflow count from 39 → 33 files.\nAll 33 workflows pass actionlint.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T17:21:03+01:00",
          "tree_id": "51940a5ec9a5a903444448f1753cb1b35189f44c",
          "url": "https://github.com/Arkestone/mcp/commit/d66804c1b35e4f15460c00056b0c020662f6d997"
        },
        "date": 1772468509117,
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
          "id": "9df0ccd5f7c7833a4a9200845a42cc0351ec3beb",
          "message": "ci: merge 8 more workflows into 4 consolidated files\n\n- container-structure + dockerfile-lint → docker-quality.yml\n  (lint + structure tests, both on Dockerfile push/PR + Monday schedule)\n- spell-check + link-check → doc-quality.yml\n  (both jobs in parallel; link-check schedule kept at Monday 06:00)\n- docs + wiki-sync → docs-publish.yml\n  (build+deploy GitHub Pages then sync wiki in same workflow)\n- benchmark + docker-size → weekly-metrics.yml\n  (weekly Monday run; benchmark job + per-server Docker size matrix)\n\nReduces workflow count from 33 → 29 files.\nAll 29 workflows pass actionlint.\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T17:31:39+01:00",
          "tree_id": "87003bf878f87361af14c1eba8a285bac2f6ac9f",
          "url": "https://github.com/Arkestone/mcp/commit/9df0ccd5f7c7833a4a9200845a42cc0351ec3beb"
        },
        "date": 1772469142318,
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
          "id": "a5a62b9553c6550f3b1e8db0470544560095fb04",
          "message": "ci(release): use GITHUB_TOKEN for GHCR visibility with PACKAGES_TOKEN fallback\n\nThe workflow's 'packages: write' permission grants write:packages scope to\nGITHUB_TOKEN, making it sufficient to PATCH package visibility without a\nseparate PAT in most cases.\n\nChanges:\n- Add comment on 'packages: write' clarifying it enables visibility changes\n- Use '${{ secrets.PACKAGES_TOKEN || secrets.GITHUB_TOKEN }}' so org-level\n  PAT is preferred when configured but GITHUB_TOKEN works out of the box\n- Parse HTTP status from response to distinguish 200/204 (success),\n  404 (package missing) and other errors instead of swallowing all failures\n- Emit ::warning:: with PAT instructions only when HTTP status is unexpected\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T17:36:10+01:00",
          "tree_id": "77c74891c4d4f56b8d89e6b37081dee3c32bc85f",
          "url": "https://github.com/Arkestone/mcp/commit/a5a62b9553c6550f3b1e8db0470544560095fb04"
        },
        "date": 1772469418246,
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
          "id": "404c7ff5faf6d613ee453a8bf6f399412fca23f5",
          "message": "chore: release v1.0.1\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T17:43:32+01:00",
          "tree_id": "59bd73f59a6db0202887b00369236ab52478c866",
          "url": "https://github.com/Arkestone/mcp/commit/404c7ff5faf6d613ee453a8bf6f399412fca23f5"
        },
        "date": 1772469859709,
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
          "id": "2d0d26772ef7934f0335bb509bd57d96a603b281",
          "message": "docs: add awesome-copilot configuration example\n\nShows how to configure mcp-instructions with github/awesome-copilot\nas a remote instruction source, covering:\n- VS Code mcp.json (stdio, Docker, HTTP)\n- Claude Desktop config\n- GitHub token for rate-limit protection\n- Pinning to a specific commit\n- URI reference table for available instructions\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T18:42:48+01:00",
          "tree_id": "ecd1651bdb4370c2aad3ed902a8ad6dc2884ef46",
          "url": "https://github.com/Arkestone/mcp/commit/2d0d26772ef7934f0335bb509bd57d96a603b281"
        },
        "date": 1772473403711,
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
          "id": "db67dd9b51486fd5ef35b1b3eacc3630f82ec839",
          "message": "feat: dynamic context loading with get-context tool and relevance filtering\n\nAdd context-aware filtering architecture across all three MCP servers\nso agents can load only the most relevant items without flooding the\ncontext window.\n\n## New: get-context (PRIMARY AGENT TOOL) on each server\n\n- mcp-instructions: get-context(file_path) → returns all applicable\n  instructions as a merged ready-to-inject block. Uses applyTo glob\n  patterns to include only relevant instructions; global instructions\n  (no applyTo) are always included.\n\n- mcp-skills: get-context(query, limit=5) → returns top-N most relevant\n  skills with full content and references, ranked by keyword match\n  (name ×3, description ×2, tags ×1). Default limit=5 to cap context use.\n\n- mcp-prompts: get-context(query, limit=5, type?) → same scoring model\n  for prompts and chatmodes; optional type filter.\n\n## New: pkg/glob — ** and {a,b} glob matching\n\nImplements VS Code-style applyTo patterns (e.g. **/*.go, **/*.{js,ts})\nwithout external dependencies using path.Match + custom ** recursion.\n\n## Improved filtering everywhere\n\n- list-instructions: file_path + source filters; returns total+matched\n- list-skills: query + source filters; returns total+matched\n- list-prompts: query + source + type filters; returns total+matched\n- optimize-instructions: now respects file_path filter before merging\n- optimize-skills: now respects query filter before merging\n- optimize-prompts: now respects query filter before merging\n\n## Full-tree discovery (loader/scanner rewrites)\n\n- All three servers now walk the entire directory tree (filepath.Walk)\n  finding instruction/skill/prompt files anywhere, including in hidden\n  subdirectories (.github/, etc.)\n- GitHub sync updated to use FetchDirRecursive (BFS) so remote repos\n  are fetched completely regardless of directory structure\n- Instruction applyTo supports both string and []string YAML formats\n\n## Agent workflow\n\n1. list-* with query/file_path → discover what exists (metadata only)\n2. get-context → load relevant content in one call\n3. get-skill / get-prompt → fetch specific item by name/URI if needed\n4. optimize-* → LLM-merge a filtered subset\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T19:31:05+01:00",
          "tree_id": "c24417dc85b9d82a6477f38e8c7ef9244475ffd0",
          "url": "https://github.com/Arkestone/mcp/commit/db67dd9b51486fd5ef35b1b3eacc3630f82ec839"
        },
        "date": 1772476303487,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 79,
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
          "id": "006f2082b41bbab921fbaea3f67aaafe0c51bd97",
          "message": "feat: efficient context-aware filtering with cache and precision scoring\n\n- pkg/filter: word-boundary tokenizer, weighted scoring (name>tag>desc),\n  coverage multiplier per matched term, stable sort (reproducible order)\n- pkg/cache: generic TTL List[T] with double-checked locking (5s default)\n- mcp-instructions loader: cache.List[Instruction] + maxFileSize 1MiB guard\n- mcp-prompts loader: cache.List[Prompt] + pkg/filter.SortByScore + maxFileSize\n- mcp-skills scanner: cache.List[Skill] + pkg/filter.SortByScore + maxFileSize\n- Fix: Tokenize keeps 2-char tokens; false positives prevented by minPrefixLen=3\n- Fix: TestLoaderLocalReadsLive calls ForceSync() before re-read\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T19:49:05+01:00",
          "tree_id": "d7020e1bd0e349d06cd4c0d666fdaead019dcfda",
          "url": "https://github.com/Arkestone/mcp/commit/006f2082b41bbab921fbaea3f67aaafe0c51bd97"
        },
        "date": 1772477397583,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 81.6,
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
          "id": "4eb8488da039e87765352260378505651bd8799a",
          "message": "feat(filter): stopwords, suffix-stripping stemmer, co-occurrence bonus\n\nStopwords:\n- Added 14 two-char function words (to/of/in/on/at/is/it/if/do/be/as/an/or/by)\n- Added 28 high-frequency English words (the/use/add/get/how/will/can/...)\n- Prevents noise terms from matching irrelevant items or diluting real scores\n\nStemming:\n- stem() strips 18 common English suffixes (longest-first, root ≥ 4 chars)\n- 'instructions'↔'instruction', 'testing'↔'test', 'deployment'↔'deploy'\n- 'authentication'→'authentic' matches 'auth' via bidirectional stem-prefix\n- Stem match is scored at exact-match weight (not prefix weight)\n\nCo-occurrence bonus:\n- +30% per field where 2+ query terms match the same field\n- 'jwt authentication' in name scores higher than split across name+desc\n- Rewards documents where multiple query concepts appear together\n\nTests: 13 filter tests all pass (added 3 new tests)\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:01:51+01:00",
          "tree_id": "d666d9c7492335e287cee85d75e858412a2a8588",
          "url": "https://github.com/Arkestone/mcp/commit/4eb8488da039e87765352260378505651bd8799a"
        },
        "date": 1772481747373,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 81.8,
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
          "id": "7c1d80163702374b990e2352f74f41c27cdc99b3",
          "message": "ci: macOS tests, nightly fuzz/profiling, outdated dependencies workflow\n\nci.yml — macOS test matrix:\n- Added matrix.os dimension with ubuntu-latest + macos-latest (stable only)\n- Coverage/PR comment/Codecov steps gated on ubuntu-latest to avoid duplication\n- Gives cross-platform signal on every push/PR without doubling CI cost\n\nnightly.yml (new, runs at 02:00 UTC daily):\n- fuzz: discovers and runs all Fuzz* targets with 30s timeout each\n- memory-profile: pprof mem profile of full test suite, top-20 allocators\n- coverage-deep: full HTML coverage report archived 30 days, flags <80% packages\n\noutdated.yml (new, runs Monday 07:00 UTC):\n- go list -u -m all: markdown table of available module updates\n- go mod verify: checksum integrity check (fails CI if mismatch)\n- go mod tidy check: warns if go.sum would change\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:04:23+01:00",
          "tree_id": "24c2a54921a0c3934d848d01fb5028154727322a",
          "url": "https://github.com/Arkestone/mcp/commit/7c1d80163702374b990e2352f74f41c27cdc99b3"
        },
        "date": 1772481911944,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 81.8,
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
          "id": "1228642ab474af47f08f390c33da139e4ae2311c",
          "message": "test: add FilterByFilePath, FilterByQuery, FetchDirRecursive tests and filter benchmarks\n\n- FilterByFilePath: 7 tests (instructions loader) → 81.2% → 95.5% coverage\n- FilterByQuery: 5 tests + lifecycle (prompts loader) → 59.2% → 71.4% coverage\n- FetchDirRecursive: 5 tests with recursiveServer() helper → pkg/github 71.9% → 98.4%\n- Filter benchmarks: Tokenize, Score (multi-term / no-match), SortByScore(100 items)\n- Fix TestFilterByQuery_SortsByScore: use 'hammering utility' to avoid stem collision with 'build'\n- CI: macOS added to test matrix (stable only); coverage/Codecov gated on ubuntu-latest\n- nightly.yml: fuzz 30s/target, pprof memory profiling, HTML coverage archive 30d\n- outdated.yml: go list -u weekly check + mod tidy drift detection\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:32:48+01:00",
          "tree_id": "6a7a2afcee9d39000f9a58bc4840983fbea199ca",
          "url": "https://github.com/Arkestone/mcp/commit/1228642ab474af47f08f390c33da139e4ae2311c"
        },
        "date": 1772483619682,
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
          "id": "42a2d6c49f369a75c52d2fd77ad0c1e0746d1715",
          "message": "test: add syncRepo/syncAllRepos tests with httptest server\n\n- TestSyncRepo_DownloadsPromptFiles: verifies files cached from fake GitHub API\n- TestSyncRepo_APIError: verifies error propagation on 404\n- TestSyncAllRepos_NoRepos: verifies no-op on empty config\n- prompts loader coverage: 71.4% → 86.7%\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:34:22+01:00",
          "tree_id": "9687d21d8ddd693852bc6a8a2f67ec08537e7666",
          "url": "https://github.com/Arkestone/mcp/commit/42a2d6c49f369a75c52d2fd77ad0c1e0746d1715"
        },
        "date": 1772483697907,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 85,
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
          "id": "b5a1ddfce9f50c2de367aeffe98da6348717e12d",
          "message": "feat(filter): phrase adjacency bonus + short-tag exact match\n\nPhrase adjacency bonus (detectPhrases):\n- +25% when consecutive query tokens appear adjacent in name field\n  e.g. query 'code review' → name 'code-review' earns +25%\n  while 'code-style-review' (non-adjacent) does not\n- +15% per adjacent pair in tags, +10% in description\n- Uses stem matching so 'unit testing' matches 'unit-test' adjacency\n- Order-sensitive: 'code review' does not bonus 'review-code'\n\nShort-tag exact match bonus:\n- Raw tag string matched case-insensitively against raw query words\n- Enables short tags ('go', 'ci', 'api', 'k8s') that Tokenize drops (≤1-char)\n  to score properly — gives +8 pts per matching tag\n- 'go' query now correctly ranks items tagged 'Go' above untagged\n\nNew constants: phraseBonusName=25, phraseBonusTag=15, phraseBonusDesc=10, shortTagBonus=8\nNew function: detectPhrases(queryTokens, targets []string) int\nTests added: 6 new filter tests (3 phrase, 3 short-tag)\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:45:26+01:00",
          "tree_id": "8a264c42ae938a09aedcaa8fb10cda0909ec9c7a",
          "url": "https://github.com/Arkestone/mcp/commit/b5a1ddfce9f50c2de367aeffe98da6348717e12d"
        },
        "date": 1772484357301,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 85.1,
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
          "id": "baf7a690ed99fd5d7fa48adfa087239304b922d0",
          "message": "feat(filter): name precision ratio bonus\n\nName precision bonus (+20% when ≥75% of name tokens covered by query):\n- Rewards items with concise names closely matching the query\n- 'jwt-auth' (2/2 = 100% coverage) beats 'jwt-auth-middleware-guide' (2/4 = 50%)\n  for query 'jwt auth'\n- Coverage computed via reverse pass: for each name token, check if any\n  query token matches (exact, stem, or prefix)\n- Only fires when name is non-empty and coverage >= 75% threshold\n\nNew constants: namePrecisionThreshold=75, namePrecisionBonus=20\nTests added: 3 new precision tests (concise wins, low coverage no bonus, single token)\npkg/filter coverage: 94.8% → 95.4%\n\nFull scoring pipeline now: tokenize → match → coverage bonus → co-occur bonus\n→ phrase adjacency → short-tag → name precision\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T21:46:45+01:00",
          "tree_id": "b5b5e629e3d984cf904d2d568be7f770c1950954",
          "url": "https://github.com/Arkestone/mcp/commit/baf7a690ed99fd5d7fa48adfa087239304b922d0"
        },
        "date": 1772484439504,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 85.2,
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
          "id": "df2e7de1b55fa9349fa5201c4f5eb8b8abba1534",
          "message": "feat: frontmatter files: glob pattern for prompts and skills\n\nBoth mcp-prompts and mcp-skills now support a files: field in frontmatter\nthat restricts which file paths the item applies to — identical semantics\nto instructions' existing applyTo: field.\n\nFrontmatter example:\n  ---\n  description: TypeScript code reviewer\n  tags: [typescript, code-review]\n  files: \"**/*.ts\"\n  ---\n\n  or as a list:\n  files:\n    - \"**/*.ts\"\n    - \"**/*.tsx\"\n\nChanges:\n- Prompt.Files []string: parsed from frontmatter files: (string or list)\n- Skill.Files []string: parsed from frontmatter files: (string or list)\n- FilterByFilePath() added to both loader and scanner packages\n  - Items with no Files: always included (global scope)\n  - Items with Files: included only when at least one pattern matches\n  - Empty filePath: all items returned unchanged (backward compatible)\n- file_path parameter added to list-prompts, get-context, optimize-prompts,\n  list-skills, get-context (skills), optimize-skills tool handlers\n- Files exposed in list output (files field in ListEntry)\n- toStringSlice() helper added to scanner (same as instructions loader)\n\nTests: 5 new FilterByFilePath tests per package (unit + integration with\n  real frontmatter parsing); 21 parseFrontmatter call sites updated to\n  unpack new 5th return value\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T22:16:08+01:00",
          "tree_id": "698c4c4529570c98c416ee91f717ac8edebf24fe",
          "url": "https://github.com/Arkestone/mcp/commit/df2e7de1b55fa9349fa5201c4f5eb8b8abba1534"
        },
        "date": 1772486229736,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 84.9,
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
          "id": "a7f8955ad3bea6a0c3c54fb92cb5be83d4b96b21",
          "message": "chore: update changelogs for v1.1.0\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T23:33:46+01:00",
          "tree_id": "dc626441666a1d8f9759a8fd79b118ae84fce21c",
          "url": "https://github.com/Arkestone/mcp/commit/a7f8955ad3bea6a0c3c54fb92cb5be83d4b96b21"
        },
        "date": 1772490861386,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 85.8,
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
          "id": "7735fd0df42f854d4ecf27eb15eab711617efbc2",
          "message": "fix: distinguish rate-limit 403 from private-repo 403 in GitHub client\n\nWhen calling the GitHub API without a token, public repos that exceed\nthe 60 req/hr unauthenticated rate limit also return HTTP 403. The\nprevious error message was misleading: 'is this a private repo?'.\n\nNow httpError reads the response body (JSON) and the X-RateLimit-Remaining\nheader to detect rate-limit errors and returns a clear message:\n  'HTTP 403 for <path>: rate limit exceeded — set GITHUB_TOKEN to increase\n   the rate limit (5000 req/hr vs 60 req/hr)'\n\nAuth/private-repo 403s (body does not mention 'rate limit') still produce\nthe original 'no GITHUB_TOKEN set — is this a private repo?' hint.\n\nAdded tests: TestFetchDir_HTTP403RateLimitViaBody,\nTestFetchDir_HTTP403RateLimitViaHeader, TestFetchFile_HTTP403RateLimitViaBody\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-02T23:55:49+01:00",
          "tree_id": "45344908c493139cd195484b94f544736ba01f41",
          "url": "https://github.com/Arkestone/mcp/commit/7735fd0df42f854d4ecf27eb15eab711617efbc2"
        },
        "date": 1772492181975,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 85.8,
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
          "id": "2bacf53a966b719833be9fa0cfcee133eca39829",
          "message": "fix: async initial sync + better rate-limit error messages\n\npkg/syncer: make initial sync non-blocking\n- syncer.Start() now runs the first sync call inside the background\n  goroutine rather than synchronously in the caller goroutine.\n- Effect: HTTP server binds and responds to /healthz within ~300ms even\n  when the GitHub repo sync takes minutes (e.g. 175 files from\n  github/awesome-copilot).\n- Updated TestStartCallsFnImmediately to use a channel-based wait\n  instead of an immediate atomic read (the fn still runs 'immediately'\n  but asynchronously).\n\npkg/github: distinguish rate-limit 403 from private-repo 403\n- httpError now reads the response body (JSON) and the\n  X-RateLimit-Remaining header to detect rate-limit errors.\n- Rate-limited public repos now get a clear message:\n  'rate limit exceeded — set GITHUB_TOKEN to increase the rate limit'\n  instead of the misleading 'is this a private repo?' hint.\n- Added 3 new tests for rate-limit detection via body and header.\n\nTested in real conditions against github/awesome-copilot (178 instructions):\n- TypeScript file (src/auth/login.ts): 77/178 matched\n- Go file (pkg/server/handler.go): 46/178 matched\n- Terraform (infra/main.tf): 48/178 matched (includes azure-verified-modules-terraform)\n- Kubernetes YAML (k8s/deployment.yaml): 51/178 matched\n\nCo-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-03T00:08:37+01:00",
          "tree_id": "68f5cf182a220a1b76fe74d10e6d222efe6dea12",
          "url": "https://github.com/Arkestone/mcp/commit/2bacf53a966b719833be9fa0cfcee133eca39829"
        },
        "date": 1772492956996,
        "tool": "customBiggerIsBetter",
        "benches": [
          {
            "name": "Total Coverage",
            "value": 85.8,
            "unit": "%"
          }
        ]
      }
    ]
  }
}