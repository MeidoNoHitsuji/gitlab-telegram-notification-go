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
	//issue, _, err := client.Issue.Get("ES-3867", &jira.GetQueryOptions{})
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//if issue.Fields == nil {
	//	fmt.Println("Fields not found")
	//	return
	//}
	//
	//if issue.Fields.TimeTracking == nil {
	//	fmt.Println("TimeTracking not found")
	//} else {
	//	out, err := json.Marshal(issue.Fields.TimeTracking)
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		fmt.Println(string(out))
	//	}
	//}
	//
	//if issue.Fields.Worklog == nil {
	//	fmt.Println("Worklog not found")
	//} else {
	//	out, err := json.Marshal(issue.Fields.Worklog)
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		fmt.Println(string(out))
	//	}
	//}
}
