# tag

[![Build Status](https://github.com/gin-gonic/gin/workflows/Run%20Tests/badge.svg?branch=master)](https://github.com/gin-gonic/gin/actions?query=branch%3Amaster)
[![GoDoc](https://pkg.go.dev/badge/github.com/gin-gonic/gin?status.svg)](https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc)

A simple library for passing groups of key/value tags around.

## Install

```
$ go get -u github.com/haleyrc/tag
```

## Usage

The basics of the `tag` package are pretty simple. Everything revolves around the idea of a tag `Group`. You create a `Group` by passing an optionally empty list of key/value pairs to `NewGroup`:

```go
g := tag.NewGroup(tag.Dict{
    "env": "prod",
    "service": "database",
})
```

With a `Group` in hand, you can add new tags, merge tags from another `Group`, retrieve the values of tags by their key, and convert to other representations (more on that one below). But a `Group` isn't very useful on its own.

A standard use of tags in a server context is to grab some important values from the request (path, verb, user agent, etc.) and tag various observability objects with them to allow engineers to troubleshoot, build reliability dashboards, and monitor SLIs. Since tags may be added at any layer of the stack, you need to be able to pass a `Group` all the way through.

To do this with the `tag` package, use the `WithTag` and `WithGroup` helpers. A
typical middleware to do this might look something like the following:

```go
func RequestMetrics(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWrite, r *http.Request) {
        ctx := tag.WithGroup(ctx, tag.NewGroup(tag.Dict{
            "method": r.Method,
            "path": r.URL.Path,
        }))

        next(w, r.WithContext(ctx))
    }
}
```

When you're ready to to use your tags, you can get the entire `Group` with `FromContext`. Since the `Group` type itself is unlikely to be used by your third-party observability package, a `Group` includes methods to convert the structure to the most common representations: `Slice` and `Map`.

### Example

If you're using the DataDog client, you'll see methods defined like:

```go
func (c *Client) Incr(name string, tags []string, rate float64) error
```

DataDog takes tags as a flat slice of strings with alternating keys and values. With a `Group`, this is easy to achieve:

```go
tags := tag.FromContext(ctx)
_ = c.Incr("my.metric", tags.Slice(), 1.0)
```
