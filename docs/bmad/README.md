# BMAD — Portfolio Monitoring Bot

Documentação de processo usando o **BMAD Method** (Breakthrough Method of Agile AI-driven Development).

## Como usar

Ative uma persona dizendo ao assistente:

| Comando | O que acontece |
|---------|----------------|
| `"atua como PM"` | Escreve/refina stories com ACs |
| `"atua como Architect"` | Propõe design técnico para uma story |
| `"atua como Dev"` | Implementa o design aprovado |
| `"atua como QA"` | Valida cobertura e ACs |

**Você decide** se implanta ou não antes de qualquer código ser escrito.

## Fluxo

```
PM escreve story
       ↓
  Você aprova?
  ├── Não → ajusta escopo
  └── Sim → Architect propõe design
                   ↓
             Você aprova?
             ├── Não → ajusta design
             └── Sim → Dev implementa
                              ↓
                         QA valida
```

## Estrutura

```
docs/bmad/
├── README.md              ← este arquivo
├── product-brief.md       ← visão do produto
├── architecture.md        ← arquitetura atual
├── personas/
│   ├── pm.md
│   ├── architect.md
│   ├── dev.md
│   └── qa.md
└── stories/
    ├── done/              ← stories concluídas
    │   ├── 001-multiusuario.md
    │   ├── 002-alertas-preco.md
    │   ├── 003-resumo-semanal-mensal.md
    │   └── 004-benchmark.md
    └── backlog/           ← aguardando aprovação/implementação
        ├── 005-alocacao-setor.md
        ├── 006-dividendos.md
        ├── 007-resumo-llm.md
        ├── 008-postgresql.md
        └── 009-simulacao-e-se.md
```

## Stories do backlog

| ID | Story | Prioridade |
|----|-------|------------|
| 005 | Alocação por setor | Alta |
| 006 | Notificação de dividendos | Média |
| 007 | Resumo LLM (Claude API) | Média |
| 008 | Migração PostgreSQL | Alta |
| 009 | Simulação "e se" | Baixa |
