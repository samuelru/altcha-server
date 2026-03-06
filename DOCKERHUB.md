# Altcha Server

A lightweight, standalone server implementation for [Altcha](https://altcha.org/) (an open-source CAPTCHA alternative) written in Go.

This server provides endpoints to generate challenges and verify solutions, allowing for proof-of-work (PoW) protection on web forms.

## Usage

### Run the Container

The container is configured using environment variables.

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

### Persistence

The server tracks verified solutions and logs events to files. It's recommended to mount a volume to persist this data:

```bash
docker run -p 3947:3947 \
  -v /path/to/data:/app \
  -e ALTCHA_SECRET=your_secret \
  altcha-server:latest
```

## Configuration

The server is configured using the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `ALTCHA_SECRET` | Secret key used for signing challenges. | *Generated randomly if not set* |
| `ALTCHA_TTL` | Time-to-live for challenges (e.g., `1h`, `30m`). | `1h` |
| `ALTCHA_CORS_ORIGIN` | Allowed CORS origins (comma-separated). | `*` |
| `IS_DEV` | If `true`, allows all origins (useful for local development). | `false` |

## API Endpoints

- `GET /challenge`: Generates a new Altcha challenge.
- `POST /verify`: Verifies a solution payload.

## Source

For more information and source code, visit the [GitHub Repository](https://github.com/your-username/altcha-server).
