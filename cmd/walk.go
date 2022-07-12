/*
Copyright © 2022 Iori Mizutani <iori.mizutani@gmail.com>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/spf13/cobra"
)

func init() {
	walkCmd.Flags().StringVarP(&Output, "out", "o", "", "output file")
	rootCmd.AddCommand(walkCmd)
}

var Output string

var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Walk the dir",
	Long:  `Walk the directory and output the result`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Stat(args[0])
		if err != nil {
			return fmt.Errorf("failed to open directory, error: %w", err)
		}
		if !dir.IsDir() {
			return fmt.Errorf("%q is not a directory", dir.Name())
		}
		err = godirwalk.Walk(dir, &godirwalk.Options{
			Callback: func(osPathname string, de *godirwalk.Dirent) error {
				if strings.Contains(osPathname, ".git") {
					return godirwalk.SkipThis
				}
				fmt.Printf("%s %s\n", de.ModeType(), osPathname)
				return nil
			},
			Unsorted: true,
		})
		return err
	},
}
