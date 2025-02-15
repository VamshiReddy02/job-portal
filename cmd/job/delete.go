package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a job listing",
	Run: func(cmd *cobra.Command, args []string) {
		data := map[string]interface{}{
			"query": `mutation DeleteJob($id: ID!) {
				deleteJobListing(id: $id) {
					deletedJobId
				}
			}`,
			"variables": map[string]string{
				"id": jobID,
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

		fmt.Println("Job deleted successfully!")
	},
}

func init() {
	DeleteCmd.Flags().StringVarP(&jobID, "id", "i", "", "Job ID")
	DeleteCmd.MarkFlagRequired("id")
}
