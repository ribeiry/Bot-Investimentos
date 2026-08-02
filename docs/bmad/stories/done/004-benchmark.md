# Story 004: Comparação com Benchmark (IBOV / S&P500)

**Como** investidor
**Quero** comparar o retorno da minha carteira com o IBOV e o S&P500
**Para** saber se estou batendo o mercado ou ficando abaixo

### Critérios de Aceitação
- [x] AC1: `GET /portfolio/benchmark?period=weekly|monthly` retorna comparação
- [x] AC2: Retorna portfolio_return_percent da carteira no período
- [x] AC3: Retorna por benchmark: name, ticker, price_start, price_current, return_percent, relative_performance
- [x] AC4: relative_performance = portfolio_return - benchmark_return
- [x] AC5: Sem histórico do benchmark → no_history: true
- [x] AC6: Sem histórico da carteira → fallback para average_price
- [x] AC7: Período inválido → 400 Bad Request

### Status: ✅ Concluído
