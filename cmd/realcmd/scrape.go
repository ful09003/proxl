package realcmd

import (
	"github.com/ful09003/cards/config"
	"github.com/ful09003/cards/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func scrape(cmd *cobra.Command, args []string) error {
	// Set up our configured scorers
	var scorers []*internal.CardsScoringProcessor

	for _, m := range config.Cfg.Scorers {
		c, err := config.ConfigToScorer(m)
		if err != nil {
			return err
		}
		scorers = append(scorers, c)
	}

	log.WithField("scorers_len", len(scorers)).Debug("configured scorers")
	
	scraper := internal.NewCardsHttpScraper(target, 1)

	mFam, err := scraper.ScrapeTarget()
	if err != nil {
		return err
	}

	for fName, fam := range mFam {
		// For each metric family, run each scorer against it and return results
		log.WithField("family_name", fName).Debug("evaluating metric family")
		for _, s := range scorers {
			flagged, err := s.Score(fam)
			if err != nil {
				log.WithField("family_name", fName).WithError(err).Error("error during scoring family")
			}
			if flagged {
				log.WithFields(log.Fields{
					"family_name": fName,
					"scorer_name": s.Name,
					"hint": s.Purpose,
				}).Info("Family failed")
			}
		}
	}

	return nil

}

/*func scrape(cmd *cobra.Command, args []string) error {
	scraper := internal.NewCardsHttpScraper(target, 1)
	tWriter := tabwriter.NewWriter(os.Stdout, 20, 8, 1, ' ', tabwriter.AlignRight|tabwriter.TabIndent)
	defer tWriter.Flush()

	mF, err := scraper.ScrapeTarget()
	if err != nil {
		return err
	}

	fmt.Fprintln(tWriter, "Metric\tScore\t")
	for mName, mData := range mF {
		// TODO: Figure out how to make this _way_ more dynamic instead of proof-of-concepty
		// BUG: This is also probably an excellent opportunity to figure out wtf channels and concurrency can do...
		scorer := internal.NewModel(mData).WithFns(
			evals.NaiveUntypedScorer,
			evals.NaiveLabelScorer)
		finalScore, err := scorer.Evaluate()

		if err != nil {
			return fmt.Errorf("error evaluating metric family: %w", err)
		}

		buildLn := func(s string, i int) string {
			return fmt.Sprintf("%s\t%d\t", s, i)
		}

		fmt.Fprintln(tWriter, buildLn(mName, finalScore))
	}

	return nil
}*/