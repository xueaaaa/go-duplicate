package cmd

import (
	"context"
	"fmt"
	"go-duplicate/internal/finder"
	"go-duplicate/internal/output"
	"go-duplicate/internal/units"
	"os"

	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <directory> [-s <sample-size>]",
	Short: "Searches for identical files in the specified directory.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		info, err := os.Stat(dir)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !info.IsDir() {
			fmt.Println(dir, "is not a directory")
			return
		}

		sampleSize, err := cmd.Flags().GetInt64("sample-size")
		if err != nil {
			fmt.Println(err)
			return
		}

		groups, err := finder.Find(context.Background(), dir, sampleSize)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = output.PlainOutput(os.Stdout, groups, dir)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().Int64P("sample-size", "s", 8*units.KiB,
		"The size of the initial and final parts used to calculate the partial hash (fingerprint)")
}
