package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(bumpUpCmd())
}

func bumpUpCmd() *cobra.Command {
	return &cobra.Command{
		Use: "bump_up",
		Short: "Bump up a ruby gem version",
		Run: BumpUp,
	}
}

func BumpUp(cmd *cobra.Command, args []string) {
	rootCmd.Println("I did my best to bump up the ruby gem version!")
}