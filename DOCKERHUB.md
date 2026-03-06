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
