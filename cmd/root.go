package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)
var myVersion string
var mySourceDate string

var Verbose bool
var Version bool
var SelfUpdate bool

func Execute(appVersion string, appSourceDate string) {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&Version, "version", "", false, "version infomation")

	rootCmd.Flags().BoolVarP(&SelfUpdate, "selfupdate", "", false, "Update this command to the latest release from Github")

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

		if(SelfUpdate){
			if(Verbose){
				selfupdate.EnableLog()
			}
			v := semver.MustParse(strings.Trim(myVersion,"v"))
			latest, err := selfupdate.UpdateSelf(v, "addshore/dockerit")
			if err != nil {
				log.Println("Binary update failed:", err)
				return
			}
			if latest.Version.Equals(v) {
				// latest version is the same as current version. It means current binary is up to date.
				log.Println("Current binary is the latest version", myVersion)
			} else {
				log.Println("Successfully updated to version", latest.Version)
				log.Println("Release note:\n", latest.ReleaseNotes)
			}
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
