# Story 015: Fix alucinação de totais no resumo LLM

**Como** usuário do bot
**Quero** que o resumo em linguagem natural traga totais numericamente corretos
**Para** não ser induzido a decisões erradas por dados inventados

### Contexto

Durante validação manual do Path D (2026-09-17), o LLM alucinou totais em 2 de 4 chamadas reais contra Groq:

- **Reportado:** "investimento total R$ 8.550, valor atual R$ 2.430,50, retorno -41,53%"
- **Real:** invested R$ 8.550 ✓, current **R$ 12.330,25**, retorno **+44,2%**

O modelo pegou o valor de um único ativo (PETR4) e reportou como total geral, inclusive invertendo o sinal do retorno.

### Causa raiz

O `portfolioSnapshot` enviado ao prompt contém:
- `summary` — agregado por mercado (`TotalperMarket`), sem total global
- `performance` — array por ativo, com `invested_value`/`current_value`/`profit_loss` individuais
- `benchmark`, `allocation`

**Nenhum campo com totais consolidados globais** (`total_invested`, `total_current`, `total_profit_loss`, `total_return_percentage`). O prompt diz "nunca invente valores" mas o LLM é forçado a somar sozinho e às vezes erra.

O validador atual (`narrative_validator.go`) só checa tickers, não números.

### Critérios de Aceitação

- [x] AC1: `portfolioSnapshot` ganha campo `Totals portfolioTotals` com `Invested`, `Current`, `ProfitLoss`, `ReturnPercentage`
- [x] AC2: Totais são calculados em Go a partir de `performance` (função `computeTotals`)
- [x] AC3: Se `Invested = 0` → `ReturnPercentage = 0` (evita divisão por zero)
- [x] AC4: Prompt (`config/llm_prompt.txt`) instrui o LLM a usar EXCLUSIVAMENTE `totals` (regras 3a/3b)
- [x] AC5: `narrative_aggregator_test.go` — 4 casos cobrem soma, prejuizo, edge cases
- [x] AC6: Performance vazio + invested zero cobertos por testes específicos
- [x] AC7: Testes existentes de `get_narrative_test.go` continuam passando
- [x] AC8: `go test ./... -race` limpo

### Status: ✅ Concluído

**Validação end-to-end com Groq real:** LLM passou a honrar o campo `totals` fielmente. Nenhuma alucinação de valor agregado nas chamadas de teste.

**Descoberta colateral:** durante o teste manual, identificado bug independente no market provider — chamadas concorrentes ao Twelve Data (4 goroutines em `aggregate()`) causam retornos zerados intermitentes para tickers US. Aberto como **Story 016**.

### Fora de escopo

- Validador numérico heurístico (rejeitar narrativa com número que não bate) — vira Story 016 se necessário
- Cache write-through (evitar race de 2 chamadas consecutivas) — vira Story 017

### Estimativa

**Baixa** — ~30 linhas de código + 2 testes + 1 ajuste no prompt.
