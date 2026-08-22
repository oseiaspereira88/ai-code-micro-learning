---
slug: mcp-stdio-foundation
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: go-runtime-foundation, learning-domain-model, catalog-schema-loader, local-event-store
priority: 50
components: mcp-server
delivers:
---

# Spec: mcp-stdio-foundation

## 1. Intent

### Goal
Entregar o servidor MCP local por stdio, seu envelope, erros, instructions e adaptadores iniciais de catálogo e sessão.

### Business value
Criar a primeira superfície estruturada utilizável pelo agente sem construir uma UI própria.

### Constraints
- Usar o SDK MCP oficial para Go compatível com Go 1.25.
- Reservar stdout exclusivamente ao protocolo.
- Não depender de resources para fluxos obrigatórios.
- Não incorporar LLM.

### Non-goals
- Transporte HTTP, OAuth, nuvem ou tools de edição.
- Implementar todos os casos de uso pedagógicos nesta spec.

## 2. Requirements

### Functional
- R1: Iniciar por ailearn serve e encerrar de forma limpa quando stdin fechar ou contexto cancelar.
- R2: Publicar nome, versão e server instructions com as invariantes centrais nos primeiros 512 caracteres.
- R3: Retornar envelope consistente com status, request, sessão, nó, efeito, disclosure, ações, data e warnings.
- R4: Mapear erros de domínio para os códigos estáveis definidos em PROJECT.md sem vazar detalhes internos.
- R5: Expor catalog_search, catalog_get, session_start, session_get e instruction_get em uma fatia mínima.
- R6: Anotar tools de leitura e escrita e aplicar schemas estritos.
- R7: Garantir que resources opcionais não sejam requisito para o host.
- R8: Rejeitar chamadas mutáveis com expected_revision obsoleta e suportar retry idempotente.

### Non-functional
- Startup deve atender ao limite de um segundo após cache aquecido.
- Contratos devem possuir golden tests e fixture de cliente real.

### Security
- Nenhuma tool aceita comando, path irrestrito ou conteúdo executável.
- Erros e logs não podem expor root, ambiente ou payload sensível.

### Compatibility
- Testar Codex CLI e extensão IDE posteriormente sem acoplar o núcleo ao host.

## 3. Technical Plan

### Affected areas
- cmd/ailearn/, internal/mcpserver/, internal/application/

### Artifacts
- modified: cmd/ailearn/main.go
- created: internal/mcpserver/server.go
- created: internal/mcpserver/instructions.go
- created: internal/mcpserver/envelope.go
- created: internal/mcpserver/errors.go
- created: internal/mcpserver/catalog_tools.go
- created: internal/mcpserver/session_tools.go
- created: internal/mcpserver/server_test.go
- created: internal/mcpserver/contract_test.go
- created: internal/application/catalog.go
- created: internal/application/session.go

### Delivery targets
O alvo planejado é o contrato `ailearn-mcp-v1`, com módulo
`internal/mcpserver`, perfil `api-contract` e entrypoint
`cmd/ailearn/main.go`. Registrá-lo como alvo tipado quando esses caminhos
forem materializados, antes do closeout desta spec.

### API/contract changes
- Criar contrato MCP v1 e schemas das cinco tools iniciais.

### Data/storage changes
- Sessões usam o event store; nenhuma nova forma de persistência.

### Technical risks
- Drift entre schemas, domínio e documentação.
- Escrita acidental em stdout invalida todo o transporte.
- API do SDK pode evoluir.

## 4. Tasks

### Planning
- [ ] Fixar versão do SDK e revisar licença e vulnerabilidades.
- [ ] Definir schemas JSON e códigos de erro.

### Implementation
- [ ] Implementar lifecycle stdio e logging em stderr.
- [ ] Implementar envelope e mapper de erros.
- [ ] Implementar as cinco tools verticais.
- [ ] Publicar instructions e resources opcionais sem dependência.
- [ ] Criar cliente de contrato para testes.

### Validation
- [ ] Executar unit, contrato e integração por stdio real.
- [ ] Testar stdout contaminado, JSON inválido, retry e cancelamento.
- [ ] Executar pose assess integrate para o contrato MCP.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Tools granulares evitam uma ação genérica ambígua.
- Options considered: execute_action; tools semânticas por caso de uso.
- Decision: expor tools pequenas com efeitos declarados.
- Rationale: melhora segurança, roteamento e auditabilidade.
- Consequences: toolset maior exige allowed_actions e skill de roteamento.

## 6. Validation

### Strategy
Validar protocolo por processo real, golden schemas, erros negativos e integração com stores temporários.

### Deterministic checks
- Test: go test ./internal/mcpserver/... ./internal/application/...
- Lint: gofmt -l internal/mcpserver internal/application
- Typecheck: go vet ./internal/mcpserver/... ./internal/application/...
- Build: go build ./cmd/ailearn
- Security / Contract: pose assess integrate; govulncheck; teste de stdout e schema fuzz.

### Execution log
- Pendente; contrato ainda não implementado.

### Results summary
- Nenhuma tool MCP entregue.

### Requirement trace
- Mapear R1–R8 a contract tests e resultado estruturado de integração.

### Known gaps
- Tools pedagógicas adicionais serão entregues pelas specs dependentes.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em cmd/ailearn, internal/mcpserver e internal/application.

### Validation executed
- Command: pose lint-spec mcp-stdio-foundation --ready-check
- Result: registrar após validação.

### Residual risks
- Compatibilidade de host só será comprovada em tutor-skill-host-integration.

### Follow-ups
- [covered: tutor-skill-host-integration] Validar o contrato MCP com hosts Codex reais.
