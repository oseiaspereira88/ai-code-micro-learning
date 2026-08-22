---
slug: assistance-hints-detours
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: session-orchestration-disclosure
priority: 70
components: assistance, sessions, mcp-server
delivers:
---

# Spec: assistance-hints-detours

## 1. Intent

### Goal
Implementar a escada de auxílio, lembrança sintática, conteúdo conceitual e desvios pedagógicos sem avanço implícito.

### Business value
Oferecer a menor ajuda suficiente e permitir aprendizagem contextual sem tomar o teclado do aluno.

### Constraints
- Subir no máximo um nível por solicitação.
- Nível de solução exige intenção explícita e registro.
- Exemplos devem evitar resolver o desafio ativo quando possível.

### Non-goals
- Redigir explicações abertas dentro do MCP.
- Avaliar ou concluir passos.

## 2. Requirements

### Functional
- R1: Entregar pistas ordenadas nos níveis 1 a 6, respeitando a política da sessão.
- R2: Tratar objetivo e critérios como nível 0 sem consumo de pista.
- R3: Bloquear níveis não permitidos com erro estável e zero mudança de progresso.
- R4: Registrar nível, contexto, timestamp e efeito de cada pista aceita.
- R5: Fornecer conteúdo canônico de conceito e lembrança sintática com custo configurável.
- R6: Abrir, consultar e fechar desvio sem alterar o nó instrucional original.
- R7: Marcar solução revelada e agendar uma variante futura sem promover autonomia.
- R8: Expor hint_request, concept_content_get, syntax_recall_get, learning_detour_start e learning_detour_finish.

### Non-functional
- Pistas da mesma versão devem ser determinísticas.
- Conteúdo conceitual deve ser acessível sem workspace.

### Security
- Conteúdo autorado é dado; não pode invocar tools nem ampliar permissões.
- Disclosure deve filtrar gabaritos e fragmentos acima do nível.

### Compatibility
- Novos níveis podem ser adicionados sem reinterpretar registros existentes.

## 3. Technical Plan

### Affected areas
- internal/assistance/, internal/session/, internal/mcpserver/

### Artifacts
- created: internal/assistance/service.go
- created: internal/assistance/ladder.go
- created: internal/assistance/detour.go
- created: internal/assistance/service_test.go
- created: internal/assistance/disclosure_test.go
- modified: internal/session/service.go
- created: internal/mcpserver/assistance_tools.go
- created: internal/mcpserver/assistance_contract_test.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar cinco tools e efeito solution_revealed.

### Data/storage changes
- Persistir hint_requested, detour_started, detour_finished e solution_revealed.

### Technical risks
- Pistas podem crescer sem diferença real de intensidade.
- Conteúdo conceitual pode vazar a solução indiretamente.

## 4. Tasks

### Planning
- [ ] Definir rubrica mensurável para intensidade de pista.
- [ ] Definir quando syntax recall conta como pista por modo.

### Implementation
- [ ] Implementar policy e ladder.
- [ ] Implementar conteúdo canônico e filtragem.
- [ ] Implementar lifecycle de detour.
- [ ] Implementar registro de solução e review futura.
- [ ] Expor tools e contract tests.

### Validation
- [ ] Testar cada transição de nível e todos os bloqueios.
- [ ] Executar golden tests de não revelação.
- [ ] Testar retomada do mesmo passo após detour.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Proibição absoluta de código prejudica alguns alunos.
- Options considered: sem código; ajuda livre; escada explícita.
- Decision: usar escada com solução somente no nível máximo.
- Rationale: controla custo pedagógico sem bloquear aprendizagem.
- Consequences: catálogo precisa autorar pistas distintas e auditáveis.

## 6. Validation

### Strategy
Usar golden disclosure, matriz de políticas e replay de eventos.

### Deterministic checks
- Test: go test ./internal/assistance/... ./internal/session/...
- Lint: gofmt -l internal/assistance internal/session
- Typecheck: go vet ./internal/assistance/... ./internal/session/...
- Build: go test ./internal/assistance/... -run ^$
- Security / Contract: testes negativos de leak e schemas MCP.

### Execution log
- Pendente.

### Results summary
- Nenhum auxílio entregue.

### Requirement trace
- Mapear R1–R8 a matriz de nível, política, evento e tool.

### Known gaps
- Qualidade linguística final depende da skill e do agente host.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em assistance, session e mcpserver.

### Validation executed
- Command: pose lint-spec assistance-hints-detours --ready-check
- Result: registrar após validação.

### Residual risks
- Revisão humana das pistas permanece obrigatória.

### Follow-ups
- [covered: tutor-skill-host-integration] Validar adaptação das explicações pelo tutor.
