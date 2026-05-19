# Compose e imagens Docker

O repositorio possui tres arquivos Compose com responsabilidades separadas.

## Arquivos disponiveis

- `docker-compose-monolito.yaml`: banco + monolito com build local.
- `docker-compose-dev.yaml`: banco + API + frontend com build local.
- `docker-compose-prod.yaml`: monolito ou stack desacoplada usando imagens publicadas no GHCR.

## Subindo o monolito localmente

```bash
docker compose -f docker-compose-monolito.yaml up --build
```

Servicos expostos:

- Monolito: `http://localhost:8080`
- PostgreSQL: `localhost:5432`

## Subindo a stack desacoplada localmente

```bash
docker compose -f docker-compose-dev.yaml up --build
```

Servicos expostos:

- Frontend: `http://localhost:8082`
- API: `http://localhost:8081`
- PostgreSQL: `localhost:5432`

## Variaveis uteis no Compose

| Variavel | Uso |
| --- | --- |
| `DB_NAME` | Nome do banco |
| `DB_USER` | Usuario do banco |
| `DB_PASSWORD` | Senha do banco |
| `CORS_ALLOW_ORIGIN` | Origem permitida pela API |
| `GOTODOLIST_API_BASE_URL` | Endereco da API usado pelo frontend |
| `GHCR_OWNER` | Dono das imagens no GHCR |
| `IMAGE_TAG` | Tag das imagens publicadas |

Exemplo com API externa ao frontend:

```bash
GOTODOLIST_API_BASE_URL=http://192.168.0.20:8081 docker compose -f docker-compose-dev.yaml up --build
```

Use um endereco que o navegador do usuario consiga resolver. Nome de container so funciona se o navegador tambem enxergar essa rede, o que normalmente nao acontece fora do Docker.

## Subindo imagens publicadas no GHCR

### Monolito

```bash
IMAGE_TAG=latest docker compose -f docker-compose-prod.yaml --profile monolito up -d
```

### Stack desacoplada

```bash
IMAGE_TAG=latest docker compose -f docker-compose-prod.yaml --profile desacoplado up -d
```

### Usando outro owner no GHCR

```bash
GHCR_OWNER=elton-bt IMAGE_TAG=1.2.3 docker compose -f docker-compose-prod.yaml --profile desacoplado up -d
```