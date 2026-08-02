# Product Brief — Portfolio Monitoring Bot

## Visão do Produto

Bot de monitoramento de carteira de investimentos multiusuário, integrado ao Telegram via n8n, construído em Go com Clean Architecture.

O produto permite que investidores acompanhem seus ativos, recebam alertas de preço e visualizem performance em tempo real diretamente no Telegram — sem precisar abrir corretoras ou planilhas.

## Usuário-alvo

Investidor pessoa física que opera B3, NYSE e NASDAQ e quer monitoramento passivo automatizado via Telegram.

## Problema resolvido

- Dispersão de informação entre corretoras, planilhas e apps
- Falta de alertas proativos quando um ativo atinge stop gain ou stop loss
- Ausência de contexto histórico (semana, mês) para avaliar se a carteira está performando bem

## Stack atual

| Componente | Tecnologia |
|---|---|
| API | Go 1.23 + Gin |
| Banco | SQLite → PostgreSQL (planejado) |
| Harness | n8n |
| Mensageria | Telegram Bot API |
| Mercado B3 | Brapi.dev |
| Mercado NYSE/NASDAQ | Twelve Data |

## Estado atual

| Fase | Status |
|---|---|
| MVP (API + Telegram + SQLite) | ✅ Concluído |
| Multiusuário (API key por user) | ✅ Concluído |
| Alertas stop gain/loss | ✅ Concluído |
| Resumo semanal/mensal | ✅ Concluído |
| Comparação com benchmark | ✅ Concluído |
| PostgreSQL | 📋 Backlog |
| Cloud deployment | 📋 Backlog |
| Alocação por setor | 📋 Backlog |
| Notificação de dividendos | 📋 Backlog |
| Resumo LLM (Claude API) | 📋 Backlog |
| Simulação "e se" | 📋 Backlog |

## Restrições

- Cada usuário possui carteira isolada por `user_id`
- Toda resposta da API inclui `telegram_id` para o n8n rotear mensagens
- n8n não tem memória entre chamadas — contexto deve vir na resposta da API
- Cobertura de testes obrigatória para cada feature
