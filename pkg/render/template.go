package render

import (
	"io"
	"text/template"

	"github.com/pkg/errors"
	"github.com/twz123/oidc-token-ferry/pkg/oidc"
)

// Prepare some data to insert into the template.
type Recipient struct {
	Name, Gift string
	Attended   bool
}

type templateRenderer struct {
	template *template.Template
}

func NewTemplateRenderer(text string) (Renderer, error) {

	// Create a new template and parse the letter into it.
	template, err := template.New("credentials").Parse(text)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	return &templateRenderer{template}, nil
}

func (renderer *templateRenderer) Render(w io.Writer, credentials *oidc.Credentials) error {
	return renderer.template.Execute(w, credentials)
}
