package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "your-command",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		// Your command logic goes here
	},
}

func NewCMD() *cobra.Command {
	return rootCmd
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

func genDefaultConfig(path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0664)
	if err != nil {
		return err
	}

	defaultConfig := map[string]interface{}{
		"ignore_map": map[string]string{},
	}

	viper.SetDefault("ignore_map", defaultConfig["ignore_map"])
	err = viper.SafeWriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func initConfig() {
	homeDir, err := getHomeDir()
	if err != nil {
		fmt.Printf("Error getting homedir: %s", err.Error())
		os.Exit(1)
	}

	configPath := filepath.Join(homeDir, ".icarus", ".config")

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	// Read the config file
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok || os.IsNotExist(err) {
			err := genDefaultConfig(configPath)
			if err != nil {
				fmt.Printf("Error writing default config: %s", err.Error())
				os.Exit(1)
			}

		} else {
			fmt.Printf("Error reading config: %s", err.Error())
			os.Exit(1)
		}
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}
