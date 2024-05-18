package page

import (
	"bytes"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var mapLock sync.Mutex

// Render is the main type for this package. Create a variable of this type
// and specify its fields, and you have access to the Show and String functions.
type Render struct {
	TemplateDir string                        // The path to templates.
	Functions   template.FuncMap              // A map of functions we want to pass to our templates.
	UseCache    bool                          // If true, use the template cache, stored in TemplateMap.
	TemplateMap map[string]*template.Template // Our template cache.
	Partials    []string                      // A list of partials.
	Debug       bool                          // Prints debugging info when true.
}

// New returns a Render type populated with sensible defaults.
func New() *Render {
	return &Render{
		Functions:   template.FuncMap{},
		UseCache:    true,
		TemplateMap: make(map[string]*template.Template),
		Partials:    []string{},
		Debug:       false,
	}
}

// Show generates a page of html from our template file(s).
func (ren *Render) Show(w http.ResponseWriter, t string, td any) error {
	// Call buildTemplate to get the template, either from the cache or by building it from disk.
	tmpl, err := ren.buildTemplate(t)
	if err != nil {
		log.Println("error building", err)
		return err
	}

	// Execute the template.
	if err := tmpl.ExecuteTemplate(w, t, td); err != nil {
		log.Println("error executing", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

// String renders a template and returns it as a string.
func (ren *Render) String(t string, td any) (string, error) {
	// Call buildTemplate to get the template, either from the cache or by building it
	// from disk.
	tmpl, err := ren.buildTemplate(t)
	if err != nil {
		return "", err
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

// GetTemplate attempts to get a template from the cache, builds it if it does not find it, and returns it.
func (ren *Render) GetTemplate(t string) (*template.Template, error) {
	// Call buildTemplate to get the template, either from the cache or by building it
	// from disk.
	tmpl, err := ren.buildTemplate(t)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
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
		log.Println("t", t)
		newTemplate, err := ren.buildTemplateFromDisk(t)
		if err != nil {
			log.Println("Error building from disk")
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
	templateSlice = append(templateSlice, ren.Partials...)

	// Append the template we want to render to the slice. Use path.Join to make it os agnostic.
	templateSlice = append(templateSlice, path.Join(ren.TemplateDir, t))

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

// LoadLayoutsAndPartials accepts a slice of strings which should consist of the types of files
// that are either layouts or partials for templates. For example, if a layout file is named
// `base.layout.gohtml` and a partial is named `footer.partial.gohtml`, then we would pass
//
//	[]string{".layout", ".partial"}
//
// Files anywhere in TemplateDir will be added the the Partials field of the Render type.
func (ren *Render) LoadLayoutsAndPartials(fileTypes []string) error {
	var templates []string
	for _, t := range fileTypes {
		files, err := addTemplate(ren.TemplateDir, t)
		if err != nil {
			return err
		}
		templates = append(templates, files...)
	}

	ren.Partials = templates
	return nil
}

func addTemplate(path, fileType string) ([]string, error) {
	files, err := find(path, ".gohtml")
	if err != nil {
		return nil, err
	}
	var templates []string
	for _, x := range files {
		if strings.Contains(x, fileType) {
			templates = append(templates, x)
		}
	}
	return templates, nil
}

func find(root, ext string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			files = append(files, s)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
