package models

import (
	"database/sql"
	"errors"
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
	// keep in mind this statement depends on db as well, postgres will use $N instead of ?
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
	query := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(query, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// Return 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	// this is required to make sure that all connections are now used up if something goes wrong.
	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// dont assume successful iteration was completed over entire result set
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
