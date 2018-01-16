package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/twz123/oidc-token-ferry/pkg/oidc"
)

const (
	xOK = iota
	xGeneralError
	xCLIUsage
)

func main() {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGTERM, os.Interrupt)

	code, msg := run(osSignals)
	switch code {
	case xOK:
		return

	case xCLIUsage:
		fmt.Fprintln(os.Stderr, msg)
		flag.Usage()

	case xGeneralError:
		fmt.Fprintln(os.Stderr, msg)
	}

	os.Exit(code)
}

func run(osSignals <-chan os.Signal) (int, string) {
	var config oidc.Config
	flag.StringVar(&config.IssuerURL, "issuer-url", "https://accounts.google.com", "")
	flag.StringVar(&config.ClientID, "client-id", "", "")
	flag.StringVar(&config.ClientSecret, "client-secret", "", "")
	flag.Parse()

	if config.IssuerURL == "" {
		return xCLIUsage, "-issuer-url missing"
	}
	if config.ClientID == "" {
		return xCLIUsage, "-client-id missing"
	}
	if config.ClientSecret == "" {
		return xCLIUsage, "-client-secret missing"
	}

	credentials, err := performChallenge(&config)
	if err != nil {
		return xGeneralError, err.Error()
	}

	fmt.Println("IDToken: ", credentials.IDToken)
	fmt.Println("RefreshToken: ", credentials.RefreshToken)
	fmt.Println("EMail: ", credentials.EMail)

	return xOK, ""
}

func performChallenge(config *oidc.Config) (*oidc.Credentials, error) {

	flow, err := oidc.NewOpenIDConnectFlow(config)
	if err != nil {
		return nil, err
	}

	challenge, err := flow.InitiateChallenge()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initiate OpenID Connect challenge")
	}

	fmt.Println("Open URL: ", challenge.RedirectURL())

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the code: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)

	return challenge.ExchangeCode(code)
}
