package main

import (
	"context"
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"adotkaya.playground/internal/models"
)

// =============================================================================
// Application Structure
// =============================================================================

// application holds the application-wide dependencies and configuration
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

// =============================================================================
// Main Function
// =============================================================================

func main() {
	// -------------------------------------------------------------------------
	// Load Environment Configuration
	// -------------------------------------------------------------------------
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// -------------------------------------------------------------------------
	// Initialize Loggers
	// -------------------------------------------------------------------------
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// -------------------------------------------------------------------------
	// Load and Validate Configuration
	// -------------------------------------------------------------------------
	cfg, err := LoadConfig()
	if err != nil {
		errorLog.Fatal("Configuration error:", err)
	}

	// -------------------------------------------------------------------------
	// Initialize Database Connection
	// -------------------------------------------------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		errorLog.Fatal("Unable to connect to database:", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		errorLog.Fatal("Unable to ping database:", err)
	}
	infoLog.Println("Database connection established")

	// -------------------------------------------------------------------------
	// Initialize Template Cache
	// -------------------------------------------------------------------------
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// -------------------------------------------------------------------------
	// Initialize Form Decoder
	// -------------------------------------------------------------------------
	formDecoder := form.NewDecoder()

	// -------------------------------------------------------------------------
	// Initialize Session Manager
	// -------------------------------------------------------------------------
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// -------------------------------------------------------------------------
	// Create Application Instance
	// -------------------------------------------------------------------------
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: pool},
		users:          &models.UserModel{DB: pool},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// -------------------------------------------------------------------------
	// Configure TLS
	// -------------------------------------------------------------------------
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// -------------------------------------------------------------------------
	// Configure HTTP Server
	// -------------------------------------------------------------------------
	// Note: Always set IdleTimeout to prevent connections from being held open
	// indefinitely. ReadTimeout and WriteTimeout should also be set to protect
	// against slow-client attacks.
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// -------------------------------------------------------------------------
	// Start HTTPS Server
	// -------------------------------------------------------------------------
	infoLog.Printf("Starting server on :%s", cfg.Server.Port)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
