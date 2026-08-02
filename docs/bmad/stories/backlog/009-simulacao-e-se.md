# Story 009: Simulação "E Se?" (Recálculo Hipotético)

**Como** investidor
**Quero** simular o impacto de trocar um ativo por outro na minha carteira
**Para** tomar decisões mais informadas antes de executar uma operação real

### Critérios de Aceitação

- [ ] AC1: `POST /portfolio/simulate` recebe lista de operações hipotéticas (venda X, compra Y)
- [ ] AC2: Retorna performance simulada da carteira com as operações aplicadas
- [ ] AC3: Compara resultado simulado vs carteira atual (delta de retorno, delta de valor)
- [ ] AC4: Operações não são persistidas — apenas simuladas
- [ ] AC5: Resposta inclui `telegram_id`

### Edge Cases
- Vender mais do que possui → erro de validação
- Comprar ativo já existente → soma à posição atual no cálculo
- Ativo inválido (sem cotação disponível) → erro com identificação do ticker

### Fora de escopo
- Simulação de rebalanceamento automático
- Otimização de carteira (Markowitz, etc.)
- Persistência da simulação

### Notas
- Feature especulativa — prioridade baixa até validação com usuários reais
- Usar preços atuais da API de mercado para o cálculo
- Não chamar APIs externas em excesso — reutilizar cotações já buscadas na mesma requisição
