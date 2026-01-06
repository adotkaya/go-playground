package main

import (
	"context"
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"adotkaya.playground/internal/models"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Load and validate configuration
	cfg, err := LoadConfig()
	if err != nil {
		errorLog.Fatal("Configuration error:", err)
	}

	// Connect to database using the config
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		errorLog.Fatal("Unable to connect to database:", err)
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		errorLog.Fatal("Unable to ping database:", err)
	}
	infoLog.Println("Database connection established")

	//cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	//load decoder
	formDecoder := form.NewDecoder()
	//load sessions manager
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: pool},
		users:          &models.UserModel{DB: pool},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	//TLS
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//Always set IdleTimeout, if not
	//ReadTimeout, WriteTimeout or any Timeout
	//values will be set as default and would cut connection
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	infoLog.Printf("Starting server on :%s", cfg.Server.Port)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
