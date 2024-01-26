package main

import "snippetbox.levlaz.org/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
