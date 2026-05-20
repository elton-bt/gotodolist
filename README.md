# gotodolist

Aplicacao didatica de lista de tarefas para demonstrar deploy de uma aplicacao web em Go com PostgreSQL. O repositorio cobre dois cenarios:

- `monolito/`: aplicacao Go com renderizacao server-side.
- `desacoplado/backend/` + `desacoplado/frontend/`: API REST em Go com frontend estatico.

## Estrutura do repositorio

```text
.
├── desacoplado/
│   ├── backend/
│   └── frontend/
├── docs/
├── internal/
├── monolito/
├── docker-compose-dev.yaml
├── docker-compose-monolito-dev.yaml
├── docker-compose-monolito-prod.yaml
├── docker-compose-prod.yaml
├── example.env
└── .github/workflows/
```

## Requisitos

- Git para clonar o repositorio.
- Go na versao definida em `go.mod` para executar pela raiz do codigo.
- PostgreSQL 16+ ou 17+ para execucao sem Docker.
- Docker Engine e Docker Compose plugin para execucao conteinerizada.

## Como executar

1. Clone o repositorio:

```bash
git clone https://github.com/elton-bt/gotodolist.git
```

2. Entre na pasta do projeto:

```bash
cd gotodolist
```

3. Copie o arquivo de exemplo e ajuste as variaveis conforme o ambiente:

```bash
cp example.env .env
```

4. Escolha um dos modos abaixo.

### Monolito com Docker Compose em desenvolvimento

```bash
docker compose --env-file .env -f docker-compose-monolito-dev.yaml up --build
```

Aplicacao: `http://localhost:8080`

### Frontend + API com Docker Compose

```bash
docker compose --env-file .env -f docker-compose-dev.yaml up --build
```

Frontend: `http://localhost:8082`

API: `http://localhost:8081`

Se o frontend precisar apontar para outra API, ajuste `GOTODOLIST_API_BASE_URL` no arquivo `.env` com um endereco acessivel pelo navegador.

### Imagens do GHCR em ambiente de deploy

```bash
docker compose --env-file .env -f docker-compose-prod.yaml up -d
```

### Monolito com imagem publicada no GHCR

```bash
docker compose --env-file .env -f docker-compose-monolito-prod.yaml up -d
```

Os arquivos de producao nao usam mais `profiles`. Cada Compose representa um unico modo de execucao.

## Documentacao

- [Guia de execucao local](docs/getting-started.md)
- [Compose e imagens do GHCR](docs/docker.md)
- [Referencia da API](docs/api.md)
- [CI/CD e releases](docs/ci-cd.md)

## Testes

```bash
go test ./...
```

## Seguranca

- Nao registre segredos reais em codigo, compose, Dockerfile ou README.
- As mensagens de erro para indisponibilidade do banco sao genericas por desenho.
- O projeto nao imprime a senha do banco.
