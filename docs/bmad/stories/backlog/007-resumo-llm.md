# Story 007: Resumo em linguagem natural via LLM (Groq)

**Como** usuário do bot no Telegram
**Quero** um resumo em linguagem natural do meu portfólio
**Para** entender rapidamente como está minha carteira sem interpretar tabelas

### Provider

**Groq** — plano free (30 req/min, 14.400 req/dia).
Modelo default: `llama-3.3-70b-versatile`.

### Critérios de Aceitação

- [ ] AC1: `GET /portfolio/summary/narrative` autenticado
- [ ] AC2: Agrega dados de 4 use cases (summary + performance + benchmark + allocation) em paralelo
- [ ] AC3: Chama Groq com prompt carregado em memória no boot
- [ ] AC4: Envelope `{"telegram_id": "...", "data": {"text": "..."}}`
- [ ] AC5: Cache in-memory 1h por `user_id`
- [ ] AC6: Rate limit por usuário: 5/dia → HTTP **429** loud
- [ ] AC7: Rate limit global: 25/min + 12.000/dia → fallback silencioso
- [ ] AC8: Deadline de resposta 500ms — se LLM demora mais, retorna fallback amigável mas processa em background e cacheia
- [ ] AC9: Validação heurística: resposta com ticker fora do portfólio → descartada + fallback
- [ ] AC10: Carteira vazia → texto fixo, sem consumir LLM
- [ ] AC11: Prompt em `config/llm_prompt.txt`, carregado 1x no boot (`log.Fatal` se ausente)
- [ ] AC12: Testes unitários com mocks (padrão limpo, sem `r0/r1/rf`)
- [ ] AC13: Log estruturado por chamada (cache hit/miss, timeout, validação)

### Env vars

```
GROQ_API_KEY=gsk_...                     # obrigatória
GROQ_MODEL=llama-3.3-70b-versatile       # opcional
LLM_RESPONSE_DEADLINE_MS=500             # opcional
LLM_USER_RATE_LIMIT_PER_DAY=5            # opcional
LLM_RATE_LIMIT_PER_MINUTE=25             # opcional
LLM_RATE_LIMIT_PER_DAY=12000             # opcional
```

### Fluxo do use case

```
1. userLimiter.Allow(userID)  → false → ErrRateLimitExceeded (429)
2. cache.Get(userID)          → hit → retorna
3. agrega dados (4 goroutines)
4. tickers vazios             → texto fixo "Sua carteira está vazia..."
5. globalLimiter.Allow()      → false → fallback silencioso
6. dispara LLM em goroutine (contexto Background, timeout 10s)
7. select:
   ├── LLM responde <500ms → valida + cacheia (se válido) + retorna
   └── deadline 500ms → retorna "⏳ Estamos processando..."
                        (goroutine continua, cacheia no fim)
```

### Camadas de defesa contra alucinação

| Camada | Contra | Custo |
|---|---|---|
| Prompt restritivo | Alucinação, off-topic | 0 |
| `temperature=0.3`, `max_tokens=500` | Prolixidade, criatividade | 0 |
| Validação heurística de tickers | Ticker fabricado | ~30 linhas |
| Rate limit user (5/dia loud) | Abuso individual | ~40 linhas |
| Rate limit global (silent) | Estouro de quota | ~50 linhas |
| Cache 1h in-memory | Consumo desnecessário | ~40 linhas |
| Carteira vazia sem LLM | Consumo desnecessário | ~5 linhas |

### Prompt template (`config/llm_prompt.txt`)

```
Você é um assistente financeiro conciso que fala português brasileiro.
Sua ÚNICA função é gerar um resumo do portfólio do usuário.

Regras OBRIGATÓRIAS:
1. Use APENAS os números fornecidos abaixo. NUNCA invente valores.
2. Não faça recomendações de compra ou venda.
3. Não comente sobre economia, política, notícias ou eventos externos.
4. Não responda perguntas — apenas gere o resumo.
5. Se os dados estiverem incompletos, mencione apenas o que está disponível.
6. Máximo 150 palavras. Tom neutro e informativo.
7. Não invente tickers ou setores que não estejam nos dados.
8. Use emojis moderadamente (📊 📈 📉 💰).
9. Destaque valores em **negrito**.

Dados do portfólio:
{{PORTFOLIO_DATA}}

Gere o resumo agora, apenas com base nos dados acima.
```

### Estrutura de arquivos

**Novos:**
```
portifolio-api/
├── config/
│   └── llm_prompt.txt
└── internal/
    ├── domain/
    │   └── llm.go                                  # LLMProvider + ErrRateLimitExceeded
    ├── infra/
    │   ├── llm/groq.go
    │   ├── cache/narrative_cache.go
    │   └── ratelimit/
    │       ├── daily_limiter.go                    # por user
    │       └── global_limiter.go                   # global (min+dia)
    ├── usecase/portifolio/
    │   ├── get_narrative.go
    │   ├── get_narrative_test.go
    │   ├── narrative_validator.go
    │   └── narrative_validator_test.go
    └── mocks/
        ├── LLMProvider.go
        ├── NarrativeCache.go
        ├── RateLimiter.go
        └── GlobalRateLimiter.go
```

**Modificados:**
- `cmd/api/server.go` — carrega prompt + wire das novas deps
- `internal/adapter/http/portfolio_handler.go` — handler `GetNarrative`
- `internal/adapter/http/portfolio_handler_test.go` — testes handler

### Handler

```go
func (h *PortfolioHandler) GetNarrative(c *gin.Context) {
    userID := c.GetInt64("userID")
    telegramID := c.GetString("telegramID")

    text, err := h.getNarrativeUseCase.Execute(c.Request.Context(), userID)
    if errors.Is(err, domain.ErrRateLimitExceeded) {
        c.JSON(http.StatusTooManyRequests, gin.H{
            "telegram_id": telegramID,
            "error":       "limite diário atingido — tente novamente amanhã",
        })
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "telegram_id": telegramID,
            "error":       err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "telegram_id": telegramID,
        "data":        gin.H{"text": text},
    })
}
```

### Testes obrigatórios (`get_narrative_test.go`)

1. Success_CacheMiss_LLMRapido — narrativa real, cache gravado
2. Success_CacheHit — não chama LLM
3. UserRateLimit_Retorna429
4. CarteiraVazia_TextoFixo_SemLLM
5. GlobalRateLimit_FallbackSilent
6. LLMTimeout500ms_FallbackProcessando_CacheiaDepois
7. LLMErro_Fallback
8. RespostaComTickerInvalido_Fallback_NaoCacheia
9. Aggregate falha (uma das 4 fontes) → Fallback

### Ordem de implementação

1. Prompt + config loader (`os.ReadFile` no boot)
2. Interface `LLMProvider` + Groq client (teste manual curl)
3. Cache in-memory + `narrative_validator` + testes
4. `dailyLimiter` (user) + `globalLimiter` + testes
5. Mocks
6. `GetNarrativeUseCase` + testes (todos os 9 cenários)
7. Handler + rota
8. Wire em `server.go`
9. `go test ./...` + `-race`
10. Teste manual end-to-end (2-3 chamadas)
11. Atualiza README

### Fora de escopo

- Streaming
- Multi-linguagem
- Recomendações de investimento
- Input livre do usuário
- Redis (in-memory serve para MVP)

### Estimativa

**Alta** — ~10 arquivos novos, cache, rate limits, goroutines, testes.
