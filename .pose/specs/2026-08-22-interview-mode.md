---
slug: interview-mode
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: feedback-evaluation-progression, safe-check-executor, tutor-skill-host-integration, learning-practice-debug-modes
priority: 160
components: interview-mode, tutor-skill
delivers:
---

# Spec: interview-mode

## 1. Intent

### Goal
Implementar simulação técnica genérica com briefing completo, tempo opcional, auxílio controlado e revisão somente ao final.

### Business value
Praticar autonomia, comunicação e solução sob condições próximas de live coding sem alegar vínculo com terceiros.

### Constraints
- Não reproduzir perguntas, marcas ou processos proprietários.
- Não prometer antifraude em máquina local.
- Pistas ficam bloqueadas ou registradas conforme política escolhida antes do início.

### Non-goals
- Proctoring, gravação de tela, vigilância ou ranking.
- Avaliação de contratação.

## 2. Requirements

### Functional
- R1: Iniciar simulação com desafio, duração, checks permitidos e política de auxílio fixada.
- R2: Exibir apenas briefing e critérios públicos durante execução.
- R3: Registrar tempo ativo, pausas autorizadas e término sem manipular relógio do sistema.
- R4: Bloquear ou registrar pedidos de pista sem revelar conteúdo proibido.
- R5: Permitir finalização antecipada ou por timeout com estado explícito.
- R6: Avaliar ao final código, testes, raciocínio, complexidade, comunicação e trade-offs.
- R7: Produzir relatório descritivo com evidências, gaps e exercícios recomendados.
- R8: Declarar limites de integridade e ausência de associação com terceiros.

### Non-functional
- Clock deve ser injetável e determinístico em testes.
- Relatório não deve reduzir desempenho a um score único.

### Security
- Não coletar áudio, vídeo, dados pessoais ou atividade fora do workspace.

### Compatibility
- Simulação deve funcionar sem recursos MCP opcionais.

## 3. Technical Plan

### Affected areas
- internal/session/, internal/assessment/, .agents/skills/codinho/, packs/go-interviews/

### Artifacts
- created: internal/session/interview.go
- created: internal/session/interview_test.go
- created: internal/assessment/interview_report.go
- created: internal/assessment/interview_report_test.go
- modified: .agents/skills/codinho/references/session-modes.md
- created: .agents/skills/codinho/assets/interview-report-template.md
- created: testdata/host/interview-transcripts/

### Delivery targets
Nenhum novo; amplia a capability de tutoria planejada.

### API/contract changes
- Adicionar propriedades de timer e policy lock à sessão.

### Data/storage changes
- Persistir interview_started, interview_finished e report_generated.

### Technical risks
- Timer pode criar pressão sem valor pedagógico.
- Agente pode oferecer ajuda por conversa fora das tools.

## 4. Tasks

### Planning
- [ ] Definir rubrica e linguagem neutra do relatório.
- [ ] Definir políticas de pausa, timeout e auxílio.

### Implementation
- [ ] Implementar timer e lock de policy.
- [ ] Integrar bloqueio de disclosure.
- [ ] Implementar avaliação e relatório final.
- [ ] Atualizar skill e templates.
- [ ] Criar transcripts de abuso e encerramento.

### Validation
- [ ] Testar clock, timeout, pausa e restart.
- [ ] Testar tentativas de obter pistas por linguagem indireta.
- [ ] Revisar privacidade e neutralidade do relatório.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Simulação local não consegue garantir integridade forte.
- Options considered: fingir antifraude; proctoring; transparência e honor system.
- Decision: declarar limites e focar prática, não certificação.
- Rationale: evita vigilância e promessas falsas.
- Consequences: resultados são evidência pedagógica, não credencial.

## 6. Validation

### Strategy
Combinar clocks falsos, policy tests, transcripts adversariais e revisão de relatório.

### Deterministic checks
- Test: go test ./internal/session/... ./internal/assessment/...
- Lint: gofmt -l internal/session internal/assessment
- Typecheck: go vet ./internal/session/... ./internal/assessment/...
- Build: go build ./cmd/codinho
- Security / Contract: privacy review e leak tests.

### Execution log
- Pendente.

### Results summary
- Modo entrevista não entregue.

### Requirement trace
- Mapear R1–R8 a timer tests, transcripts e golden report.

### Known gaps
- Validação humana de realismo e utilidade será parte do aceite.

## 7. Final Report

### Delivered scope
Nenhum; planejamento.

### Files and modules changed
- Planejados em session, assessment, skill e testdata.

### Validation executed
- Command: pose lint-spec interview-mode --ready-check
- Result: registrar após gate.

### Residual risks
- Variações de host podem afetar bloqueio conversacional.

### Follow-ups
- [covered: go-interviews-pack] Fornecer quatro simulações curadas.
