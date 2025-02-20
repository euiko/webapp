package cli

import (
	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/settings"
	"github.com/spf13/cobra"
)

func Server(app api.App) func(settings *settings.Settings) api.Module {
	return func(settings *settings.Settings) api.Module {
		return api.NewModule(api.ModuleWithCLI(func(cmd *cobra.Command) {
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
}
