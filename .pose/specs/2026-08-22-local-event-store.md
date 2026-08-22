---
slug: local-event-store
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: architecture-decision-baseline, go-runtime-foundation, learning-domain-model
priority: 40
components: event-store, evidence
delivers:
---

# Spec: local-event-store

## 1. Intent

### Goal
Persistir sessões, eventos, snapshots e evidências localmente com recuperação, auditabilidade e idempotência.

### Business value
Permitir retomada confiável sem depender da memória do chat ou de serviço externo.

### Constraints
- Usar JSONL append-only e snapshots JSON na V1.
- Manter estado runtime fora do Git por padrão.
- Nunca reescrever evidência confirmada.

### Non-goals
- Sincronização em nuvem, multiusuário ou SQLite.
- Criptografia de dados em repouso.

## 2. Requirements

### Functional
- R1: Anexar eventos imutáveis com schema, ID, revisão, timestamp e payload tipado.
- R2: Reconstruir todas as projeções a partir do log após snapshot ausente ou inválido.
- R3: Substituir snapshots atomicamente somente depois de sincronizar o evento.
- R4: Aplicar expected_revision para concorrência otimista e request_id para idempotência.
- R5: Armazenar evidências grandes por SHA-256 e verificar integridade na leitura.
- R6: Detectar linha truncada, corrupção e versão incompatível sem perder o prefixo válido.
- R7: Aplicar lock exclusivo por workspace e liberar lock em cancelamento ou encerramento.
- R8: Exportar histórico e progresso sem expor conteúdo excluído por política.

### Non-functional
- Leituras de sessão devem usar snapshot e permanecer reconstruíveis.
- Falha entre append e snapshot não pode perder evento confirmado.

### Security
- Usar permissões restritivas, paths confinados e redaction antes de persistir output.
- Nunca persistir environment completo ou secrets conhecidos.

### Compatibility
- Eventos possuem versão e upcasters explícitos.

## 3. Technical Plan

### Affected areas
- internal/eventstore/, internal/evidence/, .gitignore

### Artifacts
- created: internal/eventstore/event.go
- created: internal/eventstore/store.go
- created: internal/eventstore/snapshot.go
- created: internal/eventstore/recovery.go
- created: internal/eventstore/lock.go
- created: internal/eventstore/store_test.go
- created: internal/eventstore/recovery_test.go
- created: internal/evidence/store.go
- created: internal/evidence/store_test.go
- modified: .gitignore

### Delivery targets
Nenhum; persistência interna.

### API/contract changes
- Definir portas de EventStore, SnapshotStore e EvidenceStore no lado consumidor.

### Data/storage changes
- Criar layout runtime .ailearn/state, .ailearn/evidence e .ailearn/exports.

### Technical risks
- Lock portátil entre sistemas operacionais.
- Crescimento do log e outputs grandes.
- Upcasters podem mascarar incompatibilidades se não forem estritos.

## 4. Tasks

### Planning
- [ ] Definir schema de envelope, revisão e matriz evento versus projeção.
- [ ] Definir política de recuperação e retenção.

### Implementation
- [ ] Implementar append durável e idempotente.
- [ ] Implementar snapshot atômico e rebuild.
- [ ] Implementar store de evidência content-addressed.
- [ ] Implementar lock e falhas injetáveis.
- [ ] Criar export seguro.

### Validation
- [ ] Testar falhas em cada fronteira de escrita.
- [ ] Testar corrida de revisões e retries.
- [ ] Executar race detector e testes em filesystem temporário.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: O estado é local, pequeno e precisa ser auditável.
- Options considered: arquivo mutável; SQLite; JSONL mais snapshot.
- Decision: log append-only com projeções reconstruíveis.
- Rationale: simplicidade operacional e recuperação explícita.
- Consequences: compactação e lock precisam de implementação cuidadosa.

## 6. Validation

### Strategy
Usar testes de falha injetada, propriedades de replay e concorrência controlada.

### Deterministic checks
- Test: go test ./internal/eventstore/... ./internal/evidence/...
- Lint: gofmt -l internal/eventstore internal/evidence
- Typecheck: go vet ./internal/eventstore/... ./internal/evidence/...
- Build: go test ./internal/eventstore/... -run ^$
- Security / Contract: go test -race e inspeção de permissões/redaction.

### Execution log
- Pendente; spec draft.

### Results summary
- Nenhum estado persistente entregue.

### Requirement trace
- Mapear R1–R8 a testes de replay, falha, integridade e export.

### Known gaps
- Housekeeping e compatibilidade longa serão aprofundados em reliability-observability-compatibility.

## 7. Final Report

### Delivered scope
Nenhum; planejamento completo sem implementação.

### Files and modules changed
- Planejados em internal/eventstore, internal/evidence e .gitignore.

### Validation executed
- Command: pose lint-spec local-event-store --ready-check
- Result: registrar após o gate.

### Residual risks
- Portabilidade do lock exige prova por plataforma.

### Follow-ups
- [covered: reliability-observability-compatibility] Validar recuperação e retenção em cenários prolongados.
