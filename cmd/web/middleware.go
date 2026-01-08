package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// =============================================================================
// Security Middleware
// =============================================================================

// secureHeaders adds security headers to all HTTP responses
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content Security Policy: Restricts where resources can be loaded from
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		// Referrer Policy: Controls referrer information
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// X-Content-Type-Options: Prevents MIME-type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options: Prevents clickjacking by denying framing
		w.Header().Set("X-Frame-Options", "deny")

		// X-XSS-Protection: Disable legacy XSS filter (rely on CSP instead)
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// noSurf provides CSRF protection for all state-changing requests
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true, // Prevent JavaScript access
		Path:     "/",
		Secure:   true, // HTTPS only
	})
	return csrfHandler
}

// =============================================================================
// Logging and Error Recovery Middleware
// =============================================================================

// logRequest logs details about each HTTP request
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

// recoverPanic recovers from panics and returns a 500 Internal Server Error
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Deferred function will run in the event of a panic
		defer func() {
			if err := recover(); err != nil {
				// Set connection close header to trigger Go's HTTP server
				// to automatically close the current connection
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// =============================================================================
// Authentication Middleware
// =============================================================================

// authenticate checks if a user is authenticated and adds info to request context
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve authenticated user ID from session
		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			// User not authenticated
			next.ServeHTTP(w, r)
			return
		}

		// Check if user still exists in database
		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// If user exists, add isAuthenticated flag to request context
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// requireAuthentication redirects unauthenticated users to the login page
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Set Cache-Control header to prevent browsers from caching pages
		// that require authentication
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}
