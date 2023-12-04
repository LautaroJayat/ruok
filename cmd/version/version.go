package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of ruok",
	Long:  `Print the version of ruok (which uses semver)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ruok service monitor v0.1")
	},
}
