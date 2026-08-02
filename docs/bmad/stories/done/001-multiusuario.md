# Story 001: Suporte Multiusuário

**Como** investidor
**Quero** ter minha carteira isolada dos outros usuários
**Para** que meus dados sejam privados e cada pessoa gerencie só os seus ativos

### Critérios de Aceitação
- [x] AC1: `POST /users` cria usuário e retorna `api_key` única
- [x] AC2: Toda rota protegida exige `X-API-Key` válida
- [x] AC3: Operações de portfolio filtram por `user_id`
- [x] AC4: Dois usuários podem ter o mesmo ticker sem conflito

### Status: ✅ Concluído
