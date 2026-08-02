# Story 006: Notificação de Dividendos

**Como** investidor
**Quero** ser notificado quando um ativo da minha carteira tiver data de ex-dividendo próxima
**Para** não perder o direito ao provento e planejar minha posição

### Critérios de Aceitação

- [ ] AC1: `GET /portfolio/dividends` retorna próximos dividendos dos ativos da carteira (próximos 30 dias)
- [ ] AC2: Retorna por ativo: ticker, ex_date, payment_date, dividend_per_share, estimated_total (quantity × dividend_per_share)
- [ ] AC3: Apenas ativos com ex_date futura são retornados
- [ ] AC4: Resposta inclui `telegram_id`
- [ ] AC5: Sem dividendos próximos → `data: []`

### Edge Cases
- Ativo não tem dados de dividendo disponíveis na API → ignorado silenciosamente
- API de dividendos indisponível → 500 com mensagem de erro

### Fora de escopo
- Histórico de dividendos recebidos
- Cálculo de yield on cost
- Reinvestimento automático (DRIP)

### Notas
- Depende de disponibilidade na Brapi (B3) e Twelve Data (NYSE/NASDAQ)
- Validar suporte antes de implementar — pode precisar de API alternativa
- Dados de dividendos podem não estar disponíveis para todos os ativos
