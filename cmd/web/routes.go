package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	/* 	If you’re not using the justinas/alice package to help manage your
	middleware chains, then you’d need to use the http.HandlerFunc()
	adapter to convert your handler functions like app.home to a
	http.Handler, and then wrap that with session middleware instead.
	Like this:
	router := httprouter.New()
	router.Handler(http.MethodGet, "/",
	app.sessionManager.LoadAndSave(http.HandlerFunc(app.home)))
	router.Handler(http.MethodGet, "/snippet/view/:id",
	app.sessionManager.LoadAndSave(http.HandlerFunc(app.snippetView)))
	*/

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	/*router.Handler(http.MethodPost, "/snippet/create", app.sessionManager.LoadAndSave(app.requireAuthentication(http.HandlerFunc(app.snippetCreate))))
	 * Without justinas/alice package
	 */

	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standart := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standart.Then(router)
}
