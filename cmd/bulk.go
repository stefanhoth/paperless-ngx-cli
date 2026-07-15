package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stefanhoth/paperless-ngx-cli/api"
)

func init() {
	rootCmd.AddCommand(bulkCmd)
}

var bulkCmd = &cobra.Command{
	Use:   "bulk <operation> <ids> [param]",
	Short: "Bulk operations on documents",
	Long: `Operations:
  reprocess         <ids>                   Re-process (OCR etc.)
  delete            <ids>                   Delete
  merge             <ids>                   Merge into one document
  rotate            <ids> <90|180|270>      Rotate
  add-tag           <ids> <tag_id>          Add tag
  remove-tag        <ids> <tag_id>          Remove tag
  set-correspondent <ids> <id>              Set correspondent
  set-type          <ids> <id>              Set document type

ids: comma-separated, e.g. 1,2,3`,
	Args: cobra.MinimumNArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		op := args[0]
		ids := parseIDs(args[1])
		if len(ids) == 0 {
			fmt.Fprintln(os.Stderr, "no valid IDs provided")
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
				fmt.Fprintln(os.Stderr, "usage: bulk rotate <ids> <degrees>")
				os.Exit(1)
			}
			deg, _ := strconv.Atoi(args[2])
			method = api.MethodEnumRotate
			params["degrees"] = deg
		case "add-tag":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "usage: bulk add-tag <ids> <tag_id>")
				os.Exit(1)
			}
			tagID, _ := strconv.Atoi(args[2])
			method = api.MethodEnumAddTag
			params["tag"] = tagID
		case "remove-tag":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "usage: bulk remove-tag <ids> <tag_id>")
				os.Exit(1)
			}
			tagID, _ := strconv.Atoi(args[2])
			method = api.MethodEnumRemoveTag
			params["tag"] = tagID
		case "set-correspondent":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "usage: bulk set-correspondent <ids> <id>")
				os.Exit(1)
			}
			corrID, _ := strconv.Atoi(args[2])
			method = api.MethodEnumSetCorrespondent
			params["correspondent"] = corrID
		case "set-type":
			if len(args) < 3 {
				fmt.Fprintln(os.Stderr, "usage: bulk set-type <ids> <id>")
				os.Exit(1)
			}
			typeID, _ := strconv.Atoi(args[2])
			method = api.MethodEnumSetDocumentType
			params["document_type"] = typeID
		default:
			fmt.Fprintf(os.Stderr, "unknown operation: %s\n", op)
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
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if resp.StatusCode() >= 400 {
			fmt.Fprintf(os.Stderr, "API error %d: %s\n", resp.StatusCode(), string(resp.Body))
			os.Exit(1)
		}
		fmt.Printf("OK — %d document(s), operation: %s\n", len(ids), op)
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
