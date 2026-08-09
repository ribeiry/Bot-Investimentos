# Story 013: GET /users/me

**Como** usuário autenticado
**Quero** consultar meus próprios dados
**Para** validar minha identidade, saber meu `id`, `telegram_id` e data de cadastro

### Critérios de Aceitação

- [ ] AC1: `GET /users/me` retorna dados do usuário autenticado
- [ ] AC2: Rota protegida — exige `X-API-Key` válido
- [ ] AC3: Response inclui: `id`, `telegram_id`, `name`, `created_at`
- [ ] AC4: **NÃO** retorna `api_key` (usuário já tem a própria; não precisa expor)
- [ ] AC5: Response segue envelope `{telegram_id, data}`
- [ ] AC6: Sem `X-API-Key` → 401 unauthorized (comportamento do middleware existente)

### Request / Response

```
GET /users/me
X-API-Key: <sua-key>
```

```json
{
  "telegram_id": "123456",
  "data": {
    "id": 1,
    "telegram_id": "123456",
    "name": "João",
    "created_at": "2026-08-09T19:00:00Z"
  }
}
```

### Fora de escopo
- `GET /users` (listar todos) — exporia api_keys de outros usuários
- `PATCH /users/me` (atualizar dados) — story separada se necessário
- `DELETE /users/me` (excluir conta) — story separada com cascade em portfolio/alerts

### Notas
- Reaproveita middleware de auth existente (contexto já tem `user_id`)
- Repository já tem `FindByAPIKey` — só falta expor via handler
- Novo use case: `GetUserUseCase.Execute(userID)`
- Um novo método no repo pode ser necessário: `FindByID(userID)` ou reaproveitar dados do middleware

### Estimativa
- Baixa — endpoint simples, sem lógica de negócio nova
