package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"paperless-cli/api"
)

func init() {
	docsCmd.Flags().IntP("number", "n", 10, "Anzahl Dokumente")
	rootCmd.AddCommand(docsCmd)
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Letzte Dokumente auflisten",
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := cmd.Flags().GetInt("number")
		c, _ := mustClient()

		ordering := "-created"
		params := &api.DocumentsListParams{
			PageSize: &n,
			Ordering: &ordering,
		}
		resp, err := c.DocumentsListWithResponse(ctx(), params)
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%-6s  %-12s  %-25s  %s\n", "ID", "Datum", "Korrespondent", "Titel")
		fmt.Println(strings.Repeat("─", 90))
		for _, d := range resp.JSON200.Results {
			date := "—"
			if d.CreatedDate != nil {
				date = d.CreatedDate.String()[:10]
			}
			corr := "—"
			if d.Correspondent != nil {
				corr = fmt.Sprintf("%d", *d.Correspondent)
			}
			title := ""
			if d.Title != nil && len(*d.Title) > 50 {
				title = (*d.Title)[:50]
			} else if d.Title != nil {
				title = *d.Title
			}
			fmt.Printf("%-6d  %-12s  %-25s  %s\n", derefInt(d.Id), date, corr, title)
		}
		fmt.Printf("\n%d von %d Dokumenten\n", len(resp.JSON200.Results), resp.JSON200.Count)
	},
}
