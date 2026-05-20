# Compose e imagens Docker

O repositorio possui quatro arquivos Compose com responsabilidades separadas.

## Arquivos disponiveis

- `docker-compose-monolito-dev.yaml`: banco + monolito com build local.
- `docker-compose-monolito-prod.yaml`: banco + monolito usando imagem publicada no GHCR.
- `docker-compose-dev.yaml`: banco + API + frontend com build local.
- `docker-compose-prod.yaml`: banco + API + frontend usando imagens publicadas no GHCR.

## Arquivo de variaveis

Use `example.env` como base para criar o arquivo de variaveis do ambiente:

```bash
cp example.env .env
```

Depois suba os servicos sempre com `--env-file .env` para deixar os valores explicitos e repetir o mesmo comando em dev ou deploy.

## Subindo o monolito localmente

```bash
docker compose --env-file .env -f docker-compose-monolito-dev.yaml up --build
```

Servicos expostos:

- Monolito: `http://localhost:8080`
- PostgreSQL: `localhost:5432`

## Subindo a stack desacoplada localmente

```bash
docker compose --env-file .env -f docker-compose-dev.yaml up --build
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
docker compose --env-file .env -f docker-compose-dev.yaml up --build
```

Use um endereco que o navegador do usuario consiga resolver. Nome de container so funciona se o navegador tambem enxergar essa rede, o que normalmente nao acontece fora do Docker. Para esse caso, ajuste `GOTODOLIST_API_BASE_URL` no `.env`, por exemplo para `http://192.168.0.20:8081`.

## Subindo imagens publicadas no GHCR

### Monolito

```bash
docker compose --env-file .env -f docker-compose-monolito-prod.yaml up -d
```

### Stack desacoplada

```bash
docker compose --env-file .env -f docker-compose-prod.yaml up -d
```

### Usando outro owner no GHCR

```bash
docker compose --env-file .env -f docker-compose-prod.yaml up -d
```

Altere `GHCR_OWNER` e `IMAGE_TAG` dentro do `.env` quando precisar trocar o repositorio ou a versao da imagem.

## Sobre `profiles`

Os arquivos de producao nao usam mais `profiles`. Eles eram necessarios apenas enquanto `docker-compose-prod.yaml` reunia monolito e stack desacoplada no mesmo arquivo. Agora cada Compose define um unico conjunto de servicos, entao basta escolher o arquivo correto.