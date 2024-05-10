package page

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

var mapLock sync.Mutex

// Render is the main type for this package. Create a variable of this type
// and specify its fields, and you have access to the Show and String functions.
type Render struct {
	TemplateDir string                        // The directory where go templates are stored.
	Functions   template.FuncMap              // A map of functions we want to pass to our templates.
	UseCache    bool                          // If true, use the template cache, stored in TemplateMap.
	TemplateMap map[string]*template.Template // Our template cache.
	Partials    []string                      // A list of partials; these are assumed to be stored in TemplateDir.
	Debug       bool                          // Prints debugging info when true.
}

// Data is a struct to hold any data that is to be passed to a template.
type Data struct {
	Data map[string]any
}

// New returns a Render type populated with sensible defaults.
func New() *Render {
	return &Render{
		TemplateDir: "./templates",
		Functions:   template.FuncMap{},
		UseCache:    true,
		TemplateMap: make(map[string]*template.Template),
		Partials:    []string{},
		Debug:       false,
	}
}

// Show generates a page of html from our template file(s).
func (ren *Render) Show(w http.ResponseWriter, t string, td *Data) error {
	// Call buildTemplate to get the template, either from the cache or by building it
	// from disk.
	tmpl, err := ren.buildTemplate(t)
	if err != nil {
		return err
	}

	// if we don't have template data, just use an empty struct.
	if td == nil {
		td = &Data{}
	}

	// execute the template
	if err := tmpl.ExecuteTemplate(w, t, td); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

// String renders a template and returns it as a string.
func (ren *Render) String(t string, td *Data) (string, error) {
	// Call buildTemplate to get the template, either from the cache or by building it
	// from disk.
	tmpl, err := ren.buildTemplate(t)
	if err != nil {
		return "", err
	}

	// if we don't have template data, just use an empty struct.
	if td == nil {
		td = &Data{}
	}

	// Execute the template, storing the result in a bytes.Buffer variable.
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, td); err != nil {
		return "", err
	}

	// Return a string from the bytes.Buffer.
	result := tpl.String()
	return result, nil
}

// buildTemplate is a utility function that creates a template, either from the cache, or from
// disk. The template is ready to accept functions & data, and then get rendered.
func (ren *Render) buildTemplate(t string) (*template.Template, error) {
	// tmpl is the variable that will hold our template.
	var tmpl *template.Template

	// If we are using the cache, get try to get the pre-compiled template from our
	// map templateMap, stored in the receiver.
	if ren.UseCache {
		if templateFromMap, ok := ren.TemplateMap[t]; ok {
			if ren.Debug {
				log.Println("Reading template", t, "from cache")
			}
			tmpl = templateFromMap
		}
	}

	// At this point, tmpl will be nil if we do not have a value in the map (our template
	// cache). In this case, we build the template from disk.
	if tmpl == nil {
		newTemplate, err := ren.buildTemplateFromDisk(t)
		if err != nil {
			return nil, err
		}
		tmpl = newTemplate
	}

	return tmpl, nil
}

// buildTemplateFromDisk builds a template from disk.
func (ren *Render) buildTemplateFromDisk(t string) (*template.Template, error) {
	// templateSlice will hold all the templates necessary to build our finished template.
	var templateSlice []string

	// Read in the partials, if any.
	for _, x := range ren.Partials {
		// We use filepath.Join to make this OS-agnostic.
		path := filepath.Join(ren.TemplateDir, x)
		templateSlice = append(templateSlice, path)
	}

	// Append the template we want to render to the slice.
	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", ren.TemplateDir, t))

	// Create a new template by parsing all the files in the slice.
	tmpl, err := template.New(t).Funcs(ren.Functions).ParseFiles(templateSlice...)
	if err != nil {
		return nil, err
	}

	// Add the template to the template map stored in our receiver.
	// Note that this is ignored in development, but does not hurt anything.
	mapLock.Lock()
	ren.TemplateMap[t] = tmpl
	mapLock.Unlock()

	if ren.Debug {
		log.Println("Reading template", t, "from disk")
	}

	return tmpl, nil
}
