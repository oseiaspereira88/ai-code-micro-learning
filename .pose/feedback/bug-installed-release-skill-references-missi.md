---
title: "Installed release skill references missing policy file"
kind: bug
engine_version: 1.6.0
reported_at: 2026-08-22T01:50:28-03:00
---

# POSE Engine Report: Installed release skill references missing policy file

## Description
Reproduction: run pose skills-check --strict in a fresh schema-v1 instance installed by pose 1.6.0. Observed: .agents/skills/pose-release-closeout/SKILL.md links to ../../../.pose/release-policy.json, which is absent, so skills-check fails. Expected: the installed skill references an existing canonical release policy or the installer seeds the declared file.

## Upstream issue

- https://github.com/oseiaspereira88/pose/issues/32

---
### System Context (Auto-generated)
- **POSE Engine Version:** 1.6.0
- **OS/Arch:** linux/amd64
- **Go Version:** go1.26.5-X:nodwarf5
- **Reported At:** 2026-08-22T01:50:28-03:00
- **Kind:** bug
