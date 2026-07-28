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
	chatCmd.Flags().IntP("doc", "d", 0, "Restrict chat to this document ID")
	rootCmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(chatCmd)
}

var suggestCmd = &cobra.Command{
	Use:   "suggest <id>",
	Short: "AI-generated suggestions for a document (server-side AI must be enabled)",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: invalid id")
			os.Exit(1)
		}

		c, _ := mustClient()
		resp, err := c.DocumentsAiSuggestionsRetrieveWithResponse(ctx(), id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if apiErr := apiError(resp.StatusCode(), resp.Body); apiErr != nil {
			fmt.Fprintln(os.Stderr, "error:", apiErr)
			os.Exit(1)
		}

		s := resp.JSON200
		if s == nil {
			fmt.Println("No suggestions.")
			return
		}

		fmt.Printf("Title:              %s\n", derefStr(s.Title))
		printSuggestionField("Correspondents", s.Correspondents, s.SuggestedCorrespondents)
		printSuggestionField("Tags", s.Tags, s.SuggestedTags)
		printSuggestionField("Document types", s.DocumentTypes, s.SuggestedDocumentTypes)
		printSuggestionField("Storage paths", s.StoragePaths, s.SuggestedStoragePaths)
		if len(s.Dates) > 0 {
			fmt.Printf("%-20s%s\n", "Dates:", strings.Join(s.Dates, ", "))
		}
	},
}

// printSuggestionField prints one suggestion row: existing entity IDs the AI
// matched, plus free-text names for entities that don't exist yet.
func printSuggestionField(label string, ids []int, suggestedNew []string) {
	if len(ids) == 0 && len(suggestedNew) == 0 {
		return
	}
	parts := make([]string, 0, len(ids)+len(suggestedNew))
	for _, id := range ids {
		parts = append(parts, strconv.Itoa(id))
	}
	for _, s := range suggestedNew {
		parts = append(parts, fmt.Sprintf("%q (new)", s))
	}
	fmt.Printf("%-20s%s\n", label+":", strings.Join(parts, ", "))
}

var chatCmd = &cobra.Command{
	Use:   "chat <question>",
	Short: "Ask the AI about your documents (server-side AI must be enabled)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		docID, _ := cmd.Flags().GetInt("doc")
		question := strings.Join(args, " ")

		body := api.ChatStreamingRequest{Q: question}
		if docID != 0 {
			body.DocumentId = &docID
		}

		c, _ := mustClient()
		resp, err := c.DocumentsChatCreateWithResponse(ctx(), body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if apiErr := apiError(resp.StatusCode(), resp.Body); apiErr != nil {
			fmt.Fprintln(os.Stderr, "error:", apiErr)
			os.Exit(1)
		}

		fmt.Println(string(resp.Body))
	},
}
