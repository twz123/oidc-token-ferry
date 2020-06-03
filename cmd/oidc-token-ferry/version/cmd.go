package version

import (
	"fmt"
	"runtime"
)

type Version struct {
	Version    string
	GoVersion  string
	GoCompiler string
	GoOs       string
	GoArch     string
}

func NewVersion() Version {
	return Version{
		Version:   VERSION,
		GoVersion: runtime.Compiler,
		GoOs:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
	}
}

type VersionCmd struct {
}

func (cmd *VersionCmd) Execute(args []string) error {
	_, err := fmt.Printf("oidc-token-ferry: %+v\n", NewVersion())
	return err
}
