# Week 2 HTTP Reverse Proxy Notes

## What I learned

- `http.Handler` is an interface with one method: `ServeHTTP(ResponseWriter, *Request)`.
- `http.HandlerFunc` is a type that converts a plain function into an `http.Handler` — it is the shortcut so you don't need a struct.
- `proxyMux.HandleFunc("/", proxy)` registers the proxy function for every incoming path.
- You cannot reuse `req` as the outgoing request — the body is a stream (reading it consumes it), headers need to be filtered, and `RequestURI` cannot be set on outgoing requests. Use `http.NewRequest`.
- `req.URL.Path` drops the query string. `req.RequestURI` keeps path + query string. Always use `RequestURI` when forwarding.
- `X-Forwarded-For` goes on the **outgoing request** to the backend, not on the response writer. It tells the backend the real client IP.
- `http.Client` must be created once at package level — the connection pool inside `Transport` is only useful if the client is reused across requests.
- `http.Client` is safe for concurrent use because it has an internal `sync.Mutex` protecting the connection pool. Goroutines alone do not prevent data races — the mutex does.
- `w.WriteHeader(statusCode)` must be called before `io.Copy(w, body)` — once headers are flushed to the socket they cannot be changed.
- A 404 from the backend is NOT a proxy error. `client.Do()` returns `err == nil` for any HTTP response. The `err != nil` branch only fires when the backend is unreachable (timeout, connection refused, DNS failure).
- Hop-by-hop headers describe the current TCP connection, not the request. They must be stripped before forwarding: `Connection`, `Keep-Alive`, `Transfer-Encoding`, `Upgrade`, `Proxy-Authorization`, `Proxy-Authenticate`, `Te`, `Trailers`.
- `httptest.NewRecorder()` simulates the client's connection — it captures status code, headers, and body written by the handler. Read via `rec.Result()`.
- `httptest.NewServer` starts a real HTTP server on a random port. Its URL is in `ts.URL` (read-only). Point `backendURL` at it before calling `proxy()` in tests.
- `httputil.DumpRequest` / `DumpResponse` print the HTTP wire format (useful for debugging). They buffer the body and replace `req.Body` so the stream can still be read downstream. Do not use in production — it loads the full body into memory.

## Connection pool mental model

- Idle connections = parked cars waiting to be reused.
- `MaxIdleConns` = parking lot size (total across all hosts).
- `MaxIdleConnsPerHost` = spots per host (default 2 — too low for a proxy, set higher).
- Active connections (on the road) are not capped by default.
- When a request finishes: if the lot has space → park the connection; if full → destroy it.
- `IdleConnTimeout` = if a parked car hasn't been used in N seconds, destroy it.
- Building a new connection is expensive (TCP handshake + TLS). Reuse saves latency under load.

## Quick self-check answers

1. `http.Handler` is an interface; `http.HandlerFunc` is a type that adapts a function to satisfy it.
2. `HandleFunc("/", proxy)` registers `proxy` to be called for every request to that pattern.
3. Can't reuse `req` — body is a stream, headers need filtering, `RequestURI` is forbidden on outgoing requests.
4. Use `req.RequestURI` — it includes the query string. `URL.Path` drops it.
5. `X-Forwarded-For` belongs on the outgoing request; setting it on `w` sends it back to the client, which is wrong.
6. `http.Client` at package level — the connection pool is reused. Inside the function a new pool is created and destroyed every request.
7. Goroutines don't prevent races. `http.Client` is safe because it has an internal mutex protecting the pool map.
8. `WriteHeader` flushes the status line and headers to the socket. After that they are sealed — `io.Copy` writes the body after headers.
9. 404 from backend → `err == nil`, `res.StatusCode == 404` → forward it to client with `w.WriteHeader(404)`. Not a proxy error.
10. Hop-by-hop headers describe the TCP connection, not the request. Three examples: `Connection`, `Transfer-Encoding`, `Upgrade`.
11. `NewRecorder` captures what the handler writes. After `proxy(rec, req)` you read `rec.Result()` and assert on status, headers, body.
12. Without `backendURL = ts.URL` the proxy tries `localhost:2000` which has nothing listening → 502.
