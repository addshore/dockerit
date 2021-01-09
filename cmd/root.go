package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)

var rootCmd = &cobra.Command{
	Use:   "dt",
}

var Verbose bool

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
