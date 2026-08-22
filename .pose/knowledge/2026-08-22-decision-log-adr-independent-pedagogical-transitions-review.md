---
type: decision-log
slug: adr-independent-pedagogical-transitions-review
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

# decision-log: adr-independent-pedagogical-transitions-review

## Context

Rastreia o gatilho de revisão do ADR
[2026-08-22-independent-pedagogical-transitions](../adr/2026-08-22-independent-pedagogical-transitions.md),
que separa feedback, avaliação, conclusão e avanço como operações
independentes do MCP.

## Current state

Decisão aceita. `session-orchestration-disclosure` e
`feedback-evaluation-progression` ainda não foram implementados.

## Next checks

- Ao implementar essas specs, cobrir com teste que uma chamada isolada de
  feedback ou avaliação nunca produz `step_completed`/`step_advanced`.

## Risks

- Um atalho de implementação que infira conclusão a partir de avaliação
  violaria esta decisão silenciosamente se não houver teste dedicado.

## Next owner

Mesmo owner.

## References

- ADR: `.pose/adr/2026-08-22-independent-pedagogical-transitions.md`
- Spec: `.pose/specs/2026-08-22-architecture-decision-baseline.md`
- PROJECT.md §15.9, §32.8
