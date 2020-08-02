package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sluggishhackers/go-realopen/app"
	"github.com/sluggishhackers/go-realopen/crawler"
	"github.com/sluggishhackers/go-realopen/rmtstor"
	"github.com/sluggishhackers/go-realopen/statusmanager"
	"github.com/sluggishhackers/go-realopen/store"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run real.open.go.kr",
	Run: func(cmd *cobra.Command, args []string) {
		sm := statusmanager.New()
		rs := rmtstor.New()
		store := store.New()
		c := crawler.New(store, sm)
		newApp := app.New(c, rs, store, sm)

		newApp.RunDailyCrawler()
	},
}
