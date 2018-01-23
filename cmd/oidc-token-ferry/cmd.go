package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/twz123/oidc-token-ferry/pkg/api"
	"github.com/twz123/oidc-token-ferry/pkg/oidc"
)

type tokenFerryCmd struct {
	OIDCConfig oidc.Config `group:"OpenID Connect Options"`
}

func (cmd *tokenFerryCmd) PerformChallenge() (*api.TokenFerry, error) {
	flow, err := oidc.NewOpenIDConnectFlow(&cmd.OIDCConfig)
	if err != nil {
		return nil, err
	}

	challenge, err := flow.InitiateChallenge()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initiate OpenID Connect challenge")
	}

	fmt.Fprintln(os.Stderr, "Open URL: ", challenge.RedirectURL())

	reader := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, "Enter the code: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)

	credentials, err := challenge.ExchangeCode(code)
	if err != nil {
		return nil, err
	}

	return &api.TokenFerry{
		ClientID:     cmd.OIDCConfig.ClientID,
		ClientSecret: cmd.OIDCConfig.ClientSecret,
		Creds:        credentials,
	}, nil
}
