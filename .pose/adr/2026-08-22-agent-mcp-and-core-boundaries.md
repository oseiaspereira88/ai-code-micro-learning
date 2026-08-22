# ADR: Agent, MCP and core boundaries

## Status
Accepted

## Context

A V1 conecta um agente tutor (ex.: Codex CLI ou extensão de IDE), uma skill
pedagógica, um servidor MCP local e um núcleo de casos de uso independente de
protocolo. `PROJECT.md` §11 (Arquitetura) e §15.1–15.2 (Contrato MCP) já
descrevem essas responsabilidades e o transporte, mas nenhuma decisão
estrutural fica registrada fora do documento de visão. Sem uma fronteira
explícita, specs executáveis podem reintroduzir lógica pedagógica dentro do
servidor MCP ou acoplar o núcleo a um host específico.

Módulos afetados: agente tutor, skill, servidor MCP (`mcp-stdio-foundation`),
CLI, núcleo de casos de uso e domínio.

## Decision

1. O núcleo de casos de uso e o modelo de domínio não conhecem Codex,
   protocolo MCP, terminal ou formato de catálogo; adaptadores traduzem esses
   meios (regra arquitetural de `PROJECT.md` §11.3).
2. O servidor MCP é a fonte de verdade de catálogo, sessão e progresso;
   aplica políticas de revelação, executa checks declarados e valida
   transições, mas não chama LLM nem redige feedback aberto por conta própria
   (`PROJECT.md` §11.2, §15.2).
3. O transporte da V1 é MCP local por `stdio`, com stdout reservado ao
   protocolo e logs/diagnóstico em stderr (`PROJECT.md` §15.1).
4. A skill define o workflow pedagógico e a divulgação progressiva; o agente
   tutor conversa, interpreta intenção e não é fonte de verdade de estado.

### Alternativas rejeitadas

- **Embutir um LLM no servidor MCP**: rejeitado porque duplicaria o papel do
  agente tutor e acoplaria o núcleo a um provedor de modelo específico.
- **Adotar Streamable HTTP como transporte inicial**: rejeitado por adicionar
  superfície de rede antes de existir um fluxo local ponta a ponta validado;
  fica reservado para versão posterior.
- **Implementar a lógica pedagógica como tools do servidor MCP**: rejeitado
  porque acoplaria decisões de workflow pedagógico ao protocolo, dificultando
  reuso por outros hosts.

## Consequences

- Specs que implementam o servidor (`mcp-stdio-foundation`,
  `session-orchestration-disclosure`, `tutor-skill-host-integration`) devem
  expressar a integração como adaptadores de entrada, nunca como lógica de
  domínio.
- O núcleo permanece testável sem dependência de protocolo, permitindo
  transporte HTTP futuro sem reescrever casos de uso (requisito de
  compatibilidade da spec).
- Qualquer proposta de mover julgamento pedagógico para dentro do MCP exige
  ADR próprio antes da implementação.

### Gatilho de revisão

Revisar esta decisão se um segundo host de produção exigir um modelo de
transporte incompatível com `stdio`, ou se surgir necessidade comprovada de
raciocínio semântico no servidor sem intermediação do agente tutor.
