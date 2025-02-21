package api

import "github.com/spf13/cobra"

type (
	CLI interface {
		Command(cmd *cobra.Command)
	}
)
