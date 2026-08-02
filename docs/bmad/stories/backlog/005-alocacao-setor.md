# Story 005: Alocação por Classe/Setor

**Como** investidor
**Quero** ver quanto do meu patrimônio está alocado em cada setor
**Para** tomar decisões de rebalanceamento e evitar concentração excessiva

### Critérios de Aceitação

- [ ] AC1: `POST /portfolio/assets` aceita campo opcional `sector`
  - Se `sector` vier preenchido → usa o valor informado
  - Se `sector` vier vazio → busca automaticamente via Twelve Data (NYSE/NASDAQ) ou Brapi (B3)
  - Se a API não retornar → grava `NULL`, aparece em "Outros"
- [ ] AC2: `GET /portfolio/assets` retorna `sector` de cada ativo
- [ ] AC3: `GET /portfolio/allocation` retorna alocação agrupada por setor
  - Antes de calcular, verifica ativos com `sector = NULL`
  - Para cada um, busca setor via API e atualiza o banco
  - Calcula alocação com dados atualizados
- [ ] AC4: Cada grupo de setor retorna: `sector`, `total_value`, `percentage`, `tickers`
- [ ] AC5: `percentage` calculada sobre valor atual (quantity × current_price)
- [ ] AC6: Ativos sem setor após tentativa de busca → agrupados em `"Outros"`
- [ ] AC7: Resposta inclui `telegram_id`

### Edge Cases
- Carteira vazia → `data: []`
- API de perfil indisponível → ativo fica em "Outros", sem quebrar a chamada
- Usuário informa setor manualmente → prevalece sobre o da API (não sobrescreve)
- Todos os ativos sem setor e API offline → 100% em "Outros"

### Fora de escopo
- Sub-setores ou hierarquia de setores
- Atualização periódica automática de setores (apenas lazy on demand)
- Tradução dos nomes de setor da API (ex: "Technology" → "Tech")

### Notas
- Setor é buscado **lazy**: apenas quando `NULL` no banco, na chamada de `GET /portfolio/allocation`
- Brapi: `GET /api/quote/{ticker}?fundamental=true` retorna setor
- Twelve Data: `GET /profile?symbol={ticker}` retorna setor
- Setor informado pelo usuário nunca é sobrescrito pela API
