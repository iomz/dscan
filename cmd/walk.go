// Copyright © 2022 Iori Mizutani <iori.mizutani@gmail.com>

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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
		path := args[0]
		dir, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to open directory, error: %w", err)
		}
		if !dir.IsDir() {
			return fmt.Errorf("%q is not a directory", dir.Name())
		}
		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		err = godirwalk.Walk(path, &godirwalk.Options{
			Callback: func(osPathname string, de *godirwalk.Dirent) error {
				// don't print dirs
				if de.IsDir() {
					if strings.HasPrefix(filepath.Base(osPathname), ".") {
						return godirwalk.SkipThis
					}
					return nil
				}

				st, err := os.Stat(osPathname)
				switch err {
				case nil:
					if strings.HasPrefix(st.Name(), ".") {
						// don't print anything to a hidden file
						return nil
					}
					_, err = fmt.Printf("%v\t% 12d\t%s\n", st.ModTime().Format("2006-01-02 15:04:05"), st.Size(), osPathname)
				default:
					// ignore the error and just show the mode type
					_, err = fmt.Printf("%s\n", osPathname)
				}
				return nil
			},
			Unsorted: true,
		})
		return err
	},
}
