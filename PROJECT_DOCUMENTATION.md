# Snippetbox - Comprehensive Project Documentation

## Table of Contents
- [Project Overview](#project-overview)
- [Architecture Overview](#architecture-overview)
- [Database Schemas](#database-schemas)
- [Data Models](#data-models)
- [Application Structure](#application-structure)
- [Request Flow & Workflows](#request-flow--workflows)
- [Component Dependencies](#component-dependencies)
- [Security Implementation](#security-implementation)
- [API Routes](#api-routes)
- [Configuration](#configuration)
- [Testing Strategy](#testing-strategy)
- [Deployment](#deployment)

---

## Project Overview

**Snippetbox** is a full-featured web application for sharing and managing code snippets built with Go 1.25.1.

**Module Name**: `adotkaya.playground`

**Key Features**:
- User authentication and session management
- Create, view, and manage code snippets
- Automatic snippet expiration (1 day, 1 week, 1 year)
- Secure HTTPS/TLS communication
- CSRF protection
- Responsive web interface
- PostgreSQL database backend

**Technology Stack**:
- **Language**: Go 1.25.1
- **Database**: PostgreSQL
- **Web Server**: HTTPS/TLS (port 4000)
- **Template Engine**: Go html/template
- **Session Store**: PostgreSQL-backed sessions
- **Password Hashing**: bcrypt (cost 12)

---

## Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT (Browser)                         │
│                         HTTPS/TLS Layer                          │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                    MIDDLEWARE CHAIN                              │
│  recoverPanic → logRequest → secureHeaders                       │
│      ↓                                                            │
│  LoadAndSave (session) → noSurf (CSRF) → authenticate            │
│      ↓                                                            │
│  requireAuthentication (for protected routes)                    │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                      HTTP ROUTER                                 │
│                  (httprouter)                                    │
│                                                                   │
│  Routes: /home, /snippet/*, /user/*, /static/*, /ping           │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                      HANDLERS LAYER                              │
│  (cmd/web/handlers.go)                                           │
│                                                                   │
│  • home                    • snippetView                         │
│  • snippetCreate           • snippetCreatePost                   │
│  • userSignup              • userSignupPost                      │
│  • userLogin               • userLoginPost                       │
│  • userLogoutPost          • ping                                │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                    VALIDATION LAYER                              │
│  (internal/validator)                                            │
│                                                                   │
│  • Form validation         • Field validation                    │
│  • Email regex             • Custom validators                   │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                      MODELS LAYER                                │
│  (internal/models)                                               │
│                                                                   │
│  SnippetModel              UserModel                             │
│  • Insert                  • Insert                              │
│  • Get                     • Authenticate                        │
│  • Latest                  • Exists                              │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                   DATABASE LAYER                                 │
│                   PostgreSQL                                     │
│  (pgx/v5 connection pool)                                        │
│                                                                   │
│  Tables: snippets, users, sessions                               │
└──────────────────────────────────────────────────────────────────┘
```

### Directory Structure

```
go-playground/
├── cmd/web/                    # Application entry point
│   ├── main.go                 # Main entry, server setup
│   ├── config.go               # Configuration management
│   ├── routes.go               # Route definitions
│   ├── handlers.go             # HTTP handlers
│   ├── middleware.go           # Middleware functions
│   ├── helpers.go              # Helper utilities
│   ├── templates.go            # Template management
│   ├── context.go              # Request context keys
│   └── *_test.go               # Handler tests
│
├── internal/
│   ├── models/                 # Data access layer
│   │   ├── snippet.go          # Snippet model
│   │   ├── users.go            # User model
│   │   ├── errors.go           # Custom errors
│   │   ├── *_test.go           # Model tests
│   │   ├── mocks/              # Mock implementations
│   │   │   ├── snippets.go
│   │   │   └── users.go
│   │   └── testdata/           # Test SQL scripts
│   │       ├── setup.sql
│   │       └── teardown.sql
│   │
│   ├── validator/              # Validation utilities
│   │   └── validator.go        # Form validators
│   │
│   └── assert/                 # Testing assertions
│       └── assert.go
│
├── ui/                         # User interface assets
│   ├── html/                   # HTML templates
│   │   ├── base.tmpl           # Base layout
│   │   ├── pages/              # Page templates
│   │   │   ├── home.tmpl
│   │   │   ├── view.tmpl
│   │   │   ├── create.tmpl
│   │   │   ├── signup.tmpl
│   │   │   └── login.tmpl
│   │   └── partials/           # Reusable partials
│   │       └── nav.tmpl
│   │
│   ├── static/                 # Static assets
│   │   ├── css/main.css
│   │   ├── js/main.js
│   │   └── img/
│   │
│   └── efs.go                  # Embedded filesystem
│
├── tls/                        # TLS certificates
│   ├── cert.pem
│   └── key.pem
│
├── .env                        # Environment variables
├── .air.toml                   # Hot reload config
├── go.mod                      # Go module definition
└── go.sum                      # Dependency checksums
```

---

## Database Schemas

### Tables Overview

The application uses 3 main tables:
1. **snippets** - Stores code snippets
2. **users** - Stores user accounts
3. **sessions** - Stores session data (managed by scs library)

### Schema: `snippets`

**Purpose**: Store user-created code snippets with expiration

```sql
CREATE TABLE snippets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);
```

**Columns**:
- `id` (SERIAL PRIMARY KEY): Auto-incrementing unique identifier
- `title` (VARCHAR(100) NOT NULL): Snippet title, max 100 characters
- `content` (TEXT NOT NULL): Snippet code content, unlimited length
- `created` (TIMESTAMP NOT NULL): Creation timestamp
- `expires` (TIMESTAMP NOT NULL): Expiration timestamp

**Indexes**:
- `idx_snippets_created`: B-tree index on `created` for efficient sorting

**Business Rules**:
- Snippets are soft-deleted (filtered by expires < NOW())
- Latest snippets query uses index for performance
- No foreign key to users (snippets can be anonymous)

### Schema: `users`

**Purpose**: Store registered user accounts with authentication

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
```

**Columns**:
- `id` (SERIAL PRIMARY KEY): Auto-incrementing unique identifier
- `name` (VARCHAR(255) NOT NULL): User's display name
- `email` (VARCHAR(255) NOT NULL): User's email address (unique)
- `hashed_password` (CHAR(60) NOT NULL): bcrypt hash (always 60 chars)
- `created` (TIMESTAMP NOT NULL): Account creation timestamp

**Constraints**:
- `users_uc_email`: UNIQUE constraint on email (enforces one account per email)

**Business Rules**:
- Email must be unique across all users
- Passwords are hashed with bcrypt cost 12
- Password hash is always exactly 60 characters

### Schema: `sessions`

**Purpose**: Store server-side session data (created by scs/pgxstore)

**Note**: This table is automatically managed by the `github.com/alexedwards/scs/pgxstore` library

```sql
CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);
```

**Columns**:
- `token` (TEXT PRIMARY KEY): Session token (stored in cookie)
- `data` (BYTEA NOT NULL): Serialized session data
- `expiry` (TIMESTAMPTZ NOT NULL): Session expiration time

**Indexes**:
- `sessions_expiry_idx`: Cleanup expired sessions efficiently

**Business Rules**:
- Sessions expire after 12 hours of inactivity
- Expired sessions cleaned up automatically by library
- Token stored in secure, httpOnly cookie

---

## Data Models

### Application Structure

```go
type application struct {
    errorLog       *log.Logger                    // Error logging
    infoLog        *log.Logger                    // Info logging
    snippets       models.SnippetModelInterface   // Snippet data access
    users          models.UserModelInterface      // User data access
    templateCache  map[string]*template.Template  // Compiled templates
    formDecoder    *form.Decoder                  // Form decoder
    sessionManager *scs.SessionManager            // Session manager
}
```

### Snippet Model

**File**: `internal/models/snippet.go`

```go
type Snippet struct {
    ID      int
    Title   string
    Content string
    Created time.Time
    Expires time.Time
}

type SnippetModelInterface interface {
    Insert(title string, content string, expires int) (int, error)
    Get(id int) (*Snippet, error)
    Latest() ([]*Snippet, error)
}
```

**Methods**:

1. **Insert(title, content, expires) → (id, error)**
   - Creates new snippet
   - `expires`: Days until expiration (1, 7, or 365)
   - Returns: Snippet ID
   - SQL: `INSERT INTO snippets ... RETURNING id`

2. **Get(id) → (*Snippet, error)**
   - Retrieves single snippet by ID
   - Filters out expired snippets
   - Returns: `ErrNoRecord` if not found or expired
   - SQL: `SELECT ... WHERE id = $1 AND expires > NOW()`

3. **Latest() → ([]*Snippet, error)**
   - Fetches 10 most recent non-expired snippets
   - Ordered by creation date (newest first)
   - SQL: `SELECT ... WHERE expires > NOW() ORDER BY id DESC LIMIT 10`

### User Model

**File**: `internal/models/users.go`

```go
type User struct {
    ID             int
    Name           string
    Email          string
    HashedPassword []byte
    Created        time.Time
}

type UserModelInterface interface {
    Insert(name, email, password string) error
    Authenticate(email, password string) (int, error)
    Exists(id int) (bool, error)
}
```

**Methods**:

1. **Insert(name, email, password) → error**
   - Registers new user account
   - Hashes password with bcrypt (cost 12)
   - Returns: `ErrDuplicateEmail` if email exists
   - SQL: `INSERT INTO users (name, email, hashed_password, created) VALUES ...`

2. **Authenticate(email, password) → (userID, error)**
   - Validates user credentials
   - Compares password with bcrypt hash
   - Returns: User ID on success
   - Returns: `ErrInvalidCredentials` on failure
   - SQL: `SELECT id, hashed_password FROM users WHERE email = $1`

3. **Exists(id) → (bool, error)**
   - Checks if user ID exists
   - Used for session validation
   - SQL: `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

### Custom Errors

**File**: `internal/models/errors.go`

```go
var (
    ErrNoRecord           = errors.New("models: no matching record found")
    ErrInvalidCredentials = errors.New("models: invalid credentials")
    ErrDuplicateEmail     = errors.New("models: this email is already signed up")
)
```

### Form Structures

**File**: `cmd/web/handlers.go`

```go
type SnippetCreateForm struct {
    Title               string `form:"title"`
    Content             string `form:"content"`
    Expires             int    `form:"expires"`
    validator.Validator        `form:"-"`
}

type userSignupForm struct {
    Name                string `form:"name"`
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator        `form:"-"`
}

type userLoginForm struct {
    Email               string `form:"email"`
    Password            string `form:"password"`
    validator.Validator        `form:"-"`
}
```

**Validation Rules**:

**SnippetCreateForm**:
- `Title`: Required, max 100 characters
- `Content`: Required
- `Expires`: Must be 1, 7, or 365

**userSignupForm**:
- `Name`: Required, max 255 characters
- `Email`: Required, max 255 characters, valid email format
- `Password`: Required, min 8 characters

**userLoginForm**:
- `Email`: Required, valid email format
- `Password`: Required

### Template Data Structure

**File**: `cmd/web/templates.go`

```go
type templateData struct {
    CurrentYear     int
    Snippet         *models.Snippet
    Snippets        []*models.Snippet
    Form            any
    Flash           string
    IsAuthenticated bool
    CSRFToken       string
}
```

**Fields**:
- `CurrentYear`: For footer copyright
- `Snippet`: Single snippet for view page
- `Snippets`: Snippet list for home page
- `Form`: Form data with validation errors
- `Flash`: One-time success/error message
- `IsAuthenticated`: User login status
- `CSRFToken`: CSRF protection token

---

## Application Structure

### Configuration Management

**File**: `cmd/web/config.go`

```go
type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
}

type DatabaseConfig struct {
    User     string
    Password string
    Host     string
    Port     string
    Name     string
    SSLMode  string
}

type ServerConfig struct {
    Port         string
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}
```

**Environment Variables**:
- `DB_USER` (required)
- `DB_PASSWORD` (required)
- `DB_NAME` (required)
- `DB_HOST` (default: "localhost")
- `DB_PORT` (default: "5432")
- `DB_SSLMODE` (default: "disable")
- `SERVER_PORT` (default: "4000")
- `SERVER_READ_TIMEOUT` (default: "5s")
- `SERVER_WRITE_TIMEOUT` (default: "10s")
- `SERVER_IDLE_TIMEOUT` (default: "1m")

**Example .env**:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=1331
DB_NAME=snippetbox
DB_SSLMODE=disable
```

### Initialization Flow

**File**: `cmd/web/main.go`

```
1. Load .env file (godotenv.Load)
   ↓
2. Load configuration (loadConfig)
   ↓
3. Initialize loggers (log.New)
   ↓
4. Open database connection (openDB)
   ↓
5. Initialize template cache (newTemplateCache)
   ↓
6. Initialize form decoder (form.NewDecoder)
   ↓
7. Initialize session manager (scs.New + pgxstore)
   ↓
8. Create application struct
   ↓
9. Setup routes (app.routes)
   ↓
10. Configure TLS (crypto/tls.Config)
   ↓
11. Start HTTPS server (srv.ListenAndServeTLS)
```

---

## Request Flow & Workflows

### HTTP Request Flow

```
Client Request (HTTPS)
    ↓
┌───────────────────────────────────────┐
│  Standard Middleware Chain            │
│  1. recoverPanic                      │
│  2. logRequest                        │
│  3. secureHeaders                     │
└────────────┬──────────────────────────┘
             ↓
┌────────────────────────────────────────┐
│  Router (httprouter)                   │
│  • Static files → FileServer           │
│  • /ping → ping handler                │
│  • Dynamic routes → Dynamic middleware │
│  • Protected routes → Protected chain  │
└────────────┬───────────────────────────┘
             ↓
┌────────────────────────────────────────┐
│  Dynamic Middleware                    │
│  1. sessionManager.LoadAndSave         │
│  2. noSurf (CSRF protection)           │
│  3. authenticate (load user from session)│
└────────────┬───────────────────────────┘
             ↓
┌────────────────────────────────────────┐
│  Protected Middleware (if protected)   │
│  1. requireAuthentication              │
│     (redirect to /user/login if not)   │
└────────────┬───────────────────────────┘
             ↓
┌────────────────────────────────────────┐
│  Handler Function                      │
│  • Parse form data                     │
│  • Validate input                      │
│  • Interact with models                │
│  • Render template or redirect         │
└────────────┬───────────────────────────┘
             ↓
┌────────────────────────────────────────┐
│  Response                              │
│  • HTML page (200, 422)                │
│  • Redirect (303, 302)                 │
│  • Error page (404, 500)               │
└────────────────────────────────────────┘
```

### Middleware Chain Details

**File**: `cmd/web/middleware.go`, `cmd/web/routes.go`

**Standard Chain** (all routes):
```go
alice.New(app.recoverPanic, app.logRequest, secureHeaders)
```

1. **recoverPanic**: Catches panics, logs them, returns 500 error
2. **logRequest**: Logs IP, protocol, method, URI
3. **secureHeaders**: Sets security headers (CSP, X-Frame-Options, etc.)

**Dynamic Chain** (public pages):
```go
alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
```

4. **LoadAndSave**: Loads session from cookie, saves changes after response
5. **noSurf**: Generates CSRF token, validates on POST/PUT/PATCH/DELETE
6. **authenticate**: Checks if user ID in session exists in DB

**Protected Chain** (authenticated only):
```go
dynamic.Append(app.requireAuthentication)
```

7. **requireAuthentication**: Redirects to /user/login if not authenticated

### User Registration Workflow

```
User submits signup form
    ↓
POST /user/signup
    ↓
1. Decode form data (formDecoder)
    ↓
2. Validate fields:
   • Name: NotBlank, MaxChars(255)
   • Email: NotBlank, Matches(EmailRX), MaxChars(255)
   • Password: NotBlank, MinChars(8)
    ↓
3. If invalid → re-render form with errors (422)
    ↓
4. Call users.Insert(name, email, password)
    ↓
5. Hash password with bcrypt (cost 12)
    ↓
6. Insert into database
    ↓
7. If duplicate email → add error, re-render (422)
    ↓
8. Success:
   • Add flash message: "Your signup was successful. Please log in."
   • Redirect to /user/login (303)
```

**Files**:
- Handler: `cmd/web/handlers.go:userSignupPost`
- Model: `internal/models/users.go:Insert`
- Validation: `internal/validator/validator.go`

### User Login Workflow

```
User submits login form
    ↓
POST /user/login
    ↓
1. Decode form data
    ↓
2. Validate fields:
   • Email: NotBlank, Matches(EmailRX)
   • Password: NotBlank
    ↓
3. If invalid → re-render form with errors (422)
    ↓
4. Call users.Authenticate(email, password)
    ↓
5. Retrieve hashed password from database
    ↓
6. Compare with bcrypt.CompareHashAndPassword
    ↓
7. If invalid → add non-field error, re-render (422)
    ↓
8. Success:
   • Renew session token (prevent session fixation)
   • Store "authenticatedUserID" in session
   • Redirect to /snippet/create (303)
```

**Files**:
- Handler: `cmd/web/handlers.go:userLoginPost`
- Model: `internal/models/users.go:Authenticate`

### User Logout Workflow

```
User clicks logout
    ↓
POST /user/logout (protected route)
    ↓
1. Renew session token
    ↓
2. Remove "authenticatedUserID" from session
    ↓
3. Add flash message: "You've been logged out successfully!"
    ↓
4. Redirect to / (303)
```

**File**: `cmd/web/handlers.go:userLogoutPost`

### Create Snippet Workflow

```
User visits create page
    ↓
GET /snippet/create (protected route)
    ↓
1. Render form with default expires=365
    ↓
User submits form
    ↓
POST /snippet/create (protected route)
    ↓
1. Decode form data
    ↓
2. Validate fields:
   • Title: NotBlank, MaxChars(100)
   • Content: NotBlank
   • Expires: PermittedValue(1, 7, 365)
    ↓
3. If invalid → re-render form with errors (422)
    ↓
4. Call snippets.Insert(title, content, expires)
    ↓
5. Calculate expiry: NOW() + expires days
    ↓
6. Insert into database, get ID
    ↓
7. Add flash message: "Snippet successfully created!"
    ↓
8. Redirect to /snippet/view/:id (303)
```

**Files**:
- Handlers: `cmd/web/handlers.go:snippetCreate`, `snippetCreatePost`
- Model: `internal/models/snippet.go:Insert`

### View Snippet Workflow

```
User clicks snippet link
    ↓
GET /snippet/view/:id
    ↓
1. Extract ID from URL params
    ↓
2. Validate ID is positive integer
    ↓
3. If invalid → return 404
    ↓
4. Call snippets.Get(id)
    ↓
5. Query: SELECT ... WHERE id = $1 AND expires > NOW()
    ↓
6. If not found or expired → return 404
    ↓
7. Render view template with snippet data
```

**Files**:
- Handler: `cmd/web/handlers.go:snippetView`
- Model: `internal/models/snippet.go:Get`
- Template: `ui/html/pages/view.tmpl`

### Home Page Workflow

```
User visits homepage
    ↓
GET /
    ↓
1. Call snippets.Latest()
    ↓
2. Query: SELECT ... WHERE expires > NOW()
          ORDER BY id DESC LIMIT 10
    ↓
3. Render home template with snippets list
```

**Files**:
- Handler: `cmd/web/handlers.go:home`
- Model: `internal/models/snippet.go:Latest`
- Template: `ui/html/pages/home.tmpl`

### Session Flow

```
First Request (no session)
    ↓
1. LoadAndSave middleware creates new session
2. Generates random token
3. Stores in sessions table
4. Sets cookie with token
    ↓
Subsequent Requests
    ↓
1. Browser sends session cookie
2. LoadAndSave reads token from cookie
3. Loads session data from database
4. Makes available to handlers
5. Saves any changes back to database
6. Updates cookie expiry
    ↓
Session Expiry
    ↓
1. After 12 hours of inactivity
2. Session deleted from database
3. Cookie becomes invalid
4. New session created on next request
```

---

## Component Dependencies

### Dependency Graph

```
main.go
  ├─> godotenv (load .env)
  ├─> config.go (configuration)
  ├─> pgxpool (database connection)
  ├─> models/snippet.go
  ├─> models/users.go
  ├─> templates.go (template cache)
  ├─> form.Decoder
  ├─> scs.SessionManager
  │     └─> pgxstore (PostgreSQL session store)
  ├─> routes.go
  │     ├─> httprouter
  │     ├─> alice (middleware chaining)
  │     ├─> middleware.go
  │     └─> handlers.go
  └─> http.Server (TLS configuration)

handlers.go
  ├─> models (snippets, users)
  ├─> validator
  ├─> form.Decoder
  ├─> sessionManager
  └─> templateCache

models/snippet.go
  └─> pgxpool.Pool

models/users.go
  ├─> pgxpool.Pool
  └─> bcrypt (password hashing)

middleware.go
  ├─> nosurf (CSRF)
  ├─> sessionManager
  └─> models/users (Exists check)

templates.go
  └─> html/template

validator/validator.go
  └─> regexp (email validation)
```

### External Dependencies

**From go.mod**:

**Direct Dependencies**:
```
github.com/jackc/pgx/v5 v5.8.0
  - PostgreSQL driver and connection pooling
  - Used by: models layer

github.com/alexedwards/scs/pgxstore
  - PostgreSQL session store adapter
  - Used by: main.go (session configuration)

github.com/alexedwards/scs/v2 v2.9.0
  - Session management library
  - Used by: main.go, middleware, handlers

github.com/go-playground/form/v4 v4.3.0
  - Form decoding library
  - Used by: handlers (form parsing)

github.com/joho/godotenv v1.5.1
  - .env file loading
  - Used by: main.go

github.com/julienschmidt/httprouter v1.3.0
  - High-performance HTTP router
  - Used by: routes.go

github.com/justinas/alice v1.2.0
  - Middleware chaining
  - Used by: routes.go

github.com/justinas/nosurf v1.2.0
  - CSRF protection middleware
  - Used by: middleware.go

golang.org/x/crypto v0.46.0
  - bcrypt password hashing
  - Used by: models/users.go
```

**Indirect Dependencies**:
```
github.com/jackc/pgpassfile v1.0.0
github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761
github.com/jackc/puddle/v2 v2.2.2
golang.org/x/sync v0.19.0
golang.org/x/text v0.32.0
```

### Standard Library Dependencies

```
cmd/web/:
  - crypto/tls (TLS configuration)
  - database/sql (DB interfaces)
  - errors
  - fmt
  - html/template (templating)
  - log (logging)
  - net/http (HTTP server)
  - os
  - runtime/debug (panic recovery)
  - strconv (string conversion)
  - time

internal/models/:
  - context
  - database/sql
  - errors
  - time

internal/validator/:
  - regexp
  - slices
  - strings
  - unicode/utf8
```

---

## Security Implementation

### 1. Authentication Security

**Password Hashing**:
- Algorithm: bcrypt
- Cost: 12 (2^12 iterations)
- Library: `golang.org/x/crypto/bcrypt`
- Hash length: 60 characters

```go
// Hash password on signup
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

// Verify password on login
err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
```

**Session Security**:
- Storage: PostgreSQL (server-side)
- Lifetime: 12 hours
- Cookie attributes:
  - `Secure: true` (HTTPS only)
  - `HttpOnly: true` (no JavaScript access)
  - `SameSite: Lax` (CSRF protection)
- Session fixation prevention: Token regenerated on login/logout

```go
sessionManager := scs.New()
sessionManager.Store = pgxstore.New(db)
sessionManager.Lifetime = 12 * time.Hour
sessionManager.Cookie.Secure = true
```

### 2. CSRF Protection

**Implementation**: `github.com/justinas/nosurf`

```go
func noSurf(next http.Handler) http.Handler {
    csrfHandler := nosurf.New(next)
    csrfHandler.SetBaseCookie(http.Cookie{
        HttpOnly: true,
        Path:     "/",
        Secure:   true,
    })
    return csrfHandler
}
```

**Protection Mechanism**:
- Double submit cookie pattern
- Token generated for each session
- Token validated on mutating requests (POST, PUT, PATCH, DELETE)
- Token embedded in forms as hidden field

**Template Integration**:
```html
<input type='hidden' name='csrf_token' value='{{.CSRFToken}}'>
```

### 3. Security Headers

**File**: `cmd/web/middleware.go:secureHeaders`

```go
w.Header().Set("Content-Security-Policy",
    "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "deny")
w.Header().Set("X-XSS-Protection", "0")
```

**Header Breakdown**:

| Header | Value | Purpose |
|--------|-------|---------|
| Content-Security-Policy | default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com | Prevents XSS by restricting resource sources |
| Referrer-Policy | origin-when-cross-origin | Controls referrer information |
| X-Content-Type-Options | nosniff | Prevents MIME-type sniffing |
| X-Frame-Options | deny | Prevents clickjacking |
| X-XSS-Protection | 0 | Disables legacy XSS filter (CSP preferred) |

### 4. TLS/HTTPS Configuration

**File**: `cmd/web/main.go`

```go
tlsConfig := &tls.Config{
    CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
}

srv := &http.Server{
    Addr:         ":" + cfg.Server.Port,
    ErrorLog:     errorLog,
    Handler:      app.routes(),
    TLSConfig:    tlsConfig,
    IdleTimeout:  cfg.Server.IdleTimeout,
    ReadTimeout:  cfg.Server.ReadTimeout,
    WriteTimeout: cfg.Server.WriteTimeout,
}

srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
```

**TLS Configuration**:
- Enforced HTTPS (no HTTP fallback)
- Preferred curves: X25519, P256
- Certificates in `./tls/` directory
- Timeouts configured to prevent slowloris attacks

### 5. SQL Injection Prevention

**Parameterized Queries**:
All database queries use parameterized statements via pgx driver

```go
// SAFE: Parameterized query
stmt := `SELECT id, title, content, created, expires
         FROM snippets WHERE id = $1 AND expires > NOW()`
row := m.DB.QueryRow(context.Background(), stmt, id)

// UNSAFE: String concatenation (NOT USED)
// query := "SELECT * FROM users WHERE email = '" + email + "'"
```

### 6. Input Validation

**File**: `internal/validator/validator.go`

**Validation Functions**:
- `NotBlank`: Prevents empty strings
- `MinChars/MaxChars`: Length validation
- `Matches`: Regex validation (email)
- `PermittedValue`: Whitelist validation

**Email Validation**:
```go
var EmailRX = regexp.MustCompile(
    "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)
```

### 7. Error Handling

**Secure Error Messages**:
- Generic error messages to users (prevent information disclosure)
- Detailed logging to server logs
- Stack traces only in logs, never to users

```go
func (app *application) serverError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    app.errorLog.Output(2, trace)  // Log with stack trace
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)  // Generic message to user
}
```

### 8. Panic Recovery

**File**: `cmd/web/middleware.go:recoverPanic`

```go
func (app *application) recoverPanic(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                w.Header().Set("Connection", "close")
                app.serverError(w, fmt.Errorf("%s", err))
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

### Security Checklist

- [x] HTTPS enforced
- [x] Strong password hashing (bcrypt cost 12)
- [x] Secure session management (server-side, secure cookies)
- [x] CSRF protection on all forms
- [x] SQL injection prevention (parameterized queries)
- [x] XSS prevention (template auto-escaping, CSP headers)
- [x] Clickjacking prevention (X-Frame-Options: deny)
- [x] Input validation (all user inputs)
- [x] Secure error handling (no information disclosure)
- [x] Panic recovery (graceful degradation)
- [x] Session fixation prevention (token renewal)
- [x] Unique email constraint (prevents account enumeration)

---

## API Routes

### Route Table

| Method | Path | Middleware | Handler | Description |
|--------|------|-----------|---------|-------------|
| GET | /ping | Standard | ping | Health check |
| GET | /static/* | Standard | FileServer | Static assets |
| GET | / | Standard + Dynamic | app.home | Homepage (snippet list) |
| GET | /snippet/view/:id | Standard + Dynamic | app.snippetView | View single snippet |
| GET | /user/signup | Standard + Dynamic | app.userSignup | Signup form |
| POST | /user/signup | Standard + Dynamic | app.userSignupPost | Process signup |
| GET | /user/login | Standard + Dynamic | app.userLogin | Login form |
| POST | /user/login | Standard + Dynamic | app.userLoginPost | Process login |
| GET | /snippet/create | Standard + Protected | app.snippetCreate | Create snippet form |
| POST | /snippet/create | Standard + Protected | app.snippetCreatePost | Process snippet creation |
| POST | /user/logout | Standard + Protected | app.userLogoutPost | Logout |

**Middleware Chains**:
- **Standard**: recoverPanic → logRequest → secureHeaders
- **Dynamic**: Standard + LoadAndSave → noSurf → authenticate
- **Protected**: Dynamic + requireAuthentication

### Route Details

#### GET /ping
**Purpose**: Health check endpoint
**Auth**: None
**Response**: 200 OK with "OK" body

```go
func ping(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}
```

#### GET /
**Purpose**: Homepage with latest snippets
**Auth**: Optional (shows different nav if authenticated)
**Response**: HTML page with snippet table

**Query**: Fetches 10 latest non-expired snippets
**Template**: `home.tmpl`

#### GET /snippet/view/:id
**Purpose**: View single snippet
**Auth**: Optional
**URL Param**: `:id` (integer)
**Response**: HTML page with snippet details or 404

**Validation**: ID must be positive integer
**Query**: Fetches snippet if not expired

#### GET /user/signup
**Purpose**: User registration form
**Auth**: None
**Response**: HTML signup form

**Form Fields**: name, email, password, csrf_token

#### POST /user/signup
**Purpose**: Process user registration
**Auth**: None
**Content-Type**: application/x-www-form-urlencoded
**Response**: 303 redirect to /user/login or 422 with errors

**Validation**:
- name: required, max 255 chars
- email: required, valid format, max 255 chars, unique
- password: required, min 8 chars

#### GET /user/login
**Purpose**: Login form
**Auth**: None
**Response**: HTML login form

**Form Fields**: email, password, csrf_token

#### POST /user/login
**Purpose**: Authenticate user
**Auth**: None
**Content-Type**: application/x-www-form-urlencoded
**Response**: 303 redirect to /snippet/create or 422 with errors

**Validation**:
- email: required, valid format
- password: required

**Side Effects**:
- Renews session token
- Sets authenticatedUserID in session

#### GET /snippet/create
**Purpose**: Create snippet form
**Auth**: Required
**Response**: HTML create form or 302 redirect to /user/login

**Form Fields**: title, content, expires (radio: 1, 7, 365), csrf_token

#### POST /snippet/create
**Purpose**: Create new snippet
**Auth**: Required
**Content-Type**: application/x-www-form-urlencoded
**Response**: 303 redirect to /snippet/view/:id or 422 with errors

**Validation**:
- title: required, max 100 chars
- content: required
- expires: must be 1, 7, or 365

#### POST /user/logout
**Purpose**: Logout user
**Auth**: Required
**Response**: 303 redirect to /

**Side Effects**:
- Renews session token
- Removes authenticatedUserID from session

---

## Configuration

### Environment Variables

**Required**:
- `DB_USER`: PostgreSQL username
- `DB_PASSWORD`: PostgreSQL password
- `DB_NAME`: Database name

**Optional** (with defaults):
- `DB_HOST`: Database host (default: "localhost")
- `DB_PORT`: Database port (default: "5432")
- `DB_SSLMODE`: SSL mode (default: "disable")
- `SERVER_PORT`: HTTP server port (default: "4000")
- `SERVER_READ_TIMEOUT`: Read timeout (default: "5s")
- `SERVER_WRITE_TIMEOUT`: Write timeout (default: "10s")
- `SERVER_IDLE_TIMEOUT`: Idle timeout (default: "1m")

### Database Setup

**1. Create Database**:
```sql
CREATE DATABASE snippetbox;
```

**2. Create Tables**:
```sql
-- Snippets table
CREATE TABLE snippets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);
CREATE INDEX idx_snippets_created ON snippets(created);

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL
);
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

-- Sessions table (managed by scs)
CREATE TABLE sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions(expiry);
```

### TLS Certificates

**Location**: `./tls/`

**Generate Self-Signed Certificate** (development):
```bash
cd tls
openssl req -new -x509 -sha256 -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
```

**Production**: Use Let's Encrypt or other CA-signed certificates

### Development Setup

**1. Clone Repository**:
```bash
git clone <repo-url>
cd go-playground
```

**2. Create .env File**:
```bash
cat > .env << EOF
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=snippetbox
DB_SSLMODE=disable
EOF
```

**3. Install Dependencies**:
```bash
go mod download
```

**4. Setup Database**:
```bash
psql -U postgres -c "CREATE DATABASE snippetbox;"
psql -U postgres -d snippetbox -f internal/models/testdata/setup.sql
```

**5. Generate TLS Certificates** (if not exists):
```bash
mkdir -p tls
cd tls
openssl req -new -x509 -sha256 -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
cd ..
```

**6. Run Application**:
```bash
go run ./cmd/web
```

**7. Access Application**:
```
https://localhost:4000
```

### Hot Reload (Development)

**Install Air**:
```bash
go install github.com/air-verse/air@latest
```

**Run with Hot Reload**:
```bash
air
```

**Configuration**: `.air.toml`
- Watches: `*.go`, `*.tmpl`, `*.html`
- Build: `go build -o ./tmp/main.exe ./cmd/web`
- Output: `tmp/main.exe`

---

## Testing Strategy

### Test Coverage

**Test Files**:
- `cmd/web/handlers_test.go` (184 lines)
- `cmd/web/middleware_test.go` (57 lines)
- `cmd/web/templates_test.go` (41 lines)
- `cmd/web/testutils_test.go` (116 lines)
- `internal/models/users_test.go` (44 lines)
- `internal/models/testutils_test.go` (42 lines)

**Total Test Code**: ~484 lines

### Test Structure

#### 1. Handler Tests

**File**: `cmd/web/handlers_test.go`

**Tests**:
- `TestPing`: Health check endpoint
- `TestSnippetView`: Snippet viewing with various IDs
- `TestUserSignup`: Registration with validation

**Test Cases**:
```go
tests := []struct {
    name       string
    urlPath    string
    wantCode   int
    wantBody   string
}{
    {"Valid ID", "/snippet/view/1", http.StatusOK, "An old silent pond"},
    {"Non-existent ID", "/snippet/view/2", http.StatusNotFound, ""},
    {"Negative ID", "/snippet/view/-1", http.StatusNotFound, ""},
    {"Decimal ID", "/snippet/view/1.23", http.StatusNotFound, ""},
    {"String ID", "/snippet/view/foo", http.StatusNotFound, ""},
    {"Empty ID", "/snippet/view/", http.StatusNotFound, ""},
}
```

#### 2. Middleware Tests

**File**: `cmd/web/middleware_test.go`

**Tests**:
- `TestSecureHeaders`: Verifies all security headers

```go
func TestSecureHeaders(t *testing.T) {
    // Test security header presence and values
    assert.Equal(t, rr.Header().Get("Content-Security-Policy"),
        "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
    assert.Equal(t, rr.Header().Get("Referrer-Policy"), "origin-when-cross-origin")
    assert.Equal(t, rr.Header().Get("X-Content-Type-Options"), "nosniff")
    assert.Equal(t, rr.Header().Get("X-Frame-Options"), "deny")
    assert.Equal(t, rr.Header().Get("X-XSS-Protection"), "0")
}
```

#### 3. Template Tests

**File**: `cmd/web/templates_test.go`

**Tests**:
- `TestHumanDate`: Time formatting with various timezones

```go
tests := []struct {
    name string
    tm   time.Time
    want string
}{
    {"UTC", time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC), "17 Dec 2020 at 10:00"},
    {"Empty", time.Time{}, ""},
    {"CET", time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)), "17 Dec 2020 at 09:00"},
}
```

#### 4. Model Tests

**File**: `internal/models/users_test.go`

**Tests**:
- `TestUserModelExists`: User existence validation

```go
tests := []struct {
    name   string
    userID int
    want   bool
}{
    {"Valid ID", 1, true},
    {"Zero ID", 0, false},
    {"Non-existent ID", 2, false},
}
```

### Test Utilities

#### Application Mock

**File**: `cmd/web/testutils_test.go`

```go
func newTestApplication(t *testing.T) *application {
    return &application{
        errorLog:      log.New(io.Discard, "", 0),
        infoLog:       log.New(io.Discard, "", 0),
        snippets:      &mocks.SnippetModel{},
        users:         &mocks.UserModel{},
        templateCache: newTemplateCache(),
        formDecoder:   form.NewDecoder(),
    }
}
```

#### Test Server

```go
func newTestServer(t *testing.T, h http.Handler) *testServer {
    ts := httptest.NewTLSServer(h)
    jar, _ := cookiejar.New(nil)
    ts.Client().Jar = jar
    ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    }
    return &testServer{ts}
}
```

#### CSRF Token Extraction

```go
func extractCSRFToken(t *testing.T, body string) string {
    csrfTokenRX := regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)
    matches := csrfTokenRX.FindStringSubmatch(body)
    return matches[1]
}
```

### Mock Models

**File**: `internal/models/mocks/snippets.go`

```go
type SnippetModel struct{}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
    return 2, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
    if id == 1 {
        return &Snippet{
            ID:      1,
            Title:   "An old silent pond",
            Content: "An old silent pond...",
            Created: time.Now(),
            Expires: time.Now(),
        }, nil
    }
    return nil, models.ErrNoRecord
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
    return []*Snippet{/* ... */}, nil
}
```

### Test Database

**File**: `internal/models/testutils_test.go`

```go
func newTestDB(t *testing.T) *pgxpool.Pool {
    db, _ := openDB(/* test DSN */)

    // Setup
    script, _ := os.ReadFile("./testdata/setup.sql")
    db.Exec(context.Background(), string(script))

    // Teardown
    t.Cleanup(func() {
        script, _ := os.ReadFile("./testdata/teardown.sql")
        db.Exec(context.Background(), string(script))
        db.Close()
    })

    return db
}
```

### Running Tests

**All Tests**:
```bash
go test ./...
```

**Specific Package**:
```bash
go test ./cmd/web
go test ./internal/models
```

**With Coverage**:
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Verbose Output**:
```bash
go test -v ./...
```

**Run Single Test**:
```bash
go test -run TestPing ./cmd/web
go test -run TestUserModelExists ./internal/models
```

---

## Deployment

### Production Checklist

#### 1. Environment Setup
- [ ] PostgreSQL database created
- [ ] Database tables created (snippets, users, sessions)
- [ ] Environment variables configured
- [ ] Valid TLS certificates installed
- [ ] Secure database credentials

#### 2. Security Configuration
- [ ] Change default database password
- [ ] Use strong, random session secrets
- [ ] Enable database SSL (`DB_SSLMODE=require`)
- [ ] Use CA-signed TLS certificates (not self-signed)
- [ ] Configure firewall rules
- [ ] Enable database connection pooling limits
- [ ] Set appropriate CORS policies (if needed)

#### 3. Application Configuration
- [ ] Set appropriate timeouts
- [ ] Configure logging to file/service
- [ ] Disable debug mode
- [ ] Set production-appropriate session lifetime
- [ ] Configure database connection pool size

#### 4. Infrastructure
- [ ] Reverse proxy (nginx/caddy) configured
- [ ] Process supervisor (systemd/supervisor) setup
- [ ] Log rotation configured
- [ ] Backup strategy implemented
- [ ] Monitoring and alerting setup
- [ ] Health check endpoint monitored

### Systemd Service Example

**File**: `/etc/systemd/system/snippetbox.service`

```ini
[Unit]
Description=Snippetbox Web Application
After=network.target postgresql.service

[Service]
Type=simple
User=snippetbox
WorkingDirectory=/opt/snippetbox
ExecStart=/opt/snippetbox/bin/web
Restart=always
RestartSec=10

Environment="DB_HOST=localhost"
Environment="DB_PORT=5432"
Environment="DB_USER=snippetbox"
Environment="DB_PASSWORD=<secure-password>"
Environment="DB_NAME=snippetbox"
Environment="DB_SSLMODE=require"
Environment="SERVER_PORT=4000"

[Install]
WantedBy=multi-user.target
```

**Commands**:
```bash
sudo systemctl daemon-reload
sudo systemctl enable snippetbox
sudo systemctl start snippetbox
sudo systemctl status snippetbox
```

### Nginx Reverse Proxy Example

**File**: `/etc/nginx/sites-available/snippetbox`

```nginx
server {
    listen 80;
    server_name example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name example.com;

    ssl_certificate /etc/letsencrypt/live/example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers 'ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384';

    location / {
        proxy_pass https://localhost:4000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Docker Deployment

**Dockerfile**:
```dockerfile
FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/web ./cmd/web

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /bin/web .
COPY --from=builder /app/tls ./tls
COPY --from=builder /app/ui ./ui
EXPOSE 4000
CMD ["./web"]
```

**docker-compose.yml**:
```yaml
version: '3.8'

services:
  web:
    build: .
    ports:
      - "4000:4000"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=snippetbox
      - DB_PASSWORD=secretpassword
      - DB_NAME=snippetbox
      - DB_SSLMODE=disable
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=snippetbox
      - POSTGRES_PASSWORD=secretpassword
      - POSTGRES_DB=snippetbox
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./internal/models/testdata/setup.sql:/docker-entrypoint-initdb.d/setup.sql
    restart: unless-stopped

volumes:
  postgres_data:
```

### Build for Production

**Optimized Build**:
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o web ./cmd/web
```

**Flags Explained**:
- `CGO_ENABLED=0`: Disable CGO for static binary
- `GOOS=linux`: Target Linux OS
- `-a`: Force rebuild of all packages
- `-installsuffix cgo`: Add suffix to package directory
- `-ldflags="-s -w"`: Strip debug info and symbol table
- `-o web`: Output binary name

### Database Migrations

**Production Migration Strategy**:
1. Create migration files (up/down SQL)
2. Use migration tool (e.g., `golang-migrate/migrate`)
3. Run migrations before deployment
4. Test rollback procedures

**Example Migration Tool Usage**:
```bash
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/snippetbox?sslmode=disable" up
```

### Monitoring

**Health Check Endpoint**:
```
GET https://example.com/ping
Expected: 200 OK "OK"
```

**Metrics to Monitor**:
- HTTP response times
- Database connection pool usage
- Error rates (5xx responses)
- Active sessions count
- Database query performance
- Memory usage
- CPU usage
- Disk usage (logs, database)

**Logging**:
- Application logs: `infoLog`, `errorLog`
- Access logs: via nginx/reverse proxy
- Database logs: PostgreSQL logs

### Backup Strategy

**Database Backups**:
```bash
# Daily backup
pg_dump -U snippetbox -h localhost snippetbox > backup_$(date +%Y%m%d).sql

# Automated backup script
#!/bin/bash
BACKUP_DIR="/backups/snippetbox"
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -U snippetbox snippetbox | gzip > "$BACKUP_DIR/snippetbox_$DATE.sql.gz"
find $BACKUP_DIR -name "snippetbox_*.sql.gz" -mtime +30 -delete
```

**Application Backups**:
- TLS certificates
- Environment configuration
- Static assets
- Binary executable

---

## Summary

**Snippetbox** is a well-architected, secure Go web application demonstrating best practices in:

- Clean architecture (separation of concerns)
- Security-first design (HTTPS, bcrypt, CSRF, CSP)
- Proper error handling and logging
- Comprehensive testing strategy
- Database connection pooling
- Session management
- Template-based rendering
- Middleware-based request processing

**Key Metrics**:
- **Lines of Code**: ~1,500 (application) + ~500 (tests)
- **Database Tables**: 3 (snippets, users, sessions)
- **HTTP Routes**: 11 total
- **Middleware Functions**: 6
- **Dependencies**: 9 external + standard library
- **Test Coverage**: All major handlers, models, and middleware

**Production-Ready Features**:
- TLS/HTTPS enforcement
- Graceful panic recovery
- Configurable timeouts
- Database connection pooling
- Server-side session storage
- Security headers
- CSRF protection
- Input validation
- Password hashing
- Health check endpoint
