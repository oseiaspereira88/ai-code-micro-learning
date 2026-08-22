---
slug: mastery-review-scheduling
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: local-event-store, feedback-evaluation-progression
priority: 110
components: mastery, progress
delivers:
---

# Spec: mastery-review-scheduling

## 1. Intent

### Goal
Projetar domínio multidimensional por competência e agendar revisões espaçadas por regras transparentes.

### Business value
Distinguir conclusão momentânea de autonomia, retenção e transferência real.

### Constraints
- Não usar score opaco ou gamificação por volume.
- Solução revelada não prova autonomia.
- Retenção exige outra data; transferência exige contexto diferente.

### Non-goals
- Machine learning adaptativo.
- Ranking entre pessoas.

## 2. Requirements

### Functional
- R1: Projetar compreensão, sintaxe, implementação guiada, autônoma, depuração, explicação, retenção e transferência.
- R2: Representar não observada, introduzida, demonstra com ajuda, sem ajuda, retida e transferida.
- R3: Consumir tentativas, pistas, avaliação, reflexão, data e variante sem reescrever evidência.
- R4: Impedir promoção por uma única conclusão guiada ou por override.
- R5: Agendar revisões iniciais em 1, 3, 7, 14 e 30 dias, adaptadas por resultado.
- R6: Selecionar revisão vencida com prioridade explicável e sem bloquear sessão livre.
- R7: Expor progress_get, review_due e mastery_evidence_record.
- R8: Recalcular projeções deterministicamente a partir do log.

### Non-functional
- Cada mudança de domínio deve listar evidências causais.
- Regras e intervalos devem ser configuráveis e versionados.

### Security
- Progresso é dado local potencialmente sensível e não deve sair do dispositivo por padrão.

### Compatibility
- Mudanças na regra devem preservar projeção histórica por versão.

## 3. Technical Plan

### Affected areas
- internal/mastery/, internal/application/, internal/mcpserver/

### Artifacts
- created: internal/mastery/model.go
- created: internal/mastery/projector.go
- created: internal/mastery/rules.go
- created: internal/mastery/scheduler.go
- created: internal/mastery/projector_test.go
- created: internal/mastery/scheduler_test.go
- created: internal/application/progress.go
- created: internal/mcpserver/progress_tools.go
- created: internal/mcpserver/progress_contract_test.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar três tools de progresso e schemas de explicabilidade.

### Data/storage changes
- Persistir mastery_projected e review_scheduled como projeções reconstruíveis.

### Technical risks
- Regras rígidas podem recomendar prática demais.
- Variantes mal classificadas podem simular transferência.

## 4. Tasks

### Planning
- [ ] Definir tabela evidência versus dimensão e estado.
- [ ] Definir cálculo de revisão e desempate.

### Implementation
- [ ] Implementar projector versionado.
- [ ] Implementar scheduler e motivos.
- [ ] Implementar queries de progresso e revisão.
- [ ] Expor tools e eventos.
- [ ] Criar fixtures longitudinais.

### Validation
- [ ] Testar replay, mudança de regra e datas com clock fake.
- [ ] Testar que ajuda, solução e override não promovem indevidamente.
- [ ] Testar retenção e transferência em variantes.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Um número único esconde a natureza da lacuna.
- Options considered: XP; score 0–100; vetor de estados explicável.
- Decision: vetor de dimensões e estados derivados de evidência.
- Rationale: permite recomendação precisa e auditável.
- Consequences: UI textual é mais extensa, porém mais honesta.

## 6. Validation

### Strategy
Usar timelines sintéticas, clocks determinísticos e property tests de monotonicidade condicional.

### Deterministic checks
- Test: go test ./internal/mastery/... ./internal/application/...
- Lint: gofmt -l internal/mastery
- Typecheck: go vet ./internal/mastery/...
- Build: go test ./internal/mastery/... -run ^$
- Security / Contract: export redaction e isolamento entre perfis.

### Execution log
- Pendente.

### Results summary
- Nenhuma projeção entregue.

### Requirement trace
- Mapear R1–R8 a timelines e contract tests.

### Known gaps
- Calibração dos intervalos precisará de evidência de uso real.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em mastery, application e mcpserver.

### Validation executed
- Command: pose lint-spec mastery-review-scheduling --ready-check
- Result: registrar após gate.

### Residual risks
- Recomendações iniciais usarão heurísticas conservadoras.

### Follow-ups
- [covered: v1-integrated-acceptance] Validar retenção e transferência em sessões reais.
