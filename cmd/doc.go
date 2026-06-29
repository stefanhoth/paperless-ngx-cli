package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"paperless-cli/api"
)

func init() {
	docCmd.Flags().Bool("full-perms", false, "Vollständige Berechtigungen anzeigen")
	rootCmd.AddCommand(docCmd)
}

var docCmd = &cobra.Command{
	Use:   "doc <id>",
	Short: "Einzelnes Dokument mit Metadaten",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ungültige ID")
			os.Exit(1)
		}
		fullPerms, _ := cmd.Flags().GetBool("full-perms")

		c, _ := mustClient()
		params := &api.DocumentsRetrieveParams{FullPerms: &fullPerms}
		resp, err := c.DocumentsRetrieveWithResponse(ctx(), id, params)
		if err != nil || resp.StatusCode() != 200 {
			fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
			os.Exit(1)
		}

		d := resp.JSON200
		date := "—"
		if d.CreatedDate != nil {
			date = d.CreatedDate.String()[:10]
		}
		added := "—"
		if d.Added != nil {
			added = d.Added.Format("02.01.2006")
		}

		fmt.Printf("ID:             %d\n", derefInt(d.Id))
		fmt.Printf("Titel:          %s\n", derefStr(d.Title))
		fmt.Printf("Erstellt:       %s\n", date)
		fmt.Printf("Hinzugefügt:    %s\n", added)
		fmt.Printf("Korrespondent:  %v\n", nvlInt(d.Correspondent))
		fmt.Printf("Dokumenttyp:    %v\n", nvlInt(d.DocumentType))

		tags := make([]string, len(d.Tags))
		for i, t := range d.Tags {
			tags[i] = strconv.Itoa(t)
		}
		fmt.Printf("Tags:           %s\n", strings.Join(tags, ", "))
		fmt.Printf("Seiten:         %v\n", nvlInt(d.PageCount))
		fmt.Printf("Datei:          %s\n", derefStr(d.OriginalFileName))

		if fullPerms && d.Permissions != nil {
			b, _ := json.MarshalIndent(d.Permissions, "  ", "  ")
			fmt.Printf("Berechtigungen:\n  %s\n", b)
		}
	},
}

func nvlInt(i *int) string {
	if i == nil {
		return "—"
	}
	return strconv.Itoa(*i)
}
