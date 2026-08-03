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
  delete            <ids>                   Move to trash
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

		c, _ := mustClient()
		if err := runBulk(c, op, ids, args); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("OK — %d document(s), operation: %s\n", len(ids), op)
	},
}

// runBulk sends the operation to the API. Since API v10 the document actions
// (reprocess, delete, merge, rotate) have dedicated endpoints — on bulk_edit
// they are deprecated. The metadata operations still go through bulk_edit.
func runBulk(c *api.ClientWithResponses, op string, ids []int, args []string) error {
	switch op {
	case "reprocess":
		resp, err := c.DocumentsReprocessWithResponse(ctx(), api.DocumentsReprocessJSONRequestBody{Documents: &ids})
		if err != nil {
			return err
		}
		return apiError(resp.StatusCode(), resp.Body)

	case "delete":
		resp, err := c.DocumentsDeleteWithResponse(ctx(), api.DocumentsDeleteJSONRequestBody{Documents: &ids})
		if err != nil {
			return err
		}
		return apiError(resp.StatusCode(), resp.Body)

	case "merge":
		resp, err := c.DocumentsMergeWithResponse(ctx(), api.DocumentsMergeJSONRequestBody{Documents: &ids})
		if err != nil {
			return err
		}
		return apiError(resp.StatusCode(), resp.Body)

	case "rotate":
		degrees := bulkArg(args, "bulk rotate <ids> <degrees>")
		resp, err := c.DocumentsRotateWithResponse(ctx(), api.DocumentsRotateJSONRequestBody{Documents: &ids, Degrees: degrees})
		if err != nil {
			return err
		}
		return apiError(resp.StatusCode(), resp.Body)

	default:
		return bulkEdit(c, op, ids, args)
	}
}

// bulkEdit sends the metadata operations through the bulk_edit endpoint.
func bulkEdit(c *api.ClientWithResponses, op string, ids []int, args []string) error {
	var method api.MethodEnum
	params := map[string]interface{}{}

	switch op {
	case "add-tag":
		method = api.MethodEnumAddTag
		params["tag"] = bulkArg(args, "bulk add-tag <ids> <tag_id>")
	case "remove-tag":
		method = api.MethodEnumRemoveTag
		params["tag"] = bulkArg(args, "bulk remove-tag <ids> <tag_id>")
	case "set-correspondent":
		method = api.MethodEnumSetCorrespondent
		params["correspondent"] = bulkArg(args, "bulk set-correspondent <ids> <id>")
	case "set-type":
		method = api.MethodEnumSetDocumentType
		params["document_type"] = bulkArg(args, "bulk set-type <ids> <id>")
	default:
		fmt.Fprintf(os.Stderr, "unknown operation: %s\n", op)
		os.Exit(1)
	}

	resp, err := c.BulkEditWithResponse(ctx(), api.BulkEditJSONRequestBody{
		Documents:  &ids,
		Method:     &method,
		Parameters: &params,
	})
	if err != nil {
		return err
	}
	return apiError(resp.StatusCode(), resp.Body)
}

// bulkArg returns the operation's numeric parameter (args[2]), exiting with
// the given usage hint when it is missing.
func bulkArg(args []string, usage string) int {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: "+usage)
		os.Exit(1)
	}
	n, _ := strconv.Atoi(args[2])
	return n
}

// apiError turns a non-2xx response into an error, and nil otherwise.
func apiError(status int, body []byte) error {
	if status >= 400 {
		return fmt.Errorf("API error %d: %s", status, string(body))
	}
	return nil
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
