package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/twz123/oidc-token-ferry/pkg/oidc"
	"github.com/twz123/oidc-token-ferry/pkg/render"
)

const (
	xOK = iota
	xGeneralError
	xCLIUsage
)

type renderOpts struct {
	GoTemplate string `long:"go-template" description:"Go Template used to render credentials"`
}

type opts struct {
	OIDCConfig oidc.Config `group:"OpenID Connect" namespace:"oidc"`
	Rendering  renderOpts  `group:"Rendering" namespace:"render"`
}

func main() {
	opts, _, err := parseOpts()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(xOK)
		} else {
			os.Exit(xCLIUsage)
		}
	}

	if err := run(opts); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(xGeneralError)
	}

	os.Exit(xOK)
}

func parseOpts() (*opts, []string, error) {
	var opts opts
	parser := flags.NewParser(&opts, flags.Default)
	additionalArgs, err := parser.Parse()
	if err != nil {
		return nil, nil, err
	}

	return &opts, additionalArgs, nil
}

func run(opts *opts) error {

	var renderer render.Renderer
	if opts.Rendering.GoTemplate == "" {
		renderer = render.NewPlainRenderer()
	} else {
		var err error
		renderer, err = render.NewTemplateRenderer(opts.Rendering.GoTemplate)
		if err != nil {
			return err
		}
	}

	credentials, err := performChallenge(&opts.OIDCConfig)

	if err != nil {
		return err
	}

	return renderer.Render(os.Stdout, credentials)
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
