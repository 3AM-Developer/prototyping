package main

import (
	"fmt"

	"github.com/JaegyuDev/Icarus/cmd"
)

func main() {
	rootCmd := cmd.NewCMD()
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("I don't know how this happened: %s", err.Error())
	}
}
