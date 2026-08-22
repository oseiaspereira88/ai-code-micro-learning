# ADR: Learner code ownership and safe checks

## Status
Accepted

## Context

`PROJECT.md` §11.2 estabelece que o MCP não edita código do aluno; §20
(Executor seguro de checks) e §32 item 6 detalham que não existe
`run_command(command: string)` livre — o agente envia apenas um `check_id`
pertencente ao desafio fixado. `.pose/rules/security.md` exige mínimo
privilégio, validação de entrada externa e ausência de execução dinâmica sem
validação. Sem registrar isso como decisão estrutural, uma spec futura
poderia introduzir execução de comando livre por conveniência.

Módulos afetados: `safe-check-executor`, `workspace-observation-baselines`,
servidor MCP (tools de checks e workspace).

## Decision

1. O aluno é o único autor do código durante a tutoria; o MCP nunca edita
   arquivos do workspace do aluno.
2. A execução de verificação aceita somente `check_id` de um catálogo fixado
   pelo desafio; não existe parâmetro de comando livre.
3. Tipos permitidos na V1: `go_test`, `go_test_race`, `go_vet`,
   `gofmt_check`, `go_build`, `go_benchmark`, `catalog_validator` e
   verificadores internos por AST/arquivo (`PROJECT.md` §20.2).
4. Controles obrigatórios de execução: `os/exec` direto sem shell,
   executáveis em allowlist, argumentos derivados de campos validados,
   diretório contido no workspace autorizado, rejeição de symlink que escape
   da raiz, ambiente mínimo com allowlist de variáveis, timeout por check,
   encerramento do grupo de processos ao cancelar, limites de stdout/stderr/
   memória/paralelismo quando suportado, redaction de segredos e ausência de
   rede por padrão, exceto checks marcados e aprovados explicitamente
   (`PROJECT.md` §20.3).
5. Resultado de check é vinculado à revisão do workspace para evitar
   evidência obsoleta.

### Alternativas rejeitadas

- **`run_command(command: string)` genérico**: rejeitado por criar superfície
  de execução arbitrária incompatível com mínimo privilégio.
- **Execução via shell**: rejeitada pelo risco de injeção e pela regra de
  segurança que veta comandos dinâmicos sem validação e escaping adequados.
- **Rede habilitada por padrão nos checks**: rejeitada; qualquer necessidade
  de rede exige marcação e aprovação explícita, nunca comportamento padrão.

## Consequences

- `safe-check-executor` deve implementar a allowlist de executáveis, a
  contenção de path e os limites de recurso antes de expor `check_run` no
  MCP.
- `workspace-observation-baselines` deve fornecer as referências de diff e
  hashes que ligam o resultado do check à revisão do workspace, sem transitar
  por edição de arquivos do aluno.
- Qualquer novo tipo de check deve ser adicionado à allowlist declarada, não
  introduzido como exceção pontual.

### Gatilho de revisão

Revisar esta decisão se um cenário legítimo exigir capacidade incompatível
com o modelo de allowlist (por exemplo, comando fora do conjunto permitido ou
necessidade estrutural de shell), documentando o caso antes de qualquer
exceção.
