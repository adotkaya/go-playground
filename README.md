# Go Snippetbox

A web application for sharing code snippets, built with Go and PostgreSQL.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- [Air](https://github.com/cosmtrek/air) (optional, for hot reload during development)

## Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/adotkaya/go-playground.git
cd go-playground
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up PostgreSQL database

Create a new PostgreSQL database:

```sql
CREATE DATABASE snippetbox;
```

### 4. Configure environment variables

Create a `.env` file in the root directory:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=snippetbox
DB_SSLMODE=disable
```

Replace `your_password_here` with your actual PostgreSQL password.

### 5. Run the application

**Using Air (with hot reload):**

```bash
# Install Air if you haven't already
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

**Without Air:**

```bash
go run ./cmd/web
```

The application will be available at `http://localhost:4000`

## Project Structure

```
.
├── cmd/
│   └── web/            # Application entry point and HTTP handlers
├── internal/
│   └── models/         # Database models and queries
├── ui/
│   ├── html/           # HTML templates
│   └── static/         # Static assets (CSS, JS, images)
├── .air.toml           # Air configuration for hot reload
├── go.mod              # Go module dependencies
└── go.sum              # Dependency checksums
```

## Development

The application uses Air for hot reload during development. Any changes to `.go`, `.tmpl`, or `.html` files will automatically rebuild and restart the server.

Build artifacts are stored in the `tmp/` directory and are excluded from version control.

## License

MIT
