# ADR: Local versioned catalog and event state

## Status
Accepted

## Context

`PROJECT.md` §18 (Persistência) e §32 itens 9–10 definem que o catálogo é um
grafo versionado autorado em YAML e que o estado local usa eventos
append-only com snapshots derivados, sem banco de dados ou daemon externo.
Essa escolha ainda não está registrada como decisão durável, o que arrisca
reinterpretação divergente entre `catalog-schema-loader`, `local-event-store`
e specs dependentes.

Módulos afetados: `catalog-schema-loader`, `local-event-store`,
`workspace-observation-baselines` (evidência), layout `.ailearn/`.

## Decision

1. O catálogo é autorado em YAML e versionado no Git, como grafo (temas,
   conceitos, competências, desafios, packs, trilhas), não como lista linear.
2. O estado de sessão é gravado como eventos JSONL append-only
   (`PROJECT.md` §18.3 lista o conjunto mínimo de eventos); snapshots JSON
   derivados existem apenas para leitura rápida e são reconstruíveis a partir
   dos eventos.
3. Arquivos de evidência grandes são armazenados por hash, com índice
   reconstruível a partir do log de eventos.
4. Escrita é serializada por lock exclusivo; o evento é gravado e
   sincronizado antes de atualizar o snapshot; o snapshot é substituído
   atomicamente; chamadas mutáveis aceitam `expected_revision` para
   concorrência otimista e IDs de requisição garantem idempotência
   (`PROJECT.md` §18.4).
5. Nenhum banco de dados ou daemon externo é introduzido na V1.

### Alternativas rejeitadas

- **Adotar SQLite desde a V1**: rejeitado; fica reservado como migração
  futura, condicionada a medições reais de concorrência, consulta ou volume
  (critério explícito de `PROJECT.md` §18.1), não a preferência antecipada.
- **Catálogo embutido no binário/código**: rejeitado por remover a separação
  entre autoria de conteúdo e versão de software, dificultando revisão
  editorial independente.
- **Estado mutável em arquivo único sem log de eventos**: rejeitado por
  eliminar trilha de auditoria, idempotência de retries e capacidade de
  reconstrução após falha.

## Consequences

- `local-event-store` deve implementar o writer append-only, o lock exclusivo
  e a substituição atômica de snapshot antes de qualquer spec de sessão
  depender dele.
- `catalog-schema-loader` deve tratar o YAML como fonte de verdade versionada
  e não replicar seu conteúdo em outro formato canônico.
- Qualquer proposta de introduzir banco de dados externo exige ADR próprio e
  evidência mensurável do problema que a justifica.

### Gatilho de revisão

Revisar esta decisão apenas quando medições reais demonstrarem problemas de
concorrência, consulta ou volume que o layout local não suporte
(`PROJECT.md` §18.1).
