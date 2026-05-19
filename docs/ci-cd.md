# CI/CD e releases

O repositorio possui dois workflows principais em GitHub Actions.

## CI

Arquivo: `.github/workflows/ci.yml`

Responsabilidades:

- validar formatacao com `gofmt`
- rodar `golangci-lint`
- rodar `go test ./...`
- validar os arquivos Compose
- buildar as tres imagens Docker

## Publish Docker Images and Release

Arquivo: `.github/workflows/publish-docker.yml`

Responsabilidades:

- calcular a proxima versao a partir dos labels do PR
- validar codigo e Compose antes da release
- buildar imagens locais para scan
- rodar Trivy antes do push das imagens
- publicar imagens no GHCR
- criar uma release no GitHub

## Labels de versao

- `major`: incrementa major
- `minor`, `feature` ou `feat`: incrementa minor
- sem label reconhecida: incrementa patch

## Nomes das imagens publicadas

- `ghcr.io/<owner>/gotodolist-monolito:<versao>`
- `ghcr.io/<owner>/gotodolist-api:<versao>`
- `ghcr.io/<owner>/gotodolist-frontend:<versao>`