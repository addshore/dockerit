package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)
var myVersion string
var mySourceDate string

var Verbose bool
var Version bool

func Execute(appVersion string, appSourceDate string) {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&Version, "version", "", false, "version infomation")
	myVersion = appVersion
	mySourceDate = appSourceDate

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use: `dockerit [image]
  dockerit [image] [command]
  dockerit [flags] [image] [command] -- [command flags]`,
	Example:`  dockerit php -- -a
  dockerit --entry --user=root ubuntu bash
  dockerit --me --pwd --home composer:1 update -- --ignore-platform-reqs`,
	Short: "Run it in docker",
	Run: func(cmd *cobra.Command, args []string) {
		if(Version){
			fmt.Println("Version: " + myVersion)
			fmt.Println("Built at : " + mySourceDate)
			os.Exit(0)
		}

		if(len(args)== 0) {
			fmt.Println("Error: requires at least 1 arg(s), only received 0")
			fmt.Println("Use --help to see help text")
			os.Exit(0)
		}

		RunNow(RunNowOptions{
			Image: args[0],
			Pull: false,
			Cmd: args[1:],
		})
		},
	}
