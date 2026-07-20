# 📈 Portfolio Monitoring Bot

Bot de monitoramento de carteira de investimentos integrado ao Telegram, construído com Go, n8n e Clean Architecture.

---

## 🚀 Funcionalidades

- Consulta de cotações em tempo real (B3 via Brapi, NYSE/NASDAQ via Twelve Data)
- Cache de preços via SQLite para consultas instantâneas
- Relatório diário automático às 18:30 (seg-sex)
- 4 comandos via Telegram: `/status`, `/ativos`, `/resumo`, `/performance`
- Cálculo de performance por ativo (valor investido, valor atual, lucro/prejuízo, retorno %)
- Consolidação da carteira por mercado com Top Gainers e Top Losers

---

## 🛠 Stack

| Componente | Tecnologia |
|---|---|
| API | Go 1.23 + Gin |
| Banco | SQLite (Fase 1) → PostgreSQL (Fase 2) |
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
    ├── domain/          ← entidades e interfaces
    ├── usecase/         ← regras de negócio
    │   ├── portfolio/
    │   └── market/
    ├── adapter/         ← handlers HTTP e repositórios
    │   ├── http/
    │   └── repository/
    └── infra/           ← SQLite, Brapi, Twelve Data
        ├── db/
        └── market/
```

---

## ⚙️ Configuração

### 1. Pré-requisitos

- Docker e Docker Compose
- Conta no [Brapi](https://brapi.dev) — cotações B3
- Conta no [Twelve Data](https://twelvedata.com) — cotações NYSE/NASDAQ
- Bot Telegram criado via [@BotFather](https://t.me/BotFather)
- [ngrok](https://ngrok.com) — para expor o n8n via HTTPS (desenvolvimento)

### 2. Variáveis de ambiente

Cria um arquivo `.env` na raiz do projeto:

```env
API_KEY=seu-token-secreto
BRAPI_TOKEN=seu-token-brapi
TWELVE_DATA_KEY=seu-token-twelve-data
TELEGRAM_CHAT_ID=seu-chat-id
```

### 3. Subir os containers

```bash
docker-compose up --build
```

### 4. Expor o n8n via HTTPS

```bash
ngrok http 5678
```

Copie a URL gerada e configure em `docker-compose.yml`:

```yaml
n8n:
  environment:
    - WEBHOOK_URL=https://xxxx.ngrok-free.app
```

### 5. Configurar o n8n

1. Acesse `http://localhost:5678`
2. Crie a credencial do Telegram com o token do BotFather
3. Importe os workflows (consulta sob demanda + fechamento automático)
4. Ative os workflows

---

## 📡 Endpoints da API

| Método | Rota | Descrição |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/portfolio/assets` | Lista todos os ativos |
| POST | `/portfolio/assets` | Cadastra ou atualiza ativo |
| DELETE | `/portfolio/assets/:ticker` | Remove ativo |
| GET | `/portfolio/summary?mode=cached` | Resumo consolidado (cache) |
| GET | `/portfolio/summary?mode=realtime` | Resumo consolidado (tempo real) |
| GET | `/portfolio/performance` | Performance detalhada por ativo |
| GET | `/market/prices` | Preços atuais de todos os ativos |
| GET | `/market/close?market=B3` | Fechamento por mercado |

### Autenticação

Todas as rotas (exceto `/health`) exigem o header:

```
X-API-Key: seu-token-secreto
```

### Exemplo de cadastro de ativo

```bash
curl -X POST http://localhost:8080/portfolio/assets \
  -H "Content-Type: application/json" \
  -H "X-API-Key: seu-token" \
  -d '{"ticker":"BBSE3","market":"B3","quantity":100,"average_price":38.50}'
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
go test ./internal/... -v

# Rodar com coverage
go test ./internal/... -cover

# Ver coverage detalhado
go test ./internal/usecase/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Coverage atual

| Pacote | Coverage |
|---|---|
| `usecase/portfolio` | 100% |
| `usecase/market` | 100% |

---

## 🗺 Roadmap

| Fase | Status | Descrição |
|---|---|---|
| Fase 1 — MVP | ✅ Concluída | API Go + n8n + Telegram + SQLite |
| Fase 2 — Deploy | 🚧 Planejada | PostgreSQL + Multiusuário + Cloud |
| Fase 3 — Alertas | 📋 Backlog | Alertas de preço, dividendos, stop gain/loss |
| Fase 4 — LLM | 🤖 Backlog | Integração Claude API, resumos inteligentes |

---

## 📄 Documentação

- [SDD — Spec-Driven Development](./SDD.md)
- [BACKLOG.md](./BACKLOG.md)