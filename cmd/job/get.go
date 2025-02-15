package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

		body, err := io.ReadAll(resp.Body) // Using io.ReadAll instead of ioutil.ReadAll
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		fmt.Println("Response:", string(body))
	},
}

func init() {
	GetCmd.Flags().StringVarP(&jobID, "id", "i", "", "Filter job by ID")
	GetCmd.Flags().StringVarP(&jobTitle, "title", "t", "", "Filter job by title")
}
