/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/spf13/cobra"
	dac "github.com/xinsnake/go-http-digest-auth-client"
)

var aisegHost string
var listenAddr string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aiseg2-proxy",
	Short: "AiSEG2 Proxy",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		aisegUser := os.Getenv("AISEG2_USER")
		aisegPassword := os.Getenv("AISEG2_PASSWORD")

		t := dac.NewTransport(aisegUser, aisegPassword)
		t.Client = http.DefaultClient

		director := func(request *http.Request) {
			request.URL.Scheme = "http"
			request.URL.Host = aisegHost
		}
		rp := &httputil.ReverseProxy{Director: director, Transport: &t}
		server := http.Server{
			Addr:    listenAddr,
			Handler: rp,
		}

		log.Printf("Listen on %s\n", listenAddr)

		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}

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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aiseg2-proxy.yaml)")
	rootCmd.PersistentFlags().StringVar(&listenAddr, "listen", ":9000", "listen address (default is :9000)")
	rootCmd.PersistentFlags().StringVar(&aisegHost, "aiseg", "192.168.0.216", "aiseg address (default is 192.168.0.216)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
