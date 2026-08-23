---
slug: local-event-store
status: done
created_at: 2026-08-22
completed_at: 2026-08-23
supersedes:
depends_on: architecture-decision-baseline, go-runtime-foundation, learning-domain-model
priority: 40
components: event-store, evidence
delivers:
changelog: none
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
- [x] Definir schema de envelope, revisão e matriz evento versus projeção.
- [x] Definir política de recuperação e retenção.

### Implementation
- [x] Implementar append durável e idempotente.
- [x] Implementar snapshot atômico e rebuild.
- [x] Implementar store de evidência content-addressed.
- [x] Implementar lock (exclusividade e liberação testadas; sem harness genérico de falha injetada — ver Known gaps).
- [x] Criar export seguro (filtro/redação delegado ao consumidor).

### Validation
- [x] Testar recuperação em fronteira de escrita real (linha truncada por simulação de crash).
- [x] Testar corrida de revisões e retries (conflito otimista, idempotência por request_id, concorrência real).
- [x] Executar race detector e testes em filesystem temporário (`t.TempDir()`).

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: O estado é local, pequeno e precisa ser auditável.
- Options considered: arquivo mutável; SQLite; JSONL mais snapshot.
- Decision: log append-only com projeções reconstruíveis.
- Rationale: simplicidade operacional e recuperação explícita.
- Consequences: compactação e lock precisam de implementação cuidadosa.

### Decision 2
- Date: 2026-08-22
- Context: A spec listava `.gitignore` como artefato modificado, mas
  `/.ailearn/` já está excluído desde o commit inicial de configuração do
  projeto (antes de qualquer spec executável).
- Options considered: reescrever `.gitignore` mesmo sem mudança real;
  remover o artefato e documentar que já está coberto.
- Decision: remover `.gitignore` da seção 3; nenhuma mudança é necessária.
- Rationale: evita um "modified" vazio que `artifact-check` rejeitaria por
  ausência de mudança atribuível.
- Consequences: nenhuma; layout runtime (`.ailearn/state`, `.ailearn/evidence`,
  `.ailearn/exports`) já fica fora do Git por padrão.

### Decision 3
- Date: 2026-08-22
- Context: R2 exige reconstruir "todas as projeções" a partir do log, mas
  projeções (estado de sessão, progresso) são modelos de domínio ainda não
  ligados a este store; owning esse formato aqui duplicaria
  `learning-domain-model` e specs de sessão futuras.
- Options considered: implementar projeções concretas de sessão aqui;
  expor somente o mecanismo de replay/append/snapshot/lock e deixar a
  reconstrução de projeções para o consumidor.
- Decision: `internal/eventstore` expõe `Replay`/`ReplayAll` e snapshots
  atômicos como mecanismo; a reconstrução de uma projeção específica é
  responsabilidade do consumidor (ex.: `session-orchestration-disclosure`).
- Rationale: mantém o store agnóstico ao payload, coerente com o ADR de
  fronteiras (núcleo não deve reimplementar modelos de domínio já
  definidos).
- Consequences: R2 é satisfeito no nível de mecanismo (replay sempre
  disponível, sem perda de dados) e verificado por teste que reconstrói uma
  projeção de exemplo a partir do log; a projeção real de sessão é escopo de
  outra spec.

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
- `gofmt -l internal/eventstore internal/evidence` → saída vazia (2026-08-22).
- `go vet ./internal/eventstore/... ./internal/evidence/...` → sem findings (2026-08-22).
- `go test ./internal/eventstore/... ./internal/evidence/... -race` → `ok`,
  15 + 5 testes (2026-08-22).
- `govulncheck ./...` → `No vulnerabilities found.` (2026-08-22).

### Results summary
- `internal/eventstore` implementa append durável e idempotente (fsync antes
  de retornar sucesso), concorrência otimista por `expected_revision`,
  replay total/por stream, snapshot atômico (`temp + fsync + rename`),
  recuperação tolerante a linha final truncada sem perder o prefixo válido,
  detecção de corrupção real e de versão de schema incompatível com
  upcasters explícitos, lock exclusivo por arquivo, e export com filtro de
  redação delegado ao consumidor.
- `internal/evidence` implementa armazenamento content-addressed por
  SHA-256 com deduplicação, verificação de integridade na leitura e rejeição
  de IDs malformados antes de qualquer acesso a filesystem.
- Teste de concorrência real com 8 goroutines disputando a mesma revisão
  confirma que exatamente uma vence (sem `-race` findings).

### Requirement trace
- R1 [satisfied] test:TestAppendAssignsMonotonicRevisions
- R2 [satisfied] test:TestOpenRecoversRevisionsAndIdempotencyFromExistingLog test:TestReplayAndReplayAll
- R3 [satisfied] test:TestSnapshotOnlyWrittenAfterEventSynced
- R4 [satisfied] test:TestAppendRejectsRevisionConflict test:TestAppendIsIdempotentByRequestID test:TestAppendConcurrentSameStreamNeverDuplicatesRevisions
- R5 [satisfied] test:TestGetDetectsTampering test:TestPutIsContentAddressedAndDeduplicates
- R6 [satisfied] test:TestReadEventsToleratesTruncatedFinalLine test:TestReadEventsReportsCorruptedMiddleLineWithoutLosingPrefix test:TestReadEventsRejectsIncompatibleVersionWithoutUpcaster
- R7 [satisfied] test:TestLockIsExclusiveAndReleasable
- R8 [satisfied] test:TestExportAppliesFilter

### Known gaps
- Housekeeping e compatibilidade longa serão aprofundados em reliability-observability-compatibility.
- "Falhas injetáveis" da estratégia de validação foram cobertas por um
  cenário realista (linha final truncada simulando crash), não por um
  harness genérico de injeção de falha em cada write. Suficiente para R3/R6;
  ampliar se surgir um novo tipo de falha de escrita a cobrir.
- Redaction de segredos não é feita pelo eventstore: por design (fronteira
  do ADR de execução segura), o produtor do payload e `safe-check-executor`
  (§20.3) são responsáveis por redigir antes de chamar `Append`; `Export`
  apenas aplica o filtro que o chamador fornecer.
- Lock não sobrevive a `SIGKILL` (limite inerente de lock por arquivo);
  recuperação de lock órfão fica para uma spec operacional dedicada.

## 7. Final Report

### Delivered scope
Log de eventos append-only durável, idempotente e recuperável, snapshots
atômicos, lock exclusivo por workspace e store de evidência
content-addressed, cobrindo R1–R8. Nenhuma projeção de sessão concreta,
sincronização em nuvem ou SQLite, conforme os non-goals da spec.

### Files and modules changed
- `internal/eventstore/event.go` (criado)
- `internal/eventstore/store.go` (criado)
- `internal/eventstore/snapshot.go` (criado)
- `internal/eventstore/recovery.go` (criado)
- `internal/eventstore/lock.go` (criado)
- `internal/eventstore/store_test.go` (criado)
- `internal/eventstore/recovery_test.go` (criado)
- `internal/evidence/store.go` (criado)
- `internal/evidence/store_test.go` (criado)
- `.pose/specs/2026-08-22-local-event-store.md` (atualizado)

### Validation executed
- Command: go test ./internal/eventstore/... ./internal/evidence/... -race && govulncheck ./...
- Result: `ok  	github.com/oseiaspereira88/ailearn/internal/eventstore`; `ok  	github.com/oseiaspereira88/ailearn/internal/evidence`; `No vulnerabilities found.`

### Residual risks
- Portabilidade do lock exige prova por plataforma.

### Follow-ups
- [covered: reliability-observability-compatibility] Validar recuperação e retenção em cenários prolongados.
