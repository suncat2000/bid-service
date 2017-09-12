### Bid service

Implementation of a Bid service in Golang.

### Running

```
$ go run main.go [-addr=:8080]
```

### Usage

```
GET http://localhost:{run_addr_port}/winner?s=http://localhost:8081/primes&s=http://localhost:8081/fibo&s=http://localhost:8081/rand&s=http://localhost:8081/fact
```
> s - []string of price sources

```bash
$ curl -i -H "Accept: application/json" http://localhost:{run_addr_port}/winner?s=http://example.com/primes&s=http://example.com/fibo&s=http://example.com/rand&s=http://example.com/fact
Content-Type: application/json
Date: Tue, 12 Sep 2017 21:22:26 GMT
Content-Length: 52

{"price":23,"source":"http://example.com/primes"}
```

### Response

##### Success
```json
{
    "price": 23,
    "source": "http://example.com/primes"
}
```
##### Error
```json
{
    "code": "{40x|50x}",
    "error": "Error message"
}
```

### Test

```
$ go test requester_test.go requester.go -cpu 1
```
