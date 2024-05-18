package page

import (
	"html/template"
	"net/http/httptest"
	"strings"
	"testing"
)

type Data struct {
	Data map[string]any
}

func Test_New(t *testing.T) {
	page := New()
	if page.Debug {
		t.Error("expected page.Debug to be false, but it's true")
	}
}

var showTests = []struct {
	name          string
	template      string
	useData       bool
	useCache      bool
	errorExpected bool
}{
	{name: "valid", template: "home.page.gohtml", useData: true, useCache: true, errorExpected: false},
	{name: "valid from cache", template: "home.page.gohtml", useData: true, useCache: true, errorExpected: false},
	{name: "valid: no data", template: "nodata.page.gohtml", useData: false, useCache: false, errorExpected: false},
	{name: "invalid: no template", template: "x.page.gohtml", useData: false, errorExpected: true},
	{name: "invalid: bad template", template: "bad.page.gohtml", useData: false, errorExpected: true},
}

func TestRender_Show(t *testing.T) {
	p := New()
	p.TemplateDir = "./testdata/templates"
	p.Debug = true
	p.Partials = []string{"./testdata/templates/base.layout.gohtml"}

	for _, e := range showTests {
		rr := httptest.NewRecorder()
		data := make(map[string]any)
		p.UseCache = e.useCache
		data["payload"] = "This is passed data."

		if e.useData {
			err := p.Show(rr, e.template, &Data{Data: data})
			if err != nil && !e.errorExpected {
				t.Errorf("%s: failed to render template: %v", e.name, err)
			}
			if err == nil && e.errorExpected {
				t.Errorf("%s: expected an error but did not get one", e.name)
			}
		} else {
			err := p.Show(rr, e.template, nil)
			if err != nil && !e.errorExpected {
				t.Errorf("%s: failed to render template: %v", e.name, err)
			}
			if err == nil && e.errorExpected {
				t.Errorf("%s: expected an error but did not get one", e.name)
			}
		}
	}
}

func TestRender_GetTemplate(t *testing.T) {
	p := New()
	p.TemplateDir = "./testdata/templates"
	p.Debug = true
	p.Partials = []string{"./testdata/templates/base.layout.gohtml"}

	_, err := p.GetTemplate("home.page.gohtml")
	if err != nil {
		t.Error("error getting template:", err.Error())
	}

	_, err = p.GetTemplate("bad.page.gohtml")
	if err == nil {
		t.Error("expected error but did not get one")
	}
}

var stringTests = []struct {
	name          string
	template      string
	useData       bool
	useCache      bool
	errorExpected bool
}{
	{name: "valid", template: "home.page.gohtml", useData: true, useCache: true, errorExpected: false},
	{name: "valid from cache", template: "home.page.gohtml", useData: true, useCache: true, errorExpected: false},
	{name: "valid: no data", template: "nodata.page.gohtml", useData: false, useCache: false, errorExpected: false},
	{name: "invalid: no template", template: "x.page.gohtml", useData: false, errorExpected: true},
	{name: "invalid: bad template", template: "bad.page.gohtml", useData: false, errorExpected: true},
}

func TestRender_String(t *testing.T) {
	p := New()
	p.TemplateDir = "./testdata/templates"
	p.Debug = true
	p.Partials = []string{"./testdata/templates/base.layout.gohtml"}

	for _, e := range stringTests {
		data := make(map[string]any)
		p.UseCache = e.useCache
		data["payload"] = "This is passed data."

		if e.useData {
			s, err := p.String(e.template, &Data{Data: data})
			if err != nil && !e.errorExpected {
				t.Errorf("%s: failed to render template: %v", e.name, err)
			}
			if err == nil && e.errorExpected {
				t.Errorf("%s: expected an error but did not get one", e.name)
			}
			if len(s) == 0 && !e.errorExpected {
				t.Errorf("%s: no html returned", e.name)
			}
		} else {
			s, err := p.String(e.template, nil)
			if err != nil && !e.errorExpected {
				t.Errorf("%s: failed to render template: %v", e.name, err)
			}
			if err == nil && e.errorExpected {
				t.Errorf("%s: expected an error but did not get one", e.name)
			}
			if len(s) == 0 && !e.errorExpected {
				t.Errorf("%s: no html returned", e.name)
			}
		}
	}
}

func Test_withFuncMap(t *testing.T) {
	p := New()
	p.TemplateDir = "./testdata/templates"
	p.Partials = []string{"./testdata/templates/base.layout.gohtml"}
	fm := template.FuncMap{
		"foo": func() string {
			return "bar"
		},
	}
	p.Functions = fm

	s, err := p.String("with_func.page.gohtml", nil)
	if err != nil {
		t.Error("error rendering string:", err)
	}

	if !strings.Contains(s, "bar") {
		t.Error("did not find bar in rendered template:\n", s)
	}
}

func Test_addTemplate(t *testing.T) {
	files, err := addTemplate("./testdata/templates", ".layout")
	if err != nil {
		t.Error("error calling addTemplate:", err)
	}

	if len(files) != 1 {
		t.Error("wrong number of files in slice")
	}

	files, err = addTemplate("./nonexistent/templates", ".layout")
	if err == nil {
		t.Error("expected error but did not get one")
	}
}

func TestRender_LoadLayoutsAndPartials(t *testing.T) {
	p := New()
	p.TemplateDir = "./testdata/templates"
	p.Debug = true

	err := p.LoadLayoutsAndPartials([]string{".layout"})
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if len(p.Partials) != 1 {
		t.Error("wrong number of files in partials")
	}

	p.TemplateDir = "./nonexistent/templates"
	err = p.LoadLayoutsAndPartials([]string{".layout"})
	if err == nil {
		t.Error("expected error but did not get one")
	}

}
