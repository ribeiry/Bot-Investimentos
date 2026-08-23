# 📈 Portfolio Monitoring Bot

Bot de monitoramento de carteira de investimentos integrado ao Telegram, construído com Go, n8n e Clean Architecture.

---

## 🚀 Funcionalidades

- **Multiusuário** — cada usuário com carteira isolada via API key
- Consulta de cotações em tempo real (B3 via Brapi, NYSE/NASDAQ via Twelve Data)
- Cache de preços via SQLite para consultas instantâneas
- Relatório diário automático às 18:30 (seg-sex)
- 4 comandos via Telegram: `/status`, `/ativos`, `/resumo`, `/performance`
- Cálculo de performance por ativo (valor investido, valor atual, lucro/prejuízo, retorno % e **setor**)
- Consolidação da carteira por mercado ou por setor com Top Gainers e Top Losers
- **Alertas de preço** — stop gain e stop loss por ativo, verificação sob demanda
- **Resumo semanal/mensal** — variação da carteira na semana ou no mês via price_history
- **Comparação com benchmark** — carteira vs IBOV e S&P500 por período
- **Alocação por setor** — percentual do patrimônio em cada setor (setor informado manualmente no cadastro do ativo)
- **Simulação "e se"** — simula operações hipotéticas (venda/compra) e compara carteira atual vs simulada com delta de retorno
- **Concorrência no provider** — B3 e NYSE/NASDAQ buscados em goroutines paralelas

---

## 🛠 Stack

| Componente | Tecnologia |
|---|---|
| API | Go 1.23 + Gin |
| Banco | PostgreSQL 16 (via `jackc/pgx/v5`) |
| Harness | n8n |
| Mensageria | Telegram Bot API |
| Mercado B3 | Brapi.dev |
| Mercado NYSE/NASDAQ | Twelve Data (`/quote`) |
| Cache | price_history (Postgres) |
| Orquestração | Docker Compose |

---

## 🏗 Arquitetura

```
Telegram → n8n (Harness) → API Go → PostgreSQL
                               ↓
                         Brapi (B3)
                         Twelve Data (NYSE/NASDAQ)
```

### Clean Architecture

```
portifolio-api/
├── cmd/api/
│   ├── main.go
│   └── server.go
└── internal/
    ├── domain/           ← entidades e interfaces
    ├── usecase/          ← regras de negócio
    │   ├── alert/
    │   ├── market/
    │   ├── portifolio/
    │   └── user/
    ├── adapter/          ← handlers HTTP e repositórios
    │   ├── http/
    │   │   └── middleware/
    │   └── repository/
    ├── infra/            ← PostgreSQL, Brapi, Twelve Data
    │   ├── db/
    │   └── market/
    └── mocks/            ← mocks para testes
```

---

## ⚙️ Configuração

### 1. Pré-requisitos

- Go 1.23+
- Docker e Docker Compose (obrigatório para o PostgreSQL)
- Conta no [Brapi](https://brapi.dev) — cotações B3
- Conta no [Twelve Data](https://twelvedata.com) — cotações NYSE/NASDAQ
- Bot Telegram criado via [@BotFather](https://t.me/BotFather)
- [ngrok](https://ngrok.com) — para expor o n8n via HTTPS (desenvolvimento)

### 2. Variáveis de ambiente

Crie um arquivo `.env` na raiz de `portifolio-api/`:

```env
BRAPI_TOKEN=seu-token-brapi
TWELVE_DATA_KEY=seu-token-twelve-data
TELEGRAM_CHAT_ID=seu-chat-id
DATABASE_URL=postgres://portfolio:portfolio@localhost:5432/portfolio?sslmode=disable
```

> **Nota:** não há `API_KEY` global — cada usuário tem sua própria chave gerada via `POST /users`.
> A variável `DATABASE_URL` é **obrigatória** — a API falha ao iniciar se não estiver definida.

### 3. Rodar localmente

Suba o PostgreSQL via Docker:

```bash
cd portifolio-api
docker-compose up -d postgres
```

Rode a API:

```bash
go run ./cmd/api
```

As migrations rodam automaticamente na inicialização e criam as tabelas `users`, `portfolio`, `price_history` e `alerts`.

### 4. Subir via Docker

```bash
docker-compose up --build
```

### 5. Expor o n8n via HTTPS

```bash
ngrok http 5678
```

Copie a URL gerada e configure em `docker-compose.yml`:

```yaml
n8n:
  environment:
    - WEBHOOK_URL=https://xxxx.ngrok-free.app
```

---

## 📡 Endpoints da API

### Público

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/health` | Health check |
| `POST` | `/users` | Registra usuário e recebe API key |

### Protegido — header `X-API-Key: <sua-key>`

> Todas as respostas incluem `telegram_id` para o n8n saber para qual chat enviar.

#### Usuários

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/users/me` | Retorna dados do usuário autenticado (sem `api_key`) |

#### Portfolio

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/portfolio/assets` | Lista ativos (inclui `sector`) |
| `POST` | `/portfolio/assets` | Cadastra ou atualiza ativo (`sector` opcional — informado manualmente) |
| `DELETE` | `/portfolio/assets/:ticker` | Remove ativo |
| `GET` | `/portfolio/summary?mode=cached\|realtime` | Resumo consolidado por mercado |
| `GET` | `/portfolio/summary?mode=realtime&group_by=sector` | Resumo consolidado por setor |
| `GET` | `/portfolio/performance` | Performance por ativo (inclui `sector`) |
| `GET` | `/portfolio/period-summary?period=weekly\|monthly` | Variação da carteira no período |
| `GET` | `/portfolio/benchmark?period=weekly\|monthly` | Carteira vs IBOV e S&P500 |
| `GET` | `/portfolio/allocation` | Alocação por setor com percentual do patrimônio |
| `PATCH` | `/portfolio/assets/:ticker/sector` | Atualiza ou limpa o setor de um ativo |
| `POST` | `/portfolio/simulate` | Simula operações hipotéticas — retorna current, simulated e delta |

#### Market

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/market/prices` | Preços atuais de todos os ativos |
| `GET` | `/market/close?market=B3\|NYSE\|NASDAQ` | Fechamento por mercado |

#### Alertas

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/alerts` | Lista alertas ativos |
| `POST` | `/alerts` | Cria ou atualiza alerta de stop gain/loss |
| `DELETE` | `/alerts/:ticker` | Remove alerta |
| `GET` | `/alerts/check` | Verifica preços e retorna alertas disparados |

> `GET /alerts/check` retorna `200` com lista de triggers ou `200` com `data: []` se nenhum disparado.

---

## 🔑 Autenticação

O sistema é multiusuário. Cada usuário tem sua própria API key isolada.

### 1. Registrar usuário

```bash
curl -s -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"telegram_id": "123456", "name": "João"}' | jq
```

```json
{ "id": 1, "api_key": "a3f8c2d1...", "name": "João" }
```

### 2. Usar a API key

```bash
curl http://localhost:8080/portfolio/assets \
  -H "X-API-Key: a3f8c2d1..."
```

---

## 📋 Exemplos de uso

### Cadastrar ativo com setor

```bash
curl -X POST http://localhost:8080/portfolio/assets \
  -H "Content-Type: application/json" \
  -H "X-API-Key: <sua-key>" \
  -d '{"ticker":"BBSE3","market":"B3","quantity":100,"average_price":38.50,"sector":"Financeiro"}'
```

> `sector` é opcional. Se não informado, o ativo aparece em `"Outros"` no `/allocation`. Nenhuma API retorna setor automaticamente — deve ser informado pelo usuário.

### Performance por ativo (com setor)

```bash
curl http://localhost:8080/portfolio/performance \
  -H "X-API-Key: <sua-key>" | jq
```

```json
{
  "telegram_id": "123456",
  "data": [
    {
      "ticker": "BBSE3",
      "market": "B3",
      "sector": "Financeiro",
      "quantity": 100,
      "average_price": 38.50,
      "current_price": 40.00,
      "invested_value": 3850.00,
      "current_value": 4000.00,
      "profit_loss": 150.00,
      "return_percentage": 3.89
    }
  ]
}
```

### Alocação por setor

```bash
curl http://localhost:8080/portfolio/allocation \
  -H "X-API-Key: <sua-key>" | jq
```

```json
{
  "telegram_id": "123456",
  "data": [
    { "sector": "Financeiro", "total_value": 6400.00, "percentage": 80.0, "tickers": ["BBSE3", "ITSA4"] },
    { "sector": "Technology", "total_value": 1600.00, "percentage": 20.0, "tickers": ["AAPL"] }
  ]
}
```

### Resumo por setor

```bash
curl "http://localhost:8080/portfolio/summary?mode=realtime&group_by=sector" \
  -H "X-API-Key: <sua-key>" | jq
```

### Comparação com benchmark

```bash
curl "http://localhost:8080/portfolio/benchmark?period=weekly" \
  -H "X-API-Key: <sua-key>" | jq
```

```json
{
  "telegram_id": "123456",
  "data": {
    "period": "weekly",
    "portfolio_return_percent": 5.0,
    "benchmarks": [
      { "name": "IBOV",   "return_percent": 4.0,  "relative_performance": 1.0 },
      { "name": "S&P500", "return_percent": 10.0, "relative_performance": -5.0 }
    ]
  }
}
```

### Criar alerta de stop gain e stop loss

```bash
curl -X POST http://localhost:8080/alerts \
  -H "Content-Type: application/json" \
  -H "X-API-Key: <sua-key>" \
  -d '{"ticker":"BBSE3","market":"B3","stop_gain":45.00,"stop_loss":30.00}'
```

### Verificar alertas disparados

```bash
curl http://localhost:8080/alerts/check \
  -H "X-API-Key: <sua-key>"
```

---

## 💬 Comandos Telegram

| Comando | Descrição |
|---|---|
| `/status` | Posição consolidada da carteira |
| `/ativos` | Lista de ativos cadastrados |
| `/resumo` | Resumo por mercado com Top Gainers e Top Losers |
| `/performance` | Rentabilidade detalhada por ativo |

---

## 🧪 Testes

```bash
# Rodar todos os testes
cd portifolio-api
go test ./... -v

# Com coverage
go test ./... -cover
```

### Coverage atual

| Pacote | Testes |
|---|---|
| `usecase/alert` | 17 testes |
| `usecase/market` | 8 testes |
| `usecase/portifolio` | 44 testes |
| `usecase/user` | 4 testes |
| `adapter/http` | 37 testes |
| `adapter/http/middleware` | 4 testes |
| **Total** | **143 testes** |

---

## 🗺 Roadmap

| Fase | Status | Descrição |
|---|---|---|
| Fase 1 — MVP | ✅ Concluída | API Go + n8n + Telegram + SQLite |
| Fase 2 — Deploy | 🚧 Em andamento | Multiusuário ✅ · PostgreSQL ✅ · Cloud 📋 |
| Fase 3 — Alertas | ✅ Concluída | Stop gain/loss ✅ · Resumo semanal/mensal ✅ · Benchmark ✅ · Alocação por setor ✅ · Simulação e-se ✅ |
| Fase 4 — Performance | ✅ Concluída | Concorrência B3+NYSE ✅ · PATCH setor ✅ |
| Fase 5 — LLM | 📋 Backlog | Integração Claude API, resumos inteligentes · PostgreSQL |

Detalhes e priorização: [BACKLOG.md](./BACKLOG.md)
