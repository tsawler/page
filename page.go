package page

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Render struct {
	TemplateDir string
	Functions   template.FuncMap
	UseCache    bool
	TemplateMap map[string]*template.Template
	Partials    []string
}

type templateData struct {
	Data map[string]any
}

// Render generates a page of html from our template files
func (ren *Render) Render(w http.ResponseWriter, t string, td *templateData) {
	// declare a variable to hold the ready to execute template.
	var tmpl *template.Template

	// if we are using the cache, get try to get the pre-compiled template from our
	// map templateMap, stored in the receiver.
	if ren.UseCache {
		if templateFromMap, ok := ren.TemplateMap[t]; ok {
			//log.Println("getting template from map")
			tmpl = templateFromMap
		}
	}

	if tmpl == nil {
		newTemplate, err := ren.buildTemplateFromDisk(t)
		if err != nil {
			log.Println("Error building template:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl = newTemplate
	}

	// if we don't have template data, just use an empty struct.
	if td == nil {
		td = &templateData{}
	}

	// execute the template
	if err := tmpl.ExecuteTemplate(w, t, td); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (ren *Render) buildTemplateFromDisk(t string) (*template.Template, error) {

	templateSlice := append(ren.Partials, fmt.Sprintf("./%s/%s", ren.TemplateDir, t))

	tmpl, err := template.New(t).Funcs(ren.Functions).ParseFiles(templateSlice...)
	if err != nil {
		return nil, err
	}

	// Add the template to the template map stored in our receiver.
	// Note that this is ignored in development, but does not hurt anything.
	ren.TemplateMap[t] = tmpl

	return tmpl, nil
}
