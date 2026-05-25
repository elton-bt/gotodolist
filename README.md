# gotodolist

Aplicação didática de lista de tarefas para demonstrar deploy de uma aplicação web em Go com PostgreSQL. O repositório cobre dois cenários:

- `monolito/`: aplicação Go com renderização server-side.
- `desacoplado/backend/` + `desacoplado/frontend/`: API REST em Go com frontend estático.

## Estrutura do repositório

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

- Git para clonar o repositório.
- Go na versão definida em `go.mod` (1.26.0) para executar pela raiz do código.
- PostgreSQL 16+ ou 17+ para execução sem Docker.
- Docker Engine e Docker Compose plugin para execução conteinerizada.

## Como executar

1. Clone o repositório:

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

### Monolito com imagem publicada no GHCR

```bash
docker compose --env-file .env -f docker-compose-monolito-prod.yaml up -d
```
Aplicação: `http://localhost:8080` por padrão. Ajuste `MONOLITO_HOST_PORT` no `.env` se quiser trocar a porta publicada no host.

### Frontend + API com Docker Compose em desenvolvimento

```bash
docker compose --env-file .env -f docker-compose-dev.yaml up --build
```

### Frontend + API com imagem publicada no GHCR

```bash
docker compose --env-file .env -f docker-compose-prod.yaml up -d
```

Frontend: `http://localhost:8082` por padrão.

API: `http://localhost:8081` por padrão.

Se o frontend precisar apontar para outra API, ou se `API_HOST_PORT` for alterada, ajuste `GOTODOLIST_API_BASE_URL` no arquivo `.env` com um endereço acessível pelo navegador.

## Documentação

- [Guia de execução local](docs/getting-started.md)
- [Compose e imagens do GHCR](docs/docker.md)
- [Referência da API](docs/api.md)
- [CI/CD e releases](docs/ci-cd.md)

## Testes

```bash
go test ./...
```

## Segurança

- Nao registre segredos reais em código, compose, Dockerfile ou README.
- As mensagens de erro para indisponibilidade do banco são genéricas por design.
- O projeto não imprime a senha do banco.
