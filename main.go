package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "howoldis <domain>",
	Short: "howoldis - how old is this website?",
	Long: `howoldis
See the first date for which the Wayback Machine has records for a given domain.
    `,
	Version: "v0",
	Run: func(cmd *cobra.Command, args []string) {
		site := args[0]

		params := url.Values{}
		params.Add("url", site)
		params.Add("collection", "web")
		params.Add("output", "json")

		resp, err := http.Get("https://web.archive.org/__wb/sparkline?" + params.Encode())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: "+err.Error())
		}

		var result struct {
			FirstTs string `json:"first_ts"`
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: "+err.Error())
		}

		if result.FirstTs == "" {
			fmt.Fprintf(os.Stderr, "Failed to find %s.", site)
		}

		date, _ := time.Parse("20060102", result.FirstTs[:8])
		fmt.Printf("%s is available at least since %s, %s.",
			site, date.Format("Jan 2006"), humanize.Time(date))
	},
}
