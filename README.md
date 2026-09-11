# API Service

Go + PocketBase + Temporal backend service.

## Directory Layout

- `cmd/server/`: Application entry point.
- `config/`: Configuration setup.
- `internal/`: Internal business logic, models, services, handlers, repositories, routes, and middleware.
- `pb_data/`: PocketBase local data storage directory.
- `pkg/`: Shared utility packages and logger.

## Development

Run hot-reload development server:
```bash
task dev
```

Build binary:
```bash
task build
```
