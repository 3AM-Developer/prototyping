package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ignoreCmd = &cobra.Command{
	Use:   "ignore",
	Short: "drops in .gitignore files",
	Run: func(cmd *cobra.Command, args []string) {
		inputString := viper.GetString("input")
		inputList := parseInputString(inputString)
		fmt.Println(inputList)
	},
}

func parseInputString(inputString string) []string {
	inputString = strings.TrimSpace(inputString)
	if strings.HasPrefix(inputString, `"`) && strings.HasSuffix(inputString, `"`) {
		inputString = inputString[1 : len(inputString)-1]
	}

	inputList := strings.Split(inputString, ",")
	for i := range inputList {
		inputList[i] = strings.TrimSpace(inputList[i])
	}

	return inputList
}

func init() {
	rootCmd.AddCommand(ignoreCmd)
}
