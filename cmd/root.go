package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/open-pomodoro/go-openpomodoro"
	"github.com/open-pomodoro/openpomodoro-cli/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pomodoro",
	Short: "Pomodoro CLI",
	Long:  "A simple Pomodoro command-line client for the Open Pomodoro format",
}

var (
	client   *openpomodoro.Client
	settings *openpomodoro.Settings

	directoryFlag string
	formatFlag    string
	waitFlag      bool
)

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(
		&directoryFlag, "directory", "", ``,
		"directory to read/write Open Pomodoro data (default is ~/.config/pomodoro/)")

	RootCmd.PersistentFlags().StringVarP(
		&formatFlag, "format", "f", format.DefaultFormat,
		"format to display Pomodoros in")

	RootCmd.PersistentFlags().BoolVarP(
		&waitFlag, "wait", "w", false,
		"wait for the Pomodoro to end before exiting")

	viper.AutomaticEnv()

	directoryFlag = resolveDirectory(directoryFlag)

	RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		var err error

		client, err = openpomodoro.NewClient(directoryFlag)
		if err != nil {
			log.Fatalf("Could not create client: %v", err)
		}

		settings, err = client.Settings()
		if err != nil {
			log.Fatalf("Could not retrieve settings: %v", err)
		}
	}
}

func initConfig() {
	viper.AutomaticEnv()
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Customize built-in flag descriptions once at startup
	if helpFlag := RootCmd.Flags().Lookup("help"); helpFlag != nil {
		helpFlag.Usage = "Show help"
	}
	if versionFlag := RootCmd.Flags().Lookup("version"); versionFlag != nil {
		versionFlag.Usage = "Show version info"
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// resolveDirectory determines the directory to use for storing Pomodoro data.
//
// The resolution logic is as follows:
//  1. If the `flag` argument is non-empty, it is returned as the directory.
//  2. Otherwise, if a legacy directory `~/.pomodoro` exists in the user's home directory, it is used.
//  3. If neither of the above applies, the function falls back to the XDG config directory
//     (typically `$XDG_CONFIG_HOME/pomodoro` or `~/.config/pomodoro` if XDG_CONFIG_HOME is not set).
func resolveDirectory(directoryFlag string) string {
	if directoryFlag != "" {
		return directoryFlag
	}

	home, err := os.UserHomeDir()
	if err == nil {
		legacy := filepath.Join(home, ".pomodoro")
		if f, err := os.Stat(legacy); err == nil && f.IsDir() {
			return legacy
		}
	}

	return filepath.Join(xdg.ConfigHome, "pomodoro")
}
