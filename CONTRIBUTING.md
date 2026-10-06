# Guia de contribuição

Este projeto usa o GitHub Flow e a convenção de commits do tipo Conventional Commits. O objetivo é manter a qualidade, rastreabilidade e segurança das alterações.

## Fluxo de trabalho

1. Crie uma branch a partir da `main` atualizada.
2. Use nomes descritivos com prefixos como:
   - `feat/` para novas funcionalidades
   - `fix/` para correções de bugs
   - `docs/` para documentação
   - `ci/` para ajustes de pipeline
   - `chore/` para manutenção e governança
3. Faça commits pequenos e executáveis.
4. Abra um Pull Request com contexto técnico claro.
5. Espere a revisão e os checks de CI antes do merge.

## Padrão de commits

As mensagens devem seguir o formato:

```text
<tipo>(<escopo opcional>): <descrição curta em imperativo>
```

Exemplos válidos:

```text
feat(api): adiciona endpoint de listagem de tarefas
fix(docker): corrige versão do compose para ambiente local
ci(actions): ajusta permissões do workflow de validação
chore(governance): adiciona templates de issue e pull request
```

Tipos esperados:

- `feat`
- `fix`
- `docs`
- `refactor`
- `ci`
- `chore`

## Requisitos antes do PR

- Execute os testes locais relevantes.
- Verifique se o código está formatado com `gofmt`.
- Valide os arquivos Docker Compose com `docker compose config`.
- Não inclua segredos, tokens, senhas ou arquivos sensíveis.
- Adicione ou atualize a documentação quando a alteração impactar uso ou deploy.

## Pull Request

O template do PR está em `.github/pull_request_template.md` e deve ser preenchido com:

- descrição da alteração
- issue relacionada
- tipo de mudança
- evidências de validação
- checklist de segurança e governança

## Revisão e merge

- O merge só deve ocorrer após CI verde.
- Alterações críticas em infraestrutura, Docker, Compose e workflows exigem revisão cuidadosa.
- Mantenha o histórico limpo e a comunicação técnica clara.

## Segurança

- Ninguém deve commitar credenciais, .env reais ou secrets em texto puro.
- Antes de abrir o PR, confirme que as dependências e imagens utilizadas não introduziram vulnerabilidades críticas conhecidas.
- Sempre prefira o princípio do menor privilégio em workflows e automações.
