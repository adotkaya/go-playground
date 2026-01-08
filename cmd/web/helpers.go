package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// =============================================================================
// Template Data Helpers
// =============================================================================

// newTemplateData creates a templateData struct populated with common data
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

// =============================================================================
// Error Handlers
// =============================================================================

// serverError logs the error with a stack trace and sends a 500 response
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends a specific HTTP status code and corresponding description
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound is a convenience wrapper around clientError which sends a 404
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// =============================================================================
// Template Rendering
// =============================================================================

// render renders a template with the given data and status code
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the appropriate template from the cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Write template to a buffer first to catch any errors before writing to response
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the status code and buffered content to the response
	w.WriteHeader(status)
	buf.WriteTo(w)
}

// =============================================================================
// Authentication Helpers
// =============================================================================

// isAuthenticated checks if the current request is from an authenticated user
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

// =============================================================================
// Form Handling
// =============================================================================

// decodePostForm decodes POST form data into a destination struct
//
// Note: app.formDecoder.Decode() requires non-nil pointers. If a nil pointer
// is passed, it will return form.InvalidDecodeError which we panic on since
// this indicates a developer error rather than a user error.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Decode the form data into the destination struct
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// Check if the error is due to an invalid decode operation (developer error)
		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
		return err
	}

	return nil
}
