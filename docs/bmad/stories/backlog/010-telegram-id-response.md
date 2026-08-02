# Story 010: Devolução do telegram_id em Toda Interação

**Como** sistema de automação (n8n)
**Quero** receber o `telegram_id` do usuário em toda resposta da API
**Para** saber para qual chat Telegram enviar a mensagem, sem precisar armazenar estado entre chamadas

### Critérios de Aceitação

- [x] AC1: Toda resposta de endpoint protegido inclui `telegram_id` no body
- [x] AC2: O envelope é consistente em todos os cenários: sucesso, erro e mensagem
- [x] AC3: Respostas de sucesso com dados seguem: `{ "telegram_id": "...", "data": ... }`
- [x] AC4: Respostas de mensagem seguem: `{ "telegram_id": "...", "message": "..." }`
- [x] AC5: Respostas de erro seguem: `{ "telegram_id": "...", "error": "..." }`
- [x] AC6: Endpoint público `POST /users` **não** inclui `telegram_id` (usuário ainda não autenticado)

### Edge Cases
- Usuário autenticado com `telegram_id` vazio → retorna string vazia (não quebra)
- Endpoint de health check → sem `telegram_id` (público)

### Fora de escopo
- Cache de `telegram_id` no cliente
- Múltiplos chats por usuário

### Notas
- Motivação: n8n é **efêmero** — cada execução de workflow começa do zero sem memória de execuções anteriores
- O `telegram_id` é injetado no contexto pelo middleware de autenticação e extraído pelos response helpers
- Implementado via `internal/adapter/http/response.go` com helpers `respond`, `respondMessage`, `respondError`
- `middleware/auth.go` injeta `telegramID` no gin context após validar a API key

### Status: ⚠️ Implementado no código — aguarda validação com n8n em ambiente real
