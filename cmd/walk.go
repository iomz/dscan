// Copyright © 2022 Iori Mizutani <iori.mizutani@gmail.com>

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/spf13/cobra"
)

func init() {
	//walkCmd.Flags().StringVarP(&Output, "out", "o", "", "output file")
	rootCmd.AddCommand(walkCmd)
}

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
		err = godirwalk.Walk(dir.Name(), &godirwalk.Options{
			Callback: func(osPathname string, de *godirwalk.Dirent) error {
				if strings.HasPrefix(osPathname, ".git") {
					return godirwalk.SkipThis
				}
				if de.IsDir() {
					return nil
				}
				st, err := os.Stat(osPathname)
				switch err {
				case nil:
					_, err = fmt.Printf("% 12d %v %s\n", st.Size(), st.ModTime(), osPathname)
				default:
					// ignore the error and just show the mode type
					_, err = fmt.Printf("%s % 12d %s\n", de.ModeType(), st.Size(), osPathname)
				}
				return nil
			},
			Unsorted: true,
		})
		return err
	},
}
