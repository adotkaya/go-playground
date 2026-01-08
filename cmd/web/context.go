package main

// =============================================================================
// Request Context Keys
// =============================================================================

// contextKey is a custom type for request context keys to avoid collisions
type contextKey string

// isAuthenticatedContextKey is used to store/retrieve authentication status
// from the request context
const isAuthenticatedContextKey = contextKey("isAuthenticated")
