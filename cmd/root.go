package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "note <note-name>",
	Short: "Quick and easy note taking",
	Long: `Really quick and easy note taking, uses $EDITOR to edit markdown files and fzf for searching.

Uses the $NOTES_DIRECTORY environment variable to determine where to store notes, defaults to ~/.notes.`,
	Run: func(cmd *cobra.Command, args []string) {
		noteName := ""
		if len(args) > 0 {
			noteName = strings.Join(args, " ")
		}
		err := note(noteName)
		if err != nil {
			// If error message is "not selection cancelled", we exit gracefully
			if err.Error() == "note selection cancelled" {
				os.Exit(0)
			}
			fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred: %v\n", err)
		os.Exit(1)
	}
}
