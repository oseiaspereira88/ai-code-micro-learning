---
title: "Review bundle seal unreachable for doc-only specs: change-set not persisted, validate has no precondition"
type: "bug"
module: ""
created_at: "2026-08-22T11:16:13Z"
status: staged
upstream: "https://github.com/oseiaspereira88/pose"
privacy: sanitized-synthetic
---

# Contribution Draft: Review bundle seal unreachable for doc-only specs: change-set not persisted, validate has no precondition

POSE 1.7.4 schema-v1 reproduction: create a spec whose Technical Plan declares 'Delivery targets: Nenhum' (no code module, no build/typecheck applicable) and whose only artifacts are Markdown (e.g. ADRs under .pose/adr/). Commit those artifacts with the 'POSE-Spec: <slug>' trailer. Then: (1) run 'pose artifact-check --spec <slug> --strict' — it reconciles cleanly (0 findings, prints a change_set id), but running 'pose review bundle spec:<slug> --explain' immediately after still reports '[ERROR] no immutable attributed change set exists for spec:<slug>', so artifact-check's reconciliation is not what review bundle consults, or it is not persisted anywhere review bundle reads. (2) run 'pose validate --strict --json .pose/results/delivery-validation.json' — since no stack/module exists yet in the repo, it produces 'checks: []', 'counts.executed: 0', but 'outcome: pass' and a 'scope_provenance' entry keyed by the spec slug. 'pose review bundle --explain' still reports '[ERROR] no passed structured validation evidence is attributed to the review scope', because zero executed checks never counts as 'passed' evidence. Root cause for (2): .pose/review-profiles/spec-closeout.json marks the 'validate' tool as requiredness:'required' with no 'preconditions' field, unlike frontend-review@1's 'surface-check' tool which declares preconditions:['delivery-target-declared'] to skip itself when not applicable. A spec that legitimately has no code artifacts and no delivery target (governance/ADR-only specs, the first milestone of a fresh project before any stack exists) can therefore never produce validation evidence and never seal a review bundle, regardless of how correct and complete its content is. Expected: either (a) artifact-check's reconciled change-set is attributed/persisted so a subsequent review bundle --seal recognizes it, and/or (b) the spec-closeout@1 'validate' tool gains a precondition (e.g. 'delivery-target-declared' or 'module-exists') so it is skipped, not required, for specs with no code module — mirroring how surface-check already degrades gracefully. No project code or private paths are required to reproduce; a fresh 'pose init' plus one Markdown-only spec is sufficient.

## Upstream issue

- https://github.com/oseiaspereira88/pose/issues/36
