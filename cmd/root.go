package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "alphadns",
	Short: "alphadns is a simple and straight forward dns server",
	Long:  `Who doesn't like YAML? alphadns does! simple, straight forward, and yaml configurable!`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("welcome to alphadns!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
