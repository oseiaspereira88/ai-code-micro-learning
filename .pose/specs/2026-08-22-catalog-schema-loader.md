---
slug: catalog-schema-loader
status: done
created_at: 2026-08-22
completed_at: 2026-08-22
supersedes:
depends_on: architecture-decision-baseline, go-runtime-foundation
priority: 30
components: curriculum
delivers:
changelog: none
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
- created: internal/curriculum/model.go
- created: internal/curriculum/loader.go
- created: internal/curriculum/validator.go
- created: internal/curriculum/index.go
- created: internal/curriculum/loader_test.go
- created: internal/curriculum/testdata/valid/manifest.yaml
- created: internal/curriculum/testdata/valid/pack.yaml
- created: internal/curriculum/testdata/duplicate_id/manifest.yaml
- created: internal/curriculum/testdata/duplicate_id/pack.yaml
- created: internal/curriculum/testdata/missing_reference/manifest.yaml
- created: internal/curriculum/testdata/missing_reference/pack.yaml
- created: internal/curriculum/testdata/prerequisite_cycle/manifest.yaml
- created: internal/curriculum/testdata/prerequisite_cycle/pack.yaml
- created: internal/curriculum/testdata/criteria_without_evidence/manifest.yaml
- created: internal/curriculum/testdata/criteria_without_evidence/pack.yaml
- created: internal/curriculum/testdata/incompatible_schema_version/manifest.yaml
- created: internal/curriculum/testdata/incompatible_schema_version/pack.yaml
- created: internal/curriculum/testdata/path_escape/manifest.yaml
- created: packs/manifest.yaml
- created: packs/go-first-steps.yaml
- modified: go.mod
- created: go.sum

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
- [x] Mapear todos os campos do exemplo das seções 14.8 e 14.9 de PROJECT.md.
- [x] Definir política de IDs, versões, remoção e depreciação.

### Implementation
- [x] Implementar modelos de autoria separados dos modelos de domínio.
- [x] Implementar loader ordenado, limites e diagnósticos.
- [x] Implementar validação de referências e ciclos.
- [x] Implementar índice e view com disclosure.
- [x] Criar fixtures positivas e negativas.

### Validation
- [x] Executar testes unitários e fuzz do loader (sem golden tests — ver Known gaps).
- [x] Validar um pack mínimo de referência.
- [x] Executar dependency scan (govulncheck).

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Conteúdo precisa ser legível por humanos e estrito para o runtime.
- Options considered: JSON; YAML livre; YAML validado e convertido.
- Decision: YAML de autoria validado e convertido para tipos internos imutáveis.
- Rationale: equilibrar ergonomia editorial e determinismo.
- Consequences: exige dependency YAML avaliada e schemas mantidos.

### Decision 2
- Date: 2026-08-22
- Context: A spec original listava `schemas/event.schema.json` como artefato,
  mas nenhum requisito ou seção do Technical Plan desta spec menciona
  eventos; o formato de evento pertence a `local-event-store` (§18.3 lista os
  eventos mínimos). Manter aqui duplicaria responsabilidade.
- Options considered: implementar um schema de evento de qualquer forma;
  remover o artefato desta spec.
- Decision: remover `schemas/event.schema.json` da seção 3; o schema de
  evento será definido por `local-event-store`, que já depende desta spec.
- Rationale: evita um contrato paralelo e mantém uma única fonte de verdade
  por domínio de dado.
- Consequences: `local-event-store` deve declarar seu próprio artefato de
  schema, se optar por JSON Schema para eventos.

### Decision 3
- Date: 2026-08-22
- Context: R1 exige "validar schemas", mas introduzir uma engine de
  validação JSON Schema em runtime seria a primeira dependência não-YAML do
  binário, contrariando `PROJECT.md` §24.1 ("biblioteca padrão... não
  antecipar frameworks") para uma necessidade que a validação nativa em Go
  já cobre.
- Options considered: validar via biblioteca de JSON Schema em runtime;
  publicar os `.schema.json` como contrato documentado/lint externo e
  validar estruturalmente em Go no loader.
- Decision: os `.schema.json` sob `schemas/` são o contrato de autoria
  (referência para autores e ferramentas externas); o loader valida
  estruturalmente em Go (tipos, presença, formato, referências, ciclos),
  sem depender de uma engine de JSON Schema em runtime.
- Rationale: entrega as mesmas garantias de R1/R3 sem dependência nova;
  reavaliar se autoria externa exigir validação client-side via JSON Schema.
- Consequences: manter os `.schema.json` e a validação Go sincronizados
  manualmente até existir um gate automatizado de consistência.

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
- `gofmt -l internal/curriculum` → saída vazia (2026-08-22).
- `go vet ./internal/curriculum/...` → sem findings (2026-08-22).
- `go test ./internal/curriculum/...` → `ok`, 11 testes + 1 fuzz target (2026-08-22).
- `go test ./internal/curriculum/ -fuzz=FuzzMalformedPackNeverPanics -fuzztime=15s` →
  430.474 execuções, 0 crashes, `PASS` (2026-08-22).
- `govulncheck ./...` → `No vulnerabilities found.` (2026-08-22).
- Carga manual de `packs/` (pack de referência real, não fixture de teste):
  0 diagnósticos, desafio consultável por ID, `Variants` oculto na view
  pública (2026-08-22).

### Results summary
- `internal/curriculum` carrega `packs/manifest.yaml` + packs YAML dentro de
  limites de tamanho/profundidade/alias, valida schema_version, IDs
  duplicados, referências ausentes, ciclos de pré-requisito e critério sem
  evidência, e constrói um `Catalog` imutável consultável por ID/tipo/tema
  com campos reservados (`Variants`) ocultos por padrão.
- Um pack mínimo de referência (`packs/go-first-steps.yaml`) carrega sem
  diagnósticos, cobrindo os exemplos de PROJECT.md §14.8–§14.9.
- Contratos `.schema.json` documentam a forma de autoria (Decisão 3); a
  aplicação real das regras é nativa em Go, não via engine de JSON Schema.

### Requirement trace
- R1 [satisfied] report:schemas/pack.schema.json report:schemas/challenge.schema.json report:internal/curriculum/model.go
- R2 [satisfied] test:TestLoadIsDeterministicRegardlessOfManifestOrder test:TestLoadRejectsDuplicateID
- R3 [satisfied] test:TestLoadRejectsMissingReference test:TestLoadDetectsPrerequisiteCycle test:TestLoadRejectsCriteriaWithoutEvidence
- R4 [satisfied] report:internal/curriculum/validator.go
- R5 [satisfied] test:TestLoadValidPackBuildsQueryableCatalog
- R6 [satisfied] test:TestChallengeHidesReservedVariants test:TestLoadValidPackBuildsQueryableCatalog
- R7 [satisfied] test:TestLoadRejectsIncompatibleSchemaVersion

### Known gaps
- Gates editoriais semânticos pertencem a catalog-authoring-quality.
- Não há golden tests (comparação byte-a-byte de diagnósticos contra
  arquivos de referência); os testes de tabela verificam apenas o código do
  diagnóstico (`DiagnosticCode`), não o texto completo. Suficiente para R4
  (código estável); golden tests ficam como follow-up se o formato de
  `Detail` precisar de garantia de compatibilidade textual.
- `schemas/*.schema.json` não são validados automaticamente contra o loader
  Go (Decisão 3); dessincronia é possível até existir um gate de
  consistência.

## 7. Final Report

### Delivered scope
Loader, validador e índice de catálogo em `internal/curriculum`, contratos
`.schema.json` de autoria, manifesto de packs e um pack mínimo de
referência. Nenhuma recomendação de trilha, domínio ou pack completo
implementado, conforme os non-goals da spec.

### Files and modules changed
- `internal/curriculum/model.go` (criado)
- `internal/curriculum/loader.go` (criado)
- `internal/curriculum/validator.go` (criado)
- `internal/curriculum/index.go` (criado)
- `internal/curriculum/loader_test.go` (criado)
- `internal/curriculum/testdata/{valid,duplicate_id,missing_reference,prerequisite_cycle,criteria_without_evidence,incompatible_schema_version,path_escape}/*.yaml` (criado, 13 arquivos em 7 cenários de fixture)
- `schemas/catalog.schema.json` (criado)
- `schemas/pack.schema.json` (criado)
- `schemas/challenge.schema.json` (criado)
- `packs/manifest.yaml` (criado)
- `packs/go-first-steps.yaml` (criado)
- `go.mod` (modificado) e `go.sum` (criado): primeira dependência externa, `gopkg.in/yaml.v3`
- `.pose/specs/2026-08-22-catalog-schema-loader.md` (atualizado)

### Validation executed
- Command: go test ./internal/curriculum/... -race && govulncheck ./...
- Result: `ok  	github.com/oseiaspereira88/ailearn/internal/curriculum`; `No vulnerabilities found.`

### Residual risks
- Migração entre schemas será testada quando existir uma segunda versão.

### Follow-ups
- [covered: catalog-authoring-quality] Adicionar regras editoriais além do schema.
