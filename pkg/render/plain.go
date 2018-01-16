package render

const plainTemplate = `IDToken: {{ .IDToken.Value }}
RefreshToken: {{ .RefreshToken }}
`

func NewPlainRenderer() Renderer {
	renderer, err := NewTemplateRenderer(plainTemplate)
	if err != nil {
		panic(err)
	}
	return renderer
}
