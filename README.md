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
| `--listen`      | `DAMAP_INTEGRATION_LISTEN`      | HTTP listen port                                            | 127.0.0.1:8088 |
| `--listen-tls`  | `DAMAP_INTEGRATION_LISTEN_TLS`  | HTTPS listen port (requires --certificate and --key)        |                |
| `--certificate` | `DAMAP_INTEGRATION_CERTIFICATE` | PEM-formatted certificate or path to a file containing one. |                |
| `--key`         | `DAMAP_INTEGRATION_KEY`         | PEM-formatted private key or path to a file containing one. |                |
| `--backend`     | `DAMAP_INTEGRATION_BACKEND`     | Backend to use.                                             | dummy          |

## Backends

### Dummy

The `dummy` backend returns test data and has no configuration options.