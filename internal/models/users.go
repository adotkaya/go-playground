package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// =============================================================================
// User Model - Type Definitions
// =============================================================================

// User represents a registered user account
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModelInterface defines the interface for user operations
type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

// UserModel wraps a database connection pool
type UserModel struct {
	DB *pgxpool.Pool
}

// =============================================================================
// User Model - Methods
// =============================================================================

// Insert creates a new user account in the database
//
// The password will be hashed using bcrypt (cost 12) before storage.
// Returns ErrDuplicateEmail if the email address is already in use.
func (m *UserModel) Insert(name, email, password string) error {
	// Hash the plain-text password using bcrypt with cost factor 12
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
             VALUES ($1, $2, $3, CURRENT_TIMESTAMP)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Attempt to insert the user record
	_, err = m.DB.Exec(ctx, stmt, name, email, string(hashedPassword))
	if err != nil {
		// Check if the error is a PostgreSQL unique constraint violation
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			// Error code 23505 is unique_violation
			// Check if it's specifically for the email constraint
			if pgError.Code == "23505" && strings.Contains(pgError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// Authenticate verifies user credentials and returns the user ID
//
// Returns ErrInvalidCredentials if the email doesn't exist or the password
// doesn't match. On success, returns the user's ID.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	// Retrieve the user ID and hashed password for the given email
	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No user found with this email
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			// Password doesn't match
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// Authentication successful
	return id, nil
}

// Exists checks whether a user with the given ID exists in the database
//
// Returns true if the user exists, false otherwise
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = $1)"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, stmt, id).Scan(&exists)
	return exists, err
}
