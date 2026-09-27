package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/xueaaaa/go-duplicate/internal/finder"
	"github.com/xueaaaa/go-duplicate/internal/output"
	"github.com/xueaaaa/go-duplicate/internal/params"
	"github.com/xueaaaa/go-duplicate/internal/units"

	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <directory>",
	Short: "Searches for identical files in the specified directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := args[0]
		info, err := os.Stat(dir)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("%s is not a directory", dir)
		}

		formatString, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}

		formatString = strings.TrimSpace(strings.ToLower(formatString))
		format := params.Format(formatString)
		if !format.IsValid() {
			err = fmt.Errorf("unknown output format")
			return err
		}

		sampleSize, err := cmd.Flags().GetInt64("sample-size")
		if err != nil {
			return err
		}

		del, err := cmd.Flags().GetBool("delete")
		if err != nil {
			return err
		}

		hardlink, err := cmd.Flags().GetBool("hardlink")
		if err != nil {
			return err
		}

		skipConfirm, err := cmd.Flags().GetBool("yes")
		if err != nil {
			return err
		}

		silent, err := cmd.Flags().GetBool("silent")
		if err != nil {
			return err
		}

		p := params.Params{
			Format:     format,
			SampleSize: sampleSize,
			Delete:     del,
			Hardlink:   hardlink,
			Confirm:    !skipConfirm,
			Silent:     silent,
		}

		filesScanned, groups, err := finder.Find(context.Background(), dir, p)
		if err != nil {
			return err
		}

		stats := output.NewStats(dir, filesScanned, groups)

		switch format {
		case params.Plain:
			err = output.PlainOutput(os.Stdout, groups, stats)
			if err != nil {
				return err
			}
		case params.JSON:
			err = output.JSONOutput(os.Stdout, groups, stats)
			if err != nil {
				return err
			}
		}

		if len(groups) > 0 && p.Delete {
			if p.Confirm {
				confirm, err := output.PlainConfirm(os.Stderr, os.Stdin,
					fmt.Sprintf("Delete duplicates in %d groups?", len(groups)))
				if err != nil {
					return err
				}

				if !confirm {
					fmt.Fprintln(os.Stderr, "Aborted")
					return nil
				}
			}

			count, err := finder.Delete(groups)
			if err != nil {
				fmt.Errorf("%d files deleted before error occurred: %s\n", count, err)
				return err
			}

			fmt.Fprintf(os.Stderr, "Deleted %d files\n", count)
		}

		if len(groups) > 0 && p.Hardlink {
			if p.Confirm {
				confirm, err := output.PlainConfirm(os.Stderr, os.Stdin,
					fmt.Sprintf("Convert duplicates into hard links in %d groups?", len(groups)))
				if err != nil {
					return err
				}

				if !confirm {
					fmt.Fprintln(os.Stderr, "Aborted")
					return nil
				}
			}

			count, err := finder.Hardlink(groups)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%d files converted into hard links before error occurred: %s\n", count, err)
				return err
			}

			fmt.Fprintf(os.Stderr, "%d files converted into hard links\n", count)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().StringP("format", "f", string(params.Plain),
		"Output format of the program's results, plain and json output are available")
	findCmd.Flags().Int64P("sample-size", "s", 8*units.KiB,
		"size in bytes of the initial and final parts used to calculate the partial hash/fingerprint")
	findCmd.Flags().BoolP("delete", "d", false,
		"delete duplicate files, keeping only the first file in each group")
	findCmd.Flags().BoolP("hardlink", "l", false,
		"converts duplicate files into hardlinks to the first file in the group")
	findCmd.Flags().BoolP("yes", "y", false,
		"skip confirmation prompt (only applies with --delete or --hardlink)")
	findCmd.Flags().Bool("silent", false,
		"disables the display of progress")

	findCmd.MarkFlagsMutuallyExclusive("delete", "hardlink")
}
