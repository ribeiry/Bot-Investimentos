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

## Fase 3 — Alertas 📋 Backlog

Priorizadas por ICE Score — revisar a cada sprint.

| # | Feature | Impact | Confidence | Ease | Score | Justificativa |
|---|---|---|---|---|---|---|
| 1 | **Alertas de preço (stop gain/loss, variação %)** | 9 | 5 | 7 | **315** | Uso diário, reaproveita price_history existente, já está no roadmap |
| 2 | **Resumo semanal/mensal consolidado** | 5 | 4 | 8 | **160** | Extensão natural do relatório diário, reaproveita lógica com cron diferente |
| 3 | **Comparação com benchmark (IBOV/S&P500)** | 7 | 3 | 5 | **105** | Contextualiza retorno absoluto, requer chamada de índice extra |
| 4 | **Alocação por classe/setor** | 6 | 3 | 6 | **108** | Ajuda decisão de realocação, exige tag setor no ativo |
| 5 | **Notificação de dividendos** | 6 | 3 | 4 | **72** | Reforça avaliação de retorno, depende de Brapi/Twelve Data |

### Detalhamento Fase 3

**1. Alertas de preço (stop gain/loss, variação %)**
- Notificação via Telegram quando ativo atinge thresholds definidos pelo usuário
- Requer: novo comando `/alerta` para cadastrar, storage de limites no BD, check na tarefa de atualização de preços
- Risco: spam se alertas forem muito frequentes

**2. Resumo semanal/mensal consolidado**
- Agregar dados de múltiplos dias num só relatório de tendência
- Requer: job adicional no n8n (segunda-feira para semana, 1º dia mês para mês)
- Custo: baixo, reaproveita queries de `portfolio/summary` existentes

**3. Comparação com benchmark**
- Buscar IBOV (B3) e S&P500 (NYSE/NASDAQ) nas mesmas APIs
- Calcular retorno relativo: (carteira - benchmark) / benchmark × 100
- Novo comando Telegram: `/benchmark`

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

**2. Simulação "e se" (recálculo hipotético)**
- Simular venda de ativo X e compra de ativo Y, refazer todo o cálculo de performance
- Muito especulativo — com 0 validação externa
- Prioridade baixa até feedback real de usuário

---

## Priorização Consolidada

**Próximo (Fase 2):**
1. PostgreSQL — pré-requisito para produção real
2. Cloud deployment — após migração de banco

**Fase 3 — ordem recomendada:**
1. Alertas de preço (maior ROI, reaproveita infra)
2. Resumo semanal/mensal (fácil, incrementa valor)
3. Comparação com benchmark (contexto importante)
4. Alocação por setor (requer catalogação)
5. Dividendos (validar disponibilidade de dados)

**Fase 4 — após validação com usuários reais:**
1. Resumo em LLM
2. Simulação "e se"

---

## Como usar este backlog

- Atualizar ICE Scores a cada sprint com feedback real de usuários
- Confidence cresce quando o bot é usado; começar com 2-3 features de Fase 3
- Revisar Fase 4 após primeira geração de uso real
