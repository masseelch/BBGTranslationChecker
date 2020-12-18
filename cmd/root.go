package main

import (
	"github.com/davecgh/go-spew/spew"
	checker "github.com/masseelch/bbg-translation-checker"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "bbg-translation-checker -s /path/to/lang/dir -o /path/to/result/file",
		Short: "cli tool to check CPL's BBG translation for possible errors",
		Run: func(cmd *cobra.Command, args []string) {
			dirFlagValue, err := cmd.Flags().GetString("source")
			if err != nil {
				panic(err)
			}

			truthFlagValue, err := cmd.Flags().GetString("truth")
			if err != nil {
				panic(err)
			}

			t, ts, err := checker.Parse(truthFlagValue, dirFlagValue)
			if err != nil {
				panic(err)
			}

			// Check translations.
			rs, err := checker.Check(t, ts)
			if err != nil {
				panic(err)
			}

			spew.Dump(rs)
		},
	}
)

func init() {
	rootCmd.Flags().StringP("source", "s", "", "/path/to/lang/dir")
	rootCmd.Flags().StringP("output", "o", "", "/path/to/result/file")
	rootCmd.Flags().String("truth", "english.xml", "file of truth")

	_ = rootCmd.MarkFlagRequired("source")
}
