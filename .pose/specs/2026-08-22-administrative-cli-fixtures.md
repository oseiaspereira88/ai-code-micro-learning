---
slug: administrative-cli-fixtures
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: go-runtime-foundation, catalog-schema-loader, local-event-store, workspace-observation-baselines, safe-check-executor
priority: 140
components: cli, fixtures
delivers:
---

# Spec: administrative-cli-fixtures

## 1. Intent

### Goal
Completar a CLI administrativa e preparar fixtures de desafios com consentimento, preview e proteção contra sobrescrita.

### Business value
Tornar instalação, autoria, diagnóstico e suporte utilizáveis sem duplicar a conversa pedagógica.

### Constraints
- CLI é determinística e não chama LLM.
- workspace prepare é acionado diretamente pelo usuário, não por tool MCP.
- Nenhum comando sobrescreve por padrão.

### Non-goals
- TUI interativa ou chat.
- Gerenciar sessões como interface pedagógica primária.

## 2. Requirements

### Functional
- R1: Expor serve, init, doctor, version e ajuda.
- R2: Expor catalog validate, list e show com saída humana e JSON estável.
- R3: Expor session inspect, progress show e progress export sem mutação pedagógica.
- R4: Implementar workspace prepare com destination explícito e plano de escrita.
- R5: Recusar destino inseguro, path traversal, symlink escape e overwrite por padrão.
- R6: Materializar fixture reproduzível e registrar manifesto com hashes.
- R7: Fornecer dry-run para toda operação que escrever múltiplos arquivos.
- R8: Usar códigos de saída e mensagens consistentes, com logs em stderr.

### Non-functional
- Comandos de leitura devem funcionar offline.
- Saída JSON deve ser versionada e testada por golden files.

### Security
- Nunca incluir secrets em export ou doctor.
- Validar todas as entradas de path e permissões.

### Compatibility
- Comandos mantêm aliases somente quando documentados; mudanças quebráveis exigem versão maior.

## 3. Technical Plan

### Affected areas
- internal/cli/, internal/fixtures/, cmd/codinho/

### Artifacts
- modified: internal/cli/root.go
- created: internal/cli/catalog.go
- created: internal/cli/session.go
- created: internal/cli/progress.go
- created: internal/cli/workspace.go
- created: internal/cli/output.go
- created: internal/fixtures/manifest.go
- created: internal/fixtures/materialize.go
- created: internal/fixtures/materialize_test.go
- created: internal/cli/cli_test.go

### Delivery targets
O alvo planejado é a surface `codinho-cli`, com módulo `cmd/codinho`, perfil
`cli-surface` e entrypoint `cmd/codinho/main.go`. Registrá-lo como alvo
tipado quando esses caminhos forem materializados, antes do closeout desta
spec.

### API/contract changes
- Estabilizar comandos, flags, JSON e códigos de saída da CLI V1.

### Data/storage changes
- Manifestos de fixture e exports explícitos.

### Technical risks
- Fixture pode sobrescrever trabalho por erro de resolução.
- CLI pode duplicar lógica dos casos de uso.

## 4. Tasks

### Planning
- [ ] Definir árvore de comandos e schemas JSON.
- [ ] Definir protocolo de preview, confirmação e overwrite.

### Implementation
- [ ] Implementar comandos de catálogo, sessão e progresso.
- [ ] Implementar materialização transacional de fixture.
- [ ] Implementar outputs e códigos de saída.
- [ ] Reutilizar casos de uso do núcleo.
- [ ] Criar testes CLI end-to-end.

### Validation
- [ ] Testar destinos vazios, não vazios, symlinks e falha intermediária.
- [ ] Testar goldens humana/JSON e compatibilidade.
- [ ] Executar surface-check com reachability.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Debug challenges precisam de código inicial, mas o tutor não pode editar o aluno.
- Options considered: tool MCP de escrita; download manual; comando CLI explícito.
- Decision: workspace prepare somente pela CLI com preview e confirmação.
- Rationale: separa consentimento material de tutoria.
- Consequences: o fluxo de onboarding deve ensinar o comando.

## 6. Validation

### Strategy
Usar CLI real em diretórios temporários e golden outputs.

### Deterministic checks
- Test: go test ./internal/cli/... ./internal/fixtures/...
- Lint: gofmt -l internal/cli internal/fixtures
- Typecheck: go vet ./internal/cli/... ./internal/fixtures/...
- Build: go build ./cmd/codinho
- Security / Contract: traversal suite; pose surface-check --spec administrative-cli-fixtures --strict.

### Execution log
- Pendente.

### Results summary
- CLI completa e fixtures ainda não entregues.

### Requirement trace
- Mapear R1–R8 a e2e CLI, goldens e security tests.

### Known gaps
- Instaladores e distribuição pertencem a installation-documentation-ci.

## 7. Final Report

### Delivered scope
Nenhum; planejamento.

### Files and modules changed
- Planejados em cli, fixtures e cmd/codinho.

### Validation executed
- Command: pose lint-spec administrative-cli-fixtures --ready-check
- Result: registrar após gate.

### Residual risks
- Windows path semantics exigem CI futuro.

### Follow-ups
- [covered: installation-documentation-ci] Documentar e testar instalação da CLI.
