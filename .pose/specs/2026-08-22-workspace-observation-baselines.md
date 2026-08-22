---
slug: workspace-observation-baselines
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: go-runtime-foundation, local-event-store
priority: 90
components: workspace, evidence
delivers:
---

# Spec: workspace-observation-baselines

## 1. Intent

### Goal
Observar mudanças relevantes no workspace com baseline, contenção de paths e evidências imutáveis, sem editar o código do aluno.

### Business value
Permitir feedback e avaliação baseados no que o aluno realmente produziu, preservando trabalho preexistente.

### Constraints
- Leitura apenas; nenhuma API de escrita.
- Suportar repositórios Git e diretórios sem Git.
- Distinguir fixture, baseline e trabalho da etapa.

### Non-goals
- Executar testes, formatar arquivos ou corrigir código.
- Garantir segredo forte de testes em máquina controlada pelo usuário.

## 2. Requirements

### Functional
- R1: Autorizar uma raiz real e rejeitar paths ou symlinks que escapem dela.
- R2: Registrar commit, status, diff preexistente, hashes, manifesto de fixture e revisão ao ativar passo.
- R3: Coletar somente arquivos e símbolos dentro dos globs declarados pelo desafio.
- R4: Produzir diff relevante sem atribuir mudanças anteriores à etapa.
- R5: Comparar hashes e metadados quando Git não estiver disponível.
- R6: Analisar Go por parser e AST quando o critério for estrutural.
- R7: Criar Observation e Evidence sem implicar tentativa ou avaliação.
- R8: Expor workspace_observe e evidence_get com redaction, tamanho e escopo.
- R9: Detectar evidência obsoleta quando o workspace mudar após a coleta.

### Non-functional
- Observação não pode modificar mtime, índice Git ou working tree.
- Outputs devem ser limitados e content-addressed.

### Security
- Excluir secrets, .env, credenciais e paths sensíveis por padrão.
- Tratar nomes, comentários e conteúdo como dados não confiáveis.

### Compatibility
- Implementar abstração de Git e fallback portable.

## 3. Technical Plan

### Affected areas
- internal/workspace/, internal/evidence/, internal/mcpserver/

### Artifacts
- created: internal/workspace/root.go
- created: internal/workspace/baseline.go
- created: internal/workspace/git.go
- created: internal/workspace/hashes.go
- created: internal/workspace/golang.go
- created: internal/workspace/redaction.go
- created: internal/workspace/workspace_test.go
- created: internal/workspace/security_test.go
- modified: internal/evidence/store.go
- created: internal/mcpserver/workspace_tools.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar workspace_observe e evidence_get.

### Data/storage changes
- Persistir baseline, observation e referências content-addressed.

### Technical risks
- Symlink races e paths com diferenças entre plataformas.
- Diffs grandes podem consumir memória ou vazar dados.

## 4. Tasks

### Planning
- [ ] Definir modelo de autorização e exclusões padrão.
- [ ] Definir limite de arquivo, diff e AST.

### Implementation
- [ ] Implementar resolução real e contenção.
- [ ] Implementar baseline Git e fallback por hash.
- [ ] Implementar diff, redaction e evidência.
- [ ] Implementar inspeção AST limitada.
- [ ] Expor tools somente leitura.

### Validation
- [ ] Testar repo sujo, sem Git, rename, symlink e traversal.
- [ ] Verificar índice e arquivos byte a byte antes e depois.
- [ ] Testar prompt injection em comentários como dado inerte.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Regex é frágil para critérios estruturais Go.
- Options considered: regex; go/parser e go/ast; language server.
- Decision: usar AST padrão para estrutura e texto somente quando semântico.
- Rationale: determinismo sem daemon externo.
- Consequences: critérios devem aceitar soluções equivalentes.

## 6. Validation

### Strategy
Combinar fixtures Git/não Git, ataques de path e comparação de filesystem.

### Deterministic checks
- Test: go test ./internal/workspace/... ./internal/evidence/...
- Lint: gofmt -l internal/workspace internal/evidence
- Typecheck: go vet ./internal/workspace/... ./internal/evidence/...
- Build: go test ./internal/workspace/... -run ^$
- Security / Contract: traversal, symlink race, redaction e secret scan.

### Execution log
- Pendente.

### Results summary
- Nenhuma observação entregue.

### Requirement trace
- Mapear R1–R9 a fixtures de contenção, baseline, AST e obsolescência.

### Known gaps
- Limites de recurso do subprocesso pertencem a safe-check-executor.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em workspace, evidence e mcpserver.

### Validation executed
- Command: pose lint-spec workspace-observation-baselines --ready-check
- Result: registrar após validação.

### Residual risks
- Portabilidade de symlink requer CI por sistema operacional.

### Follow-ups
- [covered: security-privacy-hardening] Realizar revisão adversarial completa da fronteira de workspace.
