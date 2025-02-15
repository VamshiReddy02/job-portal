package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	title       string
	description string
	company     string
	url         string
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new job listing",
	Run: func(cmd *cobra.Command, args []string) {
		data := map[string]interface{}{
			"query": `mutation CreateJob($input: CreateJobListingInput!) {
				createJobListing(input: $input) {
					_id
					title
					company
				}
			}`,
			"variables": map[string]interface{}{
				"input": map[string]string{
					"title":       title,
					"description": description,
					"company":     company,
					"url":         url,
				},
			},
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Request failed:", err)
			return
		}
		defer resp.Body.Close()

		fmt.Println("Job created successfully!")
	},
}

func init() {
	CreateCmd.Flags().StringVarP(&title, "title", "t", "", "Job title")
	CreateCmd.Flags().StringVarP(&description, "description", "d", "", "Job description")
	CreateCmd.Flags().StringVarP(&company, "company", "c", "", "Company name")
	CreateCmd.Flags().StringVarP(&url, "url", "u", "", "Job URL")
	CreateCmd.MarkFlagRequired("title")
	CreateCmd.MarkFlagRequired("description")
	CreateCmd.MarkFlagRequired("company")
	CreateCmd.MarkFlagRequired("url")
}
