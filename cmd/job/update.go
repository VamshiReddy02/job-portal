package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	JobID string
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing job listing",
	Run: func(cmd *cobra.Command, args []string) {
		data := map[string]interface{}{
			"query": `mutation UpdateJob($id: ID!, $input: UpdateJobListingInput!) {
				updateJobListing(id: $id, input: $input) {
					_id
					title
					company
				}
			}`,
			"variables": map[string]interface{}{
				"id": JobID,
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

		fmt.Println("Job updated successfully!")
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&JobID, "id", "i", "", "Job ID")
	UpdateCmd.Flags().StringVarP(&title, "title", "t", "", "New job title")
	UpdateCmd.Flags().StringVarP(&description, "description", "d", "", "New job description")
	UpdateCmd.Flags().StringVarP(&company, "company", "c", "", "New company name")
	UpdateCmd.Flags().StringVarP(&url, "url", "u", "", "New job URL")
	UpdateCmd.MarkFlagRequired("id")
}
