package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stefanhoth/paperless-ngx-cli/api"
)

func init() {
	docCmd.Flags().Bool("full-perms", false, "Show full permission info")
	rootCmd.AddCommand(docCmd)
}

var docCmd = &cobra.Command{
	Use:   "doc <id>",
	Short: "Single document with metadata",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid id")
			os.Exit(1)
		}
		fullPerms, _ := cmd.Flags().GetBool("full-perms")

		c, _ := mustClient()
		params := &api.DocumentsRetrieveParams{FullPerms: &fullPerms}
		resp, err := c.DocumentsRetrieveWithResponse(ctx(), id, params)
		exitOnAPIError(err, resp)

		d := resp.JSON200
		date := "—"
		if d.CreatedDate != nil { //nolint:staticcheck // SA1019: no documented replacement for this deprecated upstream field
			date = d.CreatedDate.String()[:10] //nolint:staticcheck // SA1019: same as above
		}
		added := "—"
		if d.Added != nil {
			added = d.Added.Format("2006-01-02")
		}

		fmt.Printf("ID:             %d\n", derefInt(d.Id))
		fmt.Printf("Title:          %s\n", derefStr(d.Title))
		fmt.Printf("Created:        %s\n", date)
		fmt.Printf("Added:          %s\n", added)
		fmt.Printf("Correspondent:  %v\n", nvlInt(d.Correspondent))
		fmt.Printf("Document type:  %v\n", nvlInt(d.DocumentType))

		tags := make([]string, len(d.Tags))
		for i, t := range d.Tags {
			tags[i] = strconv.Itoa(t)
		}
		fmt.Printf("Tags:           %s\n", strings.Join(tags, ", "))
		fmt.Printf("Pages:          %v\n", nvlInt(d.PageCount))
		fmt.Printf("File:           %s\n", derefStr(d.OriginalFileName))

		if fullPerms && d.Permissions != nil {
			b, _ := json.MarshalIndent(d.Permissions, "  ", "  ")
			fmt.Printf("Permissions:\n  %s\n", b)
		}
	},
}

func nvlInt(i *int) string {
	if i == nil {
		return "—"
	}
	return strconv.Itoa(*i)
}
