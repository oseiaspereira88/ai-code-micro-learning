---
title: "Installed release-closeout skill references missing release policy"
type: "bug"
module: ""
created_at: "2026-08-22T06:00:13Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: Installed release-closeout skill references missing release policy

POSE 1.6.0 schema-v1 reproduction: initialize a fresh instance and run 'pose skills-check --strict'. The check fails because .agents/skills/pose-release-closeout/SKILL.md links to ../../../.pose/release-policy.json, but the installer does not seed that file. Expected: every shipped skill link resolves, either by referencing the canonical policy under .pose/policy/ or by installing the declared release-policy.json. No project code or private paths are required to reproduce.

## Upstream issue

- https://github.com/oseiaspereira88/pose/issues/32
