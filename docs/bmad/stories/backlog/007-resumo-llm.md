# Story 007: Resumo em Linguagem Natural (Claude API)

**Como** investidor
**Quero** receber um resumo da minha carteira escrito em linguagem natural
**Para** entender rapidamente o que aconteceu sem interpretar números brutos

### Critérios de Aceitação

- [ ] AC1: `GET /portfolio/summary/narrative?period=daily|weekly|monthly` retorna texto em prosa
- [ ] AC2: O texto destaca: top 3 ativos positivos, ativos em prejuízo, variação total da carteira
- [ ] AC3: Tom objetivo e direto (ex: "Sua carteira valorizou 2,3% na semana. BBSE3 foi o destaque com +5,1%...")
- [ ] AC4: Resposta inclui `telegram_id` e `narrative` (string)
- [ ] AC5: Falha na API Claude → retorna erro 500, não silencia

### Edge Cases
- Carteira vazia → narrativa indica que não há ativos cadastrados
- Todos os ativos em queda → destaca os menos negativos
- API Claude com timeout → propagar erro com mensagem clara

### Fora de escopo
- Recomendações de compra/venda
- Previsões de mercado
- Narrativa em outros idiomas (PT-BR apenas)

### Notas
- Requer `ANTHROPIC_API_KEY` nas variáveis de ambiente
- Prompt deve ser construído com dados reais da carteira (period-summary + performance)
- Custo por chamada à API Claude deve ser monitorado
- Validar com usuário real antes de investir em otimização de prompt
