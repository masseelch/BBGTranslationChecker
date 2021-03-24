package main

import (
	checker "github.com/masseelch/bbg-translation-checker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "bbg-translation-checker -s /path/to/lang/dir",
		Short: "cli tool to check CPL's BBG translation for possible errors",
		Run: func(cmd *cobra.Command, args []string) {
			t, ts, err := checker.Parse(viper.GetString("truth"), viper.GetString("source"), viper.GetString("only"))
			if err != nil {
				panic(err)
			}

			// Check translations.
			rs, err := checker.Check(t, ts)
			if err != nil {
				panic(err)
			}

			// Dump summary to a reports.txt file.
			f, err := os.Create(filepath.Join(viper.GetString("output"), "reports.txt"))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if err := rs.DumpSummary(f); err != nil {
				panic(err)
			}

			// If the 'overwrite' flag is set to true we change the original file and add our report
			// to the entries in form of a comment.
			if err := rs.DumpWithComments(t, viper.GetString("output")); err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().String("source", ".", "/path/to/lang/dir")
	rootCmd.Flags().String("output", ".", "/path/to/output/dir")
	rootCmd.Flags().String("truth", "english.xml", "file of truth")
	rootCmd.Flags().String("only", "", "Filename to check. If empty all files are checked")

	// Bind database flags to viper.
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile("config.yaml")

	_ = viper.ReadInConfig()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
