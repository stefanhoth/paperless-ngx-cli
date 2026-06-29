package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"paperless-cli/api"
)

func init() {
	rootCmd.AddCommand(tagsCmd)
	rootCmd.AddCommand(correspondentsCmd)
	rootCmd.AddCommand(typesCmd)
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Alle Tags auflisten",
	Run: func(cmd *cobra.Command, args []string) {
		c, _ := mustClient()
		n := 200
		name := "name"
		resp, err := c.TagsListWithResponse(ctx(), &api.TagsListParams{PageSize: &n, Ordering: &name})
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
			os.Exit(1)
		}
		for _, t := range resp.JSON200.Results {
			count := ""
			if t.DocumentCount != nil {
				count = fmt.Sprintf("  (%d Dok.)", *t.DocumentCount)
			}
			fmt.Printf("%4d  %s%s\n", derefInt(t.Id), t.Name, count)
		}
		fmt.Printf("\n%d Einträge\n", resp.JSON200.Count)
	},
}

var correspondentsCmd = &cobra.Command{
	Use:   "correspondents",
	Short: "Alle Korrespondenten auflisten",
	Run: func(cmd *cobra.Command, args []string) {
		c, _ := mustClient()
		n := 200
		name := "name"
		resp, err := c.CorrespondentsListWithResponse(ctx(), &api.CorrespondentsListParams{PageSize: &n, Ordering: &name})
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
			os.Exit(1)
		}
		for _, r := range resp.JSON200.Results {
			count := ""
			if r.DocumentCount != nil {
				count = fmt.Sprintf("  (%d Dok.)", *r.DocumentCount)
			}
			fmt.Printf("%4d  %s%s\n", derefInt(r.Id), r.Name, count)
		}
		fmt.Printf("\n%d Einträge\n", resp.JSON200.Count)
	},
}

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "Alle Dokumenttypen auflisten",
	Run: func(cmd *cobra.Command, args []string) {
		c, _ := mustClient()
		n := 200
		name := "name"
		resp, err := c.DocumentTypesListWithResponse(ctx(), &api.DocumentTypesListParams{PageSize: &n, Ordering: &name})
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
			os.Exit(1)
		}
		for _, t := range resp.JSON200.Results {
			count := ""
			if t.DocumentCount != nil {
				count = fmt.Sprintf("  (%d Dok.)", *t.DocumentCount)
			}
			fmt.Printf("%4d  %s%s\n", derefInt(t.Id), t.Name, count)
		}
		fmt.Printf("\n%d Einträge\n", resp.JSON200.Count)
	},
}
