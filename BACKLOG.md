# 🗂 Backlog — Portfolio Monitoring Bot

Detalhamento das fases de desenvolvimento com priorização por ICE Score (Impact × Confidence × Ease, 1-10 cada).

> **Nota de confiança:** priorização baseada em análise de discovery a partir do README/roadmap — sem validação externa ainda. Confidence tende a ser baixo até haver uso real de usuários.

---

## Fase 1 — MVP ✅ Concluída

| Feature | Status |
|---|---|
| API Go + n8n + Telegram + SQLite | ✅ Live |
| Consulta de cotações em tempo real (B3 + NYSE/NASDAQ) | ✅ Live |
| Relatório diário automático (seg-sex, 18:30) | ✅ Live |
| Comandos Telegram: `/status`, `/ativos`, `/resumo`, `/performance` | ✅ Live |
| Performance por ativo (valor, lucro/prejuízo, retorno %) | ✅ Live |

---

## Fase 2 — Deploy 🚧 Em Andamento

| Feature | Status | Descrição |
|---|---|---|
| Suporte multiusuário | ✅ Concluído | Carteira isolada por usuário via API key; `POST /users` para registro; auth middleware com lookup no DB |
| PostgreSQL | 📋 Backlog | Migração de SQLite para Postgres |
| Cloud deployment | 📋 Backlog | Rodar em produção (AWS/GCP/Heroku) |

---

## Fase 3 — Alertas 🚧 Em Andamento

Priorizadas por ICE Score — revisar a cada sprint.

| # | Feature | Impact | Confidence | Ease | Score | Status |
|---|---|---|---|---|---|---|
| 1 | **Alertas de preço (stop gain/loss)** | 9 | 5 | 7 | **315** | ✅ Concluído |
| 2 | **Resumo semanal/mensal consolidado** | 5 | 4 | 8 | **160** | ✅ Concluído |
| 3 | **Comparação com benchmark (IBOV/S&P500)** | 7 | 3 | 5 | **105** | 📋 Backlog |
| 4 | **Alocação por classe/setor** | 6 | 3 | 6 | **108** | 📋 Backlog |
| 5 | **Notificação de dividendos** | 6 | 3 | 4 | **72** | 📋 Backlog |

### Detalhamento Fase 3

**1. Alertas de preço (stop gain/loss)** ✅ Concluído
- `POST /alerts` — criar/atualizar alerta (stop_gain e/ou stop_loss)
- `GET /alerts` — listar alertas ativos do usuário
- `DELETE /alerts/:ticker` — remover alerta
- `GET /alerts/check` — verifica preços atuais e retorna alertas disparados (chamar via n8n)
- Tabela `alerts` com isolamento por `user_id`

**2. Resumo semanal/mensal consolidado** ✅ Concluído
- `GET /portfolio/period-summary?period=weekly` — variação da carteira desde segunda-feira
- `GET /portfolio/period-summary?period=monthly` — variação da carteira desde dia 1 do mês
- Busca preço histórico via tabela `price_history`; fallback para `average_price` se sem dados
- Retorna por ativo: `price_start`, `price_current`, `change_value`, `change_percent`
- Retorna totais: `total_value_start`, `total_value_current`, `total_change_value`, `total_change_percent`

**3. Comparação com benchmark**
- Buscar IBOV (B3) e S&P500 (NYSE/NASDAQ) nas mesmas APIs
- Calcular retorno relativo: (carteira - benchmark) / benchmark × 100
- Novo endpoint: `GET /portfolio/benchmark`

**4. Alocação por classe/setor**
- Campo novo no schema do ativo: `sector` (enum: "Tech", "Financeiro", "Energia", etc.)
- Novo endpoint: `GET /portfolio/allocation?group_by=sector`
- Requer: catalogar todos os ativos suportados com setor correspondente

**5. Notificação de dividendos**
- Verificar se Brapi/Twelve Data expõem calendar de proventos
- Se sim: job diário que checa datas de ex-dividendo
- Se não: manual/feedback loop

---

## Fase 4 — LLM 🤖 Backlog

| # | Feature | Impact | Confidence | Ease | Score | Justificativa |
|---|---|---|---|---|---|---|
| 1 | **Resumo em linguagem natural (Claude API)** | 8 | 5 | 5 | **200** | Já está no roadmap, traduz dados em insight, reaproveita dados existentes |
| 2 | **Simulação "e se"** | 8 | 2 | 3 | **48** | Alto valor percebido, mas muito especulativo — validar com usuários primeiro |

### Detalhamento Fase 4

**1. Resumo em linguagem natural (Claude API)**
- Integração com Claude API para gerar resumo em prosa
- Exemplo de prompt: "Resuma a performance de hoje em 2-3 frases. Destaque os top 3 ativos e aqueles em prejuízo."
- Requer: chave da API Anthropic, middleware de chamada na rotina de relatório
- Benefício: usuários recebem insight direto em vez de só números

**2. Simulação "e se" (recálculo hipotético)** ✅ Concluído
- `POST /portfolio/simulate` com lista de operações hipotéticas (sell/buy)
- Preço buscado do `price_history` primeiro; API só para tickers sem histórico
- Retorna `current` + `simulated` + `delta` + `skipped` (com reason)
- Operações inválidas (insufficient quantity, quote unavailable) são skipped com aviso

---

## Priorização Consolidada

**Próximo (Fase 2):**
1. PostgreSQL — pré-requisito para produção real
2. Cloud deployment — após migração de banco

**Fase 3 — próximo:**
1. ✅ Alertas de preço — concluído
2. ✅ Resumo semanal/mensal — concluído
3. Comparação com benchmark (próxima prioridade)
4. Alocação por setor (requer catalogação)
5. Dividendos (validar disponibilidade de dados)

**Fase 4 — após validação com usuários reais:**
1. Resumo em LLM
2. ~~Simulação "e se"~~ ✅ Concluído
3. ~~Concorrência no market provider~~ ✅ Concluído (story 011)

---

## Como usar este backlog

- Atualizar ICE Scores a cada sprint com feedback real de usuários
- Confidence cresce quando o bot é usado; começar com 2-3 features de Fase 3
- Revisar Fase 4 após primeira geração de uso real
