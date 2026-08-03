# Story 011: Concorrência no Market Provider

**Como** sistema
**Quero** buscar cotações de B3 e NYSE/NASDAQ em paralelo
**Para** reduzir o tempo de resposta em endpoints que buscam ativos de múltiplos mercados

### Critérios de Aceitação

- [ ] AC1: `marketProviderWithFallback` busca B3 (Brapi) e USA (Twelve Data) em goroutines paralelas
- [ ] AC2: Resultados são mergeados corretamente após ambas as goroutines finalizarem
- [ ] AC3: Erro em qualquer uma das goroutines é propagado corretamente
- [ ] AC4: Sem race conditions — acesso aos dados compartilhados protegido com `sync.Mutex`
- [ ] AC5: Todos os testes existentes continuam passando
- [ ] AC6: Sleep de 61s entre lotes do Twelve Data permanece (limite da API)

### Edge Cases
- Apenas ativos B3 → só goroutine Brapi, sem overhead
- Apenas ativos USA → só goroutine Twelve Data, sem overhead
- Erro no Brapi com sucesso no Twelve Data → propaga erro do Brapi
- Erro no Twelve Data com sucesso no Brapi → propaga erro do Twelve Data

### Fora de escopo
- Paralelismo dentro do Brapi (per-ticker) — risco de rate limit
- Paralelismo dentro do Twelve Data — respeita rate limit existente
- Cache de cotações entre requests

### Notas
- Mudança em um único arquivo: `internal/infra/market/provider.go`
- Beneficia todos os endpoints com ativos em múltiplos mercados
- Maior impacto visível em `/portfolio/simulate` (mais tickers por chamada)
- Ganho estimado: elimina tempo sequencial — resposta cai de B3+USA para max(B3, USA)
