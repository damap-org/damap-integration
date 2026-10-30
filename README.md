# DAMAP integration server

This is an integration server for the DAMAP project written in Go. Its architecture lends itself to adding multiple possible backends, the dummy backend being included.

## Building this server

You can build the server by running the following command:

```
go build -o damap-integration github.com/damap-org/damap-integration/cmd/damap-integration
```

## Configuration

The server has the following CLI and environment configuration options:

| CLI             | Environment                     | Description                                                 | Default        |
|-----------------|---------------------------------|-------------------------------------------------------------|----------------|
| `-h` `--help`   |                                 | Prints all options                                          |                |
| `--listen`      | `DAMAP_INTEGRATION_LISTEN`      | HTTP listen port                                            | 127.0.0.1:8080 |
| `--listen-tls`  | `DAMAP_INTEGRATION_LISTEN_TLS`  | HTTPS listen port (requires --certificate and --key)        |                |
| `--certificate` | `DAMAP_INTEGRATION_CERTIFICATE` | PEM-formatted certificate or path to a file containing one. |                |
| `--key`         | `DAMAP_INTEGRATION_KEY`         | PEM-formatted private key or path to a file containing one. |                |
| `--backend`     | `DAMAP_INTEGRATION_BACKEND`     | Backend to use.                                             | dummy          |

## Backends

### Dummy

The `dummy` backend returns test data and has no configuration options.

## Running with Docker

Build and start the server:

```sh
docker build -t damap-integration:local .
docker run --rm -p 127.0.0.1:8080:8080 damap-integration:local
```

The server uses the dummy backend by default. To check that it is running:

```sh
curl --fail 'http://127.0.0.1:8080/api/projects?query=test'
```

Pass configuration through environment variables (`docker run -e KEY=value`)
or CLI options after the image name. CLI options take precedence. The container
listens on `0.0.0.0:8080`; the standalone binary defaults to `127.0.0.1:8080`.

The container runs as UID/GID 65532. Mounted TLS certificates and keys must be
readable by that user.
