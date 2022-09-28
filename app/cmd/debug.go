package cmd

import (
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use: "debug",
	Run: debug,
}

func init() {
	rootCmd.AddCommand(debugCmd)
}

func debug(cmd *cobra.Command, args []string) {

}
