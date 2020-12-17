package main

import (
	"github.com/davecgh/go-spew/spew"
	checker "github.com/masseelch/bbg-translation-checker"
	"github.com/spf13/cobra"
	"path/filepath"
)

var (
	rootCmd = &cobra.Command{
		Use:   "bbg-translation-checker -s /path/to/lang/dir -o /path/to/result/file",
		Short: "cli tool to check CPL's BBG translation for possible errors",
		Run: func(cmd *cobra.Command, args []string) {
			dir, err := cmd.Flags().GetString("source")
			if err != nil {
				panic(err)
			}

			t, err := cmd.Flags().GetString("truth")
			if err != nil {
				panic(err)
			}

			// Truth
			truth, err := checker.ParseFile(filepath.Join(dir, t))
			if err != nil {
				panic(err)
			}

			// Other files
			

			// Check translations.
			r := checker.Check(truth, nil)

			spew.Dump(r)
		},
	}
)

func init() {
	rootCmd.Flags().StringP("source", "s", "", "/path/to/lang/dir")
	rootCmd.Flags().StringP("output", "o", "", "/path/to/result/file")
	rootCmd.Flags().String("truth", "english.xml", "file of truth")
}
