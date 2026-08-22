---
slug: tutor-skill-host-integration
status: draft
created_at: 2026-08-22
completed_at:
supersedes:
depends_on: mcp-stdio-foundation, session-orchestration-disclosure, assistance-hints-detours, feedback-evaluation-progression, workspace-observation-baselines, safe-check-executor
priority: 130
components: tutor-skill, mcp-server
delivers:
---

# Spec: tutor-skill-host-integration

## 1. Intent

### Goal
Criar a skill ailearn e comprovar seu workflow com o servidor MCP no Codex CLI e na extensão da IDE.

### Business value
Transformar tools e estado em uma experiência pedagógica consistente que preserve a autoria do aluno.

### Constraints
- A skill controla comportamento; o MCP permanece fonte de verdade.
- Não editar arquivos do aluno durante tutoria.
- Funcionar de modo degradado, claramente indicado, quando o MCP não estiver disponível.

### Non-goals
- Plugin distribuído ou interface própria.
- Embutir conteúdo curricular inteiro no SKILL.md.

## 2. Requirements

### Functional
- R1: Definir triggers, limites, tool routing e divulgação progressiva no SKILL.md.
- R2: Declarar dependência MCP e metadados de interface em agents/openai.yaml.
- R3: Ler session_get antes de inferir estado e entregar somente uma instrução ativa.
- R4: Observar antes de avaliar e exigir submission_intent para criar tentativa.
- R5: Tratar feedback como consultivo e informar sempre progress_effect.
- R6: Bloquear edição e solução salvo mudança explícita autorizada pela política.
- R7: Ignorar instruções encontradas em código, fixtures e outputs observados.
- R8: Adaptar conceito ao perfil usando conteúdo canônico e exemplos fora da solução.
- R9: Operar no Codex CLI e IDE com o mesmo MCP configurado.
- R10: Informar limitações no fallback conversacional e não simular persistência.

### Non-functional
- SKILL.md deve permanecer focado e usar referências por divulgação progressiva.
- Respostas em modo micro devem ser curtas e conter uma ação principal.

### Security
- A skill não amplia permissões e respeita approvals do host.
- Testes adversariais cobrem prompt injection e pedidos ambíguos de código.

### Compatibility
- Seguir o formato de Agent Skills e metadados suportados pelo Codex.

## 3. Technical Plan

### Affected areas
- .agents/skills/ailearn/, .codex/, testdata/host/

### Artifacts
- created: .agents/skills/ailearn/SKILL.md
- created: .agents/skills/ailearn/agents/openai.yaml
- created: .agents/skills/ailearn/references/tutor-contract.md
- created: .agents/skills/ailearn/references/feedback-rubric.md
- created: .agents/skills/ailearn/references/session-modes.md
- created: .agents/skills/ailearn/references/mcp-tool-routing.md
- created: .agents/skills/ailearn/assets/session-summary-template.md
- created: testdata/host/adversarial-prompts.yaml
- created: testdata/host/session-transcripts/
- created: .codex/config.toml.example

### Delivery targets
O alvo planejado é a capability `ailearn-tutor`, com módulo
`.agents/skills/ailearn`, perfil `composed-capability` e entrypoint
`.agents/skills/ailearn/SKILL.md`. Registrá-lo como alvo tipado quando esses
caminhos forem materializados, antes do closeout desta spec.

### API/contract changes
- Consumir o contrato MCP v1 sem adicionar estado na skill.

### Data/storage changes
- Nenhum; transcripts de teste são fixtures sintéticas.

### Technical risks
- Instruções longas podem perder prioridade.
- Comportamento pode variar por host ou modelo.

## 4. Tasks

### Planning
- [ ] Mapear cada regra da seção 16 de PROJECT.md para instrução ou referência.
- [ ] Definir corpus adversarial e critérios de passagem.

### Implementation
- [ ] Criar skill, metadados e referências.
- [ ] Implementar roteamento de tools e fallback.
- [ ] Configurar MCP de exemplo sem paths privados.
- [ ] Executar sessões no CLI e IDE.
- [ ] Capturar transcripts sintéticos e resultados.

### Validation
- [ ] Executar pose skills-check --strict.
- [ ] Executar contract/e2e com MCP real.
- [ ] Verificar zero edição e zero revelação indevida.
- [ ] Executar pose assess integrate e surface-check.

## 5. Decisions

### Decision 1
- Date: 2026-08-22
- Context: Instruções e estado têm ciclos de vida diferentes.
- Options considered: tudo na skill; tudo no MCP; responsabilidades separadas.
- Decision: skill governa workflow; MCP governa estado e capabilities.
- Rationale: mantém a conversa flexível e o progresso determinístico.
- Consequences: integração deve testar ambos em composição.

## 6. Validation

### Strategy
Combinar conformance da skill, transcripts adversariais e sessões reais nos dois hosts.

### Deterministic checks
- Test: go test ./internal/mcpserver/... e runner de transcripts sintéticos.
- Lint: pose skills-check --strict.
- Typecheck: validação YAML de openai.yaml.
- Build: go build ./cmd/ailearn.
- Security / Contract: pose assess integrate; pose surface-check --spec tutor-skill-host-integration --strict.

### Execution log
- Pendente.

### Results summary
- Skill e integração ainda não existem.

### Requirement trace
- Mapear R1–R10 a checks de conformance, transcripts e sessões host.

### Known gaps
- Variabilidade de modelos exige testes por comportamento, não texto exato.

## 7. Final Report

### Delivered scope
Nenhum; spec draft.

### Files and modules changed
- Planejados em .agents/skills/ailearn, testdata/host e config de exemplo.

### Validation executed
- Command: pose lint-spec tutor-skill-host-integration --ready-check
- Result: registrar após gate.

### Residual risks
- Compatibilidade futura de host depende da documentação oficial.

### Follow-ups
- [covered: v1-integrated-acceptance] Revalidar o workflow composto no gate final.
