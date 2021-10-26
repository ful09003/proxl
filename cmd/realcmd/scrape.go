package realcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ful09003/cards/config"
	"github.com/ful09003/cards/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type flaggedEntity struct {
	ScorerName, FamilyName, ScorerDesc string
	ScorerCriticality int
}

func (f *flaggedEntity) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%d\t", f.FamilyName, f.ScorerName, f.ScorerDesc, f.ScorerCriticality)
}

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

	var flaggedEntities []flaggedEntity

	for fName, fam := range mFam {
		// For each metric family, run each scorer against it and return results
		log.WithField("family_name", fName).Debug("evaluating metric family")
		for _, s := range scorers {
			flagged, err := s.Score(fam)
			if err != nil {
				log.WithField("family_name", fName).WithError(err).Error("error during scoring family")
			}
			if flagged {
				flaggedEntities = append(flaggedEntities, flaggedEntity{
					ScorerName: s.Name,
					FamilyName: fName,
					ScorerDesc: s.Purpose,
					ScorerCriticality: s.Criticality,
				})

				log.WithFields(log.Fields{
					"family_name": fName,
					"scorer_name": s.Name,
					"hint": s.Purpose,
				}).Debug("family flagged")
			}
		}
	}

	// Write out flagged entities, if any
	if len(flaggedEntities) > 0 {
		switch config.Cfg.OutputType {
		case "table":
			tWriter := tabwriter.NewWriter(os.Stdout, 20, 8, 1, ' ', tabwriter.AlignRight|tabwriter.TabIndent)
			defer tWriter.Flush()

			fmt.Fprintln(tWriter, "Metric\tScorer\tScorer Description\tCriticality\t")

			for _, f := range flaggedEntities {
				fmt.Fprintln(tWriter, f.String())
			}
		case "json":
		default:
			// JSON marshal and write, then we're done
			marshalled, err := json.Marshal(flaggedEntities)
			if err != nil {
				log.WithError(err).Error("failed marshalling results")
			}
			fmt.Println(string(marshalled))

		}
	}

	return nil

}
