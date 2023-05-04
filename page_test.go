package page

import (
	"net/http/httptest"
	"testing"
)

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
	{name: "valid: no data", template: "home.page.gohtml", useData: false, useCache: false, errorExpected: false},
	{name: "invalid: no template", template: "x.page.gohtml", useData: false, errorExpected: true},
	{name: "invalid: bad template", template: "bad.page.gohtml", useData: false, errorExpected: true},
}

func TestRender_Show(t *testing.T) {
	p := New()
	p.TemplateDir = "./testdata/templates"
	p.Debug = true
	p.Partials = []string{"base.layout.gohtml"}

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

var stringTests = []struct {
	name          string
	template      string
	useData       bool
	useCache      bool
	errorExpected bool
}{
	{name: "valid", template: "home.page.gohtml", useData: true, useCache: true, errorExpected: false},
	{name: "valid from cache", template: "home.page.gohtml", useData: true, useCache: true, errorExpected: false},
	{name: "valid: no data", template: "home.page.gohtml", useData: false, useCache: false, errorExpected: false},
	{name: "invalid: no template", template: "x.page.gohtml", useData: false, errorExpected: true},
	{name: "invalid: bad template", template: "bad.page.gohtml", useData: false, errorExpected: true},
}

func TestRender_String(t *testing.T) {
	p := New()
	p.TemplateDir = "./testdata/templates"
	p.Debug = true
	p.Partials = []string{"base.layout.gohtml"}

	for _, e := range stringTests {
		rr := httptest.NewRecorder()
		data := make(map[string]any)
		p.UseCache = e.useCache
		data["payload"] = "This is passed data."

		if e.useData {
			s, err := p.String(rr, e.template, &Data{Data: data})
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
			s, err := p.String(rr, e.template, nil)
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
