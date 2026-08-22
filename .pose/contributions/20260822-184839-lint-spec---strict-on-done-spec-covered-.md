---
title: "lint-spec --strict on done spec: [covered: <slug>] disposition falsely reports missing spec"
type: "bug"
module: ""
created_at: "2026-08-22T18:48:39Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: lint-spec --strict on done spec: [covered: <slug>] disposition falsely reports missing spec

POSE 1.7.5 schema-v1 reproduction: have spec A (status: done, completed_at set) with a Final Report follow-up '- [covered: <slug-of-B>] <text>', where spec B exists at .pose/specs/<date>-<slug-of-B>.md with matching frontmatter 'slug: <slug-of-B>' (flat dated layout, not the legacy .pose/specs/<slug>/spec.md nested layout). Running 'pose lint-spec <slug-of-A> --strict' reports: '[ERROR] <slug-of-A>: follow-up lacks a valid disposition: disposition [covered: <slug-of-B>] points to a missing spec'. This happens even after running 'pose index' to regenerate .pose/indexes/ (spec-graph.json already lists <slug-of-B> correctly, including edges from/to it), and even though 'pose lint-spec <slug-of-B> --required-only' succeeds and 'pose list-specs'/pose_list_specs resolves it fine. So the done-lifecycle disposition-target existence check used by lint-spec's strict lifecycle gate resolves spec slugs differently than the rest of the engine (list-specs, spec-graph, dependency resolution), and fails specifically for specs stored in the flat dated filename layout ('.pose/specs/YYYY-MM-DD-<slug>.md', the default produced by 'pose new-spec' in current POSE) rather than the legacy nested layout ('.pose/specs/<slug>/spec.md'). This may be the same root cause as a previously reported issue about new-spec defaulting to the flat layout while some other subsystem still assumes the legacy nested one. Expected: the lifecycle disposition-target check should resolve slugs the same way every other command does (by frontmatter 'slug:' across .pose/specs/**, not by a hardcoded path shape), so a [covered:]/[spawned:]/[duplicate:] disposition pointing at a real, existing spec never fails lint-spec --strict on close. No project code or private paths are required to reproduce; a fresh 'pose init', one flat-layout target spec, and one done spec with a covered follow-up pointing at it is sufficient.

## Upstream issue

- https://github.com/oseiaspereira88/pose/issues/37
