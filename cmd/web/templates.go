package main

import "github.com/alekslesik/snippetbox.learn/pkg/models"

type templateData struct {
    Snippet  *models.Snippet
    Snippets []*models.Snippet
}