# gentleman/tls [![Build Status](https://travis-ci.org/h2non/gentleman.png)](https://travis-ci.org/h2non/gentleman) [![GitHub release](https://img.shields.io/github/tag/h2non/gentleman.svg)](https://github.com/h2non/gentleman/releases) [![GoDoc](https://godoc.org/github.com/h2non/gentleman/plugins/tls?status.svg)](https://godoc.org/github.com/h2non/gentleman/plugins/tls) [![API](https://img.shields.io/badge/api-beta-green.svg?style=flat)](https://godoc.org/github.com/h2non/gentleman/plugins/tls) [![Go Report Card](https://goreportcard.com/badge/github.com/h2non/gentleman)](https://goreportcard.com/report/github.com/h2non/gentleman)

gentleman's plugin to easily define TLS config used by `http.Transport`/`RoundTripper` interface.

## Installation

```bash
go get -u gopkg.in/h2non/gentleman.v0/plugins/tls
```

## API

See [godoc](https://godoc.org/github.com/h2non/gentleman/plugins/tls) reference.

## Example

```go
package main

import (
  "fmt"
  "crypto/tls"
  "gopkg.in/h2non/gentleman.v0"
  "gopkg.in/h2non/gentleman.v0/plugins/tls"
)

func main() {
  // Create a new client
  cli := gentleman.New()

  // Define a custom header
  cli.Use(tls.Config(&tls.Config{ServerName: "foo.com"}))

  // Perform the request
  res, err := cli.Request().URL("http://httpbin.org/headers").End()
  if err != nil {
    fmt.Printf("Request error: %s\n", err)
    return
  }
  if !res.Ok {
    fmt.Printf("Invalid server response: %d\n", res.StatusCode)
    return
  }

  fmt.Printf("Status: %d\n", res.StatusCode)
  fmt.Printf("Body: %s", res.String())
}
```

## License

MIT - Tomas Aparicio