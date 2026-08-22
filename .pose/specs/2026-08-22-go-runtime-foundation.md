---
slug: go-runtime-foundation
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: architecture-decision-baseline
priority: 10
components: runtime, cli
delivers:
---

# Spec: go-runtime-foundation

## 1. Intent

### Goal
Criar o módulo Go, o binário ailearn e as fundações de configuração, versão e lifecycle usadas por todas as entregas.

### Business value
Transformar o repositório documental em uma base executável pequena, reproduzível e validável.

### Constraints
- Usar Go 1.25 ou superior e fixar dependências em go.mod.
- Reservar stdout ao MCP quando serve estiver ativo.
- Preferir biblioteca padrão e não antecipar frameworks.

### Non-goals
- Implementar domínio pedagógico, catálogo, persistência ou tools MCP.
- Entregar instalador ou release.

## 2. Requirements

### Functional
- R1: O repositório deve possuir módulo Go válido e reproduzível.
- R2: O binário deve expor version, doctor e ajuda determinística.
- R3: O processo deve suportar cancelamento por contexto e códigos de saída estáveis.
- R4: A configuração deve resolver defaults sem depender de segredos ou rede.
- R5: Logs devem ir para stderr e nunca contaminar uma futura sessão MCP.
- R6: O build deve incorporar versão e commit sem exigir Git em runtime.

### Non-functional
- Inicialização dos comandos administrativos deve ocorrer em menos de um segundo em ambiente típico.
- Pacotes de entrada devem permanecer finos e testáveis.

### Security
- Não ler .env implicitamente nem registrar valores de ambiente.
- Não executar arquivos do projeto durante diagnóstico.

### Compatibility
- Suportar Linux e macOS; evitar decisões que impeçam Windows.

## 3. Technical Plan

### Affected areas
- go.mod, go.sum, cmd/ailearn/, internal/cli/, internal/config/

### Artifacts
- created: go.mod
- created: go.sum
- created: cmd/ailearn/main.go
- created: internal/cli/root.go
- created: internal/cli/version.go
- created: internal/cli/doctor.go
- created: internal/config/config.go
- created: internal/cli/root_test.go
- modified: .pose/indexes/module-metadata.json
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
Nenhum público nesta etapa; o CLI completo será entregue por administrative-cli-fixtures.

### API/contract changes
- Definir códigos de saída e interfaces internas para relógio, filesystem e version info.

### Data/storage changes
- Nenhum estado persistente além de configuração somente leitura.

### Technical risks
- Uma CLI genérica prematura pode crescer sem necessidade.
- Logs acidentais em stdout quebrariam o transporte MCP futuro.

## 4. Tasks

### Planning
- [ ] Confirmar module path e versão mínima do Go.
- [ ] Registrar ADR adicional se a estrutura divergir da baseline.

### Implementation
- [ ] Inicializar módulo e binário.
- [ ] Implementar roteamento mínimo com biblioteca padrão.
- [ ] Implementar version e diagnóstico sem efeitos colaterais.
- [ ] Registrar o módulo Go na matriz POSE.

### Validation
- [ ] Executar gofmt, go test ./..., go vet ./... e build.
- [ ] Verificar stdout vazio em caminhos de erro do futuro comando serve.
- [ ] Executar pose validate --strict --module . --report.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: A V1 precisa de poucos comandos administrativos.
- Options considered: framework CLI; biblioteca padrão.
- Decision: iniciar com flag e roteamento explícito.
- Rationale: minimizar dependências e manter contratos visíveis.
- Consequences: migrar exige decisão se a superfície crescer significativamente.

## 6. Validation

### Strategy
Cobrir parsing, códigos de saída, cancelamento, build reproduzível e ausência de output indevido.

### Deterministic checks
- Test: go test ./...
- Lint: gofmt -l . deve produzir saída vazia.
- Typecheck: go vet ./...
- Build: go build ./cmd/ailearn
- Security / Contract: secret scan do diff e teste de stdout/stderr.

### Execution log
- Pendente; nenhum módulo Go existe enquanto a spec está draft.

### Results summary
- Planejamento concluído; implementação e evidências pendentes.

### Requirement trace
- Mapear R1–R6 para testes de CLI, build e relatório POSE no closeout.

### Known gaps
- O module path deve ser confirmado antes de status in-progress.

## 7. Final Report

### Delivered scope
Nenhum; spec de implementação ainda não iniciada.

### Files and modules changed
- Planejados nos artefatos da seção 3.

### Validation executed
- Command: pose lint-spec go-runtime-foundation --ready-check
- Result: registrar após o gate de planejamento.

### Residual risks
- Compatibilidade Windows ficará sem prova até existir CI correspondente.

### Follow-ups
- [covered: installation-documentation-ci] Comprovar instalação e plataformas suportadas.
