package cmd

import (
	"context"
	"fmt"
	"go-duplicate/internal/finder"
	"go-duplicate/internal/output"
	"go-duplicate/internal/params"
	"go-duplicate/internal/units"
	"os"

	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <directory>",
	Short: "Searches for identical files in the specified directory",
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

		del, err := cmd.Flags().GetBool("delete")
		if err != nil {
			fmt.Println(err)
			return
		}

		hardlink, err := cmd.Flags().GetBool("hardlink")
		if err != nil {
			fmt.Println(err)
			return
		}

		skipDryRun, err := cmd.Flags().GetBool("yes")
		if err != nil {
			fmt.Println(err)
			return
		}

		params := params.Params{
			SampleSize: sampleSize,
			Delete:     del,
			Hardlink:   hardlink,
			DryRun:     !skipDryRun,
		}

		groups, err := finder.Find(context.Background(), dir, params)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = output.PlainOutput(os.Stdout, groups, dir)
		if err != nil {
			fmt.Println(err)
			return
		}

		if params.Delete {
			if params.DryRun {
				confirm, err := output.PlainConfirm(os.Stdout, os.Stdin,
					fmt.Sprintf("Delete duplicates in %d groups?", len(groups)))
				if err != nil {
					fmt.Println(err)
					return
				}

				if !confirm {
					fmt.Println("Aborted")
					return
				}
			}

			count, err := finder.Delete(groups)
			if err != nil {
				fmt.Printf("%d files deleted before error occurred: %s\n", count, err)
				return
			}

			fmt.Printf("Deleted %d files\n", count)
		}

		if params.Hardlink {
			if params.DryRun {
				confirm, err := output.PlainConfirm(os.Stdout, os.Stdin,
					fmt.Sprintf("Convert duplicates into hard links in %d groups?", len(groups)))
				if err != nil {
					fmt.Println(err)
					return
				}

				if !confirm {
					fmt.Println("Aborted")
					return
				}
			}

			count, err := finder.Hardlink(groups)
			if err != nil {
				fmt.Printf("%d files converted into hard links before error occurred: %s\n", count, err)
				return
			}

			fmt.Printf("%d files converted into hard links", count)
		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().Int64P("sample-size", "s", 8*units.KiB,
		"size in bytes of the initial and final parts used to calculate the partial hash (fingerprint)")
	findCmd.Flags().BoolP("delete", "d", false,
		"delete duplicate files, keeping only the first file in each group")
	findCmd.Flags().BoolP("hardlink", "l", false,
		"converts duplicate files into hardlinks to the first file in the group")
	findCmd.Flags().BoolP("yes", "y", false,
		"skip confirmation prompt (only applies with --delete or --hardlink)")
}
