package cli

import (
	"fmt"
	"net/http"

	"github.com/aquaproj/aqua/pkg/config"
	"github.com/aquaproj/aqua/pkg/controller"
	"github.com/urfave/cli/v2"
)

func (runner *Runner) newListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Usage:  "List packages in Registries",
		Action: runner.listAction,
		Description: `Output the list of packages in registries.
The output format is <registry name>,<package name>

e.g.
$ aqua list
standard,99designs/aws-vault
standard,abiosoft/colima
standard,abs-lang/abs
...
`,
	}
}

func (runner *Runner) listAction(c *cli.Context) error {
	param := &config.Param{}
	logE, err := runner.setParam(c, param)
	if err != nil {
		return fmt.Errorf("parse the command line arguments: %w", err)
	}
	ctrl := controller.InitializeListCommandController(c.Context, param, http.DefaultClient)
	return ctrl.List(c.Context, param, logE) //nolint:wrapcheck
}
