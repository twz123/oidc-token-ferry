package render

import (
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/pkg/errors"

	"github.com/twz123/oidc-token-ferry/pkg/cli"
)

type goTemplateCmd struct {
	cli.OutputTarget

	Args struct {
		TemplateString string `positional-arg-name:"TEMPLATE_STRING" description:"Go Template to be rendered. An empty template indicates that the template is to be read from STDIN."`
	} `positional-args:"yes" required:"yes"`

	cli cli.CLI
}

func GoTemplateCmd(cli cli.CLI) interface{} { return &goTemplateCmd{cli: cli} }

func (cmd *goTemplateCmd) Execute(args []string) error {
	template, err := cmd.parseTemplate()
	if err != nil {
		return err
	}

	ferry, err := cmd.cli.PerformChallenge()
	if err != nil {
		return err
	}

	return cmd.OutputTarget.Write(func(writer io.Writer) error {
		return template.Execute(writer, ferry)
	})
}

func (cmd *goTemplateCmd) parseTemplate() (*template.Template, error) {
	if cmd.Args.TemplateString == "" {
		return parseTemplateFromStdin()
	}

	return parseTemplateText(cmd.Args.TemplateString)
}

func parseTemplateFromStdin() (*template.Template, error) {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from STDIN")
	}

	return parseTemplateText(string(bytes))
}

func parseTemplateText(text string) (*template.Template, error) {
	template, err := template.New("template").Parse(text)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse Go Template")
	}
	return template, nil
}
