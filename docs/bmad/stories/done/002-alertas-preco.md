# Story 002: Alertas de Preço (Stop Gain / Stop Loss)

**Como** investidor
**Quero** cadastrar alertas de stop gain e stop loss por ativo
**Para** ser notificado via Telegram quando o preço atingir meu limite

### Critérios de Aceitação
- [x] AC1: `POST /alerts` cria/atualiza alerta com stop_gain e/ou stop_loss
- [x] AC2: `GET /alerts` lista alertas ativos do usuário
- [x] AC3: `DELETE /alerts/:ticker` remove alerta
- [x] AC4: `GET /alerts/check` retorna alertas disparados (preço >= stop_gain ou <= stop_loss)
- [x] AC5: Sem triggers → retorna `data: []`
- [x] AC6: Alertas isolados por `user_id`

### Status: ✅ Concluído
