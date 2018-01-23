package render

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/twz123/oidc-token-ferry/pkg/cli"
)

type jsonCmd struct {
	cli.OutputTarget
	cli cli.CLI
}

func JsonCmd(cli cli.CLI) interface{} { return &jsonCmd{cli: cli} }

func (cmd *jsonCmd) Execute(args []string) error {
	ferry, err := cmd.cli.PerformChallenge()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(ferry)
	if err != nil {
		return errors.Wrap(err, "failed to render JSON")
	}

	return cmd.OutputTarget.WriteBytes(bytes)
}
