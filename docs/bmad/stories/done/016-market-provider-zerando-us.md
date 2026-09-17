# Story 016: Market provider retorna 0 para tickers US em chamadas concorrentes

**Como** usuário
**Quero** que preços de ativos US sejam sempre corretos, mesmo quando várias features acessam o provider em paralelo
**Para** não ver relatórios com AAPL/MSFT valendo R$ 0 e narrativas com perda de -100%

### Contexto

Descoberto durante validação da Story 015 (2026-09-17).

O endpoint `GET /portfolio/summary/narrative` dispara 4 goroutines simultâneas em `narrative_aggregator.go` chamando:
- `getSummary.Execute` → provider
- `getPerformance.Execute` → provider
- `getBenchmark.Execute` → provider (com IBOV+SPX)
- `getAllocation.Execute` → provider

Cada uma chama `marketProvider.GetByTickers` internamente. Twelve Data (US) tem rate limit apertado no plano free (8 req/min). Quando 2+ goroutines batem em ~50ms, algumas requisições retornam `0.00` para tickers US, sem erro.

### Sintoma observado

Chamada A (isolada, `GET /portfolio/performance`):
```
AAPL: current_price=337, current_value=3370 ✓
MSFT: current_price=497.75, current_value=2488.75 ✓
```

Chamada B (via narrative, 4 goroutines paralelas):
```
AAPL: current_price=0, current_value=0, return_percentage=-100 ❌
MSFT: current_price=0, current_value=0, return_percentage=-100 ❌
```

Resultado: LLM gera "carteira perdeu 78% na semana" — corretamente refletindo os dados que recebeu, mas os dados estão errados.

### Hipóteses (ordem de investigação)

1. **Rate limit do Twelve Data** → resposta vazia sem HTTP 429 → provider retorna 0 em vez de erro
2. **Timeout curto** no cliente HTTP quando concorrente
3. **Cache/dedup ausente** — cada goroutine faz sua própria chamada em vez de compartilhar resultado

### Critérios de Aceitação

- [x] AC1: Twelve Data provider agora loga status HTTP + body em caso de erro; propaga erro em vez de retornar quote vazia
- [x] AC2: Provider retorna erro explícito para HTTP != 200 (revelou que SPX/S&P500 requer plano pago — 404)
- [x] AC3: `marketProviderWithFallback` usa `singleflight.Group` com chave composta (`ticker:market`, ordenada) — 4 goroutines paralelas do narrative agora resultam em 1 request HTTP compartilhado
- [x] AC4: 6 testes em `provider_test.go` cobrem: chave estável, market na chave, dedup concorrente (10 goroutines → 1 call), chaves diferentes não compartilham, liberação após conclusão, propagação de erro
- [x] AC5: Validado end-to-end contra Groq real — narrative retorna totais corretos (R$ 12.330,25 = valor real da carteira, retorno +44,21%)

### Status: ✅ Concluído

**Solução em 3 camadas:**

1. **Fail-loud no Twelve Data** (`internal/infra/market/twelvedata.go`) — refactor de `fetchBatch` + `parseSingleQuote` + `parseBatchQuotes`; erros de HTTP status, decode e parse agora propagam com contexto (URL, body truncado).

2. **Single-flight no provider** (`internal/infra/market/provider.go`) — `NewMarketProviderWithFallback` agora é pointer receiver com `*singleflight.Group`. Função `buildAssetsKey` ordena e concatena `ticker:market`. Log `singleflight_shared` quando dedupe ocorre.

3. **Benchmark tolerante em `aggregate()`** (`internal/usecase/portifolio/narrative_aggregator.go`) — falha do benchmark não quebra a narrativa (é dado enriquecedor, não crítico). LLM gera resumo sem seção de benchmark; ainda respeita `totals` corretos.

**Descoberta colateral:** SPX (S&P500) no Twelve Data plano free retorna **HTTP 404** com mensagem sobre upgrade. Bug estava mascarado pelo `continue` que engolia erros. Solução do benchmark é tratada como "não-fatal" no narrative. Se quisermos S&P500 de volta: story 017 (trocar provider para SPY ETF ou usar Yahoo).

### Fora de escopo

- Trocar Twelve Data por outro provider
- Cache persistente de quotes (fora do escopo da story)

### Estimativa

**Média** — investigação + provider refactor + testes de race.
