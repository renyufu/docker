package image

import (
	"io"

	"golang.org/x/net/context"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type readmeOptions struct {
	image string
}

// NewSaveCommand creates a new `docker save` command
func NewReadmeCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts readmeOptions

	cmd := &cobra.Command{
		Use:   "readme IMAGE",
		Short: "Display /README.md of IMAGE",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.image = args[0]
			return runReadme(dockerCli, opts)
		},
	}

	return cmd
}

func runReadme(dockerCli *command.DockerCli, opts readmeOptions) error {

	responseBody, err := dockerCli.Client().ImageReadme(context.Background(), opts.image)
	if err != nil {
		return err
	}
	defer responseBody.Close()

	_, err = io.Copy(dockerCli.Out(), responseBody)
	return err
}
