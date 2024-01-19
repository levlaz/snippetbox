package models

import (
	"database/sql"
	"time"
)

// Snippet type holds data for an individual snippet
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel wraps sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// insert new snippet into DB
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	// note this is not supported by every DB, i.e. postgres does not have this
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

// get snippet by ID
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// Return 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
