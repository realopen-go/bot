package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/sluggishhackers/go-realopen/app"
	"github.com/sluggishhackers/go-realopen/crawler"
	"github.com/sluggishhackers/go-realopen/rmtstor"
	"github.com/sluggishhackers/go-realopen/statusmanager"
	"github.com/sluggishhackers/go-realopen/store"
)

var initCmd = &cobra.Command{
	Use:   "install",
	Short: "Install real.open.go.kr",
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Start()

		sm := statusmanager.New()
		rs := rmtstor.New()
		store := store.New()
		c := crawler.New(store, sm)
		newApp := app.New(c, rs, store, sm)

		newApp.Initialize()
		newApp.Install()

		s.Stop()
	},
}
