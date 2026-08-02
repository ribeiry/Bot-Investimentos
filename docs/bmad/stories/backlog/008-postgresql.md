# Story 008: Migração para PostgreSQL

**Como** operador do sistema
**Quero** migrar o banco de dados de SQLite para PostgreSQL
**Para** suportar múltiplos usuários simultâneos em produção com segurança e performance

### Critérios de Aceitação

- [ ] AC1: Toda a aplicação funciona com PostgreSQL sem alteração de lógica de negócio
- [ ] AC2: Migrations existentes adaptadas para PostgreSQL (sintaxe, tipos, índices)
- [ ] AC3: `docker-compose.yml` inclui serviço PostgreSQL
- [ ] AC4: Variável `DATABASE_URL` configura a conexão
- [ ] AC5: SQLite mantido como opção para desenvolvimento local (via variável de ambiente)
- [ ] AC6: Todos os 108+ testes continuam passando

### Edge Cases
- Conflito de sintaxe SQLite vs PostgreSQL (ex: `AUTOINCREMENT` vs `SERIAL`)
- `ON CONFLICT` syntax differences
- Tipos de dados: `BOOLEAN`, `DECIMAL`, `DATETIME` → equivalentes Postgres

### Fora de escopo
- Migração de dados existentes (fresh start em produção)
- Connection pooling avançado (pgBouncer)
- Read replicas

### Notas
- Usar `lib/pq` ou `pgx` como driver
- Abstrair conexão em `infra/db/connection.go` para suportar ambos os drivers
- Testar com Docker localmente antes de deploy
