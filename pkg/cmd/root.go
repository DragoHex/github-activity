package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/DragoHex/github-activity.git/pkg/github"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gact",
	Short: "A cli tool to pull recent github activities of a User.",
	Long: `A cli tool to pull recent github activities of a User.

Pulls latest 10 activities of a user`,
	Run: func(cmd *cobra.Command, args []string) {
		user, err := cmd.Flags().GetString("user")
		if err != nil {
			fmt.Printf("error in reading the user: %s", err)
			return
		}

		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			fmt.Printf("error in reading the limit: %s", err)
			return
		}

		githubActivities := github.NewGitHubEvents(user, limit)
		err = githubActivities.GetActivity()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(githubActivities.ProcessEvents())
	},
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
	rootCmd.Flags().IntP("limit", "l", 10, "max number of activities to be listed")
	err := rootCmd.MarkPersistentFlagRequired("user")
	if err != nil {
		fmt.Printf("error in setting flag as mandatory: %s", err)
		return
	}
}
