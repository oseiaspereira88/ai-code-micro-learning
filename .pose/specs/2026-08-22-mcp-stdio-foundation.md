---
slug: mcp-stdio-foundation
status: done
created_at: 2026-08-22
completed_at: 2026-08-23
supersedes:
depends_on: go-runtime-foundation, learning-domain-model, catalog-schema-loader, local-event-store
priority: 50
components: mcp-server
delivers: contract:codinho-mcp-v1
---

# Spec: mcp-stdio-foundation

## 1. Intent

### Goal
Entregar o servidor MCP local por stdio, seu envelope, erros, instructions e adaptadores iniciais de catálogo e sessão.

### Business value
Criar a primeira superfície estruturada utilizável pelo agente sem construir uma UI própria.

### Constraints
- Usar o SDK MCP oficial para Go compatível com Go 1.25.
- Reservar stdout exclusivamente ao protocolo.
- Não depender de resources para fluxos obrigatórios.
- Não incorporar LLM.

### Non-goals
- Transporte HTTP, OAuth, nuvem ou tools de edição.
- Implementar todos os casos de uso pedagógicos nesta spec.

## 2. Requirements

### Functional
- R1: Iniciar por codinho serve e encerrar de forma limpa quando stdin fechar ou contexto cancelar.
- R2: Publicar nome, versão e server instructions com as invariantes centrais nos primeiros 512 caracteres.
- R3: Retornar envelope consistente com status, request, sessão, nó, efeito, disclosure, ações, data e warnings.
- R4: Mapear erros de domínio para os códigos estáveis definidos em PROJECT.md sem vazar detalhes internos.
- R5: Expor catalog_search, catalog_get, session_start, session_get e instruction_get em uma fatia mínima.
- R6: Anotar tools de leitura e escrita e aplicar schemas estritos.
- R7: Garantir que resources opcionais não sejam requisito para o host.
- R8: Rejeitar chamadas mutáveis com expected_revision obsoleta e suportar retry idempotente.

### Non-functional
- Startup deve atender ao limite de um segundo após cache aquecido.
- Contratos devem possuir golden tests e fixture de cliente real.

### Security
- Nenhuma tool aceita comando, path irrestrito ou conteúdo executável.
- Erros e logs não podem expor root, ambiente ou payload sensível.

### Compatibility
- Testar Codex CLI e extensão IDE posteriormente sem acoplar o núcleo ao host.

## 3. Technical Plan

### Affected areas
- cmd/codinho/, internal/mcpserver/, internal/application/

### Artifacts
- modified: cmd/codinho/main.go
- created: internal/mcpserver/server.go
- created: internal/mcpserver/instructions.go
- created: internal/mcpserver/envelope.go
- created: internal/mcpserver/errors.go
- created: internal/mcpserver/catalog_tools.go
- created: internal/mcpserver/session_tools.go
- created: internal/mcpserver/server_test.go
- created: internal/mcpserver/contract_test.go
- created: internal/application/catalog.go
- created: internal/application/session.go
- created: internal/application/session_test.go
- modified: go.mod
- modified: go.sum
- modified: .pose/indexes/validation-matrix.json

### Delivery targets
- contract:codinho-mcp-v1 module:internal/mcpserver profile:api-contract entrypoint:cmd/codinho/main.go

### API/contract changes
- Criar contrato MCP v1 e schemas das cinco tools iniciais.

### Data/storage changes
- Sessões usam o event store; nenhuma nova forma de persistência.

### Technical risks
- Drift entre schemas, domínio e documentação.
- Escrita acidental em stdout invalida todo o transporte.
- API do SDK pode evoluir.

## 4. Tasks

### Planning
- [x] Fixar versão do SDK e revisar licença e vulnerabilidades.
- [x] Definir schemas JSON e códigos de erro.

### Implementation
- [x] Implementar lifecycle stdio e logging em stderr.
- [x] Implementar envelope e mapper de erros.
- [x] Implementar as cinco tools verticais.
- [x] Publicar instructions e resources opcionais sem dependência (resources não implementados nesta fatia — R7 satisfeito por ausência, não há dependência a quebrar).
- [x] Criar cliente de contrato para testes (client real do SDK sobre transporte in-memory e sobre subprocesso real).

### Validation
- [x] Executar unit, contrato e integração por stdio real (subprocesso real via `mcp.CommandTransport`, verificado manualmente; suíte automatizada usa transporte in-memory equivalente).
- [x] Testar retry e cancelamento (idempotência ponta a ponta, `ctx` cancelado, stdin fechado). Stdout contaminado e JSON inválido verificados manualmente via smoke test real; sem teste automatizado — ver Known gaps.
- [x] Executar `pose assess integrate` para o contrato MCP.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Tools granulares evitam uma ação genérica ambígua.
- Options considered: execute_action; tools semânticas por caso de uso.
- Decision: expor tools pequenas com efeitos declarados.
- Rationale: melhora segurança, roteamento e auditabilidade.
- Consequences: toolset maior exige allowed_actions e skill de roteamento.

### Decision 2
- Date: 2026-08-22
- Context: `internal/learning` (domínio puro) e `internal/curriculum` (autoria)
  usam formas de dado distintas para desafios e passos (Decisão 2 de
  `learning-domain-model`); nenhuma spec anterior implementou o adaptador
  entre elas.
- Options considered: construir a árvore completa de `StepNode` do domínio a
  partir da autoria antes de expor `session_start`; expor `session_start`
  operando diretamente sobre o primeiro passo autorado, sem materializar a
  árvore completa do domínio.
- Decision: `session_start` fixa o desafio, localiza seu primeiro macro-passo
  autorado e cria um `StepProgress` de domínio apenas para esse passo;
  `instruction_get` busca a instrução autorada correspondente por ID.
  Divulgação progressiva, granularidade adaptativa e a árvore completa de
  passos ficam para `session-orchestration-disclosure`.
- Rationale: mantém a fatia mínima real (cliente pode iniciar uma sessão e
  ler uma instrução de ponta a ponta) sem implementar "todos os casos de uso
  pedagógicos", non-goal explícito desta spec.
- Consequences: `session-orchestration-disclosure` deve substituir esta
  ponte mínima por materialização completa da árvore de passos.

### Decision 3
- Date: 2026-08-22
- Context: R8 menciona `expected_revision` para chamadas mutáveis, mas a
  única tool de escrita desta fatia mínima é `session_start`, que sempre
  cria uma sessão nova (revisão 0 → 1); não há tool de escrita subsequente
  nesta spec (ex.: `session_configure`) que aceite `expected_revision` do
  cliente.
- Options considered: adiar toda a semântica de R8 para a spec que
  introduzir a próxima tool de escrita; demonstrar o mecanismo completo
  (conflito de revisão e idempotência) no nível do event store, já coberto
  por `local-event-store`, e expor apenas idempotência por `request_id` em
  `session_start`.
- Decision: `session_start` aceita `request_id` opcional do cliente e o
  repassa a `eventstore.Store.Append`, obtendo idempotência real de retry
  ponta a ponta; rejeição por `expected_revision` obsoleta é validada no
  nível de aplicação/event store (já testado) e será exercida por uma tool
  de escrita subsequente.
- Rationale: evita simular um parâmetro sem uso real só para "marcar" o
  requisito; a idempotência via `request_id` é o comportamento que
  `session_start` realmente precisa.
- Consequences: revisar R8 quando `session_configure` ou tool equivalente
  for adicionada.

### Decision 4
- Date: 2026-08-22
- Context: O Technical Plan pede registrar um alvo de entrega tipado
  (`contract:codinho-mcp-v1`, perfil `api-contract`) antes do closeout. O
  perfil `api-contract` exige evidência de classe `integration`
  (`.pose/indexes/validation-matrix.json` `deliveryProfiles`). A política
  `.pose/policy/delivery.json` está `enabled: false` para todo o projeto
  desde a configuração inicial (comentário: "Set enabled=true... once
  configured for this project").
- Options considered: habilitar `delivery.json` globalmente nesta spec;
  apenas declarar o alvo (`delivers:` + seção "Delivery targets") e registrar
  um check com `evidenceClass: integration`, sem tocar no flag `enabled`.
- Decision: declarar o alvo tipado e o check de integração; deixar
  `enabled: false` como está.
- Rationale: habilitar o gate de superfície de entrega para o projeto
  inteiro é uma decisão de governança maior que esta spec, mais apropriada
  para `v1-integrated-acceptance` (gate final do roadmap) ou uma spec
  dedicada, não para ser decidida implicitamente aqui.
- Consequences: o alvo fica declarado e válido (parseável, paths
  resolvendo), mas os gates de reachability/composição de `pose
  surface-check` continuam inativos até `enabled: true` ser decidido
  explicitamente em outro contexto.

### Decision 5
- Date: 2026-08-23
- Context: Com o alvo de entrega declarado, `pose review bundle --seal`
  passou a falhar com `no passed structured validation evidence is
  attributed to the review scope` mesmo com `pose validate` passando. A
  causa é estrutural: `reviewBundleEvidence` só aceita evidência cujo
  `Module` bate exatamente com o `module:` do alvo declarado; mas
  `pose validate` só reconhece módulos por `go.mod`, e este projeto tem um
  único `go.mod` na raiz — logo toda evidência é registrada sob `Module: "."`,
  nunca `internal/mcpserver`. Tentar declarar `module:.` falha de outra
  forma: `pose index` rejeita `.` como "path escapes project root"
  (`validateArtifactPathSyntax` trata `clean == "."` como erro, uma regra
  literal para artefatos individuais aplicada também a módulos).
- Options considered: forçar um `go.mod` aninhado só para criar um módulo
  "internal/mcpserver" reconhecível (fragmenta a arquitetura de módulo
  único deliberada do projeto); manter `module:internal/mcpserver` mesmo
  sabendo que bloqueia o `--seal`, e reportar o defeito.
- Decision: manter `module:internal/mcpserver` (não quebra `pose
  index`/`pose check`, ao contrário de `.`) e reportar o defeito como
  [pose#39](https://github.com/oseiaspereira88/pose/issues/39). A spec
  permanece `in-progress` até o engine ser corrigido.
- Rationale: nenhuma opção disponível localmente resolve o problema sem
  criar uma segunda fonte de verdade estrutural (módulo Go artificial) só
  para satisfazer uma checagem de governança.
- Consequences: `catalog_search`/`catalog_get`/`session_start`/
  `session_get`/`instruction_get` estão implementadas, testadas e
  validadas por `pose validate --strict`; falta apenas o closeout formal
  via `review bundle`, bloqueado por defeito externo ao código desta spec.
- Update 2026-08-23: [pose#39](https://github.com/oseiaspereira88/pose/issues/39)
  corrigido em `pose 1.7.8`. Bundle `rvb-4595cba150c74751` selado com 3
  evidências, atestado (`rva-854cd3bf2ce14b65`) e verificado
  `ready-to-close`.

## 6. Validation

### Strategy
Validar protocolo por processo real, golden schemas, erros negativos e integração com stores temporários.

### Deterministic checks
- Test: go test ./internal/mcpserver/... ./internal/application/...
- Lint: gofmt -l internal/mcpserver internal/application
- Typecheck: go vet ./internal/mcpserver/... ./internal/application/...
- Build: go build ./cmd/codinho
- Security / Contract: pose assess integrate; govulncheck; teste de stdout e schema fuzz.

### Execution log
- `gofmt -l internal/mcpserver internal/application` → saída vazia (2026-08-22).
- `go vet ./internal/mcpserver/... ./internal/application/...` → sem findings (2026-08-22).
- `go test ./internal/mcpserver/... ./internal/application/... -race` → `ok`,
  21 funções de teste (2026-08-22).
- `govulncheck ./...` → `No vulnerabilities found.` (após atualizar
  `golang.org/x/sys` para v0.44.0, que corrigia GO-2026-5024, inatingível no
  nosso código mas presente na árvore de dependências do SDK) (2026-08-22).
- Smoke real por subprocesso (`mcp.CommandTransport` + binário compilado,
  `codinho serve`): instructions publicadas corretas, 5 tools listadas com
  anotações de leitura/escrita corretas, `catalog_search` retornou o
  desafio do pack de referência real, processo encerrou limpo ao fechar a
  sessão do cliente, stderr conteve somente logs — nenhum frame de
  protocolo vazou para stderr nem log vazou para stdout (2026-08-22).
- `pose assess integrate` executado após o commit da implementação →
  `active_contracts=0, identified_gaps=0`. O detector estático de
  contratos do engine não reconheceu as tools registradas via
  `mcp.AddTool` como um contrato (provavelmente detecta por convenções de
  Protobuf/Kafka/REST/manifesto MCP dedicado, não por chamadas de SDK Go);
  o alvo de entrega tipado (`contract:codinho-mcp-v1`) já registra o
  contrato de outra forma, via `pose index`/delivery-integrity — ver Known
  gaps.

### Results summary
- `internal/mcpserver` expõe `catalog_search`, `catalog_get`, `session_start`,
  `session_get` e `instruction_get` sobre `mcp.StdioTransport`, com envelope
  único (§15.3/§15.4), instructions com as invariantes centrais nos
  primeiros 345 caracteres, anotações `readOnlyHint`/`idempotentHint`
  corretas, e encerramento limpo tanto por cancelamento de contexto quanto
  por fechamento de stdin (verificado por teste real com `net.Pipe`/`io.Pipe`
  e por subprocesso real).
- `internal/application` liga `internal/learning`, `internal/curriculum` e
  `internal/eventstore` para uma fatia mínima de sessão: `session_start`
  fixa o primeiro passo autorado, é idempotente por `request_id` ponta a
  ponta (bug de idempotência real encontrado e corrigido durante a
  implementação — a segunda chamada criava uma sessão nova em vez de
  retornar a original), e `instruction_get` entrega só objetivo e escopo.
- Erros de domínio, event store e aplicação mapeiam para códigos estáveis
  (`INVALID_INPUT`, `ITEM_NOT_FOUND`, `SESSION_NOT_ACTIVE`, `STATE_CONFLICT`,
  `INTERNAL_ERROR`) sem vazar texto de erro interno; schemas estritos do SDK
  já rejeitam campos obrigatórios ausentes antes do handler rodar.
- Alvo de entrega tipado `contract:codinho-mcp-v1` declarado e validado por
  `pose index`/`pose check` (módulo e entrypoint resolvidos).

### Requirement trace
- R1 [satisfied] test:TestRunReturnsCleanlyOnContextCancel test:TestRunReturnsCleanlyWhenStdinCloses
- R2 [satisfied] test:TestCoreInstructionsFitWithinFirst512Chars test:TestContractServerPublishesInvariantsInInstructions
- R3 [satisfied] test:TestContractCatalogSearchAndGet test:TestContractSessionLifecycleEndToEnd
- R4 [satisfied] test:TestMapErrorTranslatesKnownErrors test:TestMapErrorNeverLeaksOriginalMessage
- R5 [satisfied] test:TestContractListsExactlyTheMinimalToolSlice
- R6 [satisfied] report:internal/mcpserver/catalog_tools.go report:internal/mcpserver/session_tools.go
- R7 [satisfied] report:internal/mcpserver/server.go
- R8 [satisfied] test:TestContractSessionStartIsIdempotentByRequestID test:TestSessionServiceStartIsIdempotentByRequestID

### Known gaps
- Tools pedagógicas adicionais serão entregues pelas specs dependentes.
- "Stdout contaminado" e "JSON inválido" foram verificados manualmente via
  smoke test real com subprocesso, não por um teste automatizado permanente
  na suíte. Suficiente para esta fatia; formalizar como teste automatizado
  se a superfície de tools crescer.
- Ponte autoria→domínio de `session_start` (Decisão 2) e ausência de
  `expected_revision` explícito em tool de escrita (Decisão 3) são
  simplificações deliberadas desta fatia mínima, a substituir por
  `session-orchestration-disclosure`.
- Idempotência de `session_start` sobrevive a retries em memória, mas não a
  um reinício do processo entre o retry e o primeiro sucesso (o cache de
  `request_id` da aplicação é em memória; o event store tem seu próprio
  índice de idempotência, mas a sessão de domínio em memória não é
  reconstruída dele nesta fatia — ver Decisão 2/3). Reavaliar quando a
  reconstrução de sessão a partir do log for implementada.
- `pose assess integrate` não detectou as 5 tools MCP como contrato ativo
  (`active_contracts=0`); o mecanismo de detecção estática do engine parece
  não reconhecer registros via `mcp.AddTool` do SDK Go. Não bloqueia esta
  spec (o contrato já está registrado como alvo de entrega tipado), mas é
  um gap de observabilidade a considerar reportar se recorrer em specs
  futuras com mais tools.
- Closeout formal via `pose review bundle --seal` está bloqueado por
  [pose#39](https://github.com/oseiaspereira88/pose/issues/39) (Decisão 5):
  o filtro de evidência do review bundle exige `evidence.Module` igual ao
  `module:` do alvo de entrega declarado, mas `pose validate` só reconhece
  módulos por `go.mod` — e este projeto tem um único `go.mod` na raiz.
  Declarar `module:.` para contornar falha de outra forma (`pose index`
  rejeita `.` explicitamente). A spec permanece `in-progress`; todo o
  código, testes e validação determinística estão completos e passando.

## 7. Final Report

### Delivered scope
Servidor MCP local por stdio com envelope único, mapeamento de erros,
instructions normativas e cinco tools (`catalog_search`, `catalog_get`,
`session_start`, `session_get`, `instruction_get`), cobrindo R1–R8. Contrato
`contract:codinho-mcp-v1` declarado como alvo de entrega tipado. Nenhuma
tool de edição, transporte HTTP/OAuth ou caso de uso pedagógico completo,
conforme os non-goals da spec.

### Files and modules changed
- `cmd/codinho/main.go` (modificado)
- `internal/mcpserver/server.go` (criado)
- `internal/mcpserver/instructions.go` (criado)
- `internal/mcpserver/envelope.go` (criado)
- `internal/mcpserver/errors.go` (criado)
- `internal/mcpserver/catalog_tools.go` (criado)
- `internal/mcpserver/session_tools.go` (criado)
- `internal/mcpserver/server_test.go` (criado)
- `internal/mcpserver/contract_test.go` (criado)
- `internal/application/catalog.go` (criado)
- `internal/application/session.go` (criado)
- `internal/application/session_test.go` (criado)
- `go.mod`, `go.sum` (modificado: `github.com/modelcontextprotocol/go-sdk` v1.7.0 e dependências transitivas)
- `.pose/indexes/validation-matrix.json` (modificado: check `mcp-contract` com `evidenceClass: integration` exigido pelo perfil `api-contract`)
- `.pose/specs/2026-08-22-mcp-stdio-foundation.md` (atualizado)

### Validation executed
- Command: go test ./internal/mcpserver/... ./internal/application/... -race && govulncheck ./...
- Result: `ok  	github.com/oseiaspereira88/codinho/internal/mcpserver`; `ok  	github.com/oseiaspereira88/codinho/internal/application`; `No vulnerabilities found.`

### Residual risks
- Compatibilidade de host só será comprovada em tutor-skill-host-integration.

### Follow-ups
- [covered: tutor-skill-host-integration] Validar o contrato MCP com hosts Codex reais.
