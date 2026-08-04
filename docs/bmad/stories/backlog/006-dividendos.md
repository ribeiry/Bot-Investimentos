# Story 006: Dividendos Internacionais (NYSE/NASDAQ)

**Como** investidor
**Quero** ver os próximos dividendos dos meus ativos internacionais
**Para** planejar minha posição e não perder o direito ao provento

### Critérios de Aceitação

- [ ] AC1: `GET /portfolio/dividends` retorna próximos dividendos dos ativos NYSE/NASDAQ da carteira
- [ ] AC2: Retorna por ativo: `ticker`, `ex_date`, `amount` (por ação), `estimated_total` (quantity × amount)
- [ ] AC3: Apenas dividendos com `ex_date` futura são retornados
- [ ] AC4: Ativos B3 são ignorados silenciosamente (Brapi não suporta no plano free)
- [ ] AC5: Carteira sem ativos internacionais → `data: []`
- [ ] AC6: Resposta inclui `telegram_id`

### Request / Response

```
GET /portfolio/dividends
```

```json
{
  "telegram_id": "123456",
  "data": [
    {
      "ticker": "AAPL",
      "market": "NASDAQ",
      "quantity": 5,
      "ex_date": "2026-05-11",
      "amount_per_share": 0.27,
      "estimated_total": 1.35
    }
  ]
}
```

### Edge Cases
- Ativo NYSE/NASDAQ sem dividendos disponíveis na Twelve Data → ignorado silenciosamente
- Twelve Data indisponível → 500 com mensagem de erro
- Todos os dividendos com `ex_date` no passado → `data: []`

### Fora de escopo
- Dividendos de ativos B3 (requer plano pago Brapi)
- Histórico de dividendos recebidos
- Cálculo de dividend yield
- Reinvestimento automático (DRIP)

### Notas
- Endpoint Twelve Data: `GET /dividends?symbol=AAPL&apikey=xxx`
- Resposta já validada no plano free: retorna `ex_date` e `amount`
- Filtrar apenas `ex_date >= hoje` no use case
- Nome do endpoint deixa claro o escopo: `/portfolio/dividends` com nota na resposta futuramente
