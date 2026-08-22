---
slug: architecture-decision-baseline
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on:
priority: 0
components: governance, architecture
delivers:
---

# Spec: architecture-decision-baseline

## 1. Intent

### Goal
Registrar as decisões estruturais mínimas que permitem implementar a V1 sem reinterpretar PROJECT.md a cada spec.

### Business value
Reduzir retrabalho e impedir que contratos centrais de pedagogia, execução segura e persistência sejam alterados incidentalmente.

### Constraints
- Tratar PROJECT.md como visão, não como evidência de entrega.
- Registrar decisões duráveis em ADRs; não copiar toda a especificação do produto.
- Manter o núcleo independente de Codex, MCP, terminal e formato de catálogo.

### Non-goals
- Implementar código, tools MCP ou conteúdo curricular.
- Decidir detalhes reversíveis de organização interna.

## 2. Requirements

### Functional
- R1: Registrar os limites entre agente tutor, skill, servidor MCP, CLI e núcleo.
- R2: Registrar MCP local por stdio como transporte V1 e a ausência de LLM no servidor.
- R3: Registrar catálogo YAML versionado e estado local por JSONL append-only com snapshots.
- R4: Registrar que o MCP não edita código e só executa checks declarados por ID.
- R5: Registrar a separação entre feedback, avaliação, conclusão e avanço.
- R6: Registrar critérios objetivos para rever cada decisão.

### Non-functional
- ADRs devem ser curtos, vinculados a PROJECT.md e individualmente revisáveis.
- Nenhum ADR pode alegar comportamento implementado.

### Security
- A decisão de execução deve exigir contenção de paths, ausência de shell e mínimo privilégio.

### Compatibility
- As decisões devem preservar a possibilidade de outros hosts MCP e transporte HTTP futuro.

## 3. Technical Plan

### Affected areas
- .pose/adr/
- PROJECT.md como referência somente leitura

### Artifacts
- created: .pose/adr/2026-08-22-agent-mcp-and-core-boundaries.md
- created: .pose/adr/2026-08-22-local-versioned-catalog-and-event-state.md
- created: .pose/adr/2026-08-22-learner-code-ownership-and-safe-checks.md
- created: .pose/adr/2026-08-22-independent-pedagogical-transitions.md

### Delivery targets
Nenhum; esta spec governa decisões e não entrega superfície executável.

### API/contract changes
- Estabelecer contratos arquiteturais que specs posteriores deverão citar.

### Data/storage changes
- Nenhum dado será criado além dos ADRs versionados.

### Technical risks
- ADRs amplos demais podem congelar escolhas reversíveis.
- Duplicação com PROJECT.md pode criar duas fontes de verdade.

## 4. Tasks

### Planning
- [ ] Confirmar cada decisão contra as seções 11, 15, 18, 20 e 32 de PROJECT.md.
- [ ] Delimitar contexto, decisão, consequências e gatilhos de revisão.

### Implementation
- [ ] Criar os quatro ADRs pelo CLI POSE.
- [ ] Ligar cada ADR às specs dependentes sem duplicar seu conteúdo.
- [ ] Registrar alternativas rejeitadas e consequências negativas.

### Validation
- [ ] Executar pose check --strict.
- [ ] Revisar todos os links e garantir ausência de alegação de entrega.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: A V1 contém decisões estruturais já descritas em uma visão extensa.
- Options considered: implementar diretamente; repetir decisões em cada spec; consolidar ADRs mínimos.
- Decision: consolidar quatro ADRs temáticos antes da implementação.
- Rationale: manter rastreabilidade sem transformar o documento de visão em log de decisão.
- Consequences: specs executáveis dependem desta baseline e podem exigir ADR adicional ao mudar contratos.

## 6. Validation

### Strategy
Validar estrutura, referências e precisão editorial; não executar checks de aplicação.

### Deterministic checks
- Test: pose check --strict; todos os ADRs e links devem ser válidos.
- Lint: pose docs-check quando a governança documental estiver ativa.
- Typecheck: não aplicável; não há código.
- Build: não aplicável; não há binário.
- Security / Contract: revisão manual contra .pose/rules/security.md.

### Execution log
- Pendente; a spec permanece draft.

### Results summary
- Nenhuma decisão foi registrada ou validada ainda.

### Requirement trace
- Será preenchido no closeout com um ADR e evidência para cada R-ID.

### Known gaps
- Os nomes finais dos ADRs devem ser confirmados pelo output de pose new-adr.

## 7. Final Report

### Delivered scope
Nenhum; este documento é um plano governado em estado draft.

### Files and modules changed
- Planejados: quatro ADRs sob .pose/adr/.

### Validation executed
- Command: pose lint-spec architecture-decision-baseline --ready-check
- Result: registrar após a validação do planejamento.

### Residual risks
- Alterações arquiteturais futuras exigirão ADRs próprios.

### Follow-ups
- [covered: v1-integrated-acceptance] Reconciliar as decisões arquiteturais no aceite final.
