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
├── docker-compose-monolito.yaml
├── docker-compose-prod.yaml
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

3. Escolha um dos modos abaixo.

### Monolito com Docker Compose

```bash
docker compose -f docker-compose-monolito.yaml up --build
```

Aplicacao: `http://localhost:8080`

### Frontend + API com Docker Compose

```bash
docker compose -f docker-compose-dev.yaml up --build
```

Frontend: `http://localhost:8082`

API: `http://localhost:8081`

Se o frontend precisar apontar para outra API, passe a variavel `GOTODOLIST_API_BASE_URL` com um endereco acessivel pelo navegador:

```bash
GOTODOLIST_API_BASE_URL=http://192.168.0.20:8081 docker compose -f docker-compose-dev.yaml up --build
```

### Imagens do GHCR em ambiente de deploy

```bash
IMAGE_TAG=latest docker compose -f docker-compose-prod.yaml --profile desacoplado up -d
```

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
