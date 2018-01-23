package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"

	"github.com/twz123/oidc-token-ferry/cmd/oidc-token-ferry/kubeconfig"
	"github.com/twz123/oidc-token-ferry/cmd/oidc-token-ferry/render"
)

const (
	xOK = iota
	xGeneralError
	xCLIUsage
)

func main() {
	cli := &tokenFerryCmd{}

	parser := flags.NewParser(cli, flags.Default)

	cmd(parser, render.JsonCmd(cli), "render-json", "renders credentials as JSON")
	cmd(parser, render.GoTemplateCmd(cli), "render-go-template", "renders credentials using Go Templates")
	cmd(parser, kubeconfig.PatchCmd(cli), "patch-kubeconfig", "patches Kubernetes kubeconfig files")

	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(xOK)
			}

			os.Exit(xCLIUsage)
		}

		os.Exit(xGeneralError)
	}

	os.Exit(xOK)
}

func cmd(parser *flags.Parser, data interface{}, name, desc string) {
	if _, err := parser.AddCommand(name, desc, "", data); err != nil {
		panic(err)
	}
}
