# Referencia da API

## Endpoints do monolito

- `GET /`
- `POST /tasks`
- `POST /tasks/{id}/edit`
- `POST /tasks/{id}/toggle`
- `POST /tasks/{id}/delete`
- `GET /health`

## Endpoints da API REST

- `GET /api/tasks`
- `POST /api/tasks`
- `PUT /api/tasks/{id}`
- `DELETE /api/tasks/{id}`
- `GET /health`

## Payload para criar tarefa

```json
{
  "title": "Preparar aula sobre Docker",
  "description": "Criar roteiro para compose, health check e multi-stage build"
}
```

## Payload para atualizar tarefa

```json
{
  "title": "Preparar aula final",
  "description": "Roteiro revisado",
  "completed": true
}
```