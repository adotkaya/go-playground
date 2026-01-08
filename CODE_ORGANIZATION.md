# Code Organization Best Practices

This document explains the code organization structure used throughout this project, following Go best practices and conventions.

## File Organization Standard

### Standard Go File Structure

All `.go` files in this project follow this ordering:

```go
1. Package declaration
2. Import statements (grouped and ordered)
3. Constants
4. Variables
5. Types (structs, interfaces)
6. Constructor functions (New*, Make*)
7. Methods (grouped by receiver type)
8. Utility/helper functions
```

### Import Organization

Imports are organized in **three groups**, separated by blank lines:

```go
import (
    // 1. Standard library imports
    "context"
    "fmt"
    "net/http"
    "time"

    // 2. Third-party imports
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/justinas/alice"

    // 3. Local/internal imports
    "adotkaya.playground/internal/models"
    "adotkaya.playground/internal/validator"
)
```

**Benefits:**
- Easy to identify dependencies at a glance
- Reduces merge conflicts
- Follows Go community standards

---

## File-by-File Organization

### 1. `cmd/web/main.go`

**Purpose**: Application entry point and initialization

**Structure:**
```
Package declaration
Imports (stdlib → third-party → local)
═════════════════════════════════════
Application Structure
  - application type definition
═════════════════════════════════════
Main Function
  - Environment loading
  - Logger initialization
  - Configuration loading
  - Database connection
  - Template cache
  - Form decoder
  - Session manager
  - Application instance
  - TLS configuration
  - HTTP server setup
  - Server start
```

**Key Points:**
- Main function is clearly sectioned with comment dividers
- Each initialization step is clearly labeled
- Comments explain WHY not just WHAT

---

### 2. `cmd/web/config.go`

**Purpose**: Configuration management

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Configuration Types
  - Config struct
  - DatabaseConfig struct
  - ServerConfig struct
═════════════════════════════════════
Configuration Loading
  - LoadConfig() function
  - Validate() method
═════════════════════════════════════
Configuration Methods
  - DSN() method (connection string)
═════════════════════════════════════
Helper Functions
  - getEnvOrDefault()
  - parseDurationOrDefault()
```

**Key Points:**
- Types defined before functions
- Methods grouped with their receivers
- Private helpers at the end

---

### 3. `cmd/web/handlers.go`

**Purpose**: HTTP request handlers

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Form Types
  - SnippetCreateForm
  - userSignupForm
  - userLoginForm
═════════════════════════════════════
Public Handlers
  - ping() - health check
  - home() - homepage
═════════════════════════════════════
Snippet Handlers
  - snippetView()
  - snippetCreate()
  - snippetCreatePost()
═════════════════════════════════════
User Authentication Handlers
  - userSignup()
  - userSignupPost()
  - userLogin()
  - userLoginPost()
  - userLogoutPost()
```

**Key Points:**
- Form structs at the top (data definitions)
- Handlers grouped by functionality
- Logical ordering: public → snippets → users
- GET handlers before POST handlers

---

### 4. `cmd/web/middleware.go`

**Purpose**: HTTP middleware functions

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Security Middleware
  - secureHeaders() - security headers
  - noSurf() - CSRF protection
═════════════════════════════════════
Logging and Error Recovery Middleware
  - logRequest() - request logging
  - recoverPanic() - panic recovery
═════════════════════════════════════
Authentication Middleware
  - authenticate() - check auth status
  - requireAuthentication() - enforce auth
```

**Key Points:**
- Grouped by purpose (security, logging, auth)
- Comments explain what each middleware does
- Inline comments explain WHY for complex logic

---

### 5. `cmd/web/helpers.go`

**Purpose**: Helper utilities for handlers

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Template Data Helpers
  - newTemplateData()
═════════════════════════════════════
Error Handlers
  - serverError()
  - clientError()
  - notFound()
═════════════════════════════════════
Template Rendering
  - render()
═════════════════════════════════════
Authentication Helpers
  - isAuthenticated()
═════════════════════════════════════
Form Handling
  - decodePostForm()
```

**Key Points:**
- Organized by functionality
- Most commonly used first (template data)
- Clear separation between sections

---

### 6. `cmd/web/templates.go`

**Purpose**: Template management and caching

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Template Data Structure
  - templateData struct
═════════════════════════════════════
Template Functions
  - humanDate()
  - functions map
═════════════════════════════════════
Template Cache
  - newTemplateCache()
```

**Key Points:**
- Data types first
- Functions second
- Complex initialization last
- Inline comments explain field purposes

---

### 7. `cmd/web/routes.go`

**Purpose**: Route configuration and middleware chains

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Route Configuration
  - routes() function
    → Custom error handlers
    → Static file server
    → Health check route
    → Dynamic middleware chain
    → Public routes
    → Protected routes
    → Standard middleware chain
```

**Key Points:**
- Single routes() function with clear sections
- Comment blocks explain middleware order
- Routes grouped by access level (public/protected)
- Comments explain what each middleware does

---

### 8. `cmd/web/context.go`

**Purpose**: Request context key definitions

**Structure:**
```
Package declaration
═════════════════════════════════════
Request Context Keys
  - contextKey type
  - isAuthenticatedContextKey constant
```

**Key Points:**
- Small, focused file
- Custom type prevents key collisions
- Constants for context keys

---

### 9. `internal/models/errors.go`

**Purpose**: Custom error definitions

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Custom Error Definitions
  - ErrNoRecord
  - ErrInvalidCredentials
  - ErrDuplicateEmail
```

**Key Points:**
- Centralized error definitions
- Descriptive error messages
- Comments explain when each error is used

---

### 10. `internal/models/snippet.go`

**Purpose**: Snippet data model and database operations

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Snippet Model - Type Definitions
  - Snippet struct
  - SnippetModelInterface
  - SnippetModel struct
═════════════════════════════════════
Snippet Model - Methods
  - Insert() - with detailed comments
  - Get() - with detailed comments
  - Latest() - with detailed comments
```

**Key Points:**
- Interface before implementation
- Methods have detailed documentation
- SQL queries use clear formatting
- Comments explain business logic

---

### 11. `internal/models/users.go`

**Purpose**: User data model and authentication

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
User Model - Type Definitions
  - User struct
  - UserModelInterface
  - UserModel struct
═════════════════════════════════════
User Model - Methods
  - Insert() - with bcrypt hashing
  - Authenticate() - with validation
  - Exists() - simple check
```

**Key Points:**
- Same structure as snippet.go for consistency
- Security considerations in comments
- Error handling clearly documented

---

### 12. `internal/validator/validator.go`

**Purpose**: Form validation utilities

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Validator Type
  - Validator struct
═════════════════════════════════════
Email Regular Expression
  - EmailRX regex
═════════════════════════════════════
Validator Methods
  - Valid()
  - AddNonFieldError()
  - AddFieldError()
  - CheckField()
═════════════════════════════════════
Validation Functions
  - NotBlank()
  - MinChars()
  - MaxChars()
  - PermittedValue()
  - Matches()
```

**Key Points:**
- Type definition first
- Methods grouped together
- Standalone validation functions last
- Alphabetically ordered functions

---

### 13. `internal/assert/assert.go`

**Purpose**: Testing assertion helpers

**Structure:**
```
Package declaration
Imports
═════════════════════════════════════
Test Assertion Helpers
  - StringContains()
  - Equal()
  - NilError()
```

**Key Points:**
- Simple, focused file
- All test helpers in one place
- Consistent naming pattern

---

## Section Dividers

The project uses two types of comment dividers:

### Major Section Dividers (80 characters wide)
```go
// =============================================================================
// Section Name
// =============================================================================
```

Use for major sections like:
- Type definitions
- Groups of related functions
- Different areas of functionality

### Minor Section Dividers (75 characters)
```go
// -------------------------------------------------------------------------
// Subsection Name
// -------------------------------------------------------------------------
```

Use for subsections within main():
- Initialization steps
- Configuration blocks
- Related route groups

---

## Comments Best Practices

### 1. **Function Comments**

Always comment exported functions with:
- What the function does
- Important parameters
- What it returns
- Any side effects or important behavior

```go
// Insert creates a new snippet in the database
//
// Parameters:
//   - title: The snippet title (max 100 characters)
//   - content: The snippet code content
//   - expires: Number of days until expiration (1, 7, or 365)
//
// Returns the ID of the newly created snippet, or an error
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
```

### 2. **Inline Comments**

Use inline comments to explain:
- WHY something is done (not what)
- Complex logic
- Security considerations
- Non-obvious behavior

```go
// Only add if an error doesn't already exist for this field
if _, exists := v.FieldErrors[key]; !exists {
    v.FieldErrors[key] = message
}
```

### 3. **Type Comments**

Document what each type represents:

```go
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
```

---

## Naming Conventions

### Variables
- **Lowercase for private**: `templateCache`, `formDecoder`
- **Uppercase for exported**: `EmailRX`, `Validator`
- **Descriptive names**: `authenticatedUserID` not `uid`

### Functions
- **Action verbs**: `Insert()`, `Authenticate()`, `Exists()`
- **Getters without Get**: `Latest()` not `GetLatest()`
- **New for constructors**: `newTemplateCache()`, `newTestServer()`

### Types
- **PascalCase**: `SnippetModel`, `UserModelInterface`
- **Descriptive**: `templateData` not `td`

### Interfaces
- **-er suffix** for single method: `Handler`, `Writer`
- **Interface suffix** for multiple methods: `SnippetModelInterface`

---

## Grouping Related Code

### Example: Handlers

Instead of:
```go
// ❌ BAD: Mixed together
func (app *application) home() { }
func (app *application) userSignup() { }
func (app *application) snippetCreate() { }
func (app *application) userLogin() { }
```

Group by functionality:
```go
// ✅ GOOD: Organized by feature
// Public Handlers
func (app *application) home() { }

// Snippet Handlers
func (app *application) snippetView() { }
func (app *application) snippetCreate() { }
func (app *application) snippetCreatePost() { }

// User Handlers
func (app *application) userSignup() { }
func (app *application) userSignupPost() { }
func (app *application) userLogin() { }
func (app *application) userLoginPost() { }
```

---

## Reading Order

Files are organized so they can be read **top to bottom** naturally:

1. **Types** (what we're working with)
2. **Constants/Variables** (fixed values)
3. **Constructors** (how to create things)
4. **Methods** (what things can do)
5. **Helpers** (supporting functions)

This matches how you'd explain code to someone:
1. "Here's the data structure"
2. "Here are some constants we use"
3. "Here's how you create one"
4. "Here's what you can do with it"
5. "Here are some helper functions"

---

## Benefits of This Organization

1. **Easier Navigation**: Find things quickly by knowing where they should be
2. **Reduced Cognitive Load**: Similar code is grouped together
3. **Better Collaboration**: Team members know where to add new code
4. **Easier Reviews**: Reviewers can scan code more efficiently
5. **Maintainability**: Changes are localized to specific sections
6. **Onboarding**: New developers can understand the structure faster

---

## Quick Reference: File Order Template

```go
// 1. Package declaration
package main

// 2. Imports (stdlib → third-party → local)
import (
    "fmt"
    "net/http"

    "github.com/some/package"

    "yourproject/internal/models"
)

// 3. Constants
const (
    defaultPort = 4000
)

// 4. Variables
var (
    version = "1.0.0"
)

// 5. Types
type MyStruct struct {
    Field string
}

type MyInterface interface {
    Method() error
}

// 6. Constructors
func NewMyStruct() *MyStruct {
    return &MyStruct{}
}

// 7. Methods (grouped by receiver)
func (m *MyStruct) Method() error {
    return nil
}

func (m *MyStruct) AnotherMethod() {
    // ...
}

// 8. Helper functions
func helperFunction() {
    // ...
}
```

---

## Additional Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)

---

## Summary

This codebase now follows Go best practices for:
- ✅ Import organization (stdlib → third-party → local)
- ✅ Code ordering (types → constructors → methods → helpers)
- ✅ Logical grouping (related code together)
- ✅ Clear section dividers
- ✅ Comprehensive comments
- ✅ Consistent naming conventions
- ✅ Top-to-bottom readability

Every file in the project follows these patterns, making the codebase **consistent**, **maintainable**, and **professional**.
