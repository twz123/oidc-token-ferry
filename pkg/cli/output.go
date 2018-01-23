package cli

import (
	"io"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type OutputTarget struct {
	Filename flags.Filename `short:"o" long:"output-file" description:"Output file to write (defaults to STDOUT if omitted)"`
}

func (target *OutputTarget) WriteBytes(bytes []byte) error {
	return target.Write(func(writer io.Writer) error {
		_, err := writer.Write(bytes)
		return err
	})
}

func (target *OutputTarget) Write(write func(io.Writer) error) error {
	out := os.Stdout
	if target.Filename != "" {
		out, err := os.Create(string(target.Filename))
		if err != nil {
			return errors.Wrap(err, "failed to create output file")
		}
		defer out.Close()
	}

	return write(out)
}
