package cli

import (
	"github.com/twz123/oidc-token-ferry/pkg/api"
)

type CLI interface {
	PerformChallenge() (*api.TokenFerry, error)
}
