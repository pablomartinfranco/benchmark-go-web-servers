# benchmark-go-web-servers

Inspired by [benchmark-http-servers](https://github.com/pablomartinfranco/benchmark-http-servers), this repository provides minimal Go targets focused on shared-hosting protocols.

## Implementations

- `fastcgi-go-stdlib`: FastCGI server using Go's `net/http/fcgi`.
- `fastcgi-go-dispatch`: FastCGI server using a non-blocking dispatch queue pattern on top of Go stdlib handlers.
- `wsgi-go-stdlib`: Go HTTP backend plus `passenger_wsgi.py` adapter example for cPanel/Passenger-style WSGI entrypoints.

## Build

```bash
go build ./...
```

## Test

```bash
go test ./...
```
