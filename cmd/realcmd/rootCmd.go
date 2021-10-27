package realcmd

import (
	"github.com/ful09003/proxl/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Flags
	target     string     // Endpoint to scrape with cards
	ll         string     // Log level, used to init rootCmd

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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetLevel(getLogLevel(ll))
			logrus.WithField("level", ll).Info("setting log level")
			logrus.SetFormatter(&logrus.JSONFormatter{})
			config.Cfg.LogConfig()

			return nil
		},
	}
)

func init() {
	// Much of init() here is taken from the Cobra docs
	cobra.OnInitialize(initConfig)

	// Set config global, weeee.
	config.Cfg = parseConfigFile()

	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "http://localhost:9100/metrics", "HTTP endpoint to evaluate")
	rootCmd.PersistentFlags().StringVar(&ll, "loglevel", "info", "log level to set")

	viper.BindPFlag("target", rootCmd.PersistentFlags().Lookup("target"))
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
}

func parseConfigFile() *config.CardsConfig {
	viper.AddConfigPath(".")
	viper.SetConfigName("cards")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic(err)
	}
	c := &config.CardsConfig{}
	if err := viper.Unmarshal(c); err != nil {
		logrus.Panic(err)
	}

	return c
}

func initConfig() {
	viper.SetEnvPrefix("cards")
	viper.AutomaticEnv()
}

func Execute() error {
	return rootCmd.Execute()
}

// getLogLevel calls logrus.ParseLevel(), logging any errors and returning info level
func getLogLevel(l string) logrus.Level {
	lvl, e := logrus.ParseLevel(l)
	if e != nil {
		logrus.WithError(e).Error("error setting log level")
		lvl = logrus.InfoLevel
	}

	return lvl
}
