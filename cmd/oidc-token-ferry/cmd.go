package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/twz123/oidc-token-ferry/pkg/api"
	"github.com/twz123/oidc-token-ferry/pkg/oidc"
)

type tokenFerryCmd struct {
	OIDCConfig oidc.Config `group:"OpenID Connect Options"`
	NoOpenURL  bool        `long:"no-open-url" description:"Don't open the redirect URL in a browser automatically"`
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

	if err := cmd.notifyUserAboutRedirectURL(challenge.RedirectURL()); err != nil {
		return nil, errors.Wrap(err, "failed to notify user about redirect URL")
	}

	code, err := obtainCodeFromUser()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain code")
	}

	credentials, err := challenge.ExchangeCode(code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange code for credentials")
	}

	return &api.TokenFerry{
		ClientID:     cmd.OIDCConfig.ClientID,
		ClientSecret: cmd.OIDCConfig.ClientSecret,
		Creds:        credentials,
	}, nil
}

func (cmd *tokenFerryCmd) notifyUserAboutRedirectURL(redirectURL string) error {
	var openCmd *exec.Cmd
	if !cmd.NoOpenURL {
		switch os := runtime.GOOS; os {
		case "darwin":
			openCmd = exec.Command("open", redirectURL)
		case "linux":
			openCmd = exec.Command("xdg-open", redirectURL)
		default:
			// unsupported OS
		}
	}

	if openCmd != nil {
		if err := openCmd.Run(); err == nil {
			return nil
		}
	}

	_, err := fmt.Fprintln(os.Stderr, "Open the following URL and authenticate with the IdP: ", redirectURL)
	return err
}

func obtainCodeFromUser() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprint(os.Stderr, "Enter the code provided by the IdP: ")
	code, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}
