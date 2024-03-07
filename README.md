# Respond

[![GoDoc](https://pkg.go.dev/badge/github.com/gookit/respond.svg)](https://pkg.go.dev/github.com/gookit/respond)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/respond)](https://goreportcard.com/report/github.com/gookit/respond)
[![Unit-Tests](https://github.com/gookit/respond/workflows/Unit-Tests/badge.svg)](https://github.com/gookit/respond/actions)

Quickly respond `Text`, `HTML`, `XML`, `JSON`, `JSONP` and more data to `http.ResponseWriter`.

## Godoc

- [godoc](https://pkg.go.dev/github.com/gookit/respond)

## Quick start

```go
package main

import (
    "net/http"
	
    "github.com/gookit/respond"
)

func main() {
    // config and init the default Responder
    respond.Initialize(func(opts *respond.Options) {
        opts.TplLayout = "two-column.tpl"
        opts.TplViewsDir = "templates"
    })
    
    http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
        respond.JSON(w, 200, map[string]string{
            "name": "tom",
        })
    })
    
    http.HandleFunc("/xml", func(w http.ResponseWriter, r *http.Request) {
        respond.XML(w, 200, map[string]string{
            "name": "tom",
        })
    })
    
    http.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
        respond.HTML(w, 200, "home.tpl", map[string]string{
            "name": "tom",
        })
    })
    
    http.ListenAndServe(":8080", nil)
}
```

## Create new

```go
package main

import (
    "net/http"
	
    "github.com/gookit/respond"
)

func main() {
    render := respond.New(func(opts *respond.Options) {
        opts.TplLayout = "two-column.tpl"
        opts.TplViewsDir = "templates"
    })
    render.Initialize()
    
    // usage
    http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
        render.JSON(w, 200, map[string]string{
            "name": "tom",
        })
    })
    http.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
        render.HTML(w, 200, "home.tpl", map[string]string{
            "name": "tom",
        })
    })
    
}
```

## Reference

- https://github.com/unrolled/render
- https://github.com/thedevsaddam/renderer

## License

**MIT**
