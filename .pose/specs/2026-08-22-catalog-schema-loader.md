---
slug: catalog-schema-loader
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: architecture-decision-baseline, go-runtime-foundation
priority: 30
components: curriculum
delivers:
---

# Spec: catalog-schema-loader

## 1. Intent

### Goal
Definir schemas versionados e carregar packs YAML em um catálogo imutável, validado e consultável.

### Business value
Permitir autoria de conteúdo rica sem incorporar lógica pedagógica ou dados curriculares no binário.

### Constraints
- Sessões devem fixar versões de packs.
- IDs publicados são estáveis e não podem ser reutilizados com outro significado.
- Conteúdo reservado não pode aparecer em consultas públicas.

### Non-goals
- Implementar recomendação de trilha, domínio ou os packs completos.
- Gerar conteúdo com LLM.

## 2. Requirements

### Functional
- R1: Validar schemas de pack, tema, conceito, competência, desafio, step, pista, check e trilha.
- R2: Carregar múltiplos packs com versão semântica e IDs globalmente não ambíguos.
- R3: Rejeitar referências ausentes, ciclos de pré-requisito, tipos inválidos e critérios sem evidência.
- R4: Produzir diagnóstico com arquivo, item, campo e código estável.
- R5: Fixar uma visão imutável do catálogo para cada sessão.
- R6: Consultar itens por ID, tipo e filtros básicos sem revelar campos reservados.
- R7: Detectar versões incompatíveis de schema antes de materializar o catálogo.

### Non-functional
- A validação do catálogo completo da V1 deve ser determinística e offline.
- A ordem de arquivos no filesystem não pode alterar o resultado.

### Security
- Limitar tamanho, profundidade e aliases YAML.
- Confinar includes e caminhos de fixture dentro do pack.

### Compatibility
- Preservar leitura de versões anteriores por migração explícita, nunca por heurística silenciosa.

## 3. Technical Plan

### Affected areas
- schemas/, internal/curriculum/, packs/

### Artifacts
- created: schemas/catalog.schema.json
- created: schemas/pack.schema.json
- created: schemas/challenge.schema.json
- created: schemas/event.schema.json
- created: internal/curriculum/model.go
- created: internal/curriculum/loader.go
- created: internal/curriculum/validator.go
- created: internal/curriculum/index.go
- created: internal/curriculum/loader_test.go
- created: internal/curriculum/testdata/
- created: packs/manifest.yaml

### Delivery targets
Nenhum; o catálogo é capacidade interna consumida por MCP e CLI.

### API/contract changes
- Definir o contrato de autoria schema v1 e a API de snapshot imutável.

### Data/storage changes
- Introduzir YAML versionado sob packs e JSON Schema sob schemas.

### Technical risks
- YAML permissivo pode aceitar conteúdo surpreendente.
- Validação estrutural não garante qualidade pedagógica.

## 4. Tasks

### Planning
- [ ] Mapear todos os campos do exemplo das seções 14.8 e 14.9 de PROJECT.md.
- [ ] Definir política de IDs, versões, remoção e depreciação.

### Implementation
- [ ] Implementar modelos de autoria separados dos modelos de domínio.
- [ ] Implementar loader ordenado, limites e diagnósticos.
- [ ] Implementar validação de referências e ciclos.
- [ ] Implementar índice e view com disclosure.
- [ ] Criar fixtures positivas e negativas.

### Validation
- [ ] Executar testes unitários, golden tests e fuzz do loader.
- [ ] Validar um pack mínimo de referência.
- [ ] Executar dependency e secret scan.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Conteúdo precisa ser legível por humanos e estrito para o runtime.
- Options considered: JSON; YAML livre; YAML validado e convertido.
- Decision: YAML de autoria validado e convertido para tipos internos imutáveis.
- Rationale: equilibrar ergonomia editorial e determinismo.
- Consequences: exige dependency YAML avaliada e schemas mantidos.

## 6. Validation

### Strategy
Usar fixtures, golden diagnostics, testes de ciclo e fuzzing de entradas malformadas.

### Deterministic checks
- Test: go test ./internal/curriculum/...
- Lint: gofmt -l internal/curriculum
- Typecheck: go vet ./internal/curriculum/...
- Build: go test ./internal/curriculum/... -run ^$
- Security / Contract: fuzz loader, limites YAML e validação JSON Schema.

### Execution log
- Pendente; implementação não iniciada.

### Results summary
- Contrato planejado, sem evidência de runtime.

### Requirement trace
- Associar R1–R7 a fixtures nomeadas e ao relatório do validador.

### Known gaps
- Gates editoriais semânticos pertencem a catalog-authoring-quality.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em schemas, internal/curriculum e packs/manifest.yaml.

### Validation executed
- Command: pose lint-spec catalog-schema-loader --ready-check
- Result: registrar após o gate.

### Residual risks
- Migração entre schemas será testada quando existir uma segunda versão.

### Follow-ups
- [covered: catalog-authoring-quality] Adicionar regras editoriais além do schema.
