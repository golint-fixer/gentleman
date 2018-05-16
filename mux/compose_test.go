package mux

import (
	"testing"

	"github.com/lytics/gentleman/context"
	"github.com/lytics/gentleman/utils"
)

func TestMuxComposeIfMatches(t *testing.T) {
	mx := New()
	mx.Use(If(Method("GET"), Host("foo.com")).UseRequest(func(ctx *context.Context, h context.Handler) {
		ctx.Request.Header.Set("foo", "bar")
		h.Next(ctx)
	}))
	ctx := context.New()
	ctx.Request.Method = "GET"
	ctx.Request.URL.Host = "foo.com"
	mx.Run("request", ctx)
	utils.Equal(t, ctx.Request.Header.Get("foo"), "bar")
}

func TestMuxComposeIfUnmatch(t *testing.T) {
	mx := New()
	mx.Use(If(Method("GET"), Host("bar.com")).UseRequest(func(ctx *context.Context, h context.Handler) {
		ctx.Request.Header.Set("foo", "bar")
		h.Next(ctx)
	}))
	ctx := context.New()
	ctx.Request.Method = "GET"
	ctx.Request.URL.Host = "foo.com"
	mx.Run("request", ctx)
	utils.Equal(t, ctx.Request.Header.Get("foo"), "")
}

func TestMuxComposeOrMatch(t *testing.T) {
	mx := New()
	mx.Use(Or(Method("GET"), Host("bar.com")).UseRequest(func(ctx *context.Context, h context.Handler) {
		ctx.Request.Header.Set("foo", "bar")
		h.Next(ctx)
	}))
	ctx := context.New()
	ctx.Request.Method = "GET"
	ctx.Request.URL.Host = "foo.com"
	mx.Run("request", ctx)
	utils.Equal(t, ctx.Request.Header.Get("foo"), "bar")
}

func TestMuxComposeOrUnMatch(t *testing.T) {
	mx := New()
	mx.Use(Or(Method("GET"), Host("bar.com")).UseRequest(func(ctx *context.Context, h context.Handler) {
		ctx.Request.Header.Set("foo", "bar")
		h.Next(ctx)
	}))
	ctx := context.New()
	ctx.Request.Method = "POST"
	ctx.Request.URL.Host = "foo.com"
	mx.Run("request", ctx)
	utils.Equal(t, ctx.Request.Header.Get("foo"), "")
}
