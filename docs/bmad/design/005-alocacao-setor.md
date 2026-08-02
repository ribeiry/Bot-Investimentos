# Design: Story 005 — Alocação por Setor

## Novos endpoints

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/portfolio/allocation` | Alocação agrupada por setor (com lazy update de sector) |

`POST /portfolio/assets` e `GET /portfolio/assets` já existem — recebem/retornam o novo campo `sector`.

---

## Request / Response

**POST /portfolio/assets** (sector opcional):
```json
{ "ticker": "BBSE3", "market": "B3", "quantity": 100, "average_price": 38.50, "sector": "Financeiro" }
{ "ticker": "AAPL",  "market": "NYSE", "quantity": 10,  "average_price": 150.00 }
```

**GET /portfolio/allocation**:
```json
{
  "telegram_id": "123456",
  "data": [
    {
      "sector": "Financeiro",
      "total_value": 10000.00,
      "percentage": 45.5,
      "tickers": ["BBSE3", "ITSA4"]
    },
    {
      "sector": "Technology",
      "total_value": 8000.00,
      "percentage": 36.4,
      "tickers": ["AAPL"]
    },
    {
      "sector": "Outros",
      "total_value": 2000.00,
      "percentage": 18.1,
      "tickers": ["XPTO3"]
    }
  ]
}
```

---

## Schema de banco

```sql
-- Migration: adiciona coluna sector (nullable)
ALTER TABLE portfolio ADD COLUMN sector VARCHAR(50);
-- Idempotente via hasColumn(db, "portfolio", "sector")
```

### Novo método no AssetRepository
```go
UpdateSector(userID int64, ticker string, sector string) error
```

---

## Novo provider: MarketProfileProvider

### Interface (domain)
```go
type MarketProfileProvider interface {
    GetSector(ticker string, market string) (string, error)
}
```

### Implementações (infra/market)
- `brapiProfileProvider` — chama `GET /api/quote/{ticker}?fundamental=true&token=...`
- `twelveDataProfileProvider` — chama `GET /profile?symbol={ticker}&apikey=...`
- `marketProfileWithFallback` — roteia B3 → brapi, NYSE/NASDAQ → twelvedata

---

## Fluxo do GET /portfolio/allocation

```
1. ReturnAllPortfolio(userID)
2. Para cada asset com sector == "":
   a. GetSector(ticker, market)   ← chama API
   b. UpdateSector(userID, ticker, sector)  ← persiste no banco
3. GetByTickers(assets)           ← preços atuais
4. Agrupa por sector, calcula percentagem
5. Retorna
```

---

## Fluxo do POST /portfolio/assets

```
1. Valida input
2. sector informado? → usa o valor do usuário
3. sector vazio?
   a. GetSector(ticker, market)
   b. Se erro ou vazio → sector = "" (NULL no banco)
4. Upsert(userID, asset)
```

---

## Novos arquivos

```
internal/domain/allocation.go                    ← SectorAllocation struct
internal/domain/marketProfileRepository.go       ← MarketProfileProvider interface
internal/infra/market/brapi_profile.go           ← brapiProfileProvider
internal/infra/market/twelvedata_profile.go      ← twelveDataProfileProvider
internal/infra/market/profile_provider.go        ← marketProfileWithFallback
internal/usecase/portifolio/get_allocation.go
internal/usecase/portifolio/get_allocation_test.go
internal/mocks/MarketProfileProvider.go
```

## Impacto em arquivos existentes

| Arquivo | Mudança |
|---|---|
| `internal/domain/asset.go` | Adiciona `Sector string` |
| `internal/domain/assetRepository.go` | Adiciona `UpdateSector(userID int64, ticker, sector string) error` |
| `internal/infra/db/migrations.go` | `migratePortfolioAddSector()` + chamada em `RunMigrations` |
| `internal/adapter/repository/portfolio_repo.go` | `sector` no INSERT/UPDATE/SELECT + implementa `UpdateSector` |
| `internal/usecase/portifolio/upsert_asset.go` | Aceita `sector`, chama profile provider se vazio |
| `internal/adapter/http/portfolio_handler.go` | Adiciona `GetAllocation` |
| `internal/mocks/AssetRepository.go` | Adiciona mock de `UpdateSector` |
| `cmd/api/server.go` | Wiring do profile provider + rota `GET /portfolio/allocation` |

---

## Decisões técnicas

- **Lazy update**: setor só é buscado quando `NULL` e na chamada de `GET /portfolio/allocation` — sem jobs em background
- **Setor do usuário é imutável via API**: se o usuário informou o setor manualmente, nunca é sobrescrito
- **Falha na API de perfil não quebra a chamada**: ativo simplesmente fica em "Outros"
- **String livre no banco**: sem enum, sem constraint — validação apenas no use case se necessário
- **Roteamento igual ao market provider**: B3 → brapi, NYSE/NASDAQ → twelvedata
