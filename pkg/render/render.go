package render

import (
	"io"

	"github.com/twz123/oidc-token-ferry/pkg/oidc"
)

type Renderer interface {
	Render(w io.Writer, credentials *oidc.Credentials) error
}
