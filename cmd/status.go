package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show system statistics",
	Run: func(cmd *cobra.Command, args []string) {
		c, _ := mustClient()
		resp, err := c.StatisticsRetrieveWithResponse(ctx())
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if resp.JSON200 == nil {
			fmt.Println("No data.")
			return
		}
		s := *resp.JSON200
		get := func(key string) interface{} { return s[key] }
		fmt.Printf("Documents:        %v\n", get("documents_total"))
		fmt.Printf("Tags:             %v\n", get("tag_count"))
		fmt.Printf("Correspondents:   %v\n", get("correspondent_count"))
		fmt.Printf("Document types:   %v\n", get("document_type_count"))
		fmt.Printf("Total characters: %v\n", get("character_count"))
		if ftRaw, ok := s["document_file_type_counts"]; ok {
			b, _ := json.Marshal(ftRaw)
			var ft []struct {
				MimeType      string `json:"mime_type"`
				MimeTypeCount int    `json:"mime_type_count"`
			}
			if json.Unmarshal(b, &ft) == nil && len(ft) > 0 {
				fmt.Println("\nFile types:")
				for _, f := range ft {
					fmt.Printf("  %-30s %5d\n", f.MimeType, f.MimeTypeCount)
				}
			}
		}
	},
}
