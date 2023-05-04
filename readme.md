# page

This is a simple package that makes rendering Go html templates as painless as possible. 

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
			log.Println(err)
		}
		log.Println(out)
		fmt.Fprint(w, "Check the console; you should see html")
	})

	log.Println("Starting on port", port)
	_ = http.ListenAndServe(port, nil)
}

```