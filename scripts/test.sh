#!/usr/bin/env sh
set -eu

go test ./...

if command -v docker >/dev/null 2>&1; then
  docker compose --profile monolito config >/dev/null
  docker compose --profile desacoplado config >/dev/null
fi
