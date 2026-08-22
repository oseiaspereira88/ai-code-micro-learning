# ADR: Independent pedagogical transitions

## Status
Accepted

## Context

`PROJECT.md` §15.9 (Tools de feedback e avaliação) e §32 item 8 estabelecem
que feedback, avaliação, conclusão e avanço são operações separadas. Os
valores de `progress_effect` em §15.3 (`feedback_recorded`,
`evaluation_recorded`, `step_completed`, `step_advanced`, entre outros)
reforçam que cada tool tem um efeito de progresso único e explícito. Sem
registrar essa separação como decisão estrutural, uma spec de sessão poderia
implementar atalhos que inferem conclusão a partir de feedback ou avaliação,
violando a divulgação progressiva.

Módulos afetados: `session-orchestration-disclosure`,
`feedback-evaluation-progression`, `mastery-review-scheduling`.

## Decision

1. `feedback_prepare`/`feedback_record` produzem e registram feedback sem
   concluir ou avaliar; sempre retornam `progress_effect: feedback_recorded`.
2. `step_evaluate` executa critérios determinísticos e recebe julgamentos
   qualitativos ligados a evidência, cria tentativa somente quando
   `submission_intent` é verdadeiro, e não conclui nem avança.
3. `step_complete` valida política e marca o passo concluído, sem ativar o
   próximo nó.
4. `step_advance` ativa o próximo nó permitido ou informa ramificação,
   como operação distinta da conclusão.
5. Nenhuma tool pode produzir, como efeito colateral implícito, o
   `progress_effect` de outra categoria (feedback ≠ avaliação ≠ conclusão ≠
   avanço).

### Alternativas rejeitadas

- **Avançar automaticamente após avaliação aprovada**: rejeitado por remover
  o ponto de decisão explícito do aluno e mascarar a intenção pedagógica da
  transição.
- **Unificar feedback e avaliação em uma única chamada**: rejeitado por
  confundir orientação aberta (feedback) com julgamento baseado em rubrica
  (avaliação), quebrando a rastreabilidade de evidência exigida por §15.9.
- **Inferir conclusão a partir de checks determinísticos sem chamada
  explícita**: rejeitado por remover o controle de política que
  `step_complete` deve aplicar antes de marcar o passo.

## Consequences

- `session-orchestration-disclosure` e `feedback-evaluation-progression`
  devem expor essas quatro operações como tools/transições distintas, cada
  uma com seu próprio `progress_effect`.
- `mastery-review-scheduling` só pode consumir eventos de avaliação e
  conclusão já registrados separadamente, nunca inferir um a partir do outro.
- Testes de aceite devem cobrir que uma chamada de feedback ou avaliação,
  isoladamente, nunca produz avanço de estado.

### Gatilho de revisão

Revisar esta decisão se evidência de uso mostrar que a separação confunde o
aluno de forma sistemática, ou se um padrão pedagógico comprovar que uma
transição automática é segura sob um controle de divulgação específico.
