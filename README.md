# Altcha Server

A lightweight, standalone server implementation for [Altcha](https://altcha.org/) (an open-source CAPTCHA alternative) written in Go.

This server provides endpoints to generate challenges and verify solutions, making it easy to integrate proof-of-work (PoW) protection into your web forms without needing a complex backend integration.

## Features

- **Lightweight & Fast**: Written in Go for optimal performance and low memory footprint.
- **Standalone**: Can run as a binary or within a Docker container.
- **Anti-Replay Protection**: Automatically tracks and prevents the reuse of verified solutions.
- **Configurable**: Easily customize through environment variables.
- **CORS Support**: Integrated CORS handling for seamless frontend integration.

## Quick Start

### Running with Docker

The easiest way to get started is by using Docker:

```bash
docker run -p 3947:3947 \
  -e ALTCHA_SECRET=your_very_secret_key \
  -e ALTCHA_CORS_ORIGIN=https://your-frontend.com \
  altcha-server:latest
```

### Running with Docker Compose (Docker Desktop)

For a more persistent setup on Docker Desktop, use a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  altcha:
    image: altcha-server:latest
    ports:
      - "3947:3947"
    environment:
      - ALTCHA_SECRET=your_very_secret_key
      - ALTCHA_CORS_ORIGIN=https://your-frontend.com
    volumes:
      - ./data:/app
```

Then run it with:

```bash
docker-compose up -d
```

### Running in Docker Swarm with Traefik

If you're using Docker Swarm and Traefik as a reverse proxy, you can deploy it with the following stack definition:

```yaml
version: '3.8'

services:
  altcha:
    image: altcha-server:latest
    networks:
      - traefik-public
    environment:
      - ALTCHA_SECRET=your_very_secret_key
      - ALTCHA_CORS_ORIGIN=https://altcha.example.com
    volumes:
      - altcha-data:/app
    deploy:
      labels:
        - "traefik.enable=true"
        - "traefik.http.routers.altcha.rule=Host(`altcha.example.com`)"
        - "traefik.http.routers.altcha.entrypoints=websecure"
        - "traefik.http.routers.altcha.tls.certresolver=myresolver"
        - "traefik.http.services.altcha.loadbalancer.server.port=3947"

volumes:
  altcha-data:

networks:
  traefik-public:
    external: true
```

Deploy the stack:

```bash
docker stack deploy -c docker-compose.yml altcha
```

### Running from Binary

If you have the binary:

```bash
export ALTCHA_SECRET=your_very_secret_key
export ALTCHA_CORS_ORIGIN=https://your-frontend.com
./altcha-server
```

## API Endpoints

### 1. Get Challenge
**Endpoint**: `GET /challenge`

Returns a new Altcha challenge object.

**Response**:
```json
{
  "algorithm": "SHA-256",
  "challenge": "7e91513ebdb...5486aaf72d35b8",
  "salt": "2e341e4918e69a71d2eadac5?expires=1772819149",
  "signature": "b2effa1a4a555f1358a1cf020f7b1dc14b996a88ae66503b0f89bacf9a1f7a3b"
}
```

### 2. Verify Solution
**Endpoint**: `POST /verify`

Verifies a solution payload submitted by the client.

**Request Body**:
```json
{
  "payload": "eyJhbGdvcml0aG0iOiJTSEEtMjU2IiwiY2hhbGxlbmdlIjoiN2U5..."
}
```

**Response**:
- `200 OK`: Verification successful.
- `400 Bad Request`: Invalid request body.
- `403 Forbidden`: Invalid Altcha payload or payload already used (replay).

## Configuration

The server is configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `ALTCHA_SECRET` | Secret key used for signing challenges. | *Generated randomly if not set* |
| `ALTCHA_TTL` | Time-to-live for challenges (e.g., `1h`, `30m`). | `1h` |
| `ALTCHA_CORS_ORIGIN` | Allowed CORS origins (comma-separated). Use `*` for all. | `*` |
| `IS_DEV` | If set to `true`, disables CORS origin checks and allows all. | `false` |

## Persistence

The server tracks verified solutions to prevent replay attacks. By default, it stores these in:
- `solutions.txt`: List of used payloads.
- `verifications.log`: Log of all verification attempts.
- `challenges.log`: Log of all generated challenges.

When running in Docker, it's recommended to mount a volume to persist these files:

```bash
docker run -p 3947:3947 \
  -v ./data:/app \
  -e ALTCHA_SECRET=your_secret \
  altcha-server:latest
```

## Development

To build the server from source:

1. Ensure you have Go 1.26+ installed.
2. Clone the repository.
3. Run `go build -o altcha-server .`.

## License

MIT
