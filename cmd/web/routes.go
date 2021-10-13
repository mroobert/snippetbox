package main

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.home)

	router.HandlerFunc(http.MethodPost, "/snippet", app.createSnippet)
	router.HandlerFunc(http.MethodGet, "/snippet", app.showSnippetForm)
	router.HandlerFunc(http.MethodGet, "/snippet/:id", app.showSnippet)

	// Use the router.ServeFiles() function to register the file server as the handler for
	// all URL paths that start with "/static/".
	router.ServeFiles("/static/*filepath", unexposedFileSystem{http.Dir("./ui/static/")})

	//without alice package
	//return app.recoverPanic(app.logRequest(secureHeaders(router)))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standardMiddleware.Then(router)
}

type unexposedFileSystem struct {
	fs http.FileSystem
}

func (ufs unexposedFileSystem) Open(path string) (http.File, error) {
	file, err := ufs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := ufs.fs.Open(index); err != nil {
			closeErr := file.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return file, nil
}
