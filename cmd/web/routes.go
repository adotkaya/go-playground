package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"adotkaya.playground/ui"
)

// =============================================================================
// Route Configuration
// =============================================================================

// routes configures all application routes and middleware chains
func (app *application) routes() http.Handler {
	// Initialize router
	router := httprouter.New()

	// -------------------------------------------------------------------------
	// Custom Error Handlers
	// -------------------------------------------------------------------------

	// Handle 404 Not Found errors
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// -------------------------------------------------------------------------
	// Static File Server
	// -------------------------------------------------------------------------

	// Serve static files (CSS, JS, images) from embedded filesystem
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// -------------------------------------------------------------------------
	// Health Check Route
	// -------------------------------------------------------------------------

	// Health check endpoint (no middleware required)
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// -------------------------------------------------------------------------
	// Dynamic Middleware Chain
	// -------------------------------------------------------------------------
	// Applied to routes that need session management, CSRF protection, and
	// authentication checking (but don't require authentication)
	//
	// Middleware order:
	//   1. LoadAndSave - Load session data and save after response
	//   2. noSurf - CSRF token generation and validation
	//   3. authenticate - Check if user is authenticated and add to context

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// -------------------------------------------------------------------------
	// Public Routes (Dynamic Middleware)
	// -------------------------------------------------------------------------

	// Homepage
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	// View snippet (by ID)
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	// User signup
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))

	// User login
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// -------------------------------------------------------------------------
	// Protected Routes (Authentication Required)
	// -------------------------------------------------------------------------
	// These routes require the user to be authenticated. If not authenticated,
	// the user will be redirected to the login page.
	//
	// Additional middleware:
	//   4. requireAuthentication - Redirect to login if not authenticated

	protected := dynamic.Append(app.requireAuthentication)

	// Create snippet
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	// User logout
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// -------------------------------------------------------------------------
	// Standard Middleware Chain
	// -------------------------------------------------------------------------
	// Applied to ALL routes for core functionality
	//
	// Middleware order:
	//   1. recoverPanic - Recover from panics and return 500 error
	//   2. logRequest - Log all incoming requests
	//   3. secureHeaders - Add security headers to all responses

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the router wrapped in the standard middleware chain
	return standard.Then(router)
}
