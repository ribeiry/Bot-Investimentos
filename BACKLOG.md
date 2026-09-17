# 🗂 Backlog — Portfolio Monitoring Bot

Detalhamento das fases de desenvolvimento com priorização por ICE Score (Impact × Confidence × Ease, 1-10 cada).

> **Fonte da verdade:** as **stories BMAD** em `docs/bmad/stories/` são o backlog operacional (com critérios de aceite, contexto e edge cases). Este `BACKLOG.md` é a visão executiva/roadmap e deve ser sincronizado a cada story concluída.
>
> **Nota de confiança:** ICE original baseado em discovery a partir do README/roadmap, sem validação externa. Confidence de features **entregues** foi recalibrado; das ainda não entregues permanece baixo até haver uso real.

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
| Suporte multiusuário | ✅ Concluído | Carteira isolada por usuário via API key; `POST /users` para registro; auth middleware com lookup no DB (story `001`) |
| PostgreSQL | ✅ Concluído | Migração SQLite → PostgreSQL via `pgx/v5`; docker-compose com healthcheck (story `008`) |
| Cloud deployment | 📋 Backlog | Rodar em produção — ver quebra em **Fase 5** |

---

## Fase 3 — Alertas & Análise ✅ Concluída

Priorizadas por ICE Score original — mantidas para histórico. Confidence recalibrada pós-entrega.

| # | Feature | Impact | Confidence | Ease | Score | Status |
|---|---|---|---|---|---|---|
| 1 | **Alertas de preço (stop gain/loss)** | 9 | 9 | 7 | **567** | ✅ Concluído |
| 2 | **Resumo semanal/mensal consolidado** | 5 | 9 | 8 | **360** | ✅ Concluído |
| 3 | **Comparação com benchmark (IBOV/S&P500)** | 7 | 9 | 5 | **315** | ✅ Concluído |
| 4 | **Alocação por classe/setor** | 6 | 9 | 6 | **324** | ✅ Concluído |
| 5 | **Notificação de dividendos** | — | — | — | — | ❌ Descartado |

### Detalhamento Fase 3

**1. Alertas de preço** ✅ (story `002`)
- `POST /alerts`, `GET /alerts`, `DELETE /alerts/:ticker`, `GET /alerts/check`
- Tabela `alerts` isolada por `user_id`.

**2. Resumo semanal/mensal** ✅ (story `003`)
- `GET /portfolio/period-summary?period=weekly|monthly`.
- Busca via `price_history`; fallback para `average_price`.
- Retorna por ativo e totais consolidados.

**3. Benchmark IBOV / S&P500** ✅ (story `004`)
- `GET /portfolio/benchmark?period=weekly|monthly`.
- Retorna `portfolio_return_percent` + benchmarks (`price_start`, `price_current`, `return_percent`, `relative_performance`).
- Sem histórico → flag `no_history: true`; período inválido → 400.

**4. Alocação por setor** ✅ (story `005`)
- `POST /portfolio/assets` aceita `sector` opcional; autopreenche via Brapi/Twelve Data quando ausente.
- `GET /portfolio/allocation` agrupa por setor com `total_value`, `percentage`, `tickers`.
- `PATCH /portfolio/assets/:ticker/sector` para correção manual.
- Ativos sem setor após tentativa de fetch → `"Outros"`.

**5. Notificação de dividendos** ❌ Descartado
- Discovery: **Brapi e Twelve Data não expõem calendar de proventos no plano gratuito**.
- Alternativas pagas fora do orçamento atual; input manual quebra a proposta de valor (automática).
- Reavaliar se surgir provider gratuito ou se houver mudança de modelo de custo.

---

## Fase 4 — LLM 🤖 Concluída

| # | Feature | Impact | Confidence | Ease | Score | Status |
|---|---|---|---|---|---|---|
| 1 | **Resumo em linguagem natural (Groq)** | 8 | 9 | 5 | **360** | ✅ Concluído |
| 2 | **Simulação "e se"** | 8 | 9 | 6 | **432** | ✅ Concluído |

### Detalhamento Fase 4

**1. Resumo LLM via Groq** ✅ (story `007`)
- **Decisão técnica:** Groq (não Claude, conforme roadmap inicial). Plano free: 30 req/min, 14.400 req/dia. Modelo default `llama-3.3-70b-versatile`.
- Endpoint: `GET /portfolio/summary/narrative`.
- Envelope: `{"telegram_id": "...", "data": {"text": "..."}}`.
- Rate limit usuário: 5/dia → **HTTP 429** loud. Rate limit global: 25/min + 12k/dia → fallback silencioso.
- Deadline 500ms: se LLM excede, retorna "⏳ Processando..." e cacheia em background.
- Validador heurístico anti-alucinação: narrativas com ticker fora do portfólio são descartadas.
- Cache in-memory 1h por `user_id`; carteira vazia → texto fixo sem consumir LLM.
- Todos os 13 ACs entregues; 11 cenários de teste (use case + handler) passando com `-race`.

**2. Simulação "e se"** ✅ (story `009`)
- `POST /portfolio/simulate` com lista de operações hipotéticas (sell/buy).
- Preço via `price_history` primeiro; API só para tickers sem histórico.
- Retorna `current` + `simulated` + `delta` + `skipped` (com `reason`).

---

## Fase 5 — Produção & Qualidade 📋 Backlog

Novos itens levantados após revisão do estado atual do repo. Sem ICE ainda — a definir.

| # | Feature | Motivação |
|---|---|---|
| 1 | **Cloud deployment** | Colocar API em produção (AWS/GCP/Fly.io/Railway). Sub-tarefas: escolha do provider, secrets manager, domínio, CI/CD, deploy pipeline |
| 2 | **`.env.example`** | Hoje `.env` é desstrackeado e não há template — friction para novos ambientes |
| 3 | **CI (GitHub Actions)** | Rodar `go test ./...` + lint em cada PR |
| 4 | **Versionamento de migrations** | Migrations rodam via `db.RunMigrations` no boot, sem histórico; migrar para `golang-migrate` ou similar |
| 5 | **Observabilidade** | Logs estruturados (zap/slog), métricas Prometheus, healthcheck detalhado |
| 6 | **Documentação de rate limits** | Rate limit já existe em código mas não está documentado como comportamento externo (README/API docs) |
| 7 | **`Makefile` / task runner** | Padronizar comandos (`make test`, `make run`, `make lint`) |

---

## Melhorias de UX/API (concluídas)

Micro-stories que polem contrato da API — histórico consolidado.

| Story | Descrição | Status |
|---|---|---|
| `010` | Envelope `{telegram_id, data\|message\|error}` em toda resposta protegida | ✅ |
| `011` | Concorrência no market provider (Brapi + Twelve Data em paralelo) | ✅ |
| `013` | `GET /users/me` (sem expor `api_key`) | ✅ |
| `014` | `telegram_id` único no `POST /users` (409 em duplicata) | ✅ |
| `015` | Fix alucinação de totais no resumo LLM (campo `totals` pré-calculado) | ✅ |
| `016` | Fix zero em tickers US no narrative (single-flight + fail-loud + benchmark não-fatal) | ✅ |

> Gap de numeração: stories `006` e `012` não existem no diretório BMAD — confirmar se foram descartadas ou renumeradas.

### Bugs em backlog

_(vazio)_

---

## Priorização Consolidada (próximas sprints)

1. **Fase 5 #1 — Cloud deployment** — pré-requisito para uso real e para validar Confidence do backlog.
2. **Fase 5 #2/3/7 — `.env.example`, CI, Makefile** — quick wins de higiene que destravam contribuições.
3. **Fase 5 #4/5 — Migrations versionadas + observabilidade** — antes ou junto do go-live.

---

## Como usar este backlog

- **Ao concluir uma story BMAD:** mover de `docs/bmad/stories/backlog/` → `done/` **e** atualizar a linha correspondente aqui.
- **A cada sprint:** revisar ICE Scores com feedback de uso real; Confidence sobe conforme validação empírica.
- **Débito técnico:** entra em Fase 5 ou em nova seção dedicada — não misturar com features de produto.
- **Fase 4 (LLM):** revisitar prioridade e prompt tuning **após** primeira geração de uso real do resumo narrativo.
