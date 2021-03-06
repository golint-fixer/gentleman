package main

import (
	"fmt"

	"github.com/lytics/gentleman"
	"github.com/lytics/gentleman/context"
	"github.com/lytics/gentleman/mux"
	"github.com/lytics/gentleman/plugins/url"
)

func main() {
	// Create a new client
	cli := gentleman.New()

	// Use a custom multiplexer for GET requests
	cli.Use(mux.New().AddMatcher(func(ctx *context.Context) bool {
		return ctx.GetString("$phase") == "request" && ctx.Request.Method == "GET"
	}).Use(url.URL("http://httpbin.org/headers")))

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
