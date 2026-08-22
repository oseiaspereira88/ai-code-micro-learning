---
title: "Doc-audit report truncates AGENTS.md filename"
kind: bug
engine_version: 1.6.0
reported_at: 2026-08-22T01:50:28-03:00
---

# POSE Engine Report: Doc-audit report truncates AGENTS.md filename

## Description
Reproduction: modify AGENTS.md and run pose report --type doc-audit --task <task> --outcome pass with pose 1.6.0. Observed: Files Changed records GENTS.md. Expected: AGENTS.md is preserved exactly.

---
### System Context (Auto-generated)
- **POSE Engine Version:** 1.6.0
- **OS/Arch:** linux/amd64
- **Go Version:** go1.26.5-X:nodwarf5
- **Reported At:** 2026-08-22T01:50:28-03:00
- **Kind:** bug

