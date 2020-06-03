package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"

	"github.com/twz123/oidc-token-ferry/cmd/oidc-token-ferry/kubeconfig"
	"github.com/twz123/oidc-token-ferry/cmd/oidc-token-ferry/render"
	"github.com/twz123/oidc-token-ferry/cmd/oidc-token-ferry/version"
)

const (
	xOK = iota
	xGeneralError
	xCLIUsage
)

type cli struct {
	VersionCmd    version.VersionCmd   `command:"version" description:"Show oidc-token-ferry version information"`
	JSONCmd       render.JSONCmd       `command:"render-json" description:"renders credentials as JSON"`
	GoTemplateCmd render.GoTemplateCmd `command:"render-go-template" description:"renders credentials using Go Templates"`
	PatchCmd      *kubeconfig.PatchCmd `command:"patch-kubeconfig" description:"patches Kubernetes kubeconfig files"`
}

func main() {
	parser := flags.NewParser(&cli{PatchCmd: kubeconfig.NewPatchCmd()}, flags.Default)

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
