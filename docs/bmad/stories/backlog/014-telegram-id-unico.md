# Story 014: Validação de telegram_id único na criação de usuário

**Como** operador do sistema
**Quero** que o `telegram_id` seja único entre usuários
**Para** evitar carteiras duplicadas para o mesmo usuário do Telegram

### Contexto

Hoje `POST /users` permite criar múltiplos usuários com o mesmo `telegram_id`. Isso gera:
- Carteiras duplicadas
- Múltiplas api_keys para o mesmo humano
- Confusão na entrega de notificações via n8n

### Critérios de Aceitação

- [ ] AC1: Schema `users` tem índice `UNIQUE` em `telegram_id`
- [ ] AC2: `POST /users` com `telegram_id` já existente → retorna **409 Conflict**
- [ ] AC3: Response de conflito: `{"error": "telegram_id já cadastrado"}`
- [ ] AC4: Validação acontece no use case (antes do INSERT), não só no DB
- [ ] AC5: Migration idempotente — pode rodar em base já populada
- [ ] AC6: Testes cobrem: criação normal, duplicata, telegram_id vazio

### Request / Response

**Sucesso (201):**
```
POST /users
{"telegram_id": "123456", "name": "João"}
```
```json
{"id": 1, "api_key": "abc...", "name": "João"}
```

**Duplicata (409):**
```
POST /users
{"telegram_id": "123456", "name": "Outro João"}
```
```json
{"error": "telegram_id já cadastrado"}
```

### Design técnico

**Migration:**
```sql
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);
```

**Use case:**
```go
func (u CreateUserUseCase) Execute(user domain.User) (*domain.User, error) {
    existing, err := u.repo.FindByTelegramID(user.TelegramID)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, domain.ErrTelegramIDAlreadyExists
    }
    return u.repo.Create(user)
}
```

**Repository:** novo método `FindByTelegramID(telegramID string) (*domain.User, error)`

**Handler:** mapeia `ErrTelegramIDAlreadyExists` → HTTP 409

### Edge Cases
- Base populada com duplicatas pré-existentes → migration precisa detectar antes de aplicar (opção: deduplicar mantendo o mais antigo, ou falhar com aviso)
- Race condition entre validação e INSERT → constraint do DB protege (retorna erro que deve ser mapeado para 409 também)

### Fora de escopo
- Deduplicação retroativa de dados existentes (fresh start, sem dados hoje)
- Atualizar `telegram_id` de usuário existente
- Merge de carteiras duplicadas

### Estimativa
- Baixa/Média — 1 migration + 1 método repo + 1 validação use case + testes
