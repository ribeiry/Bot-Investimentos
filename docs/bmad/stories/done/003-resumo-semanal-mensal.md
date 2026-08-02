# Story 003: Resumo Semanal e Mensal

**Como** investidor
**Quero** ver a variação da minha carteira na semana ou no mês
**Para** entender a tendência recente e não só o retorno total desde a compra

### Critérios de Aceitação
- [x] AC1: `GET /portfolio/period-summary?period=weekly` retorna variação desde segunda-feira
- [x] AC2: `GET /portfolio/period-summary?period=monthly` retorna variação desde dia 1 do mês
- [x] AC3: Retorna por ativo: price_start, price_current, change_value, change_percent
- [x] AC4: Retorna totais consolidados da carteira
- [x] AC5: Sem histórico → usa average_price como fallback
- [x] AC6: Período inválido → 400 Bad Request

### Status: ✅ Concluído
