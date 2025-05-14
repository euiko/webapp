package cli

import (
	"github.com/euiko/webapp/core"
	"github.com/euiko/webapp/settings"
	"github.com/spf13/cobra"
)

func Server(app core.App) core.Module {
	return core.NewModule(core.ModuleWithCLI(func(cmd *cobra.Command, _ *settings.Settings) {
		startCmd := cobra.Command{
			Use:   "start",
			Short: "Start the web application",
			RunE: func(cmd *cobra.Command, args []string) error {
				return app.Start(cmd.Context())
			},
		}
		cmd.AddCommand(&startCmd)
	}))

}
