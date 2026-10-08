# Guia de Contribuição do Projeto

Agradecemos pelo interesse em contribuir com nossa infraestrutura e aplicação! Para mantermos a qualidade técnica e a rastreabilidade do projeto, exigimos a conformidade com as regras abaixo.

## 🌿 Fluxo de Branching (GitHub Flow)
1. Nunca realize commits diretamente na branch `main`.
2. Crie uma branch a partir da `main` atualizada utilizando o padrão de nomenclatura:
   - `feat/nome-da-funcionalidade`
   - `fix/correcao-do-problema`
   - `ci/melhoria-no-pipeline`
   - `chore/tarefa-ou-governanca`
   - `docs/ajuste-de-documentacao`

## 💬 Padrão de Commits (Conventional Commits)
Todas as mensagens de commit devem seguir rigorosamente o padrão:
`<tipo>(<escopo opcional>): <descrição no imperativo e em minúsculas>`

Exemplos aceitos:
- `feat(api): adiciona endpoint para listagem de metricas`
- `fix(docker): corrige permissao de usuario non-root no dockerfile`
- `ci(actions): fixa tag imutavel na action do trivy`

## 🧪 Validação Local Antes do PR
```bash
gofmt -l .                 # não deve listar nenhum arquivo
go test ./...              # todos os testes devem passar
docker compose --env-file example.env -f docker-compose-monolito-dev.yaml up --build
curl -i http://localhost:8080/health
```

## 🚀 Submissão de Pull Requests
- Todos os PRs devem preencher o checklist do template oficial.
- O merge só será liberado após:
  1. Todos os status checks do pipeline de CI/CD passarem com sucesso (verde).
  2. Aprovação de pelo menos 1 revisor listado no `CODEOWNERS`.

## 🐛 Abertura de Issues
Utilize sempre os formulários disponíveis em **Issues → New issue** (Relato de Bug ou Sugestão de Funcionalidade). Issues em branco estão desabilitadas.
