[![Version](https://img.shields.io/badge/goversion-1.20.x-blue.svg)](https://golang.org)
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/tsawler/goblender/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tsawler/page)](https://goreportcard.com/report/github.com/tsawler/page)
![Tests](https://github.com/tsawler/page/actions/workflows/tests.yml/badge.svg)
<a href="https://pkg.go.dev/github.com/tsawler/page"><img src="https://img.shields.io/badge/godoc-reference-%23007d9c.svg"></a>
[![Go Coverage](https://github.com/tsawler/page/wiki/coverage.svg)](https://raw.githack.com/wiki/tsawler/page/coverage.html)


# page

This is a simple package that makes rendering Go html templates as painless as possible. It has no dependencies
outside of the standard library.

## Installation
Install it the usual way:

```
go get -u github.com/tsawler/page
```

## Usage
To use, import `github.com/tsawler/page`, and create a variable of type `page.Render`, either by
using the New() method, or by manually constructing the variable. 

At the minimum, you must specify the directory where your templates live with the `TemplateDir` field, and
a list of partials you want to include, with the `Partials` field.

If you set UseCache to true, then every time a template is rendered it will be stored in the `TemplateMap` field, 
which is of type `map[string]*template.Template`. This way, pre-parsed templates can be read very quickly from
memory, rather than rebuilding them from disk on every request, which is an expensive operation.


## Example

Assuming that you have two Go templates named `base.layout.gohtml` and `home.page.gohtml` which look like this...

```gotemplate
{{define "base"}}
    <!doctype html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport"
              content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
        <meta http-equiv="X-UA-Compatible" content="ie=edge">
        <title>Document</title>
    </head>
    <body>
    {{block "content" .}}

    {{end}}
    </body>
    </html>
{{end}}
```

```gotemplate
{{template "base" .}}

{{define "content"}}
    <h1>Hello, world!</h1>
    <p>{{index .Data "payload"}}</p>
{{end}}
```

... then you can use code like this:

```go
package main

import (
	"fmt"
	"github.com/tsawler/page"
	"html/template"
	"log"
	"net/http"
)

const port = ":8080"

func main() {

	render := page.Render{
		TemplateDir: "./templates",
		TemplateMap: make(map[string]*template.Template),
		Functions:   template.FuncMap{},
		Partials:    []string{"base.layout.gohtml"},
		Debug:       true,
		UseCache:    true,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]any)
		data["payload"] = "This is passed data."
		err := render.Show(w, "home.page.gohtml", &page.Data{Data: data})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		}
	})

	http.HandleFunc("/string", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]any)
		data["payload"] = "This is passed data."
		out, err := render.String(w, "home.page.gohtml", &page.Data{Data: data})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		}
		log.Println(out)
		fmt.Fprint(w, "Check the console; you should see html")
	})

	log.Println("Starting on port", port)
	_ = http.ListenAndServe(port, nil)
}

```