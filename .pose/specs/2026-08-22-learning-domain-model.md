---
slug: learning-domain-model
status: done
created_at: 2026-08-22
completed_at: 2026-08-22
supersedes:
depends_on: architecture-decision-baseline, go-runtime-foundation
priority: 20
components: learning-domain
delivers:
changelog: none
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
- [x] Mapear cada entidade e invariante às seções 7, 8, 12 e 13 de PROJECT.md.
- [x] Identificar tipos que realmente exigem identidade.

### Implementation
- [x] Implementar IDs, enums e value objects validados.
- [x] Implementar agregados e transições.
- [x] Implementar erros de domínio.
- [x] Cobrir tabelas completas de transição e casos negativos.

### Validation
- [x] Executar testes unitários e casos negativos de transições inválidas.
- [x] Executar race detector (mutabilidade compartilhada em `LearningSession`/`StepProgress`).
- [x] Revisar dependências para garantir domínio puro.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Domínio rico pode induzir hierarquias de structs rígidas.
- Options considered: herança conceitual; union por kind; tipos separados.
- Decision: usar composição e tipos discriminados somente onde a invariante exigir.
- Rationale: preservar clareza idiomática em Go.
- Consequences: conversões explícitas nos adaptadores.

### Decision 2
- Date: 2026-08-22
- Context: `Catalog` (§12.1) é listado como agregado deste domínio, mas indexar
  e validar packs é responsabilidade de `catalog-schema-loader` (non-goal
  explícito: "carregar catálogo").
- Options considered: implementar indexação completa aqui; representar
  apenas referências versionadas (`CatalogID`/`ContentVersion`) consumidas
  por desafio e sessão.
- Decision: representar somente tipos e referências versionadas; loader e
  índice ficam para `catalog-schema-loader`.
- Rationale: evita duplicar responsabilidade e mantém o domínio puro sem
  acesso a filesystem/YAML.
- Consequences: `catalog-schema-loader` deve produzir/validar instâncias
  destes mesmos tipos, não um formato paralelo.

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
- `gofmt -l internal/learning` → saída vazia (2026-08-22).
- `go vet ./internal/learning/...` → sem findings (2026-08-22).
- `go test ./internal/learning/...` → `ok`, 15 funções de teste (2026-08-22).
- `go test ./internal/learning/... -race` → sem data races (2026-08-22).
- `go list -deps ./internal/learning/...` → produção importa somente `fmt` e
  `time`; nenhuma dependência de MCP, CLI, YAML ou filesystem (2026-08-22).

### Results summary
- Domínio puro implementado sob `internal/learning`: ontologia (catálogo,
  desafio, camada, passo), agregados `LearningSession`/`StepProgress`/
  `LearningDetour` com tabelas de transição completas de §13.1–§13.3,
  políticas de sessão (§8.3) com construtores que falham na criação para
  valores inválidos, evidências imutáveis, feedback consultivo com
  bloqueio por tabela (§8.7) e avaliação sem efeito de conclusão/avanço.
- 15 funções de teste (com subtestes de tabela) cobrem as tabelas de
  transição completas, casos negativos e determinismo; nenhum teste de
  fuzzing (`go test -fuzz`) foi adicionado — ver Known gaps.

### Requirement trace
- R1 [satisfied] report:internal/learning/catalog.go report:internal/learning/challenge.go
- R2 [satisfied] report:internal/learning/evidence.go report:internal/learning/evaluation.go
- R3 [satisfied] test:TestValidSessionTransitionTable test:TestValidStepTransitionTableAndOverrideSkip
- R4 [satisfied] test:TestSetActiveInstructionRejectsSecondActive test:TestAdvanceRequiresCompletedStepUnlessOverride
- R5 [satisfied] test:TestRevealSolutionInvalidatesAutonomyEvidence
- R6 [satisfied] test:TestNewSessionPolicyRejectsInvalidValues test:TestNewSessionPolicyAcceptsValidValues
- R7 [satisfied] report:internal/learning/errors.go

### Known gaps
- Formato persistido será fechado junto de local-event-store.
- A spec citava "fuzzing de transições inválidas"; implementado como testes
  de tabela negativos (determinísticos), não `go test -fuzz` real. Fuzzing
  de fato fica como follow-up caso as tabelas de transição cresçam.

## 7. Final Report

### Delivered scope
Modelo pedagógico puro sob `internal/learning`: ontologia, agregados de
sessão e passo com tabelas de transição completas, políticas de sessão
validadas na construção, evidências imutáveis e feedback/avaliação
independentes de conclusão/avanço. Nenhuma persistência, carregamento de
catálogo ou observação de workspace, conforme os non-goals da spec.

### Files and modules changed
- `internal/learning/catalog.go` (criado)
- `internal/learning/challenge.go` (criado)
- `internal/learning/step.go` (criado)
- `internal/learning/session.go` (criado)
- `internal/learning/evidence.go` (criado)
- `internal/learning/evaluation.go` (criado)
- `internal/learning/policy.go` (criado)
- `internal/learning/errors.go` (criado)
- `internal/learning/session_test.go` (criado)
- `internal/learning/policy_test.go` (criado)
- `.pose/specs/2026-08-22-learning-domain-model.md` (atualizado)

### Validation executed
- Command: go test ./internal/learning/... -race
- Result: `ok  	github.com/oseiaspereira88/codinho/internal/learning`

### Residual risks
- Alterações posteriores no protocolo deverão adaptar, não contaminar, o domínio.

### Follow-ups
- [covered: v1-integrated-acceptance] Verificar todas as invariantes em fluxos compostos.
