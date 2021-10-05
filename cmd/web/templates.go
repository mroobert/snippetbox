package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/mroobert/snippetbox/pkg/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
    CurrentYear int
    Snippet *models.Snippet
    Snippets []*models.Snippet
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
    return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
    "humanDate": humanDate,
}

// Each and every time we render a web page, our application could read and parse the relevant template files
// using the template.ParseFiles() function. We could avoid this duplicated work by parsing the files once
// — when starting the application — and storing the parsed templates in an in-memory cache.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
    // Initialize a new map to act as the cache.
    cache := map[string]*template.Template{}

    // Use the filepath.Glob function to get a slice of all filepaths with
    // the extension '.page.html'. This essentially gives us a slice of all the
    // 'page' templates for the application.
    pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
    if err != nil {
        return nil, err
    }

    for _, page := range pages {
        // Extract the file name (like 'home.page.html') from the full file path
        name := filepath.Base(page)

  
        //!!! The html/template package automatically escapes any data that is yielded between {{ }} tags.
	    //!!! This behavior is hugely helpful in avoiding cross-site scripting (XSS) attacks,
	    //!!! and is the reason that you should use the html/template package instead of the
	    //!!! more generic text/template package that Go also provides.

        // Register the FuncMap with the template
        templateWithFuncs := template.New(name).Funcs(functions)

        // Parse the page template file in to a template set.
        templateSet, err := templateWithFuncs.ParseFiles(page)
        if err != nil {
            return nil, err
        }

        // Use the ParseGlob method to add any 'layout' templates to the
        // template set (in our case, it's just the 'base' layout at the
        // moment).
        templateSet, err = templateSet.ParseGlob(filepath.Join(dir, "*.layout.html"))
        if err != nil {
            return nil, err
        }

        // Use the ParseGlob method to add any 'partial' templates to the
        // template set (in our case, it's just the 'footer' partial at the
        // moment).
        templateSet, err = templateSet.ParseGlob(filepath.Join(dir, "*.partial.html"))
        if err != nil {
            return nil, err
        }

        // Add the template set to the cache, using the name of the page
        // (like 'home.page.html') as the key.
        cache[name] = templateSet
    }

    return cache, nil
}