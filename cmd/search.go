package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stefanhoth/paperless-ngx-cli/api"
)

func init() {
	searchCmd.Flags().IntP("limit", "l", 20, "Max results")
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Full-text search",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		query := strings.Join(args, " ")
		c, _ := mustClient()

		ordering := "-created"
		params := &api.DocumentsListParams{
			Query:    &query,
			PageSize: &limit,
			Ordering: &ordering,
		}
		resp, err := c.DocumentsListWithResponse(ctx(), params)
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if len(resp.JSON200.Results) == 0 {
			fmt.Printf("No results for: %s\n", query)
			return
		}

		fmt.Printf("%-6s  %-12s  %s\n", "ID", "Date", "Title")
		fmt.Println(strings.Repeat("─", 70))
		for _, d := range resp.JSON200.Results {
			date := "—"
			if d.CreatedDate != nil {
				date = d.CreatedDate.String()[:10]
			}
			title := derefStr(d.Title)
			if len(title) > 55 {
				title = title[:55]
			}
			fmt.Printf("%-6d  %-12s  %s\n", derefInt(d.Id), date, title)
		}
		fmt.Printf("\n%d results\n", resp.JSON200.Count)
	},
}
