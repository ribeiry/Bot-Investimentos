# Story 012: Brapi Batch Request

**Como** sistema
**Quero** buscar todas as cotações B3 numa única chamada HTTP
**Para** reduzir latência e número de requests à API do Brapi

### Critérios de Aceitação

- [ ] AC1: `brapiProvider.GetByTickers` faz **uma única** chamada HTTP com todos os tickers separados por vírgula: `/api/quote/PETR4,BBSE3,ITSA4?token=xxx`
- [ ] AC2: Response com múltiplos resultados no array `results` é parseada corretamente
- [ ] AC3: Ticker sem cotação na resposta é ignorado silenciosamente (comportamento atual mantido)
- [ ] AC4: Todos os testes existentes continuam passando
- [ ] AC5: Carteira com 1 ticker → continua funcionando (sem regressão)

### Motivação

Atualmente `brapiProvider` faz N chamadas HTTP sequenciais (1 por ticker). Com batch, reduz para 1 chamada independente do tamanho da carteira.

### Mudança

Apenas `internal/infra/market/brapi.go` — sem impacto em interfaces ou handlers.

### Fora de escopo

- Paginação (Brapi suporta até ~50 tickers por chamada — suficiente para o uso atual)
- Goroutines por ticker (batch é superior: 1 round trip vs N paralelos)
