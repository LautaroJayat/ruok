package main

import (
	"fmt"
	"os"

	migrations "github.com/back-end-labs/ruok/cmd/migrate"
	"github.com/back-end-labs/ruok/cmd/scheduler"
	"github.com/back-end-labs/ruok/cmd/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ruok",
	Short: "ruok - make any postgres database into a backend service monitor",
	Long:  `Turn your postgres database into a backend service monitor. Receive notifications via http, slack, sqs/sns, and much more!`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("working!")
	},
}

func init() {

}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(scheduler.StartScheduler)
	rootCmd.AddCommand(migrations.SetupDB)
	execute()
}
