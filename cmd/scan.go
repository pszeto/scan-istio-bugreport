package cmd

import (
	"os"

	"github.com/pszeto/scan-istio-bugreport/pkg/scan"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfg = &scan.Config{}

var rootCmd = &cobra.Command{
	Use:   "create-fake-service",
	Short: "Scans istio-bug report and performs various task.",
	Run: func(cmd *cobra.Command, args []string) {
		scan.New(cfg).Entry()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&cfg.BugReportDirector, "directory", "", "Specify bug report directory.")
	rootCmd.Flags().BoolVar(&cfg.GenerateFakeService, "generate-fake-service", false, "Specify if it's a dry run. Default false")
}
