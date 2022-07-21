package main

import "github.com/alekslesik/snippetbox.learn/pkg/models"

// Добавляем поле Snippets в структуру templateData
type templateData struct {
    Snippet  *models.Snippet
    Snippets []*models.Snippet
}