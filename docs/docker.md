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

Os Compose mantem a porta interna padrao de cada servico e expoem apenas a porta publicada no host via `.env`.

## Subindo o monolito localmente

```bash
docker compose --env-file .env -f docker-compose-monolito-dev.yaml up --build
```

Servicos expostos:

- Monolito: `http://localhost:8080` por padrao, controlado por `MONOLITO_HOST_PORT`
- PostgreSQL: `localhost:5432` por padrao, controlado por `DB_HOST_PORT`

## Subindo a stack desacoplada localmente

```bash
docker compose --env-file .env -f docker-compose-dev.yaml up --build
```

Servicos expostos:

- Frontend: `http://localhost:8082` por padrao, controlado por `FRONTEND_HOST_PORT`
- API: `http://localhost:8081` por padrao, controlado por `API_HOST_PORT`
- PostgreSQL: `localhost:5432` por padrao, controlado por `DB_HOST_PORT`

## Variaveis uteis no Compose

| Variavel | Uso |
| --- | --- |
| `DB_NAME` | Nome do banco |
| `DB_USER` | Usuario do banco |
| `DB_PASSWORD` | Senha do banco |
| `DB_HOST` | Host do PostgreSQL visto pela aplicacao |
| `DB_PORT` | Porta do PostgreSQL vista pela aplicacao |
| `DB_HOST_PORT` | Porta publicada do PostgreSQL no host |
| `MONOLITO_HOST_PORT` | Porta publicada do monolito no host |
| `API_HOST_PORT` | Porta publicada da API no host |
| `FRONTEND_HOST_PORT` | Porta publicada do frontend no host |
| `CORS_ALLOW_ORIGIN` | Origem permitida pela API |
| `GOTODOLIST_API_BASE_URL` | Endereco da API usado pelo frontend |
| `GHCR_OWNER` | Dono das imagens no GHCR |
| `IMAGE_TAG` | Tag das imagens publicadas |

Exemplo com API externa ao frontend:

```bash
docker compose --env-file .env -f docker-compose-dev.yaml up --build
```

Use um endereco que o navegador do usuario consiga resolver. Nome de container so funciona se o navegador tambem enxergar essa rede, o que normalmente nao acontece fora do Docker. Para esse caso, ajuste `GOTODOLIST_API_BASE_URL` no `.env`, por exemplo para `http://192.168.0.20:8081`. Se voce mudar `API_HOST_PORT`, atualize essa URL para a mesma porta publicada.

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

## Notas sobre portas

- `DB_HOST` e `DB_PORT` definem como a aplicacao acessa o PostgreSQL. Quando voce usa o banco do proprio Compose, mantenha `DB_HOST=db` e `DB_PORT=5432`.
- `DB_HOST_PORT`, `MONOLITO_HOST_PORT`, `API_HOST_PORT` e `FRONTEND_HOST_PORT` controlam apenas as portas publicadas no host.
- As aplicacoes ja usam internamente as portas padrao do projeto: `8080` no monolito, `8081` na API e `8080` no frontend estatico.

## Sobre `profiles`

Os arquivos de producao nao usam mais `profiles`. Eles eram necessarios apenas enquanto `docker-compose-prod.yaml` reunia monolito e stack desacoplada no mesmo arquivo. Agora cada Compose define um unico conjunto de servicos, entao basta escolher o arquivo correto.