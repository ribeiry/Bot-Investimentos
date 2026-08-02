# Architecture — Portfolio Monitoring Bot

## Visão Geral

```
Telegram ←→ n8n ←→ API Go (portifolio-api) ←→ SQLite
                           ↓
                     Brapi (B3)
                     Twelve Data (NYSE/NASDAQ)
```

## Clean Architecture

```
portifolio-api/
├── cmd/api/
│   ├── main.go          ← entrypoint
│   └── server.go        ← DI container + roteamento
└── internal/
    ├── domain/          ← entidades, interfaces (zero dependências externas)
    ├── usecase/         ← regras de negócio
    │   ├── alert/
    │   ├── market/
    │   ├── portifolio/
    │   └── user/
    ├── adapter/
    │   ├── http/        ← handlers Gin + response envelope
    │   │   └── middleware/
    │   └── repository/  ← implementações SQLite
    ├── infra/
    │   ├── db/          ← connection + migrations
    │   └── market/      ← brapi, twelvedata, cached, fallback
    └── mocks/           ← mocks testify para testes
```

## Padrões estabelecidos

### Response envelope
Toda resposta inclui `telegram_id` para o n8n rotear:
```json
{ "telegram_id": "123456", "data": { ... } }
{ "telegram_id": "123456", "message": "..." }
{ "telegram_id": "123456", "error": "..." }
```

### Autenticação
- `POST /users` → cria usuário, retorna `api_key` gerada com `crypto/rand` (64 chars hex)
- Todas as rotas protegidas validam `X-API-Key` header via `middleware.Auth`
- Middleware injeta `userID` e `telegramID` no gin context

### Isolamento multiusuário
- Toda query de portfolio filtra por `user_id`
- Alertas também isolados por `user_id`
- `price_history` é compartilhada (dados de mercado, não de usuário)

### Migrations
- Migrations rodam no boot via `db.RunMigrations`
- Idempotentes: `CREATE TABLE IF NOT EXISTS` + `PRAGMA table_info` para ALTER TABLE
- UNIQUE constraints criadas como índices separados (limitação SQLite)

### Market providers
```
marketProviderWithFallback
├── B3 → brapiProvider
├── NYSE/NASDAQ → twelveDataProvider
└── Salva tudo em price_history (usado por cachedProvider e period-summary)
```

## Banco de dados (SQLite)

```sql
users         (id, telegram_id, name, api_key, created_at)
portfolio     (id, user_id, ticker, market, quantity, average_price, created_at)
              UNIQUE INDEX (user_id, ticker)
price_history (id, ticker, price, captured_at)
alerts        (id, user_id, ticker, market, stop_gain, stop_loss, active, created_at)
              UNIQUE INDEX (user_id, ticker)
```

## Endpoints

| Método | Rota | Auth |
|---|---|---|
| POST | /users | público |
| GET | /health | público |
| GET/POST/DELETE | /portfolio/assets | sim |
| GET | /portfolio/summary | sim |
| GET | /portfolio/performance | sim |
| GET | /portfolio/period-summary | sim |
| GET | /portfolio/benchmark | sim |
| GET | /market/prices | sim |
| GET | /market/close | sim |
| GET/POST/DELETE | /alerts | sim |
| GET | /alerts/check | sim |

## Cobertura de testes

108 testes — use cases, handlers HTTP e middleware todos cobertos com mocks testify.
