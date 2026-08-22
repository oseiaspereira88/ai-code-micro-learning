# Cobertura das specs da V1

Data: 2026-08-22
Fonte de verdade do produto: `PROJECT.md`
Roadmap governado: `.pose/roadmaps/ailearn-v1.md`

## Resultado

O escopo integral da V1 foi decomposto em 26 specs `draft` prontas para
execução: uma baseline arquitetural, 24 incrementos de produto e uma spec de
aceite integrado. O grafo possui 105 relações de dependência e é acíclico.

`draft` é o estado correto: a documentação de execução está pronta, mas
nenhuma spec alega implementação, evidência ou entrega inexistente. Cada alvo
de entrega cujo caminho ainda será criado está descrito como planejado e só
será registrado como alvo tipado depois de materializado.

## Matriz de cobertura

| Escopo em `PROJECT.md` | Specs proprietárias |
|---|---|
| 1–6: produto, objetivos, perfis e princípios pedagógicos | `learning-domain-model`, `session-orchestration-disclosure`, `assistance-hints-detours`, `feedback-evaluation-progression` |
| 7–10: ontologia, modelo pedagógico, experiência e requisitos funcionais | `learning-domain-model`, `learning-practice-debug-modes`, `interview-mode`, `mastery-review-scheduling` e specs do núcleo pedagógico |
| 11, 23–24 e 32–33: arquitetura, repositório e decisões técnicas | `architecture-decision-baseline`, `go-runtime-foundation` |
| 12–13: agregados, invariantes, estados e transições | `learning-domain-model`, `local-event-store`, `session-orchestration-disclosure` |
| 14: catálogo, taxonomia, packs, trilhas, grafo e qualidade editorial | `catalog-schema-loader`, `curriculum-graph-path-recommendation`, `catalog-authoring-quality`, `go-foundations-packs`, `go-backend-packs`, `go-production-architecture-packs`, `go-interviews-pack` |
| 15: transporte, envelopes, erros e tools MCP | `mcp-stdio-foundation` e as seis specs do núcleo pedagógico que ampliam o contrato |
| 16–17: skill do tutor e CLI | `tutor-skill-host-integration`, `administrative-cli-fixtures` |
| 18: persistência local e eventos | `local-event-store` |
| 19–20: observação do workspace e checks seguros | `workspace-observation-baselines`, `safe-check-executor` |
| 21: avaliação, domínio e revisão | `feedback-evaluation-progression`, `mastery-review-scheduling` |
| 22 e 25–27: segurança, privacidade, confiabilidade, observabilidade e testes | `security-privacy-hardening`, `reliability-observability-compatibility`, `installation-documentation-ci` |
| 28–29: fases e sequência de implementação | milestones e dependências de `ailearn-v1` |
| 30–31: aceite, riscos e mitigação | `v1-integrated-acceptance` e requisitos de risco distribuídos nas specs proprietárias |
| 34: decomposição completa de desafio | `catalog-authoring-quality` e os quatro conjuntos de packs |

## Metas editoriais preservadas

- 84 desafios: 42 atômicos, 20 combinados, 8 slices, 6 de
  debugging/refatoração, 4 missões e 4 simulações.
- 160 conceitos, 100 competências, 500 nós curriculares e 12 trilhas.
- Cada meta tem uma spec autora e volta a ser reconciliada no aceite integrado.

## Ordem governada

1. Baseline arquitetural.
2. Fundação executável.
3. Núcleo pedagógico.
4. Superfícies do tutor e modos.
5. Catálogo V1.
6. Hardening de produção.
7. Aceite integrado da V1.

O roadmap mantém os sete grupos como milestones explícitos. Dependências entre
specs impedem que aceite, distribuição ou conteúdo declarem prontidão antes
dos contratos que consomem.

## Evidência de planejamento

- `pose lint-spec --all --ready-check`: 26 verificadas, 0 falhas.
- `pose check --strict`: estrutura POSE válida.
- `pose index`: índices gerados com integridade de entrega.
- `pose followups --all --json`: 26 follow-ups, 0 abertos, 0 vencidos e 0 sem disposição.
- `pose roadmap-check ailearn-v1 --tolerant --json`: roadmap ativo e íntegro; os 26 bloqueios reportados correspondem corretamente às 26 specs ainda não implementadas.
- Ordenação topológica de `.pose/indexes/spec-graph.json`: concluída sem ciclo.
- Auditoria editorial: nenhum placeholder `TODO`, `TBD` ou de template permaneceu nas specs e no roadmap.

## Limite desta entrega

Esta entrega conclui o planejamento executável do projeto. Código, conteúdo,
testes e evidências permanecem deliberadamente não implementados e devem ser
produzidos seguindo a ordem, os gates e os critérios de cada spec.
