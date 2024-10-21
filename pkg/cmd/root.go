package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gact",
	Short: "A cli tool to pull recent github activities of a User.",
	Long: `A cli tool to pull recent github activities of a User.

Pulls latest 10 activities of a user`,
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("user", "u", "", "GitHub Username")
	err := rootCmd.MarkPersistentFlagRequired("user")
	if err != nil {
		fmt.Printf("error in setting flag as mandatory: %s", err)
		return
	}
}
