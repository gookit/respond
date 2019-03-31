# Respond

respond `Text`, `HTML`, `XML`, `JSON`, `JSONP` data to http.ResponseWriter

## Godoc

- [doc on gowalker](https://gowalker.org/github.com/gookit/respond)
- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/respond.v1)
- [godoc for github](https://godoc.org/github.com/gookit/respond)

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
