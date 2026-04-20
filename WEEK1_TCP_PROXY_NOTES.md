# Week 1 TCP Proxy Notes

## What I learned

- `net.Listen("tcp", ":2000")` creates a listening socket for incoming connections.
- `net.Dial("tcp", ":2001")` creates an outgoing connection to a backend.
- `listener.Accept()` blocks until a new client connection is ready.
- `conn.Read(buf)` blocks until bytes arrive, EOF happens, or an error happens.
- `Read` returns `n, err`.
- `n` is the number of valid bytes in the buffer.
- `err` tells me whether the read ended because of EOF or another error.
- `string(buf)` is wrong after `Read` because it prints stale bytes too.
- `string(buf[:n])` is correct because it prints only the bytes that were actually read.
- A TCP server usually needs one goroutine per connection so one blocked client does not block all other clients.
- A TCP proxy needs two copy directions:
  - client -> server
  - server -> client
- Two sequential `io.Copy` calls are wrong for a full-duplex proxy because the second one waits for the first one to finish.
- `net.Dial` must happen per client connection, not once globally.
- `defer conn.Close()` closes silently; it does not log anything.
- `io.TeeReader` helps inspect bytes while still letting `io.Copy` handle the loop.

## Quick self-check answers

1. `Listen` waits for incoming connections; `Dial` starts an outgoing connection.
2. `Accept()` blocks because there is no completed incoming connection ready yet.
3. `Read()` blocks because there are no bytes ready to read yet.
4. `Read()` returns `n, err`.
5. Use `buf[:n]` because only the first `n` bytes are valid.
6. One goroutine per connection prevents one slow connection from blocking all others.
7. A proxy needs both directions because TCP is full-duplex.
8. Sequential `io.Copy` is wrong because they do not run in parallel.
9. Each client needs its own backend connection.
10. If the client disconnects, the copy that reads from the client eventually returns EOF or an error.
11. `defer Close()` does not print; check with logs or `lsof`/`ss`.
12. Check `err` before using `conn.RemoteAddr()`.
13. The accept loop lives for the server lifetime; the per-connection copy/read loop lives only for one client.
14. A shared channel for all backend connections can pair the wrong client and server together.
15. `io.TeeReader` lets me log bytes while forwarding them.

## What is still not done

- Test with `curl` and a real backend like `python3 -m http.server`.
- Clean up `cmd/proxy/main.go`.
- Add a short code comment explaining why each connection gets its own goroutine.
