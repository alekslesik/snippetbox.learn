package mysql

import (
	"database/sql"
	"errors"
	"github.com/alekslesik/snippetbox.learn/pkg/models"
)

// Determine type which wrap connect pool sql.DB
type SnippetModel struct {
	DB *sql.DB
}

// Create new snippet in database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// SQL request we wanted to execute
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use Exec() for execute SQL request
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Get the last created snippet ID from snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Return snippet data by ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// SQL request for getting data of one record
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use QueryRow() for executing SQL request passing unreliable variable ID like a placeholder
	row := m.DB.QueryRow(stmt, id)

	// Initialise the pointer to new struct Snippet
	s := &models.Snippet{}

	// Use row.Scan() to copy the value from every sql.Row field to Snippet Struct
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If all ok return Snippet object
	return s, nil
}

// Return last 10 snippets
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	// SQL request we wanted to execute
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	// Use Query() for execute SQL request
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*models.Snippet

	// Use rows.Next() to run over the result
	for rows.Next() {
		s := &models.Snippet{}
		// Use row.Scan() to copy the value from every sql.Row field to Snippet Struct
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Add the struct to slice
		snippets = append(snippets, s)
	}

	// Call rows.Err() after rows.Next() to ensure we haven't any errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If all ok return slice
	return snippets, nil
}
