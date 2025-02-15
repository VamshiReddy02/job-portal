package main

import (
	"log"
	"os"

	"github.com/vamshireddy02/job-portal/cmd"
)

func main() {
	rootCmd := cmd.RootCmd()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
