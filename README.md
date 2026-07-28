# 📈 Portfolio Monitoring Bot

Bot de monitoramento de carteira de investimentos integrado ao Telegram, construído com Go, n8n e Clean Architecture.

---

## 🚀 Funcionalidades

- **Multiusuário** — cada usuário com carteira isolada via API key
- Consulta de cotações em tempo real (B3 via Brapi, NYSE/NASDAQ via Twelve Data)
- Cache de preços via SQLite para consultas instantâneas
- Relatório diário automático às 18:30 (seg-sex)
- 4 comandos via Telegram: `/status`, `/ativos`, `/resumo`, `/performance`
- Cálculo de performance por ativo (valor investido, valor atual, lucro/prejuízo, retorno %)
- Consolidação da carteira por mercado com Top Gainers e Top Losers
- **Alertas de preço** — stop gain e stop loss por ativo, verificação sob demanda
- **Resumo semanal/mensal** — variação da carteira na semana ou no mês via price_history

---

## 🛠 Stack

| Componente | Tecnologia |
|---|---|
| API | Go 1.23 + Gin |
| Banco | SQLite (local) → PostgreSQL (produção) |
| Harness | n8n |
| Mensageria | Telegram Bot API |
| Mercado B3 | Brapi.dev |
| Mercado NYSE/NASDAQ | Twelve Data |
| Cache | price_history (SQLite) |
| Orquestração | Docker Compose |

---

## 🏗 Arquitetura

```
Telegram → n8n (Harness) → API Go → SQLite
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
    ├── infra/            ← SQLite, Brapi, Twelve Data
    │   ├── db/
    │   └── market/
    └── mocks/            ← mocks gerados para testes
```

---

## ⚙️ Configuração

### 1. Pré-requisitos

- Go 1.23+
- Docker e Docker Compose
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
```

> **Nota:** não há `API_KEY` global — cada usuário tem sua própria chave gerada via `POST /users`.

### 3. Rodar localmente

```bash
cd portifolio-api
go run ./cmd/api
```

O banco SQLite é criado automaticamente em `./data/portfolio.db`. As migrations rodam na inicialização.

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

#### Portfolio

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/portfolio/assets` | Lista ativos do usuário |
| `POST` | `/portfolio/assets` | Cadastra ou atualiza ativo |
| `DELETE` | `/portfolio/assets/:ticker` | Remove ativo |
| `GET` | `/portfolio/summary?mode=cached` | Resumo consolidado (cache) |
| `GET` | `/portfolio/summary?mode=realtime` | Resumo consolidado (tempo real) |
| `GET` | `/portfolio/performance` | Performance detalhada por ativo |
| `GET` | `/portfolio/period-summary?period=weekly` | Variação da carteira na semana |
| `GET` | `/portfolio/period-summary?period=monthly` | Variação da carteira no mês |

#### Market

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/market/prices` | Preços atuais de todos os ativos |
| `GET` | `/market/close?market=B3` | Fechamento por mercado (B3, NYSE, NASDAQ) |

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
{
  "id": 1,
  "api_key": "a3f8c2d1...",
  "name": "João"
}
```

### 2. Usar a API key

```bash
curl http://localhost:8080/portfolio/assets \
  -H "X-API-Key: a3f8c2d1..."
```

---

## 📋 Exemplos de uso

### Cadastrar ativo

```bash
curl -X POST http://localhost:8080/portfolio/assets \
  -H "Content-Type: application/json" \
  -H "X-API-Key: <sua-key>" \
  -d '{"ticker":"BBSE3","market":"B3","quantity":100,"average_price":38.50}'
```

### Resumo semanal

```bash
curl "http://localhost:8080/portfolio/period-summary?period=weekly" \
  -H "X-API-Key: <sua-key>"
```

```json
{
  "telegram_id": "123456",
  "data": {
    "period": "weekly",
    "period_start": "2026-07-21T00:00:00Z",
    "assets": [
      {
        "ticker": "BBSE3",
        "price_start": 40.00,
        "price_current": 42.00,
        "change_value": 200.00,
        "change_percent": 5.0
      }
    ],
    "total_change_value": 200.00,
    "total_change_percent": 5.0
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

# Coverage detalhado por pacote
go test ./internal/usecase/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Coverage atual

| Pacote | Testes |
|---|---|
| `usecase/alert` | 17 testes |
| `usecase/market` | 8 testes |
| `usecase/portifolio` | 26 testes |
| `usecase/user` | 4 testes |
| `adapter/http` | 32 testes |
| `adapter/http/middleware` | 4 testes |
| **Total** | **91 testes** |

---

## 🗺 Roadmap

| Fase | Status | Descrição |
|---|---|---|
| Fase 1 — MVP | ✅ Concluída | API Go + n8n + Telegram + SQLite |
| Fase 2 — Deploy | 🚧 Em andamento | Multiusuário ✅ · PostgreSQL 📋 · Cloud 📋 |
| Fase 3 — Alertas | 🚧 Em andamento | Stop gain/loss ✅ · Resumo semanal/mensal ✅ · Benchmark 📋 |
| Fase 4 — LLM | 📋 Backlog | Integração Claude API, resumos inteligentes |

Detalhes e priorização: [BACKLOG.md](./BACKLOG.md)
