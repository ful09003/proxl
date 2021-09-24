package realcmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Flags
	target     string     // Endpoint to scrape with cards
	multFactor int    = 1 // Multiply all calculations by this

	rootCmd = &cobra.Command{
		Use:   "cards",
		Short: "Runs the Prometheus Exporter CARDinality Scorer (CARDS)",
		Long:  "cards is a utility meant to help gauge the quality of a Prometheus text-format exporter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := scrape(cmd, args); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	// Much of init() here is taken from the Cobra docs
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "http://localhost:9100/metrics", "HTTP endpoint to evaluate")
	rootCmd.PersistentFlags().IntVarP(&multFactor, "mult", "m", 1, "Set to maximum estimated hosts being scraped by Prometheus")

	viper.BindPFlag("target", rootCmd.PersistentFlags().Lookup("target"))
	viper.BindPFlag("mult", rootCmd.PersistentFlags().Lookup("mult"))
}

func initConfig() {
	viper.AutomaticEnv()
}

func Execute() error {
	return rootCmd.Execute()
}
