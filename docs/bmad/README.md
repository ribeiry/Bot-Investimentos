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
    │   ├── 004-benchmark.md
    │   ├── 008-postgresql.md
    │   ├── 009-simulacao-e-se.md
    │   ├── 011-market-provider-concorrencia.md
    │   └── 014-telegram-id-unico.md
    └── backlog/           ← aguardando aprovação/implementação
        ├── 007-resumo-llm.md
        ├── 013-get-user-me.md
        ├── 015-menu-interativo-telegram.md
        ├── 016-n8n-ajuste-endpoints-response.md
        └── 017-n8n-mapping-telegram-apikey.md
```

## Stories do backlog

| ID | Story | Status |
|----|-------|--------|
| 001 | Multiusuário | ✅ Done |
| 002 | Alertas de preço | ✅ Done |
| 003 | Resumo semanal/mensal | ✅ Done |
| 004 | Benchmark IBOV/S&P500 | ✅ Done |
| 005 | Alocação por setor | ✅ Done |
| 006 | Notificação de dividendos | ❌ Cancelado — Brapi e Twelve Data exigem plano pago |
| 007 | Resumo LLM (Claude API) | 📋 Backlog — Média |
| 008 | Migração PostgreSQL | ✅ Done |
| 009 | Simulação "e se" | ✅ Done |
| 010 | Envelope `telegram_id` na response | ✅ Done |
| 011 | Concorrência no Market Provider | ✅ Done |
| 012 | Brapi Batch Request | ❌ Cancelado — plano free limita 1 ativo/request |
| 013 | GET /users/me | 📋 Backlog — Baixa |
| 014 | Validação telegram_id único | ✅ Done |
| 015 | Menu interativo no Telegram (n8n) | 📋 Backlog — Média |
| 016 | Ajuste n8n para novos endpoints + envelope | 📋 Backlog — Alta |
| 017 | Mapping telegram_id → api_key no n8n | 📋 Backlog — Alta |
