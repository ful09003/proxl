package realcmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ful09003/cards/internal"
	evals "github.com/ful09003/cards/internal/criteria"
	"github.com/spf13/cobra"
)

func scrape(cmd *cobra.Command, args []string) error {
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
			evals.NaiveUntypedScorer)
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
}
