package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/sluggishhackers/realopen.go/app"
	"github.com/sluggishhackers/realopen.go/crawler"
	"github.com/sluggishhackers/realopen.go/rmtstor"
	"github.com/sluggishhackers/realopen.go/statmanager"
	"github.com/sluggishhackers/realopen.go/store"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize real.open.go.kr",
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		sm := statmanager.New()
		rs := rmtstor.New()
		store := store.New()
		c := crawler.New(store, sm)
		newApp := app.New(c, rs, store, sm)

		newApp.Initialize()

		s.Stop()
	},
}
