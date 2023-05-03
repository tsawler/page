package page

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Render is the main type for this package. Create a variable of this type
// and specify its fields, and you have access to the Show function.
type Render struct {
	TemplateDir string                        // the directory where go templates are stored.
	Functions   template.FuncMap              // a map of functions we want to pass to our templates.
	UseCache    bool                          // if true, use the template cache, stored in TemplateMap.
	TemplateMap map[string]*template.Template // our template cache.
	Partials    []string                      // a list of partials; these are stored in TemplateDir.
}

// Data is a struct to hold any data that is to be passed to a template.
type Data struct {
	Data map[string]any
}

// Show generates a page of html from our template file(s).
func (ren *Render) Show(w http.ResponseWriter, t string, td *Data) {
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

	// tmpl will be nil if we do not have a value in the map (our template cache). In this case,
	// we build the template from disk.
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
		td = &Data{}
	}

	// execute the template
	if err := tmpl.ExecuteTemplate(w, t, td); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// buildTemplateFromDisk builds a template from disk.
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
