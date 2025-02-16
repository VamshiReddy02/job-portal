package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var jobID string
var jobTitle string

// GetCmd represents the `job get` command
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve job listings",
	Run: func(cmd *cobra.Command, args []string) {
		var data map[string]interface{}

		switch {
		case jobID != "":
			data = map[string]interface{}{
				"query": `query GetJob($id: ID!) {
					job(id: $id) {
						_id
						title
						company
						description
						url
					}
				}`,
				"variables": map[string]string{
					"id": jobID,
				},
			}
		case jobTitle != "":
			data = map[string]interface{}{
				"query": `query GetJobsByTitle($title: String!) {
					jobsByTitle(title: $title) {
						_id
						title
						company
						description
						url
					}
				}`,
				"variables": map[string]string{
					"title": jobTitle,
				},
			}
		default:
			data = map[string]interface{}{
				"query": `query {
					jobs {
						_id
						title
						company
						description
						url
					}
				}`,
			}
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error parsing response:", err)
			return
		}

		dataKey := "jobs"
		if jobID != "" {
			dataKey = "job"
		} else if jobTitle != "" {
			dataKey = "jobsByTitle"
		}

		jobs, ok := result["data"].(map[string]interface{})[dataKey]
		if !ok {
			fmt.Println("No job data found.")
			return
		}

		jobList, ok := jobs.([]interface{})
		if !ok {
			jobList = []interface{}{jobs}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Title", "Company", "Description", "URL"})
		table.SetBorder(false)

		for _, j := range jobList {
			jobMap, _ := j.(map[string]interface{})
			table.Append([]string{
				fmt.Sprintf("%v", jobMap["_id"]),
				fmt.Sprintf("%v", jobMap["title"]),
				fmt.Sprintf("%v", jobMap["company"]),
				fmt.Sprintf("%v", jobMap["description"]),
				fmt.Sprintf("%v", jobMap["url"]),
			})
		}

		table.Render()
	},
}

func init() {
	GetCmd.Flags().StringVarP(&jobID, "id", "i", "", "Filter job by ID")
	GetCmd.Flags().StringVarP(&jobTitle, "title", "t", "", "Filter job by title")
}
