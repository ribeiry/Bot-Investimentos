# Story 008: Migração para PostgreSQL

**Como** operador do sistema
**Quero** migrar o banco de dados de SQLite para PostgreSQL
**Para** suportar múltiplos usuários simultâneos em produção com segurança e performance

### Critérios de Aceitação

- [ ] AC1: Toda a aplicação funciona com PostgreSQL sem alteração de lógica de negócio
- [ ] AC2: Migrations existentes reescritas em sintaxe PostgreSQL
- [ ] AC3: `docker-compose.yml` inclui serviço PostgreSQL com healthcheck
- [ ] AC4: Variável `DATABASE_URL` configura a conexão (com default local)
- [ ] AC5: SQLite **removido** do projeto (fresh start, paridade dev/prod)
- [ ] AC6: Todos os 143 testes continuam passando sem alteração (todos usam mocks)

### Fora de escopo
- Migração de dados existentes (fresh start em produção)
- Connection pooling avançado (pgBouncer)
- Read replicas
- Testes de integração com Postgres real (não há testes de repository)

---

## 🏛 Design Técnico (Architect)

### Decisão: driver único, sem dualidade

SQLite abandonado completamente. Manter dois drivers dobra código sem ganho — objetivo é produção.

- `mattn/go-sqlite3` (CGO) sai → build vira **pure Go**
- Dev local usa Postgres via `docker-compose`
- Trade-off: perde `go run` sem Docker. Ganho: **paridade dev/prod**

### Driver escolhido: `pgx/v5/stdlib`

Mais performático que `lib/pq`, compatível com `database/sql`, ativamente mantido.

### Diferenças críticas de sintaxe

| Concept | SQLite | PostgreSQL |
|---------|--------|-----------|
| Auto-increment PK | `INTEGER PRIMARY KEY` | `SERIAL PRIMARY KEY` |
| Placeholders | `?` | `$1, $2, $3...` |
| Boolean default | `DEFAULT 1` | `DEFAULT TRUE` |
| Introspection | `PRAGMA table_info(x)` | `information_schema.columns` |
| Insert retornando ID | `LastInsertId()` | `RETURNING id` + `QueryRow().Scan(&id)` |
| Datetime default | `DEFAULT CURRENT_TIMESTAMP` | ✅ igual |
| VARCHAR / DECIMAL | ✅ igual | ✅ igual |
| ON CONFLICT | ✅ igual | ✅ igual |

### Onde SQLite está acoplado hoje

| Arquivo | Alteração |
|---|---|
| `infra/db/connection.go` | Driver hardcoded + path hardcoded — rewrite completo |
| `infra/db/migrations.go` | `PRAGMA`, `INTEGER PK`, `BOOLEAN DEFAULT 1`, workarounds de recriação — rewrite completo |
| `adapter/repository/portfolio_repo.go` | `?` → `$N` |
| `adapter/repository/alert_repo.go` | `?` → `$N` |
| `adapter/repository/user_repo.go` | `?` → `$N` + `LastInsertId()` → `RETURNING id` |
| `adapter/repository/price_history_repo.go` | `?` → `$N` |
| `docker-compose.yml` | +serviço postgres, +healthcheck, -volume sqlite |
| `go.mod` | -`mattn/go-sqlite3`, +`jackc/pgx/v5` |
| `README.md` | Novas instruções: subir via docker-compose, exige Postgres |

**13 queries afetadas** ao total nos 4 repositories.

### Novo `connection.go`

```go
package db

import (
    "database/sql"
    "os"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func InitDatabase() (*sql.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://portfolio:portfolio@localhost:5432/portfolio?sslmode=disable"
    }
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    return db, db.Ping()
}
```

### Novo `migrations.go`

- `INTEGER PRIMARY KEY` → `SERIAL PRIMARY KEY`
- `PRAGMA table_info` → query em `information_schema.columns`
- `BOOLEAN DEFAULT 1` → `BOOLEAN DEFAULT TRUE`
- Remover workarounds de recriação de tabela (Postgres suporta `ALTER TABLE ADD CONSTRAINT`)
- Como não migramos dados existentes, `migratePortfolio*` viram `CREATE TABLE` completo com todas as colunas

### `user_repo.go`: LastInsertId → RETURNING

```go
// antes
result, err := r.db.Exec("INSERT INTO users (...) VALUES (?, ?, ?)", ...)
id, err := result.LastInsertId()

// depois
var id int64
err := r.db.QueryRow(
    "INSERT INTO users (...) VALUES ($1, $2, $3) RETURNING id", ...,
).Scan(&id)
```

### `docker-compose.yml`

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=portfolio
      - POSTGRES_PASSWORD=portfolio
      - POSTGRES_DB=portfolio
    volumes:
      - postgres-data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U portfolio"]
      interval: 5s
      retries: 5
    networks:
      - portfolio-network

  portifolio-api:
    build: ...
    environment:
      - DATABASE_URL=postgres://portfolio:portfolio@postgres:5432/portfolio?sslmode=disable
      ...
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres-data:  # substitui sqlite-data
```

### Testes

**Todos os 143 testes usam mocks das interfaces** (`AssetRepository`, `UserRepository`, etc.) — nunca tocam `*sql.DB`. Zero testes de repository. Migração **não afeta testes**.

- `go test ./...` roda sem Docker antes e depois
- Docker necessário apenas para rodar a aplicação e teste manual com curl

### Ordem de implementação

1. Adicionar serviço Postgres no `docker-compose.yml`, subir e validar
2. Trocar driver em `connection.go`, adicionar `DATABASE_URL`
3. Reescrever `migrations.go`
4. Reescrever queries em `repository/` (placeholders + RETURNING)
5. Rodar `docker-compose up` e testar todos os endpoints manualmente
6. `go test ./...` — 143 testes devem passar sem alteração
7. Atualizar `README.md`
8. Remover `mattn/go-sqlite3` do `go.mod`

### Estimativa

- Complexidade: média
- Arquivos alterados: ~10
- Novo código: ~150 linhas
- Testes afetados: zero

### Riscos

| Risco | Mitigação |
|-------|-----------|
| Data em SQLite atual perdido | Aceito — story diz "fresh start" |
| Erro em query Postgres só em runtime | Teste manual completo antes de encerrar |
| Devs sem Docker travam | Documentar setup no README |
