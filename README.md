# go-time-server

Go current time server for exploring the standard library.

Specifically:
- Implement a Handler (`net/http`) that responds with the current time, using `time`.
  Response is either text (default) or json (with request header `Accept: application/json`),
  making use of `encoding/json`.
- Use structured logging (`log/slog`) to log the IP of the requester.
  Implemented as a middleware function.

## Testing

```sh
$ make run
# On other instance...
# JSON output:
$ curl -i -H "Accept: application/json" http://127.0.0.1:8080/now
# Text output:
$ curl -i http://127.0.0.1:8080/now
```
