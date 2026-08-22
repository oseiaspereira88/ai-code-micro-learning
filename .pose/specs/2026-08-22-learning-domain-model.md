---
slug: learning-domain-model
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: architecture-decision-baseline, go-runtime-foundation
priority: 20
components: learning-domain
delivers:
---

# Spec: learning-domain-model

## 1. Intent

### Goal
Implementar o modelo pedagógico puro, sua ontologia e as invariantes de sessão e progressão descritas em PROJECT.md.

### Business value
Dar a todas as superfícies uma semântica única para desafio, passos, evidências e transições.

### Constraints
- O domínio não pode importar MCP, CLI, YAML ou filesystem.
- Feedback, avaliação, conclusão e avanço permanecem independentes.
- Uma sessão possui no máximo uma instrução ativa.

### Non-goals
- Persistir estado, carregar catálogo ou observar workspace.
- Calcular domínio longitudinal.

## 2. Requirements

### Functional
- R1: Representar trilha, tema, desafio, camada, macro, meso e micropasso com IDs estáveis.
- R2: Representar conceito, competência, tentativa, observação, pista, feedback, avaliação, evidência, reflexão e desvio.
- R3: Aplicar os estados e transições válidos de sessão, nó e desvio pedagógico.
- R4: Rejeitar mais de uma instrução ativa e avanço de nó não concluído sem override.
- R5: Marcar solução revelada sem promover autonomia.
- R6: Distinguir políticas de modo, profundidade, auxílio, avaliação, avanço e tempo.
- R7: Emitir erros de domínio tipados e comparáveis sem mensagens acopladas à UI.

### Non-functional
- Operações puras devem ser determinísticas para a mesma entrada.
- Value objects inválidos devem falhar na construção.

### Security
- Texto vindo de catálogo ou aluno deve permanecer dado, nunca instrução executável.

### Compatibility
- Tipos persistidos devem possuir versões ou adaptadores de migração posteriores.

## 3. Technical Plan

### Affected areas
- internal/learning/

### Artifacts
- created: internal/learning/catalog.go
- created: internal/learning/challenge.go
- created: internal/learning/step.go
- created: internal/learning/session.go
- created: internal/learning/evidence.go
- created: internal/learning/evaluation.go
- created: internal/learning/policy.go
- created: internal/learning/errors.go
- created: internal/learning/session_test.go
- created: internal/learning/policy_test.go

### Delivery targets
Nenhum; trata-se do núcleo interno.

### API/contract changes
- Criar API interna consumida pelos casos de uso, sem interface por entidade.

### Data/storage changes
- Definir representações serializáveis sem implementar armazenamento.

### Technical risks
- Modelagem excessiva pode tornar micropassos difíceis de autorar.
- Estados permissivos podem acoplar novamente operações pedagógicas.

## 4. Tasks

### Planning
- [ ] Mapear cada entidade e invariante às seções 7, 8, 12 e 13 de PROJECT.md.
- [ ] Identificar tipos que realmente exigem identidade.

### Implementation
- [ ] Implementar IDs, enums e value objects validados.
- [ ] Implementar agregados e transições.
- [ ] Implementar erros de domínio.
- [ ] Cobrir tabelas completas de transição e casos negativos.

### Validation
- [ ] Executar testes unitários e fuzzing de transições inválidas.
- [ ] Executar race detector se houver mutabilidade compartilhada.
- [ ] Revisar dependências para garantir domínio puro.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Domínio rico pode induzir hierarquias de structs rígidas.
- Options considered: herança conceitual; union por kind; tipos separados.
- Decision: usar composição e tipos discriminados somente onde a invariante exigir.
- Rationale: preservar clareza idiomática em Go.
- Consequences: conversões explícitas nos adaptadores.

## 6. Validation

### Strategy
Testar invariantes por tabelas, transições negativas e propriedades de determinismo.

### Deterministic checks
- Test: go test ./internal/learning/...
- Lint: gofmt -l internal/learning
- Typecheck: go vet ./internal/learning/...
- Build: go test ./internal/learning/... -run ^$
- Security / Contract: inspeção de imports e fuzz de valores inválidos.

### Execution log
- Pendente; spec em estado draft.

### Results summary
- Nenhum comportamento de domínio entregue ainda.

### Requirement trace
- Mapear R1–R7 a testes nomeados e commit atribuído no closeout.

### Known gaps
- Formato persistido será fechado junto de local-event-store.

## 7. Final Report

### Delivered scope
Nenhum; somente planejamento.

### Files and modules changed
- Planejados sob internal/learning/.

### Validation executed
- Command: pose lint-spec learning-domain-model --ready-check
- Result: registrar após validação.

### Residual risks
- Alterações posteriores no protocolo deverão adaptar, não contaminar, o domínio.

### Follow-ups
- [covered: v1-integrated-acceptance] Verificar todas as invariantes em fluxos compostos.
