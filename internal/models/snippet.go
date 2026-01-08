package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// =============================================================================
// Snippet Model - Type Definitions
// =============================================================================

// Snippet represents a code snippet with metadata
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModelInterface defines the interface for snippet operations
type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

// SnippetModel wraps a database connection pool
type SnippetModel struct {
	DB *pgxpool.Pool
}

// =============================================================================
// Snippet Model - Methods
// =============================================================================

// Insert creates a new snippet in the database
//
// Parameters:
//   - title: The snippet title (max 100 characters)
//   - content: The snippet code content
//   - expires: Number of days until expiration (1, 7, or 365)
//
// Returns the ID of the newly created snippet, or an error
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
             VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + make_interval(days => $3))
             RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	err := m.DB.QueryRow(ctx, stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Get retrieves a specific snippet by ID
//
// Only returns snippets that have not expired. Returns ErrNoRecord if the
// snippet doesn't exist or has expired.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
             FROM snippets
             WHERE expires > CURRENT_TIMESTAMP AND id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s := &Snippet{}
	err := m.DB.QueryRow(ctx, stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return s, nil
}

// Latest retrieves the 10 most recently created snippets
//
// Only returns snippets that have not expired, ordered by creation date
// (most recent first).
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
             FROM snippets
             WHERE expires > CURRENT_TIMESTAMP
             ORDER BY id DESC
             LIMIT 10`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and build a slice of snippets
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
