---
title: "Review bundle seal crashes on overlay profiles referencing uninstalled extension rules"
type: "bug"
module: ""
created_at: "2026-08-22T10:46:01Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: Review bundle seal crashes on overlay profiles referencing uninstalled extension rules

POSE 1.7.3 schema-v1 reproduction: initialize a project with 'pose init' before any language stack exists (empty repo, extensions.lock.json stays {"extensions":{}}, 'pose stacks' reports no markers). .pose/policy/review.json ships with default overlay_profiles ["backend-review@1", "frontend-review@1"]. Both profiles declare criteria that reference rule IDs 'backend-go' and 'frontend-react', which are only provided by extensions (per current docs, these moved out of the core rule set into extension packages) and are never auto-detected/installed because no stack existed at init time. Running 'pose review bundle spec:<any-spec> --explain' or '--seal' then hard-fails with 'pose: unknown review rule "backend-go" in backend-review@1' and 'pose: unknown review rule "frontend-react" in frontend-review@1', blocking closeout for every spec in the project regardless of that spec's own components/selectors. Unmapped review components already degrade gracefully to a [WARN]; an overlay profile whose selectors do not match the review scope's language, or whose referenced rules are not installed, should degrade the same way (skip/warn) instead of hard-erroring the seal. Expected: 'pose review bundle --seal' succeeds for a spec with no applicable overlay profile, and/or 'pose init'/'pose doctor' warns up front that overlay_profiles reference rules with no installed extension. No project code or private paths are required to reproduce.
