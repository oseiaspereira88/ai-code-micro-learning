---
type: decision-log
slug: adr-local-versioned-catalog-and-event-state-review
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

# decision-log: adr-local-versioned-catalog-and-event-state-review

## Context

Rastreia o gatilho de revisão do ADR
[2026-08-22-local-versioned-catalog-and-event-state](../adr/2026-08-22-local-versioned-catalog-and-event-state.md),
que fixa catálogo YAML versionado e estado local por eventos JSONL
append-only com snapshots, sem banco de dados externo.

## Current state

Decisão aceita. `local-event-store` e `catalog-schema-loader` ainda não foram
implementados.

## Next checks

- Ao implementar `local-event-store`, medir concorrência de escrita e volume
  de eventos sob uso real antes de considerar qualquer alternativa a JSONL.

## Risks

- Migração prematura para SQLite sem medição violaria a decisão registrada.

## Next owner

Mesmo owner.

## References

- ADR: `.pose/adr/2026-08-22-local-versioned-catalog-and-event-state.md`
- Spec: `.pose/specs/2026-08-22-architecture-decision-baseline.md`
- PROJECT.md §18
