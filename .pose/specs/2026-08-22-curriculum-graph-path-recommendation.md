---
slug: curriculum-graph-path-recommendation
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: catalog-schema-loader, mastery-review-scheduling
priority: 120
components: curriculum, recommendations
delivers:
---

# Spec: curriculum-graph-path-recommendation

## 1. Intent

### Goal
Implementar o grafo curricular, pesquisa e recomendação explicável de trilhas e remediações.

### Business value
Escolher o que praticar com base em objetivo, dependências e evidências, em vez de uma lista linear fixa.

### Constraints
- Sugestões são consultivas e nunca iniciam trilha sem aceitação.
- Relações de precedência devem ser acíclicas.
- Recomendações não dependem de rede ou LLM.

### Non-goals
- Gerar desafios.
- Inferir conhecimento sem evidência.

## 2. Requirements

### Functional
- R1: Suportar requires, recommended_before, relates_to, contrasts_with, commonly_fails_with, applies_in, deepens_into e evidences.
- R2: Validar aciclicidade somente nas relações que implicam precedência.
- R3: Pesquisar itens por texto, tema, competência, dificuldade, tipo, duração e pré-requisito.
- R4: Recomendar caminho a partir de objetivo, tempo, perfil, progresso e revisões vencidas.
- R5: Explicar cada recomendação por evidência, gap, dependência e custo estimado.
- R6: Propor remediação menor e retorno ao desafio original.
- R7: Expor catalog_search, concept_relations_get e learning_path_recommend de forma completa.
- R8: Produzir resultado estável para mesma versão e estado.

### Non-functional
- Consultas p95 inferiores a 100 ms após indexação.
- Ranking deve ter tie-break determinístico.

### Security
- Query de texto é dado e deve ter limites de tamanho e custo.

### Compatibility
- Relações novas devem ser aditivas e desconhecidas falham de modo explícito.

## 3. Technical Plan

### Affected areas
- internal/curriculum/, internal/recommendation/, internal/mcpserver/

### Artifacts
- created: internal/curriculum/graph.go
- created: internal/curriculum/search.go
- created: internal/curriculum/graph_test.go
- created: internal/recommendation/service.go
- created: internal/recommendation/explanation.go
- created: internal/recommendation/service_test.go
- modified: internal/mcpserver/catalog_tools.go
- created: internal/mcpserver/recommendation_contract_test.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Completar filtros de busca e schemas de relações e recomendação.

### Data/storage changes
- Índices são projeções reconstruíveis do catálogo e do progresso.

### Technical risks
- Heurística pode favorecer pré-requisitos demais e atrasar prática.
- Relações editoriais inconsistentes degradam explicações.

## 4. Tasks

### Planning
- [ ] Definir semântica e direção de cada relação.
- [ ] Definir função de ranking, tie-break e limites.

### Implementation
- [ ] Implementar grafo validado e vizinhança limitada.
- [ ] Implementar índice de busca.
- [ ] Implementar recomendador e explicações.
- [ ] Implementar remediação e retorno.
- [ ] Expor tools e fixtures.

### Validation
- [ ] Testar ciclos, relações desconhecidas e ranking estável.
- [ ] Testar perfis sem evidência e objetivos conflitantes.
- [ ] Benchmarkar catálogo no tamanho-alvo da V1.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Recomendação deve ser previsível e auditável.
- Options considered: LLM; score oculto; regras determinísticas.
- Decision: ranking determinístico com decomposição de fatores.
- Rationale: permite entender e corrigir recomendações.
- Consequences: personalização sofisticada fica para versões futuras.

## 6. Validation

### Strategy
Usar grafos sintéticos, golden rankings e benchmarks no volume V1.

### Deterministic checks
- Test: go test ./internal/curriculum/... ./internal/recommendation/...
- Lint: gofmt -l internal/curriculum internal/recommendation
- Typecheck: go vet ./internal/curriculum/... ./internal/recommendation/...
- Build: go test ./internal/recommendation/... -run ^$
- Security / Contract: limites de query, profundidade e schemas.

### Execution log
- Pendente.

### Results summary
- Nenhum grafo ou recomendador entregue.

### Requirement trace
- Mapear R1–R8 a testes de grafo, ranking, remediação e contrato.

### Known gaps
- Qualidade depende da curadoria dos 160 conceitos e 100 competências.

## 7. Final Report

### Delivered scope
Nenhum; planejamento.

### Files and modules changed
- Planejados em curriculum, recommendation e mcpserver.

### Validation executed
- Command: pose lint-spec curriculum-graph-path-recommendation --ready-check
- Result: registrar após gate.

### Residual risks
- Ranking inicial deverá ser reavaliado após uso real.

### Follow-ups
- [covered: catalog-authoring-quality] Validar cobertura e coerência das relações.
