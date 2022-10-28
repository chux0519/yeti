package cmd

import (
	"os"

	"github.com/chux0519/yeti/pkg/api"
	"github.com/chux0519/yeti/pkg/config"
	logging "github.com/ipfs/go-log"
	"github.com/spf13/cobra"
)

var serverFlag = config.ServerFlag{}

var rootLog = logging.Logger("root")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yeti",
	Short: "yeti backend server",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.LoadServerConfig(serverFlag.Config)
		if config.Debug {
			logging.SetLogLevel("*", "debug")
		} else {
			logging.SetLogLevel("*", "info")
		}
		rootLog.Debugf("config loaded: %v", config)
		api.StartYetiServer(config)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	flags := rootCmd.Flags()

	flags.StringVar(&serverFlag.Config, "config", "config.toml", "Configuration of yeti server")
}
