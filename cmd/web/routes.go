package main

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes.
	dynamicMiddleware := alice.New(app.session.Enable, app.checkCSRF, app.authenticate)

	router := httprouter.New()

	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.home))

	router.Handler(http.MethodPost, "/snippet", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	router.Handler(http.MethodGet, "/snippet", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showSnippetForm))
	router.Handler(http.MethodGet, "/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	router.Handler(http.MethodGet, "/user/signup", dynamicMiddleware.ThenFunc(app.showSignupForm))
	router.Handler(http.MethodPost, "/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	router.Handler(http.MethodGet, "/user/login", dynamicMiddleware.ThenFunc(app.showLoginForm))
	router.Handler(http.MethodPost, "/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	router.Handler(http.MethodPost, "/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

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
