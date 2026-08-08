package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godup",
	Short: "Find and work with identical files based on hash",
	Long: `go-duplicate is a utility designed to find and work with duplicate files based on hashing. 
It uses the SHA-256 hashing algorithm.
By using partial hashing methods and initial grouping by size, the utility can quickly process multiple 
directories containing huge files in a short time.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
