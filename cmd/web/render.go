package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type templateData struct {
	StringMap            map[string]string
	IntMap               map[string]int
	FloatMap             map[string]float32
	Data                 map[string]interface{}
	CSRFToken            string
	Flash                string
	Warning              string
	Error                string
	IsAuthenticated      int
	API                  string
	CSSVersion           string
	StripeSecretKey      string
	StripePublishableKey string
}

// function to pass as middleware when rendering pages
var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
}

func formatCurrency(n int) string {
	f := float32(n / 100)
	// to set to 2 decimal places
	return fmt.Sprintf("$%.2f", f)
}

// Package embed provides access to files embedded in the running Go program.
// Go source files that import "embed" can use the //go:embed directive to initialize a variable of type string, []byte, or FS with the contents of files read from the package directory or subdirectories at compile time.
// allow us to compile the app into binary with all embedded file, very convenient

//go:embed templates
var templateFS embed.FS

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.API = fmt.Sprintf("%s", app.config.api)
	td.StripePublishableKey = app.config.stripe.key

	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)

	// if template exist in cache templeInMap will be true otherwise false
	_, templateInMap := app.templeCache[templateToRender]

	// we want to use the cache in production only
	if app.config.env == "production" && templateInMap {
		t = app.templeCache[templateToRender]
	} else {
		// build the template
		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	// so we check to see if templateData was passed with the call to ??? at defaultdata
	// if it was not we create an empty templateData Object
	if td == nil {
		td = &templateData{}
	}

	// then we add our default data
	td = app.addDefaultData(td, r)

	// finally execute the template
	err = t.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

func (app *application) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	// build the partials
	// prepending templates to the partials name
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.gohtml", x)
		}
	}

	// if we do have partials we have to call our parseTemplate function
	// .Funcs(functions) pass some functions available to our template
	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", strings.Join(partials, ","), templateToRender)
	} else {
		// situation we don t have our partials associated with the template
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", templateToRender)
	}

	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	// set in cache
	app.templeCache[templateToRender] = t

	return t, nil
}
