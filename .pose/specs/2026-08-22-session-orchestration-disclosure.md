---
slug: session-orchestration-disclosure
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: learning-domain-model, catalog-schema-loader, local-event-store, mcp-stdio-foundation
priority: 60
components: sessions, mcp-server
delivers:
---

# Spec: session-orchestration-disclosure

## 1. Intent

### Goal
Orquestrar lifecycle de sessão, instrução única, granularidade e divulgação progressiva com estado persistente.

### Business value
Impedir que o agente revele etapas ou soluções fora do momento pedagógico autorizado.

### Constraints
- Sessão fixa versões de conteúdo no início.
- Avanço é sempre explícito na V1.
- Ajuste de granularidade não modifica a árvore canônica.

### Non-goals
- Gerar feedback semântico, executar checks ou calcular domínio.
- Implementar todos os modos especializados.

## 2. Requirements

### Functional
- R1: Iniciar sessão com desafio ou trilha, modo, profundidade, políticas, workspace opcional e tempo opcional.
- R2: Consultar e retomar sessão com nó ativo, revisão, disclosure, ações permitidas e pendências.
- R3: Manter exatamente uma instrução ativa e omitir filhos futuros, gabarito e pistas bloqueadas.
- R4: Configurar apenas propriedades mutáveis e registrar toda mudança.
- R5: Pausar, retomar, concluir e abandonar lifecycle sem inferir conclusão dos passos.
- R6: Ajustar janela entre desafio, camada, macro, meso e micro sem reescrever conteúdo.
- R7: Detectar desatualização de revisão e retries idempotentes.
- R8: Expor session_start, session_get, session_configure, session_pause, session_resume, session_finish, instruction_get e granularity_adjust.

### Non-functional
- Mesma versão e estado produzem a mesma instrução e ações.
- Consulta de sessão deve atender p95 inferior a 100 ms após aquecimento.

### Security
- Disclosure deve ser aplicado no servidor, não apenas na skill.
- IDs de sessão não autorizam outro workspace.

### Compatibility
- Sessões antigas permanecem legíveis após atualização de packs.

## 3. Technical Plan

### Affected areas
- internal/session/, internal/application/, internal/mcpserver/

### Artifacts
- created: internal/session/service.go
- created: internal/session/disclosure.go
- created: internal/session/granularity.go
- created: internal/session/service_test.go
- created: internal/session/disclosure_test.go
- modified: internal/application/session.go
- modified: internal/mcpserver/session_tools.go
- created: internal/mcpserver/session_contract_test.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar oito tools de sessão e seus efeitos ao contrato MCP.

### Data/storage changes
- Persistir eventos de lifecycle, instrução e granularidade.

### Technical risks
- Disclosure incorreto pode vazar solução.
- Reagrupar passos pode perder posição ou evidência.

## 4. Tasks

### Planning
- [ ] Definir matriz modo versus disclosure versus ação permitida.
- [ ] Definir seleção determinística do próximo nó.

### Implementation
- [ ] Implementar lifecycle e projeção da sessão.
- [ ] Implementar disclosure e instrução única.
- [ ] Implementar ajuste de granularidade.
- [ ] Expor tools e eventos.
- [ ] Cobrir branches, retries e retomada.

### Validation
- [ ] Executar testes de tabela para toda transição.
- [ ] Executar golden tests em cada nível de disclosure.
- [ ] Reiniciar servidor no meio de uma sessão e comparar estado.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Granularidade é uma visão da árvore, não conteúdo duplicado.
- Options considered: duplicar desafios por nível; podar árvore; janela derivada.
- Decision: derivar uma janela instrucional da árvore canônica.
- Rationale: evita drift e permite adaptação.
- Consequences: o cursor precisa preservar ancestry e progresso.

## 6. Validation

### Strategy
Cobrir state machine, disclosure negativo, replay e contratos MCP.

### Deterministic checks
- Test: go test ./internal/session/... ./internal/mcpserver/...
- Lint: gofmt -l internal/session internal/mcpserver
- Typecheck: go vet ./internal/session/... ./internal/mcpserver/...
- Build: go build ./cmd/ailearn
- Security / Contract: golden leak tests e pose assess integrate.

### Execution log
- Pendente; spec draft.

### Results summary
- Nenhuma sessão orquestrada entregue.

### Requirement trace
- Mapear R1–R8 a testes de lifecycle, disclosure, replay e schemas.

### Known gaps
- Adaptação automática baseada em domínio será integrada depois.

## 7. Final Report

### Delivered scope
Nenhum; planejamento.

### Files and modules changed
- Planejados em internal/session, application e mcpserver.

### Validation executed
- Command: pose lint-spec session-orchestration-disclosure --ready-check
- Result: registrar após o gate.

### Residual risks
- Conteúdo mal classificado pode contornar disclosure e exige gate editorial.

### Follow-ups
- [covered: catalog-authoring-quality] Validar classificação de conteúdo reservado.
