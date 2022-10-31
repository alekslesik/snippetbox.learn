package mock

import (
	"errors"
	"time"

	"github.com/alekslesik/snippetbox.learn/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

// Rewrite all mysql.SnippetModel methods
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

type SnippetModelERR struct{}

// Rewrite all mysql.SnippetModel methods, return errors
func (m *SnippetModelERR) Insert(title, content, expires string) (int, error) {
	return 0, errors.New("test error Insert()")
}

func (m *SnippetModelERR) Get(id int) (*models.Snippet, error) {
	return &models.Snippet{}, errors.New("test error Get()")
}

func (m *SnippetModelERR) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{}, errors.New("test error Latest()")
}