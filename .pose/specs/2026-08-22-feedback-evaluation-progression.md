---
slug: feedback-evaluation-progression
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: session-orchestration-disclosure, local-event-store, mcp-stdio-foundation
priority: 80
components: assessment, sessions, mcp-server
delivers:
---

# Spec: feedback-evaluation-progression

## 1. Intent

### Goal
Implementar feedback consultivo, avaliação híbrida, reflexão, conclusão e avanço como operações independentes e auditáveis.

### Business value
Permitir experimentação segura e avaliação justa sem transformar opinião idiomática em reprovação.

### Constraints
- MCP prepara contexto; o agente redige feedback e julgamento semântico.
- Tentativa só existe com intenção de submissão.
- Avaliação nunca conclui; conclusão nunca avança.

### Non-goals
- Executar checks ou observar workspace.
- Calcular retenção longitudinal.

## 2. Requirements

### Functional
- R1: Preparar pacote de feedback com instrução, pergunta, rubrica, observações e evidências.
- R2: Registrar feedback com tipos e progress_effect sem criar avaliação ou tentativa.
- R3: Avaliar critérios determinísticos e julgamentos qualitativos separadamente.
- R4: Exigir referência de evidência e rubrica para todo julgamento qualitativo.
- R5: Produzir met, partially_met, not_met, unverifiable ou not_applicable por critério.
- R6: Classificar achados em blocking, important_non_blocking ou advisory.
- R7: Concluir apenas quando a política permitir e avançar somente em chamada separada.
- R8: Registrar override sem falsificar avaliação positiva ou domínio.
- R9: Registrar reflexão como evidência distinta da implementação.
- R10: Expor feedback_prepare, feedback_record, step_evaluate, reflection_record, step_complete e step_advance.

### Non-functional
- Toda mutação deve ser idempotente e revisionada.
- O histórico deve preservar avaliações múltiplas do mesmo passo.

### Security
- Evidências e feedback devem aplicar redaction e escopo da sessão.
- Texto qualitativo não pode conter instruções executáveis.

### Compatibility
- Novos tipos de feedback e critérios devem ser aditivos.

## 3. Technical Plan

### Affected areas
- internal/assessment/, internal/session/, internal/mcpserver/

### Artifacts
- created: internal/assessment/service.go
- created: internal/assessment/rubric.go
- created: internal/assessment/criteria.go
- created: internal/assessment/service_test.go
- created: internal/assessment/progression_test.go
- modified: internal/session/service.go
- created: internal/mcpserver/assessment_tools.go
- created: internal/mcpserver/assessment_contract_test.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar seis tools e envelopes de avaliação e achado.

### Data/storage changes
- Persistir feedback_recorded, attempt_submitted, evaluation_recorded, reflection_recorded, step_completed e step_advanced.

### Technical risks
- O agente pode produzir avaliação inconsistente com a rubrica.
- Completion override pode ser abusado e distorcer métricas.

## 4. Tasks

### Planning
- [ ] Definir rubricas de Go idiomático e comunicação técnica.
- [ ] Definir matriz critério versus evidência e severidade.

### Implementation
- [ ] Implementar packet de feedback e registro.
- [ ] Implementar avaliação híbrida e validação de evidências.
- [ ] Implementar reflexão, conclusão, override e avanço.
- [ ] Expor tools e allowed_actions.
- [ ] Cobrir sequências inválidas e avaliações repetidas.

### Validation
- [ ] Executar tabelas de transição e contract tests.
- [ ] Testar feedback com progress_effect sem avanço.
- [ ] Testar spoof de evidência e revisão obsoleta.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: O servidor não possui LLM, mas precisa persistir avaliação semântica.
- Options considered: fingir avaliação no servidor; não persistir; aceitar julgamento estruturado do tutor.
- Decision: receber julgamentos tipados, cada um ligado a rubrica e evidência.
- Rationale: mantém responsabilidades honestas e auditáveis.
- Consequences: a skill precisa orientar o agente e testes devem validar o envelope, não a linguagem aberta.

## 6. Validation

### Strategy
Testar operações isoladas e sequências completas, incluindo feedback sem efeito e overrides.

### Deterministic checks
- Test: go test ./internal/assessment/... ./internal/session/...
- Lint: gofmt -l internal/assessment internal/session
- Typecheck: go vet ./internal/assessment/... ./internal/session/...
- Build: go test ./internal/assessment/... -run ^$
- Security / Contract: schema tests, redaction e rejeição de evidence IDs fora da sessão.

### Execution log
- Pendente.

### Results summary
- Nenhum fluxo de avaliação entregue.

### Requirement trace
- Mapear R1–R10 a testes unitários, de contrato e de sequência.

### Known gaps
- Checks determinísticos serão conectados em safe-check-executor.

## 7. Final Report

### Delivered scope
Nenhum; planejamento.

### Files and modules changed
- Planejados em assessment, session e mcpserver.

### Validation executed
- Command: pose lint-spec feedback-evaluation-progression --ready-check
- Result: registrar após o gate.

### Residual risks
- Calibração das rubricas exige sessões humanas.

### Follow-ups
- [covered: v1-integrated-acceptance] Comprovar separação operacional em E2E.
