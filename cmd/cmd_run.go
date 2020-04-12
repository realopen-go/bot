package cmd

import (
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/cobra"

	"github.com/sluggishhackers/realopen.go/app"
	"github.com/sluggishhackers/realopen.go/crawler"
	"github.com/sluggishhackers/realopen.go/rmtstor"
	"github.com/sluggishhackers/realopen.go/statmanager"
	"github.com/sluggishhackers/realopen.go/store"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run real.open.go.kr",
	Run: func(cmd *cobra.Command, args []string) {
		sm := statmanager.New()
		rs := rmtstor.New()
		store := store.New()
		c := crawler.New(store, sm)
		newApp := app.New(c, rs, store, sm)

		newApp.RunDailyCrawler()
		<-gocron.Start()
	},
}
