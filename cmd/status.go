package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show system statistics",
	Run: func(_ *cobra.Command, _ []string) {
		c, _ := mustClient()
		resp, err := c.StatisticsRetrieveWithResponse(ctx())
		exitOnAPIError(err, resp)
		if resp.JSON200 == nil {
			fmt.Println("No data.")
			return
		}
		s := *resp.JSON200
		get := func(key string) string { return formatStat(s[key]) }
		fmt.Printf("Documents:        %s\n", get("documents_total"))
		fmt.Printf("Tags:             %s\n", get("tag_count"))
		fmt.Printf("Correspondents:   %s\n", get("correspondent_count"))
		fmt.Printf("Document types:   %s\n", get("document_type_count"))
		fmt.Printf("Total characters: %s\n", get("character_count"))
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

// formatStat prints JSON numbers (decoded as float64) without scientific notation.
func formatStat(v interface{}) string {
	if f, ok := v.(float64); ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	return fmt.Sprint(v)
}
