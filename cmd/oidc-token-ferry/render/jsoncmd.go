package render

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/twz123/oidc-token-ferry/pkg/cli"
)

type JSONCmd struct {
	cli.OutputTarget
	cli.TokenFerryCmd
}

func (cmd *JSONCmd) Execute(args []string) error {
	ferry, err := cmd.TokenFerryCmd.PerformChallenge()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(ferry)
	if err != nil {
		return errors.Wrap(err, "failed to render JSON")
	}

	return cmd.OutputTarget.WriteBytes(bytes)
}
