package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)
var appVersion string

var Verbose bool
var Version bool

func Execute(mainVersion string) {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&Version, "version", "", false, "version infomation")
	appVersion = mainVersion

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: `dockerit [image]
  dockerit [image] [command]
  dockerit [flags] [image] [command] -- [command flags]`,
	Example:`  dockerit --pwd=0 php -- -a
  dockerit --entry --user=root --pwd=0 ubuntu bash
  dockerit --pwd composer:1 update -- --ignore-platform-reqs`,
	Short: "Run it in docker",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if(Version){
			fmt.Println("Version: " + appVersion)
			os.Exit(0)
		}

		RunNow(RunNowOptions{
			Image: args[0],
			Pull: false,
			Cmd: args[1:],
		})
		},
	}
