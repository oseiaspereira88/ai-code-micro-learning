---
slug: session-orchestration-disclosure
status: done
created_at: 2026-08-22
completed_at: 2026-08-23
supersedes:
depends_on: learning-domain-model, catalog-schema-loader, local-event-store, mcp-stdio-foundation
priority: 60
components: sessions, mcp-server
delivers:
---

# Spec: session-orchestration-disclosure

## 1. Intent

### Goal
Orquestrar lifecycle de sessão, instrução única, granularidade e divulgação progressiva com estado persistente.

### Business value
Impedir que o agente revele etapas ou soluções fora do momento pedagógico autorizado.

### Constraints
- Sessão fixa versões de conteúdo no início.
- Avanço é sempre explícito na V1.
- Ajuste de granularidade não modifica a árvore canônica.

### Non-goals
- Gerar feedback semântico, executar checks ou calcular domínio.
- Implementar todos os modos especializados.

## 2. Requirements

### Functional
- R1: Iniciar sessão com desafio ou trilha, modo, profundidade, políticas, workspace opcional e tempo opcional.
- R2: Consultar e retomar sessão com nó ativo, revisão, disclosure, ações permitidas e pendências.
- R3: Manter exatamente uma instrução ativa e omitir filhos futuros, gabarito e pistas bloqueadas.
- R4: Configurar apenas propriedades mutáveis e registrar toda mudança.
- R5: Pausar, retomar, concluir e abandonar lifecycle sem inferir conclusão dos passos.
- R6: Ajustar janela entre desafio, camada, macro, meso e micro sem reescrever conteúdo.
- R7: Detectar desatualização de revisão e retries idempotentes.
- R8: Expor session_start, session_get, session_configure, session_pause, session_resume, session_finish, instruction_get e granularity_adjust.

### Non-functional
- Mesma versão e estado produzem a mesma instrução e ações.
- Consulta de sessão deve atender p95 inferior a 100 ms após aquecimento.

### Security
- Disclosure deve ser aplicado no servidor, não apenas na skill.
- IDs de sessão não autorizam outro workspace.

### Compatibility
- Sessões antigas permanecem legíveis após atualização de packs.

## 3. Technical Plan

### Affected areas
- internal/session/, internal/application/, internal/mcpserver/

### Artifacts
- created: internal/session/service.go
- created: internal/session/disclosure.go
- created: internal/session/granularity.go
- created: internal/session/service_test.go
- created: internal/session/disclosure_test.go
- modified: internal/application/session.go
- modified: internal/mcpserver/session_tools.go
- created: internal/mcpserver/session_contract_test.go
- modified: internal/mcpserver/contract_test.go
- modified: internal/eventstore/event.go
- modified: internal/mcpserver/errors.go
- modified: internal/mcpserver/server_test.go
- modified: internal/application/catalog.go

### Delivery targets
Nenhum novo; amplia o contrato MCP V1 planejado.

### API/contract changes
- Adicionar oito tools de sessão e seus efeitos ao contrato MCP.

### Data/storage changes
- Persistir eventos de lifecycle, instrução e granularidade.

### Technical risks
- Disclosure incorreto pode vazar solução.
- Reagrupar passos pode perder posição ou evidência.

## 4. Tasks

### Planning
- [x] Definir matriz modo versus disclosure versus ação permitida.
- [x] Definir seleção determinística do próximo nó.

### Implementation
- [x] Implementar lifecycle e projeção da sessão.
- [x] Implementar disclosure e instrução única.
- [x] Implementar ajuste de granularidade.
- [x] Expor tools e eventos.
- [x] Cobrir branches, retries e retomada.

### Validation
- [x] Executar testes de tabela para toda transição (sessão, lifecycle, idempotência, conflito de revisão).
- [x] Cobrir os dois níveis de disclosure suportados (briefing/instruction) com testes unitários (sem golden tests em arquivo — ver Known gaps).
- [ ] Reiniciar servidor no meio de uma sessão e comparar estado — não realizado: sessões continuam em memória apenas (ver Known gaps, mcp-stdio-foundation Decisão 2/3).

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Granularidade é uma visão da árvore, não conteúdo duplicado.
- Options considered: duplicar desafios por nível; podar árvore; janela derivada.
- Decision: derivar uma janela instrucional da árvore canônica.
- Rationale: evita drift e permite adaptação.
- Consequences: o cursor precisa preservar ancestry e progresso.

### Decision 2
- Date: 2026-08-22
- Context: R1 pede iniciar sessão "com desafio ou trilha", mas nenhuma spec
  entregue até aqui resolve uma trilha (`TrackAuthoring`) em uma sequência
  de desafios — essa é responsabilidade de
  `curriculum-graph-path-recommendation` (ainda não implementada).
- Options considered: implementar resolução mínima de trilha aqui;
  manter `session_start` aceitando somente `challenge_id`.
- Decision: `session_start` aceita somente `challenge_id` nesta spec;
  iniciar por trilha fica para quando `curriculum-graph-path-recommendation`
  existir.
- Rationale: evita duplicar/antecipar lógica de recomendação de trilha fora
  de sua spec dona.
- Consequences: `curriculum-graph-path-recommendation` deve estender
  `session_start` (ou compor com ele) para aceitar `track_id`.

### Decision 3
- Date: 2026-08-22
- Context: `mcp-stdio-foundation` Decisão 3 adiou `expected_revision`
  explícito para "quando `session_configure` ou tool equivalente for
  adicionada" — esta spec adiciona cinco tools mutáveis
  (`session_configure`, `session_pause`, `session_resume`, `session_finish`,
  `granularity_adjust`).
- Decision: todas as cinco exigem `expected_revision` do cliente (obtido via
  `session_get`) e aceitam `request_id` opcional para idempotência,
  seguindo o mesmo mecanismo de `eventstore.Store.Append` já testado em
  `local-event-store`. Uma chamada repetida com o mesmo `request_id` nunca
  reaplica a transição de domínio subjacente (checado via
  `revisão retornada == expected_revision + 1`), só retorna o resultado já
  efetivado.
- Rationale: fecha o gap de R7/R8 deixado explicitamente em aberto pela
  spec anterior, sem inventar semântica nova.
- Consequences: nenhuma; fecha a Decisão 3 de `mcp-stdio-foundation`.

### Decision 4
- Date: 2026-08-22
- Context: R5 exige `session_resume`, mas a lista de eventos mínimos de
  `PROJECT.md` §18.3 não inclui um marcador de retomada (só
  `session_paused`); inferir retomada de outro evento perderia
  auditabilidade.
- Options considered: reaproveitar `session_started` para retomada (ambíguo
  no replay); adicionar `EventSessionResumed` a `internal/eventstore`.
- Decision: adicionar `EventSessionResumed = "session_resumed"` a
  `internal/eventstore/event.go`.
- Rationale: §18.3 é explicitamente uma lista de eventos "mínimos", não
  exaustiva; retomada precisa de um marcador próprio para não perder
  auditabilidade no replay.
- Consequences: nenhuma migração necessária — evento novo, não uma mudança
  de formato de evento existente.

## 6. Validation

### Strategy
Cobrir state machine, disclosure negativo, replay e contratos MCP.

### Deterministic checks
- Test: go test ./internal/session/... ./internal/mcpserver/...
- Lint: gofmt -l internal/session internal/mcpserver
- Typecheck: go vet ./internal/session/... ./internal/mcpserver/...
- Build: go build ./cmd/codinho
- Security / Contract: golden leak tests e pose assess integrate.

### Execution log
- `gofmt -l internal/session internal/mcpserver internal/application internal/eventstore` → saída vazia (2026-08-22).
- `go vet ./...` → sem findings (2026-08-22).
- `go test ./... -race` → `ok` em todos os pacotes; 22 funções de teste
  novas (17 em `internal/session`, 5 em
  `internal/mcpserver/session_contract_test.go`), 95 no total do
  repositório (2026-08-22).
- `govulncheck ./...` → `No vulnerabilities found.` (2026-08-22).
- Smoke real por subprocesso (`mcp.CommandTransport`) contra o pack de
  referência real (`packs/go-first-steps.yaml`): `session_start` →
  `granularity_adjust` → `instruction_get` → `session_configure` →
  `session_pause` → `session_resume` → `session_finish`, revisão
  incrementando monotonicamente (1→6) a cada chamada, todas `status: ok`
  (2026-08-22).
- Bug real de idempotência encontrado e corrigido durante a implementação:
  a heurística inicial de "chamada nova vs. replay" (`ev.Revision ==
  expectedRevision+1`) não distinguia um append fresco de um replay
  idempotente, porque um retry por definição repete o mesmo
  `expected_revision` — ambos os casos satisfaziam a equação. Corrigido
  comparando a revisão do stream *antes* da chamada com `expectedRevision`.

### Results summary
- `internal/session.Service` orquestra lifecycle completo (`Start`, `Get`,
  `Instruction`, `Configure`, `Pause`, `Resume`, `Finish`,
  `GranularityAdjust`) sobre `internal/learning`, `internal/curriculum` e
  `internal/eventstore`, com idempotência por `request_id` e rejeição de
  `expected_revision` obsoleta em toda operação mutável.
- `internal/session/granularity.go` deriva a janela instrucional (desafio,
  camada, macro, meso, micro) andando pela árvore autorada existente
  (`Children` recursivos), sem reescrever conteúdo; degrada com segurança
  quando o conteúdo não desce até o nível pedido.
- `internal/application/session.go` virou um adaptador fino (type aliases +
  delegação) para `internal/session.Service`, preservando a superfície que
  `internal/mcpserver` já consumia.
- `internal/mcpserver` ganhou `session_configure`, `session_pause`,
  `session_resume`, `session_finish` e `granularity_adjust`, totalizando 10
  tools; `session_start`/`session_get`/`instruction_get` passaram a expor
  `revision` e `disclosure` reais (antes, `instruction_get` sempre
  retornava `solution_revealed: false` hardcoded).
- `EventSessionResumed` adicionado a `internal/eventstore` (Decisão 4):
  §18.3 lista eventos mínimos, não exaustivos, e retomada precisa de
  marcador próprio para auditabilidade no replay.

### Requirement trace
- R1 [satisfied] test:TestStartAssignsFirstMacroStep test:TestContractSessionLifecycleEndToEnd
- R2 [satisfied] test:TestGetReportsRevisionAndDisclosure
- R3 [satisfied] test:TestInstructionNeverExposesChildrenOrSolution
- R4 [satisfied] test:TestConfigureAppliesOnlySpecifiedFields test:TestContractSessionConfigureChangesHelpPolicy
- R5 [satisfied] test:TestPauseResumeFinishLifecycle test:TestFinishFromPausedFailsWithoutInferringResume test:TestContractSessionPauseResumeFinishLifecycle
- R6 [satisfied] test:TestGranularityAdjustWalksToRequestedDepth test:TestContractGranularityAdjustWalksToMicro
- R7 [satisfied] test:TestConfigureRejectsStaleRevision test:TestContractSessionConfigureRejectsStaleRevision
- R8 [satisfied] report:internal/mcpserver/session_tools.go

### Known gaps
- Adaptação automática baseada em domínio será integrada depois.
- Sem golden tests em arquivo para os níveis de disclosure — cobertos por
  testes unitários (`disclosureFor`). Só dois níveis são alcançáveis nesta
  spec (`briefing`, `instruction`); níveis acima (pseudocódigo, fragmento,
  solução) chegam via `hint_request`, de `assistance-hints-detours`.
- Reiniciar o servidor no meio de uma sessão não reconstrói o estado a
  partir do log: sessões continuam em memória apenas (limitação já
  documentada em `mcp-stdio-foundation` Decisão 2/3, ainda não fechada).
  `internal/eventstore` já suporta replay completo; falta o adaptador que
  reconstrói `*learning.LearningSession`/`*learning.StepProgress` a partir
  dos eventos.
- `session_start` só aceita `challenge_id`, não `track_id` (Decisão 2) —
  fica para `curriculum-graph-path-recommendation`.
- Não há teste automatizado de latência (p95 < 100ms); não verificado.

## 7. Final Report

### Delivered scope
Orquestração completa de sessão (lifecycle, instrução única, disclosure,
granularidade) sobre o contrato MCP, com 5 novas tools somando 10 ao
servidor. Nenhuma adaptação automática por domínio, reconstrução de sessão
a partir do log, ou início por trilha, conforme os non-goals/Decisões desta
spec.

### Files and modules changed
- `internal/session/service.go` (criado)
- `internal/session/disclosure.go` (criado)
- `internal/session/granularity.go` (criado)
- `internal/session/service_test.go` (criado)
- `internal/session/disclosure_test.go` (criado)
- `internal/application/session.go` (reescrito como adaptador fino)
- `internal/application/catalog.go` (modificado: `ErrNotFound` alinhado a `session.ErrNotFound`)
- `internal/mcpserver/session_tools.go` (modificado: +5 tools)
- `internal/mcpserver/session_contract_test.go` (criado)
- `internal/mcpserver/contract_test.go` (modificado: lista de tools esperada)
- `internal/mcpserver/errors.go` (modificado: +2 mapeamentos de erro)
- `internal/mcpserver/server_test.go` (modificado: +2 casos de mapeamento)
- `internal/eventstore/event.go` (modificado: `EventSessionResumed`)
- `.pose/specs/2026-08-22-session-orchestration-disclosure.md` (atualizado)

### Validation executed
- Command: go test ./... -race && govulncheck ./...
- Result: `ok` em todos os 8 pacotes; `No vulnerabilities found.`

### Residual risks
- Conteúdo mal classificado pode contornar disclosure e exige gate editorial.

### Follow-ups
- [covered: catalog-authoring-quality] Validar classificação de conteúdo reservado.
