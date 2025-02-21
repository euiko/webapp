package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/settings"
	"github.com/spf13/cobra"
)

func Config() api.Module {
	return api.NewModule(api.ModuleWithCLI(func(cmd *cobra.Command, s *settings.Settings) {
		cmd.AddCommand(configCmd(s))
	}))
}

var configWriter = map[string]func(*settings.Settings, io.Writer) error{
	"yaml": func(s *settings.Settings, w io.Writer) error {
		return settings.Write(s, settings.FormatYaml, w)
	},
	"json": func(s *settings.Settings, w io.Writer) error {
		return settings.Write(s, settings.FormatJson, w)
	},
}

func configCmd(s *settings.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage the configuration",
	}
	cmd.AddCommand(configGetCmd(s))
	cmd.AddCommand(configWriteCmd(s))
	return cmd
}

func configGetCmd(s *settings.Settings) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			writer, ok := configWriter[format]
			if !ok {
				return fmt.Errorf("unsupported format: %s", format)
			}

			if err := writer(s, os.Stdout); err != nil {
				log.Fatal("error when writing configuration", log.WithError(err))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format")
	return cmd
}

func configWriteCmd(s *settings.Settings) *cobra.Command {
	var (
		format string
		output string
	)
	cmd := &cobra.Command{
		Use:   "write [FILE]",
		Short: "Write the configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			writer, ok := configWriter[format]
			if !ok {
				return fmt.Errorf("unsupported format: %s", format)
			}

			// set default output when not being set
			if output == "" && len(os.Args) > 0 {
				execName := os.Args[0]
				output = execName + "." + format
			}

			f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()

			if err := writer(s, f); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file (default to <exec_name>.<format> when not being set)")

	return cmd
}
