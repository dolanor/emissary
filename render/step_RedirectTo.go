package render

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/benpate/derp"
)

// StepRedirectTo represents an action-step that forwards the user to a new page.
type StepRedirectTo struct {
	URL *template.Template
}

func (step StepRedirectTo) Get(renderer Renderer, buffer io.Writer) error {
	return step.redirect(renderer)
}

// Post updates the stream with approved data from the request body.
func (step StepRedirectTo) Post(renderer Renderer, buffer io.Writer) error {
	return step.redirect(renderer)
}

// Redirect returns an HTTP 307 Temporary Redirect that works for both GET and POST methods
func (step StepRedirectTo) redirect(renderer Renderer) error {

	const location = "render.StepRedirectTo.Redirect"
	var nextPage bytes.Buffer

	if err := step.URL.Execute(&nextPage, renderer); err != nil {
		return derp.Wrap(err, location, "Error evaluating 'url'")
	}

	return renderer.context().NoContent(http.StatusTemporaryRedirect)
}
