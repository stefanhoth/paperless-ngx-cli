package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stefanhoth/paperless-ngx-cli/api"
)

func init() {
	rootCmd.AddCommand(tagsCmd)
	rootCmd.AddCommand(correspondentsCmd)
	rootCmd.AddCommand(typesCmd)
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List all tags",
	Run: func(_ *cobra.Command, _ []string) {
		c, _ := mustClient()
		n := 200
		name := "name"
		resp, err := c.TagsListWithResponse(ctx(), &api.TagsListParams{PageSize: &n, Ordering: &name})
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		for _, t := range resp.JSON200.Results {
			count := ""
			if t.DocumentCount != nil {
				count = fmt.Sprintf("  (%d docs)", *t.DocumentCount)
			}
			fmt.Printf("%4d  %s%s\n", derefInt(t.Id), t.Name, count)
		}
		fmt.Printf("\n%d entries\n", resp.JSON200.Count)
	},
}

var correspondentsCmd = &cobra.Command{
	Use:   "correspondents",
	Short: "List all correspondents",
	Run: func(_ *cobra.Command, _ []string) {
		c, _ := mustClient()
		n := 200
		name := "name"
		resp, err := c.CorrespondentsListWithResponse(ctx(), &api.CorrespondentsListParams{PageSize: &n, Ordering: &name})
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		for _, r := range resp.JSON200.Results {
			count := ""
			if r.DocumentCount != nil {
				count = fmt.Sprintf("  (%d docs)", *r.DocumentCount)
			}
			fmt.Printf("%4d  %s%s\n", derefInt(r.Id), r.Name, count)
		}
		fmt.Printf("\n%d entries\n", resp.JSON200.Count)
	},
}

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "List all document types",
	Run: func(_ *cobra.Command, _ []string) {
		c, _ := mustClient()
		n := 200
		name := "name"
		resp, err := c.DocumentTypesListWithResponse(ctx(), &api.DocumentTypesListParams{PageSize: &n, Ordering: &name})
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		for _, t := range resp.JSON200.Results {
			count := ""
			if t.DocumentCount != nil {
				count = fmt.Sprintf("  (%d docs)", *t.DocumentCount)
			}
			fmt.Printf("%4d  %s%s\n", derefInt(t.Id), t.Name, count)
		}
		fmt.Printf("\n%d entries\n", resp.JSON200.Count)
	},
}
