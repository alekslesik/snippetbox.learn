package mock

import (
	"errors"

	"github.com/alekslesik/snippetbox.learn/pkg/models"
)

type SnippetModelERR struct{}

// Rewrite all mysql.SnippetModel methods
func (m *SnippetModelERR) Insert(title, content, expires string) (int, error) {
	return 0, errors.New("test error Insert()")
}

func (m *SnippetModelERR) Get(id int) (*models.Snippet, error) {
	return &models.Snippet{}, errors.New("test error Get()")
}

func (m *SnippetModelERR) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{}, errors.New("test error Latest()")
}
