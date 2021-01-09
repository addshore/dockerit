package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)

var rootCmd = &cobra.Command{
	Use: `dockerit [image] [command]
  dockerit [flags] [image] [command] -- [command flags]`,
	Example:`  dockerit --pwd=0 php -- -a
  dockerit --entry --user=root --pwd=0 ubuntu bash
  dockerit --pwd composer:1 update -- --ignore-platform-reqs`,
	Short: "Run it in docker",
	Run: func(cmd *cobra.Command, args []string) {
		RunNow(RunNowOptions{
			Image: args[0],
			Pull: false,
			Cmd: args[1:],
		})
		},
	}

var Verbose bool

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
