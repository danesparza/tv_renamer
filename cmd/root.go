package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverInterface       string
	cfgFile               string
	ProblemWithConfigFile bool
)

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tv_renamer",
	Short: "A simple tool for renaming TV episodes",
	Long: `tv_renamer renames a set of video files.  You pass it a show name to use 
and it will attempt to match your filenames using TVDB episode names, find the season and episode,  
and rename your file(s) using the format sNNeNN.ext 

==========
Example
==========
Show: 'Looney Tunes'  
Filename: 04 - Scaredy Cat.avi

Using TVDB, this will rename to s1948e33.avi `,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/tv_renamer.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.AutomaticEnv() // read in environment variables that match

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	//	Set our defaults
	/*
		viper.SetDefault("server.port", "3000")
		viper.SetDefault("server.bind", "")
		viper.SetDefault("server.allowed-origins", "*")
	*/

	viper.SetConfigName("tv_renamer") // name of config file (without extension)
	viper.AddConfigPath("$HOME")      // adding home directory as first search path
	viper.AddConfigPath(".")          // also look in the working directory

	// If a config file is found, read it in
	// otherwise, make note that there was a problem
	if err := viper.ReadInConfig(); err != nil {
		ProblemWithConfigFile = true
	}
}
