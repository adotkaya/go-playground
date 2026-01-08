package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"adotkaya.playground/internal/models"
	"adotkaya.playground/ui"
)

// =============================================================================
// Template Data Structure
// =============================================================================

// templateData holds dynamic data that we want to pass to HTML templates
type templateData struct {
	CurrentYear     int               // For copyright year in footer
	Snippet         *models.Snippet   // Single snippet for view page
	Snippets        []*models.Snippet // Multiple snippets for home page
	Form            any               // Form data with validation errors
	Flash           string            // One-time flash message
	IsAuthenticated bool              // User authentication status
	CSRFToken       string            // CSRF protection token
}

// =============================================================================
// Template Functions
// =============================================================================

// humanDate formats a time.Time object into a human-readable string
func humanDate(t time.Time) string {
	// Return empty string for zero time
	if t.IsZero() {
		return ""
	}

	// Convert to UTC and format as "DD MMM YYYY at HH:MM"
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// functions is a map of custom template functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// =============================================================================
// Template Cache
// =============================================================================

// newTemplateCache creates a cache of all templates
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Get all page templates from the embedded filesystem
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Parse each page template along with base and partials
	for _, page := range pages {
		// Extract the filename (e.g., 'home.tmpl') from the full path
		name := filepath.Base(page)

		// Define the patterns for parsing: base + partials + page
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// Parse the template files with custom functions
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the page name as the key
		cache[name] = ts
	}

	return cache, nil
}
