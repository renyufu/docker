package container

import (
	"errors"
	"strings"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/spf13/pflag"
)

// NewCopyCommand creates a new `docker cp` command
func NewReadmeCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts copyOptions
	var createOpts createOptions
	var copts *containerOptions

	cmd := &cobra.Command{
		Use: `readme IMAGE`,
		Short: "Display README.md of an image",
		Long: strings.Join([]string{
			"Display README.md of an image\n",
			"README.md should been placed at /README.md",
		}, ""),
		Args: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("image can not be empty")
			}
			opts.source = args[0]
			return runReadme(dockerCli, cmd.Flags(), opts, &createOpts, copts)
		},
	}

	flags := cmd.Flags()

	flags.SetInterspersed(false)

	flags.StringVar(&createOpts.name, "name", "", "Assign a name to the container")

	flags.Bool("help", false, "Print usage")

	command.AddTrustVerificationFlags(flags)
	copts = addFlags(flags)

	flags.BoolVarP(&opts.followLink, "follow-link", "L", false, "Always follow symbol link in SRC_PATH")

	return cmd
}

func runReadme(dockerCli *command.DockerCli, flags *pflag.FlagSet, opts copyOptions, createOpts *createOptions, copts *containerOptions) error {
	//image := opts.source
        srcPath := "/README.md"
        dstPath := "-x"
	var rmOpts rmOptions

	ctx := context.Background()

        config, hostConfig, networkingConfig, err := parse(flags, copts)
        if err != nil {
                reportError(dockerCli.Err(), "create", err.Error(), true)
                return cli.StatusError{StatusCode: 125}
        }

        config.Image = opts.source

        response, err := createContainer(ctx, dockerCli, config, hostConfig, networkingConfig, hostConfig.ContainerIDFile, createOpts.name)
        if err != nil {
                return err
        }

        containerID := response.ID

        rmOpts.containers = make([]string, 1)
        rmOpts.containers[0] = containerID
        rmOpts.rmVolumes = true

	cpParam := &cpConfig{
		followLink: opts.followLink,
	}


	copyFromContainer(ctx, dockerCli, containerID, srcPath, dstPath, cpParam)
      
	var errs []string
	options := types.ContainerRemoveOptions{
		RemoveVolumes: rmOpts.rmVolumes,
		RemoveLinks:   rmOpts.rmLink,
		Force:         rmOpts.force,
	}

	errChan := parallelOperation(ctx, rmOpts.containers, func(ctx context.Context, container string) error {
		container = strings.Trim(container, "/")
		if container == "" {
			return errors.New("Container name cannot be empty")
		}
		return dockerCli.Client().ContainerRemove(ctx, container, options)
	})

	for range rmOpts.containers {
		if err := <-errChan; err != nil {
			errs = append(errs, err.Error())
			continue
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
