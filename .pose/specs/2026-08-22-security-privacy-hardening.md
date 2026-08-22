---
slug: security-privacy-hardening
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: mcp-stdio-foundation, workspace-observation-baselines, safe-check-executor, tutor-skill-host-integration, administrative-cli-fixtures
priority: 220
components: security, privacy
delivers:
---

# Spec: security-privacy-hardening

## 1. Intent

### Goal
Executar threat modeling e hardening transversal de MCP, CLI, workspace, checks, evidências, catálogo e skill.

### Business value
Garantir que a ferramenta de aprendizagem preserve código, dados e máquina local enquanto observa e executa verificações.

### Constraints
- Operação local, zero telemetria e zero upload por padrão.
- Nenhuma tool de edição ou comando livre.
- Mínimo privilégio e fail closed nas fronteiras.

### Non-goals
- Sandbox de kernel para código deliberadamente hostil.
- Autenticação remota, multiusuário ou compliance certificado.

## 2. Requirements

### Functional
- R1: Publicar threat model com assets, atores, trust boundaries, abuso e risco residual.
- R2: Impedir traversal, symlink escape, env injection, shell injection e execução fora do registry.
- R3: Tratar código, comentários, fixtures e outputs como conteúdo não confiável.
- R4: Excluir e redigir secrets, .env, credenciais, paths sensíveis e dados pessoais.
- R5: Aplicar permissões restritivas e retenção configurável a estado e evidências.
- R6: Fornecer export e remoção local explícitos sem apagar código ou conteúdo versionado.
- R7: Negar rede por padrão e exigir aprovação visível quando necessária.
- R8: Avaliar dependências, licenças e vulnerabilidades com política de atualização.
- R9: Provar que repositório sujo, índice Git e arquivos fora de escopo são preservados.
- R10: Documentar limites de segurança, testes reservados e modo entrevista.

### Non-functional
- Todo finding crítico ou alto bloqueia closeout sem exceção governada.
- Logs não contêm código por padrão.

### Security
- Aplicar integralmente .pose/rules/security.md e revisão independente.

### Compatibility
- Controles devem degradar explicitamente quando plataforma não oferecer uma primitiva.

## 3. Technical Plan

### Affected areas
- internal/security/, internal/workspace/, internal/checks/, internal/evidence/, internal/mcpserver/, internal/cli/, docs/

### Artifacts
- created: docs/security/threat-model.md
- created: docs/security/privacy.md
- created: docs/security/executor-limitations.md
- created: internal/security/policy.go
- created: internal/security/redaction.go
- created: internal/security/policy_test.go
- created: testdata/security/
- modified: internal/workspace/security_test.go
- modified: internal/checks/security_test.go
- modified: internal/mcpserver/contract_test.go
- modified: internal/cli/cli_test.go

### Delivery targets
Nenhum novo; hardening transversal.

### API/contract changes
- Erros de autorização e redaction permanecem estáveis no MCP e CLI.

### Data/storage changes
- Adicionar política de retenção e remoção ao estado local.

### Technical risks
- Garantias podem ser superestimadas em documentação.
- Redaction pode falhar para formatos desconhecidos.

## 4. Tasks

### Planning
- [ ] Modelar ameaças e classificar riscos.
- [ ] Definir allowlists, exclusions, retenção e approvals.

### Implementation
- [ ] Centralizar políticas e redaction.
- [ ] Corrigir findings em todas as fronteiras.
- [ ] Implementar export e remoção segura.
- [ ] Criar corpus adversarial.
- [ ] Documentar garantias e limites.

### Validation
- [ ] Executar secret scan, govulncheck e suites negativas.
- [ ] Executar revisão independente de threat model e diff.
- [ ] Verificar zero mutação do workspace fora do comando explícito.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Execução local segura não equivale a sandbox completa.
- Options considered: prometer isolamento; exigir container; declarar modelo local limitado.
- Decision: controles fortes de processo e path com limites honestos.
- Rationale: atende V1 sem falsa garantia.
- Consequences: código hostil permanece fora do modelo de uso suportado.

## 6. Validation

### Strategy
Usar threat model, corpus adversarial, scanners e revisão separada.

### Deterministic checks
- Test: go test -race ./internal/security/... ./internal/workspace/... ./internal/checks/... ./internal/mcpserver/... ./internal/cli/...
- Lint: gofmt -l internal
- Typecheck: go vet ./...
- Build: go build ./cmd/ailearn
- Security / Contract: govulncheck ./...; secret scan; negative suite; pose assess integrate.

### Execution log
- Pendente.

### Results summary
- Hardening não executado.

### Requirement trace
- Mapear R1–R10 a threat cases, scanners, tests e revisão.

### Known gaps
- Isolamento forte poderá ser uma capacidade pós-V1.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados transversalmente e em docs/security.

### Validation executed
- Command: pose lint-spec security-privacy-hardening --ready-check
- Result: registrar após gate.

### Residual risks
- Riscos aceitos exigem follow-up ou ADR, nunca silêncio.

### Follow-ups
- [covered: v1-integrated-acceptance] Reexecutar suite adversarial no candidato V1.
