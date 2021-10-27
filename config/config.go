// Package config contains the structs and methods necessary to turn a Cards .yml file into the necessary structs for scoring a metric family
package config

import (
	"strconv"

	"github.com/ful09003/proxl/internal"
	log "github.com/sirupsen/logrus"
)

var Cfg *CardsConfig

type CardsConfig struct {
	Scorers []CardsScorerConfig `mapstructure:"scorers"`
	OutputType string `mapstructure:"output"`
}

type CardsScoringMethodConfig struct {
	Type string `mapstructure:"type"`
	Criteria []string `mapstructure:"criteria"`
}

type CardsScorerConfig struct {
	Name string `mapstructure:"name"`
	Purpose string `mapstructure:"description"`
	Criticality int `mapstructure:"criticality"`
	Method CardsScoringMethodConfig `mapstructure:"scoring-method"`
}

func (c *CardsConfig) LogConfig() {
	log.WithFields(log.Fields{
		"scorers_len": len(c.Scorers),
	}).Debug("parsed cards configuration")

	log.WithField("output-type", c.OutputType).Debug("choosing output type")

	for _, v := range c.Scorers {
		log.WithFields(log.Fields{
			"scorer_name": v.Name,
			"scorer_purpose": v.Purpose,
			"scorer_criticality": v.Criticality,
			"scorer_type": v.Method.Type,
			"scorer_criteria": v.Method.Criteria,
		}).Debug("scorer configuration")
	}
}

func ConfigToScorer(c CardsScorerConfig) (*internal.CardsScoringProcessor, error) {
	convType := aToScorer(c.Method.Type)

	newProcessor, err := internal.NewScoringProcessor(c.Name, c.Purpose, convType, c.Criticality)
	if err != nil {
		return &internal.CardsScoringProcessor{}, err
	}

	//TODO(mfuller): The below switch is probably not very scalable, but I'm not sure a better approach at this time
	switch convType{
	case internal.RegexScorer:
		return newProcessor.WithRegexScorer(c.Method.Criteria[0])
	case internal.LabelLengthScorer:
		i, err := strconv.Atoi(c.Method.Criteria[0])
		if err != nil {
			return newProcessor, err
		}
		return newProcessor.WithLabelLengthScorer(i)
	case internal.FamilyExcluderScorer:
		return newProcessor.WithMetricNameExclusionScorer(c.Method.Criteria)
	default:
		return newProcessor, nil
	}

}

func aToScorer(s string) internal.ScoringType {
	switch s{
	case "regex_scorer":
		return internal.RegexScorer
	case "label_length_scorer":
		return internal.LabelLengthScorer
	case "family_name_scorer":
		return internal.FamilyExcluderScorer
	default:
		return internal.OtherScorer
	}
}