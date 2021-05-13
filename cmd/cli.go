package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"redisdual/commons"
	"redisdual/core"
	"redisdual/global"
)

type CommandFlags struct {
	URL      string
	CAPath   string
	CertPath string
	KeyPath  string
	Help     bool
}

var (
	commandFlags    = CommandFlags{}
	cfgFile         string
	detach          bool
	syncserver      string
	Confignotseterr error
	interact        bool
	version         bool
)

var LivePrefixState struct {
	LivePrefix string
	IsEnable   bool
}

var query = ""

func init() {
	cobra.EnablePrefixMatching = true
	cobra.OnInitialize(initConfig)

}

func getBasicCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "redissyncer-portal",
		Short: "redissyncer-portal",
		Long:  "",
	}

	rootCmd.PersistentFlags().BoolVarP(&commandFlags.Help, "help", "h", false, "help message")

	rootCmd.AddCommand(
		NewConfigCommand(),
		NewStartCommand(),
	)

	rootCmd.Flags().ParseErrorsWhitelist.UnknownFlags = true
	rootCmd.SilenceErrors = true
	return rootCmd
}

func getMainCmd(args []string) *cobra.Command {
	rootCmd := getBasicCmd()

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/config.yaml)")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Print version information and exit.")

	rootCmd.SetArgs(args)
	rootCmd.ParseFlags(args)
	rootCmd.SetOut(os.Stdout)

	return rootCmd
}

// MainStart start main command
func MainStart(args []string) {
	startCmd(getMainCmd, args)
}

func startCmd(getCmd func([]string) *cobra.Command, args []string) {
	rootCmd := getCmd(args)

	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		if !commons.FileExists(cfgFile) {
			panic("config file not exists")
		}
		global.RSPViper = core.Viper(cfgFile)
		return
	}

	global.RSPViper = core.Viper()

	//if syncserver == "" {
	//	fmt.Println(syncserver)
	//	syncserver = os.Getenv("SYNCSERVER")
	//}
	//
	//if cfgFile != "" && commons.FileExists(cfgFile) {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	viper.AddConfigPath(".")
	//	viper.SetConfigName(".config")
	//}
	//
	//viper.ReadInConfig()
	//
	//viper.AutomaticEnv() // read in environment variables that match
	//
	//if syncserver != "" {
	//	viper.Set("SYNCSERVER", syncserver)
	//}

}
