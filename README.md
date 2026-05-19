# gotodolist

Aplicacao didatica de lista de tarefas para demonstrar deploy de uma aplicacao web em Go que conversa com PostgreSQL. O repositorio foi organizado como monorepo para cobrir dois cenarios de aula:

- `monolito/`: aplicacao Go com renderizacao server-side de HTML.
- `desacoplado/backend/`: API REST em Go.
- `desacoplado/frontend/`: frontend estatico em HTML, CSS e JavaScript consumindo a API.

O projeto nao possui autenticacao de proposito. O foco e ensinar o fluxo de dados de listar, criar, atualizar e excluir tarefas, incluindo health checks, desligamento gracioso, conteinerizacao e pipeline de entrega.

## Funcionalidades

- Listagem de tarefas.
- Criacao de nova tarefa com titulo e descricao.
- Edicao de titulo e descricao.
- Marcacao e desmarcacao de conclusao.
- Exclusao de tarefa.
- Endpoint de saude em `GET /health`.
- Tratamento generico de indisponibilidade do banco sem expor detalhes internos.
- Graceful shutdown para encerramento com `SIGINT` e `SIGTERM`.

## Estrutura do repositorio

```text
.
├── desacoplado/
│   ├── backend/
│   └── frontend/
├── internal/
│   ├── config/
│   ├── httpx/
│   ├── postgres/
│   ├── runtime/
│   └── todo/
├── monolito/
├── scripts/
├── docker-compose.yml
└── .github/workflows/
```

## Requisitos

- Go atualizado conforme `go.mod`.
- PostgreSQL 16+ ou 17+.
- Docker e Docker Compose plugin para execucao conteinerizada.

## Configuracao por variaveis de ambiente

Toda configuracao entra por variaveis de ambiente. Nao coloque segredos reais no repositorio.

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
| `HTTP_READ_TIMEOUT` | Timeout de leitura HTTP | `10s` |
| `HTTP_READ_HEADER_TIMEOUT` | Timeout de cabecalho HTTP | `5s` |
| `HTTP_WRITE_TIMEOUT` | Timeout de escrita HTTP | `15s` |
| `HTTP_SHUTDOWN_TIMEOUT` | Tempo maximo para graceful shutdown | `10s` |

## Executando a versao monolitica via codigo-fonte

1. Garanta um PostgreSQL acessivel.
2. Exporte as variaveis necessarias.
3. Rode a aplicacao:

```bash
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_NAME=gotodolist
export DB_USER=gotodolist
export DB_PASSWORD='<defina-em-runtime>'
export DB_SSLMODE=disable

go run ./monolito
```

A interface ficara disponivel em `http://localhost:8080`.

## Executando a API REST via codigo-fonte

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

A API ficara disponivel em `http://localhost:8081`.

## Executando o frontend desacoplado

O frontend desacoplado e composto apenas por arquivos estaticos. Voce pode abrir `desacoplado/frontend/index.html` diretamente no navegador ou servi-lo com qualquer servidor estatico.

Se o frontend estiver fora do host padrao, use o parametro `api` para apontar a base da API:

```text
http://localhost:8082/?api=http://localhost:8081
```

Por padrao, o frontend tenta acessar `http://<host-atual>:8081`.

## Scripts uteis

- `scripts/run-monolito.sh`: executa o monolito com `go run`.
- `scripts/run-api.sh`: executa a API REST com `go run`.
- `scripts/test.sh`: roda `go test ./...` e valida o `docker-compose.yml` quando Docker estiver disponivel.

## Endpoints principais

### Monolito

- `GET /`
- `POST /tasks`
- `POST /tasks/{id}/edit`
- `POST /tasks/{id}/toggle`
- `POST /tasks/{id}/delete`
- `GET /health`

### API REST

- `GET /api/tasks`
- `POST /api/tasks`
- `PUT /api/tasks/{id}`
- `DELETE /api/tasks/{id}`
- `GET /health`

### Exemplo de payload para criar tarefa

```json
{
	"title": "Preparar aula sobre Docker",
	"description": "Criar roteiro para compose, health check e multi-stage build"
}
```

### Exemplo de payload para atualizar tarefa

```json
{
	"title": "Preparar aula final",
	"description": "Roteiro revisado",
	"completed": true
}
```

## Executando com Docker Compose

O arquivo `docker-compose.yml` usa perfis para cobrir os dois formatos de aula.

### Perfil monolito

```bash
docker compose --profile monolito up --build
```

Servicos expostos:

- Monolito: `http://localhost:8080`
- PostgreSQL: `localhost:5432`

### Perfil desacoplado

```bash
docker compose --profile desacoplado up --build
```

Servicos expostos:

- Frontend: `http://localhost:8082`
- API: `http://localhost:8081`
- PostgreSQL: `localhost:5432`

### Variaveis no Compose

O compose aceita sobrescrita em runtime, por exemplo:

```bash
DB_NAME=gotodolist DB_USER=gotodolist DB_PASSWORD='<defina-em-runtime>' docker compose --profile monolito up --build
```

## Deployment didatico

### Em VM

- O monolito pode ser publicado diretamente em uma VM com `APP_PORT=8080`.
- A API desacoplada pode ficar em outra VM ou no mesmo host, desde que o frontend aponte para o endereco correto.
- A aplicacao foi escrita para funcionar por IP direto e tambem atras de proxy reverso.

### Em ambiente conteinerizado

- O banco e a aplicacao podem rodar na mesma VM via Docker Compose.
- Os Dockerfiles da aplicacao Go usam multi-stage build para reduzir a imagem final.
- Os endpoints `/health` podem ser usados por orquestradores, balanceadores e proxys.

## Qualidade e testes

- Testes unitarios do dominio em `internal/todo/service_test.go`.
- Testes HTTP da API em `desacoplado/backend/server_test.go`.
- Testes HTTP do monolito em `monolito/server_test.go`.

Para executar localmente:

```bash
go test ./...
```

## CI/CD

O repositorio inclui dois workflows em GitHub Actions:

- `.github/workflows/ci.yml`: roda formatacao, lint, testes, validacao do compose, build das imagens e scan com Trivy.
- `.github/workflows/publish-docker.yml`: segue a ideia da referencia fornecida, gera uma nova versao com base nos labels do PR, publica imagens no GHCR, roda Trivy nas imagens publicadas e cria uma release no GitHub.

### Labels de versao para release

- `major`: incrementa major.
- `minor`, `feature` ou `feat`: incrementa minor.
- sem label reconhecida: incrementa patch.

As imagens publicadas ficam com estes nomes:

- `ghcr.io/<owner>/gotodolist-monolito:<versao>`
- `ghcr.io/<owner>/gotodolist-api:<versao>`
- `ghcr.io/<owner>/gotodolist-frontend:<versao>`

## Seguranca

- Nao registre segredos reais em codigo, compose, Dockerfile ou README.
- As mensagens de erro para indisponibilidade do banco sao genericas por desenho.
- O projeto nao imprime a senha do banco.
