---
slug: go-runtime-foundation
status: in-progress
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
- [x] Confirmar module path e versão mínima do Go.
- [x] Registrar ADR adicional se a estrutura divergir da baseline.

### Implementation
- [x] Inicializar módulo e binário.
- [x] Implementar roteamento mínimo com biblioteca padrão.
- [x] Implementar version e diagnóstico sem efeitos colaterais.
- [x] Registrar o módulo Go na matriz POSE.

### Validation
- [x] Executar gofmt, go test ./..., go vet ./... e build.
- [x] Verificar stdout vazio em caminhos de erro (usage/comando desconhecido).
- [x] Executar pose validate --strict --module . --report.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: A V1 precisa de poucos comandos administrativos.
- Options considered: framework CLI; biblioteca padrão.
- Decision: iniciar com flag e roteamento explícito.
- Rationale: minimizar dependências e manter contratos visíveis.
- Consequences: migrar exige decisão se a superfície crescer significativamente.

### Decision 2
- Date: 2026-08-22
- Context: `go.mod` exige um module path estável; o remote atual do repositório
  (`github.com/oseiaspereira88/ai-code-micro-learning`) diverge do nome do
  produto (`ailearn`) usado em `PROJECT.md`/`README.md`.
- Options considered: usar o nome do remote atual; usar o nome do produto;
  domínio próprio.
- Decision: `module github.com/oseiaspereira88/ailearn`, confirmado pelo
  responsável do projeto.
- Rationale: o import path deve refletir a identidade do produto, não o nome
  histórico do repositório; evita reescrever imports numa renomeação futura
  do repositório.
- Consequences: se o repositório GitHub não for renomeado/criado como
  `ailearn`, o module path não corresponderá ao caminho de `go get`
  automático a partir do remote atual — aceito como risco conhecido, sem
  impacto na V1 (sem publicação de módulo externo prevista).
- Estrutura (`cmd/ailearn/`, `internal/cli/`, `internal/config/`) confirmada
  compatível com
  [ADR agent-mcp-and-core-boundaries](../adr/2026-08-22-agent-mcp-and-core-boundaries.md):
  nenhum pacote sob `internal/` importa MCP ou protocolo; a CLI é um
  adaptador de entrada isolado. Nenhum ADR adicional é necessário nesta
  etapa.

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
- `gofmt -l .` → saída vazia (2026-08-22).
- `go vet ./...` → sem findings (2026-08-22).
- `go build ./cmd/ailearn` → binário gerado com sucesso (2026-08-22).
- `go test ./... -race` → `ok internal/cli`, sem data races (2026-08-22).
- `pose validate --strict --module . --report` → `Result: SUCCESS` (test, vet, build) (2026-08-22).
- `pose check --strict` → `Result: SUCCESS` (2026-08-22).

### Results summary
- Módulo Go inicializado (`github.com/oseiaspereira88/ailearn`, `go 1.25.0`)
  com binário `cmd/ailearn`, roteamento mínimo em `internal/cli` (version,
  doctor, help, usage) e resolução de config em `internal/config`.
- Seis testes automatizados cobrem parsing, exit codes, determinismo de
  `version`, `doctor`, `help` e respeito a contexto cancelado.
- Módulo `.` registrado em `.pose/indexes/module-metadata.json` (domínio
  `runtime`) e em `moduleOverrides` de `validation-matrix.json` (check de
  build adicional).

### Requirement trace
- R1 [satisfied] test:internal/cli check:build report:.pose/reports/2026-08-22-standard-validate-native.md
- R2 [satisfied] test:TestRunVersionIsDeterministic test:TestRunDoctorReportsOK test:TestRunHelp
- R3 [satisfied] test:TestRunRespectsCancelledContext
- R4 [satisfied] report:internal/config/config.go
- R5 [satisfied] test:TestRunNoArgsPrintsUsageToStderr
- R6 [satisfied] report:internal/cli/version.go

### Known gaps
- `gofmt -l .` não está automatizado em `pose validate` porque o exit code do
  comando não reflete achados (sempre 0); verificado manualmente nesta
  entrega e deve ser revisitado quando um wrapper multiplataforma existir.

## 7. Final Report

### Delivered scope
Módulo Go executável (`ailearn`) com comandos administrativos `version`,
`doctor`, `help` e usage determinística; base de configuração e roteamento
para specs futuras. Nenhum tool MCP, catálogo ou persistência foi
implementado, conforme os non-goals da spec.

### Files and modules changed
- `go.mod` (criado)
- `cmd/ailearn/main.go` (criado)
- `internal/cli/root.go` (criado)
- `internal/cli/version.go` (criado)
- `internal/cli/doctor.go` (criado)
- `internal/cli/root_test.go` (criado)
- `internal/config/config.go` (criado)
- `.pose/indexes/module-metadata.json` (modificado)
- `.pose/indexes/validation-matrix.json` (modificado)
- `.pose/specs/2026-08-22-go-runtime-foundation.md` (atualizado)

### Validation executed
- Command: pose validate --strict --module . --report
- Result: `Result: SUCCESS` (go test, go vet, go build)

### Residual risks
- Compatibilidade Windows ficará sem prova até existir CI correspondente.

### Follow-ups
- [covered: installation-documentation-ci] Comprovar instalação e plataformas suportadas.
