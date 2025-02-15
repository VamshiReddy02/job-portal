package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vamshireddy02/job-portal/cmd/job" // Import the job subcommands
)

// RootCmd represents the base command when called without any subcommands
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "job",
		Short: "Job management CLI",
		Long:  `CLI for managing job listings using a GraphQL API backend.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Job CLI - Use 'job create', 'job update', or 'job delete'.")
		},
	}

	// Add subcommands from the `job` package
	rootCmd.AddCommand(job.CreateCmd)
	rootCmd.AddCommand(job.UpdateCmd)
	rootCmd.AddCommand(job.DeleteCmd)
	rootCmd.AddCommand(job.GetCmd)

	return rootCmd
}
