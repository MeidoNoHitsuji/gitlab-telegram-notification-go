package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab-telegram-notification-go/gitclient"
)

var debugCmd = &cobra.Command{
	Use: "debug",
	Run: debug,
}

func init() {
	rootCmd.AddCommand(debugCmd)
}

func debug(cmd *cobra.Command, args []string) {
	commits, err := gitclient.GetCommitsLastPipeline(338, "2d5a0f79e238ee704c015520909719445f33692a", "a2f7d5f5698425a4a5d8d630c11389b7c723b5e8")
	if err != nil {
		fmt.Println(err)
	}

	for _, commit := range commits {
		fmt.Println(commit.Title)
	}
}
