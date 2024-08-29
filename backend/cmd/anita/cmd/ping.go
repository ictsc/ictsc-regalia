package cmd

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Post ping to the server",
	Run: func(cmd *cobra.Command, args []string) {
		uri := "http://localhost:" + strconv.Itoa(config.Port) + "/ping"
		u, err := url.Parse(uri)
		if err != nil {
			log.Panic(err)
		}

		resp, err := http.Get(u.String())
		if err != nil {
			log.Panic(err)
		}

		if err = resp.Body.Close(); err != nil {
			log.Panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
