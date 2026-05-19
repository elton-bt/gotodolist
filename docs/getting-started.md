# Guia de execucao local

Este guia cobre a execucao do projeto a partir do codigo-fonte.

## Requisitos

- Git
- Go na versao definida em `go.mod`
- PostgreSQL 16+ ou 17+

## Clonando o repositorio

```bash
git clone https://github.com/elton-bt/gotodolist.git
cd gotodolist
```

## Variaveis de ambiente

Toda configuracao entra por variaveis de ambiente.

| Variavel | Descricao | Padrao |
| --- | --- | --- |
| `APP_PORT` | Porta HTTP da aplicacao | `8080` no monolito, `8081` na API |
| `DB_HOST` | Host do PostgreSQL | `localhost` |
| `DB_PORT` | Porta do PostgreSQL | `5432` |
| `DB_NAME` | Nome do banco | `gotodolist` |
| `DB_USER` | Usuario do banco | `gotodolist` |
| `DB_PASSWORD` | Senha do banco | `replace-me` |
| `DB_SSLMODE` | Modo SSL do Postgres | `disable` |
| `CORS_ALLOW_ORIGIN` | Origem permitida para a API REST | `*` |
| `GOTODOLIST_API_BASE_URL` | Base da API usada pelo frontend | vazio, com fallback para `http://<host>:8081` |
| `HTTP_READ_TIMEOUT` | Timeout de leitura HTTP | `10s` |
| `HTTP_READ_HEADER_TIMEOUT` | Timeout de cabecalho HTTP | `5s` |
| `HTTP_WRITE_TIMEOUT` | Timeout de escrita HTTP | `15s` |
| `HTTP_SHUTDOWN_TIMEOUT` | Tempo maximo para graceful shutdown | `10s` |

## Executando o monolito

1. Garanta um PostgreSQL acessivel.
2. Exporte as variaveis abaixo.
3. Rode a aplicacao.

```bash
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_NAME=gotodolist
export DB_USER=gotodolist
export DB_PASSWORD='<defina-em-runtime>'
export DB_SSLMODE=disable

go run ./monolito
```

Aplicacao: `http://localhost:8080`

## Executando a API REST

```bash
export APP_PORT=8081
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_NAME=gotodolist
export DB_USER=gotodolist
export DB_PASSWORD='<defina-em-runtime>'
export DB_SSLMODE=disable
export CORS_ALLOW_ORIGIN='*'

go run ./desacoplado/backend
```

API: `http://localhost:8081`

## Executando o frontend desacoplado

O frontend desacoplado e composto por arquivos estaticos.

1. Suba a API REST.
2. Abra `desacoplado/frontend/index.html` diretamente no navegador ou sirva essa pasta com um servidor estatico.
3. Se a API nao estiver no host padrao, defina `GOTODOLIST_API_BASE_URL` no arquivo `desacoplado/frontend/config.js` ou use a imagem Docker do frontend, que gera esse arquivo em runtime.

Quando o frontend nao recebe configuracao explicita, ele tenta acessar `http://<host-atual>:8081`. Se nao houver host no navegador, o fallback e `http://localhost:8081`.

## Testes locais

```bash
go test ./...
```