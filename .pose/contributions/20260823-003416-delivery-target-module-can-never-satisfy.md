---
title: "Delivery target module can never satisfy review-bundle evidence in a single-module project"
type: "bug"
module: ""
created_at: "2026-08-23T00:34:16Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: Delivery target module can never satisfy review-bundle evidence in a single-module project

POSE 1.7.7 schema-v1 reproduction: in a single-module Go (or any single-manifest-file) repo, declare a spec's '### Delivery targets' section as e.g. '- contract:my-contract module:internal/mypkg profile:api-contract entrypoint:cmd/myapp/main.go', matching the exact syntax pose expects (kind:id module:<dir> profile:<name> entrypoint:<file>). 'pose index'/'pose check' accept this and populate delivery-integrity.json's 'deliveries' with Module:'internal/mypkg'. Then run 'pose validate --json <path>' (all checks recorded under Module:'.', the only module 'pose validate' ever detects, since the repo has exactly one go.mod at the root) and 'pose index' again. Run 'pose review bundle spec:<slug> --explain' or '--seal': it now reports '[ERROR] no passed structured validation evidence is attributed to the review scope' with review_bundle.evidence=0, even though 'pose validate' just passed. Root cause: reviewBundleEvidence builds a 'modules' allow-set from graph.Deliveries where target.Spec matches the scope (here {'internal/mypkg': true}), then filters graph.ValidationResults to '!modules[evidence.Module]' — but every recorded evidence.Module is '.', never 'internal/mypkg', because pose validate has no way to attribute checks to an arbitrary subdirectory of a single go.mod module. Attempting the 'obvious' fix — declaring the delivery target's module as '.' instead of a subdirectory — fails differently: validateArtifactPathSyntax hard-rejects clean=='.' with 'path escapes project root' (delivery_integrity.go), so 'pose index' itself then fails with 'pose index: delivery integrity: <ref> module: path escapes project root'. So neither a real subdirectory nor '.' can ever produce a delivery-target module value that lines up with the only module pose validate actually reports evidence for. This makes review-bundle sealing structurally unreachable for any delivery target declared in a single-module project, regardless of how much passing validation evidence exists. Expected: either (a) reviewBundleEvidence's module filter should also accept evidence whose Module is an ancestor directory of (or equal to the git root for) the declared target Module, so root-level '.' evidence still counts for a subdirectory target, or (b) '.' should be an allowed delivery-target module value (it is a completely different concept from an artifact path 'escaping' the root — it names the whole-repo module, which is exactly what a single-module project's target belongs to), or (c) pose validate should support attributing a check's Module to a non-go.mod subdirectory the caller declares (e.g. via moduleOverrides on a path with no marker file), so evidence can be scoped narrower than the whole repo. No project code or private paths are required to reproduce: a fresh 'pose init', a single go.mod at the root, one spec declaring a delivery target for any subdirectory, and one passing 'pose validate --json' run is sufficient.

## Upstream issue

- https://github.com/oseiaspereira88/pose/issues/39
