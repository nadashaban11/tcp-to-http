# tcp-to-http 

A custom HTTP/1.1 server built entirely from scratch in Go using raw TCP sockets. No external frameworks, and no standard `net/http` package. Just pure bytes and socket.

## Usage

First, clone the repo and run the server:

```bash
go run main.go
```

The server will start listening on http://localhost:8080.

Endpoints to try:
1. Get the Home Page

```Bash
curl http://localhost:8080/
```

2. Serve a Static HTML File

```Bash
curl http://localhost:8080/html-file
```

3. Post JSON Data (creates a mock user)

```Bash
curl -X POST http://localhost:8080/api/user \
     -H "Content-Type: application/json" \
     -d '{"name":"Nada", "role":"Backend Engineer"}'
```



## This http was built to explore web internals:

**Raw Sockets**: Manages TCP connections directly using `net.Listen`.

**Parsing**: Reads exact byte lengths using the `Content-Length` header to prevent stream corruption when parsing JSON bodies.

**O(1) Routing**: Uses a Hash Map for fast endpoint matching based on Method + Path.