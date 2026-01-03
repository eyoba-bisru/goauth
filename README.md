# Google OAuth using Go

A minimal example demonstrating Google OAuth2 in Go. This project exposes a few simple endpoints and includes an OpenAPI (Swagger) description plus hosted documentation views.

## Features

- Login via Google OAuth2 (`/login`)
- OAuth callback that returns Google user profile JSON (`/callback`)
- Health check (`/health`)
- OpenAPI spec at `/swagger.json` with Swagger UI at `/docs` and ReDoc at `/redoc`

## Prerequisites

- Go 1.25.4 installed
- A Google OAuth client (Client ID, Client Secret) with an OAuth redirect URI configured

Create a `.env` file in the project root with the following variables:

```
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URI=http://localhost:8080/callback
```

## Run locally

Start the server with:

```bash
go run main.go
```

Or build and run the binary:

```bash
go build -o bin/goauth ./...
./bin/goauth
```

The server listens on port `8080` by default.

## API Docs

- OpenAPI JSON: `/swagger.json`
- Swagger UI: `/docs` (served using the Swagger UI CDN)
- ReDoc: `/redoc`

Visit `http://localhost:8080/docs` to interactively explore the API similar to FastAPI's docs.

## Endpoints

- `GET /health` — returns `Alive`
- `GET /login` — redirects to Google to start the OAuth flow
- `GET /callback?code=...` — OAuth callback; returns Google user profile JSON

## Notes

- The project serves the OpenAPI spec from `docs/swagger.json`. You can replace or regenerate that file if you change handlers or models.
