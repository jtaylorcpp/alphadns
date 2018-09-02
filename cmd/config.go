package main

import (
	"log"

	"github.com/jtaylorcpp/alphadns"
	"github.com/spf13/cobra"
)

var configFile string

func init() {
	runCmd.Flags().StringVarP(&configFile, "config-file", "f", "", "Config file to run from")

	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a alphadns server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("starting dns servers from file: ", configFile)
		dnsServer := alphadns.LoadFromFile(configFile)
		log.Println("dns server: ", dnsServer)
		dnsServer.Start()
	},
}
