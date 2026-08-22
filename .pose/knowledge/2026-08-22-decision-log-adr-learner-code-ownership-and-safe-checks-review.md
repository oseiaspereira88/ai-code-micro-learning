---
type: decision-log
slug: adr-learner-code-ownership-and-safe-checks-review
owner: @oseiaspereira
sensitivity: public-internal
created_at: 2026-08-22
last_reviewed_at: 2026-08-22
expires_at: 2026-11-20
source_refs:
  spec: "architecture-decision-baseline"
  workflow: ".pose/workflows/refactor.md"
  commands: []
  external_sources: []
---

# decision-log: adr-learner-code-ownership-and-safe-checks-review

## Context

Rastreia o gatilho de revisão do ADR
[2026-08-22-learner-code-ownership-and-safe-checks](../adr/2026-08-22-learner-code-ownership-and-safe-checks.md),
que veda edição de código do aluno pelo MCP e restringe execução a
`check_id` de uma allowlist declarada.

## Current state

Decisão aceita. `safe-check-executor` ainda não foi implementado.

## Next checks

- Ao implementar `safe-check-executor`, validar allowlist, contenção de path
  e ausência de shell contra `.pose/rules/security.md`.

## Risks

- Um novo tipo de check fora da allowlist declarada exigiria exceção
  documentada, não implementação silenciosa.

## Next owner

Mesmo owner.

## References

- ADR: `.pose/adr/2026-08-22-learner-code-ownership-and-safe-checks.md`
- Spec: `.pose/specs/2026-08-22-architecture-decision-baseline.md`
- PROJECT.md §11.2, §20
