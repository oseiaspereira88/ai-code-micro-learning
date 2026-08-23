# ailearn
Microaprendizado de código assistido por IA

O **ailearn** é um projeto de prática deliberada assistida por agente. Seu
objetivo inicial é desenvolver e recuperar fluência em Go sem permitir que a
IA substitua o raciocínio, a decomposição e a escrita de código pelo aluno,
usando o conceito de microaprendizado de código assistido por IA.

## Estado atual

O projeto está em fase de especificação e preparação operacional. A visão, o
escopo e os critérios da V1 estão definidos em [`PROJECT.md`](PROJECT.md). O
código da aplicação ainda não foi iniciado.

## Componentes planejados

- skill responsável pelo comportamento do agente tutor;
- servidor MCP local em Go para currículo, sessões e evidências;
- CLI administrativa enxuta;
- catálogo versionado de conceitos, competências e desafios de Go;
- persistência local e verificações controladas do workspace.

## Desenvolvimento

O repositório usa POSE para governar o trabalho de agentes e a evolução das
especificações. Leia [`AGENTS.md`](AGENTS.md) antes de alterar o projeto e use
[`POSE.md`](POSE.md) como manual operacional.

Os artefatos de produto e de governança têm responsabilidades distintas:

- [`PROJECT.md`](PROJECT.md) é a referência principal da visão da V1;
- [`.pose/specs/`](.pose/specs/) conterá specs incrementais de implementação;
- [`.pose/roadmaps/`](.pose/roadmaps/) conterá a sequência governada de entrega;
- [`.pose/adr/`](.pose/adr/) registrará decisões arquiteturais relevantes.
