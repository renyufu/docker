package image

import (
	"io"
        "archive/tar"

	"golang.org/x/net/context"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type readmeOptions struct {
	images []string
}

// NewSaveCommand creates a new `docker save` command
func NewReadmeCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts readmeOptions

	cmd := &cobra.Command{
		Use:   "readme IMAGE",
		Short: "Display /README.md of IMAGE",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.images = args
			return runReadme(dockerCli, opts)
		},
	}

	return cmd
}

func runReadme(dockerCli *command.DockerCli, opts readmeOptions) error {

	responseBody, err := dockerCli.Client().ImageReadme(context.Background(), opts.images)
	if err != nil {
		return err
	}
	defer responseBody.Close()

	tr := tar.NewReader(responseBody)

	// Iterate through the files in the archive.
	for {
		_, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if _, err := io.Copy(dockerCli.Out(), tr); err != nil {
			return err
		}
	}
	return nil
}
