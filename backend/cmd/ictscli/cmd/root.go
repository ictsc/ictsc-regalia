package cmd

import (
	"os"

	"github.com/ictsc/ictsc-regalia/backend/cmd/ictscli/client"
	"github.com/ictsc/ictsc-regalia/backend/cmd/ictscli/cmd/team"
	"github.com/spf13/cobra"
)

var (
	baseURL  string
	useJSON  bool
)

var rootCmd = &cobra.Command{
	Use:   "ictscli",
	Short: "ICTSC Regalia CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		client.InitClients(baseURL)
		client.SetJSONOutput(useJSON)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "base", "http://localhost:8080", "Base URL of the backend server")
	rootCmd.PersistentFlags().BoolVar(&useJSON, "json", false, "Output in JSON format")
	rootCmd.AddCommand(team.TeamCmd)
}
