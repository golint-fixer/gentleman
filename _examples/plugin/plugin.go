package main

import (
	"fmt"
	"net/url"

	"github.com/lytics/gentleman"
	"github.com/lytics/gentleman/context"
	"github.com/lytics/gentleman/plugin"
	"github.com/lytics/gentleman/plugins/headers"
)

func main() {
	// Create a new client
	cli := gentleman.New()

	// Define a custom header
	cli.Use(headers.Set("Token", "s3cr3t"))

	// Create a request plugin to define the URL
	cli.Use(plugin.NewRequestPlugin(func(ctx *context.Context, h context.Handler) {
		u, _ := url.Parse("http://httpbin.org/headers")
		ctx.Request.URL = u
		h.Next(ctx)
	}))

	// Perform the request
	res, err := cli.Request().Send()
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
