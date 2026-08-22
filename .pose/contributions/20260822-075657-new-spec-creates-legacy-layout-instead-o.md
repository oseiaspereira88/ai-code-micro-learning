---
title: "new-spec creates legacy layout instead of modern dated flat format"
type: "bug"
module: ""
created_at: "2026-08-22T07:56:57Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: new-spec creates legacy layout instead of modern dated flat format

POSE version: 1.7.0, schema 1. Synthetic reproduction: install POSE into an empty temporary directory, run pose new-spec synthetic-layout-repro, then run pose spec-format status --json. Actual: new-spec creates .pose/specs/synthetic-layout-repro/spec.md and status reports date_prefixed=0, legacy=1, conforming=false. A dry run of pose spec-format migrate --all --format flat proposes .pose/specs/2026-08-22-synthetic-layout-repro.md. Expected: new-spec should create the current modern date-prefixed layout by default and its own output should immediately be conforming. Documentation is internally inconsistent: the installed POSE manual says new-spec creates date-prefixed folders, while pose new-spec --help advertises the legacy undated folder. Workaround: run pose spec-format migrate --all --format flat after scaffolding. Suggested fix: make new-spec use the canonical modern default, align help and manual with the selected default, and add a regression test asserting that a freshly created spec is conforming without migration.

## Upstream issue

- https://github.com/oseiaspereira88/pose/issues/33
