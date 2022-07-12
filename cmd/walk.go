// Copyright © 2022 Iori Mizutani <iori.mizutani@gmail.com>

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/karrick/godirwalk"
	"github.com/spf13/cobra"
)

// Ignore contains paths to ignore
var Ignore []string

func init() {
	walkCmd.Flags().StringSliceVarP(&Ignore, "ignore", "i", []string{"^\\..*", "Thumbs\\.db", "\\.DS_Store", "\\.\\_.*"}, "regexp patterns to ignore")
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
				for _, p := range Ignore {
					if match, _ := regexp.MatchString(p, filepath.Base(osPathname)); match {
						return godirwalk.SkipThis
					}
				}

				// don't print dirs
				if de.IsDir() {
					return nil
				}

				// ignore symlinks
				if de.IsSymlink() {
					return godirwalk.SkipThis
				}

				st, err := os.Stat(osPathname)
				switch err {
				case nil:
					// ignore files without extention
					if filepath.Ext(osPathname) == "" {
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
