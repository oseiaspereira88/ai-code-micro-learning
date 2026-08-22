---
type: decision-log
slug: adr-agent-mcp-and-core-boundaries-review
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

# decision-log: adr-agent-mcp-and-core-boundaries-review

## Context

Rastreia o gatilho de revisão do ADR
[2026-08-22-agent-mcp-and-core-boundaries](../adr/2026-08-22-agent-mcp-and-core-boundaries.md),
que fixa as fronteiras entre agente tutor, skill, servidor MCP, CLI e núcleo,
e o transporte `stdio` como decisão da V1.

## Current state

Decisão aceita e sem violação conhecida. Nenhuma spec executável foi
implementada ainda sobre esta fronteira.

## Next checks

- Ao implementar `mcp-stdio-foundation` e `tutor-skill-host-integration`,
  confirmar que nenhuma lógica de domínio foi movida para dentro do servidor
  MCP.

## Risks

- Um segundo host de produção pode exigir um transporte incompatível com
  `stdio`, forçando revisão do ADR antes de expandir compatibilidade.

## Next owner

Mesmo owner.

## References

- ADR: `.pose/adr/2026-08-22-agent-mcp-and-core-boundaries.md`
- Spec: `.pose/specs/2026-08-22-architecture-decision-baseline.md`
- PROJECT.md §11, §15.1–15.2
