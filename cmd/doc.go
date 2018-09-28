package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "doc",
	Long:  "Generate markdown documentation.",
	RunE: func(c *cobra.Command, args []string) (err error) {
		if err := doc.GenMarkdownTree(FanplaneCmd, docDirectory); err != nil {
			log.WithError(err).Fatal("Couldn't generate markdown documentation.")
			os.Exit(-1)
		}
		return
	},
}

var docDirectory string

func init() {
	docCmd.Flags().StringVarP(&docDirectory, "dir", "d", "docs/", "output directory.")
}
