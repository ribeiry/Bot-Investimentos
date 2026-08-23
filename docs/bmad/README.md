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
    ├── done/              ← stories concluídas (código Go)
    │   ├── 001-multiusuario.md
    │   ├── 002-alertas-preco.md
    │   ├── 003-resumo-semanal-mensal.md
    │   ├── 004-benchmark.md
    │   ├── 005-alocacao-setor.md
    │   ├── 008-postgresql.md
    │   ├── 009-simulacao-e-se.md
    │   ├── 010-telegram-id-response.md
    │   ├── 011-market-provider-concorrencia.md
    │   ├── 013-get-user-me.md
    │   └── 014-telegram-id-unico.md
    ├── backlog/           ← aguardando aprovação/implementação (API Go)
    │   └── 007-resumo-llm.md
    └── n8n/               ← workflows do n8n (sem código Go)
        ├── 015-menu-interativo-telegram.md
        ├── 016-n8n-ajuste-endpoints-response.md
        └── 017-n8n-mapping-telegram-apikey.md
```

## Status geral

### API Go

| ID | Story | Status |
|----|-------|--------|
| 001 | Multiusuário | ✅ Done |
| 002 | Alertas de preço | ✅ Done |
| 003 | Resumo semanal/mensal | ✅ Done |
| 004 | Benchmark IBOV/S&P500 | ✅ Done |
| 005 | Alocação por setor | ✅ Done |
| 006 | Notificação de dividendos | ❌ Cancelado — APIs exigem plano pago |
| 007 | Resumo LLM (Claude API) | 📋 Backlog — Média |
| 008 | Migração PostgreSQL | ✅ Done |
| 009 | Simulação "e se" | ✅ Done |
| 010 | Envelope `telegram_id` na response | ✅ Done |
| 011 | Concorrência no Market Provider | ✅ Done |
| 012 | Brapi Batch Request | ❌ Cancelado — free limita 1 ativo/request |
| 013 | GET /users/me | ✅ Done |
| 014 | Validação telegram_id único | ✅ Done |

### Workflows n8n

| ID | Story | Status |
|----|-------|--------|
| 015 | Menu interativo no Telegram | 📋 Pendente — Média |
| 016 | Ajuste endpoints + envelope | 📋 Pendente — Alta |
| 017 | Mapping telegram_id → api_key | 📋 Pendente — Alta |
