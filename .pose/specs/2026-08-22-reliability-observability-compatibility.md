---
slug: reliability-observability-compatibility
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: local-event-store, mcp-stdio-foundation, safe-check-executor, mastery-review-scheduling, administrative-cli-fixtures, security-privacy-hardening
priority: 230
components: reliability, observability, compatibility
delivers:
---

# Spec: reliability-observability-compatibility

## 1. Intent

### Goal
Comprovar recuperação, desempenho local, diagnósticos, logging e compatibilidade do runtime V1.

### Business value
Tornar sessões retomáveis e problemas operacionais explicáveis sem depender de inspeção manual de arquivos.

### Constraints
- stdout do serve permanece exclusivo do protocolo.
- Logs são locais, estruturados, limitados e sem código por padrão.
- Compatibilidade declarada exige CI ou evidência equivalente.

### Non-goals
- Telemetria remota, dashboard, tracing distribuído ou SLO de serviço cloud.

## 2. Requirements

### Functional
- R1: Recuperar após interrupção entre evento e snapshot sem perder evento confirmado.
- R2: Detectar e diagnosticar log truncado, snapshot adulterado, lock órfão e evidência ausente.
- R3: Implementar logs em stderr com correlation ID, tool, duração, status e tamanho, sem payload sensível.
- R4: Completar ailearn doctor para binário, catálogo, estado, workspace, checks e conexão estática MCP.
- R5: Atender startup até 1 segundo e queries p95 até 100 ms no volume-alvo.
- R6: Limitar memória, output, tempo e concorrência em operações custosas.
- R7: Testar Linux e macOS; documentar Windows como suportado somente após evidência.
- R8: Testar atualização de pack sem alterar sessão fixada e migração de evento suportada.
- R9: Detectar qualquer escrita em stdout fora de frames MCP.
- R10: Fornecer version info e relatório diagnóstico sanitizado.

### Non-functional
- Testes de recuperação devem ser determinísticos e repetíveis.
- Benchmarks são informativos com thresholds documentados.

### Security
- Diagnóstico não expõe roots, ambiente completo, código ou secrets.

### Compatibility
- Fixar matriz de Go, OS, SDK MCP, schema de catálogo e schema de evento.

## 3. Technical Plan

### Affected areas
- internal/eventstore/, internal/mcpserver/, internal/cli/, internal/diagnostics/, tests/

### Artifacts
- created: internal/diagnostics/doctor.go
- created: internal/diagnostics/report.go
- created: internal/diagnostics/doctor_test.go
- created: internal/eventstore/fault_test.go
- created: internal/mcpserver/lifecycle_test.go
- created: internal/mcpserver/benchmark_test.go
- modified: internal/cli/doctor.go
- created: tests/compatibility/
- created: docs/compatibility.md
- created: docs/troubleshooting.md

### Delivery targets
Nenhum novo; qualidade transversal.

### API/contract changes
- Estabilizar schema do relatório doctor e version info.

### Data/storage changes
- Implementar política de compactação/retenção se medições exigirem, preservando log auditável.

### Technical risks
- Thresholds variam por máquina.
- Teste de crash pode ser frágil.

## 4. Tasks

### Planning
- [ ] Definir fault matrix e compatibility matrix.
- [ ] Definir benchmark dataset no volume V1.

### Implementation
- [ ] Completar doctor e relatório sanitizado.
- [ ] Implementar fault injection e recovery cases.
- [ ] Implementar logging e correlation.
- [ ] Implementar limites e benchmarks.
- [ ] Criar suites de compatibilidade e docs.

### Validation
- [ ] Repetir fault suite e benchmarks.
- [ ] Executar stdio lifecycle e stdout purity.
- [ ] Executar CI nas plataformas declaradas.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Observabilidade local não deve coletar conteúdo do aluno.
- Options considered: logs detalhados; nenhum log; metadados estruturados.
- Decision: registrar metadados operacionais e correlation IDs sem payload.
- Rationale: diagnosticar com privacidade.
- Consequences: reprodução pode exigir export explícito adicional do usuário.

## 6. Validation

### Strategy
Combinar fault injection, golden diagnostics, benchmarks e CI multi-OS.

### Deterministic checks
- Test: go test -race ./... e suites de fault/compatibility.
- Lint: gofmt -l .
- Typecheck: go vet ./...
- Build: go build ./cmd/ailearn em plataformas declaradas.
- Security / Contract: stdout purity, redaction e schema compatibility.

### Execution log
- Pendente.

### Results summary
- Nenhuma evidência de confiabilidade ou compatibilidade.

### Requirement trace
- Mapear R1–R10 a fault cases, benchmarks, CI e diagnostics goldens.

### Known gaps
- Windows permanece desejável até execução real.

## 7. Final Report

### Delivered scope
Nenhum; planejamento.

### Files and modules changed
- Planejados em diagnostics, stores, mcpserver, CLI, tests e docs.

### Validation executed
- Command: pose lint-spec reliability-observability-compatibility --ready-check
- Result: registrar após gate.

### Residual risks
- Performance será avaliada com hardware e dataset documentados.

### Follow-ups
- [covered: installation-documentation-ci] Automatizar a matriz de compatibilidade.
