# Story 009: Simulação "E Se?" (Recálculo Hipotético)

**Como** investidor
**Quero** simular o impacto de operações hipotéticas na minha carteira
**Para** tomar decisões mais informadas antes de executar uma operação real

### Critérios de Aceitação

- [ ] AC1: `POST /portfolio/simulate` recebe lista de operações hipotéticas
- [ ] AC2: Cada operação tem: `action` (sell|buy), `ticker`, `quantity`, `market` (obrigatório para buy de ativo novo)
- [ ] AC3: Preço de compra é buscado automaticamente via API de mercado (mesmo provider atual)
- [ ] AC4: Retorna carteira atual + carteira simulada + delta entre as duas
- [ ] AC5: Operações não são persistidas — cálculo apenas em memória
- [ ] AC6: Vender mais do que possui → operação marcada como `skipped` com aviso, restante da simulação continua
- [ ] AC7: Comprar ativo já existente na carteira → soma à posição atual no cálculo simulado
- [ ] AC8: Ativo sem cotação disponível → operação marcada como `skipped` com aviso
- [ ] AC9: Resposta inclui `telegram_id`

### Request

```json
POST /portfolio/simulate
{
  "operations": [
    { "action": "sell", "ticker": "BBSE3", "quantity": 50 },
    { "action": "buy",  "ticker": "AAPL",  "quantity": 5, "market": "NYSE" }
  ]
}
```

### Response

```json
{
  "telegram_id": "123456",
  "data": {
    "skipped": [
      { "ticker": "XPTO3", "reason": "insufficient quantity" }
    ],
    "current": {
      "total_invested": 5000.00,
      "total_current_value": 5500.00,
      "total_return_percent": 10.0
    },
    "simulated": {
      "total_invested": 4800.00,
      "total_current_value": 5350.00,
      "total_return_percent": 11.45
    },
    "delta": {
      "value": -150.00,
      "return_percent": 1.45
    }
  }
}
```

### Edge Cases
- Todas as operações skipped → retorna current = simulated, delta zerado, lista de skipped preenchida
- Carteira vazia + só operações de venda → todas skipped
- Buy de ativo novo sem `market` → erro de validação (400)
- Lista de operações vazia → erro de validação (400)

### Fora de escopo
- Simulação de rebalanceamento automático
- Otimização de carteira (Markowitz, etc.)
- Persistência da simulação
- Simulação de dividendos ou proventos

### Notas
- Preço atual buscado via `marketProvider.GetByTickers` — reutilizar cotações numa única chamada
- Carteira simulada é calculada em memória: aplica operações sobre cópia dos assets reais
- `total_invested` na simulação = recalculado com `average_price` dos ativos restantes + preço atual dos comprados
