package realcmd

import (
	"os"

	"github.com/ful09003/cards/internal"
	expfmt "github.com/prometheus/common/expfmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(debugCmd)
}

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Scrape Prometheus exporter",
	Long: `
Scrapes the desired Prometheus exporter target, and prints its output. Save yourself switching back to a browser!	
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := internal.NewCardsHttpScraper(target, 1)
		mF, err := s.ScrapeTarget()
		if err != nil {
			return err
		}

		for _, actualFam := range mF {
			expfmt.MetricFamilyToText(os.Stdout, actualFam)
		}

		return nil
	},
}
