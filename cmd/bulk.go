package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"paperless-cli/api"
)

func init() {
	rootCmd.AddCommand(bulkCmd)
}

var bulkCmd = &cobra.Command{
	Use:   "bulk <operation> <ids> [param]",
	Short: "Bulk-Operationen auf Dokumenten",
	Long: `Operationen:
  reprocess       <ids>                   Neu verarbeiten (OCR etc.)
  delete          <ids>                   Löschen
  merge           <ids>                   Zusammenführen
  rotate          <ids> <90|180|270>      Drehen
  add-tag         <ids> <tag_id>          Tag hinzufügen
  remove-tag      <ids> <tag_id>          Tag entfernen
  set-correspondent <ids> <id>            Korrespondent setzen
  set-type        <ids> <id>              Dokumenttyp setzen

ids: kommagetrennt, z.B. 1,2,3`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		op := args[0]
		ids := parseIDs(args[1])
		if len(ids) == 0 {
			fmt.Fprintln(os.Stderr, "Keine gültigen IDs")
			os.Exit(1)
		}

		var method api.MethodEnum
		params := map[string]interface{}{}

		switch op {
		case "reprocess":
			method = api.MethodEnumReprocess
		case "delete":
			method = api.MethodEnumDelete
		case "merge":
			method = api.MethodEnumMerge
		case "rotate":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: bulk rotate <ids> <degrees>")
				os.Exit(1)
			}
			deg, _ := strconv.Atoi(args[2])
			method = api.MethodEnumRotate
			params["degrees"] = deg
		case "add-tag":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: bulk add-tag <ids> <tag_id>")
				os.Exit(1)
			}
			tagID, _ := strconv.Atoi(args[2])
			method = api.MethodEnumAddTag
			params["tag"] = tagID
		case "remove-tag":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: bulk remove-tag <ids> <tag_id>")
				os.Exit(1)
			}
			tagID, _ := strconv.Atoi(args[2])
			method = api.MethodEnumRemoveTag
			params["tag"] = tagID
		case "set-correspondent":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: bulk set-correspondent <ids> <id>")
				os.Exit(1)
			}
			corrID, _ := strconv.Atoi(args[2])
			method = api.MethodEnumSetCorrespondent
			params["correspondent"] = corrID
		case "set-type":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "Usage: bulk set-type <ids> <id>")
				os.Exit(1)
			}
			typeID, _ := strconv.Atoi(args[2])
			method = api.MethodEnumSetDocumentType
			params["document_type"] = typeID
		default:
			fmt.Fprintf(os.Stderr, "Unbekannte Operation: %s\n", op)
			os.Exit(1)
		}

		c, _ := mustClient()
		body := api.BulkEditJSONRequestBody{
			Documents:  &ids,
			Method:     &method,
			Parameters: &params,
		}
		resp, err := c.BulkEditWithResponse(ctx(), body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fehler: %v\n", err)
			os.Exit(1)
		}
		if resp.StatusCode() >= 400 {
			fmt.Fprintf(os.Stderr, "API-Fehler %d: %s\n", resp.StatusCode(), string(resp.Body))
			os.Exit(1)
		}
		fmt.Printf("OK — %d Dokument(e), Operation: %s\n", len(ids), op)
	},
}

func parseIDs(s string) []int {
	parts := strings.Split(s, ",")
	var ids []int
	for _, p := range parts {
		if id, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
