package internal

import (
	"errors"

	"github.com/ful09003/proxl/internal/criteria"
	log "github.com/sirupsen/logrus"

	dto "github.com/prometheus/client_model/go"
)

type ScoringType int

func (f ScoringType) String() string {
	switch f {
	case SimpleScorer:
		return "simple scorer"
	case RegexScorer:
		return "regex-capable scorer"
	case FamilyExcluderScorer:
		return "metric family name scorer"
	default:
		return "unknown scorer type"
	}
}

const (
	SimpleScorer ScoringType = iota // A simplistic scorer
	LabelLengthScorer // Label length scoring
	FamilyExcluderScorer // Family name exclusion
	RegexScorer // A regex-capable scorer
	OtherScorer // Everything else
)

type CardsScoringProcessor struct {
	Name string // Name of this processor
	Purpose string // Purpose of this processor
	ScorerType ScoringType // Type of scorer this is
	Criticality int // How important this scorer is to the user

	Evaluator CardsEvaluator
}

// NewScoringProcessor creates and returns a pointer to a Cards Scoring Processor
func NewScoringProcessor(name, purpose string, pT ScoringType, criticality int) (*CardsScoringProcessor, error) {
	
	if name == "" || purpose == "" {
		return &CardsScoringProcessor{}, errors.New("scoring processor must have name and purpose")
	}

	return &CardsScoringProcessor{
		Name: name,
		Purpose: purpose,
		Criticality: criticality,
		ScorerType: pT,
	}, nil
}

// WithRegexScorer accepts a regex-compatible string and ensures that the Processor has a regex parser which uses the input pattern. 
func (sp *CardsScoringProcessor) WithRegexScorer(s string) (*CardsScoringProcessor, error) {

	if sp.ScorerType != RegexScorer {
		return sp, errors.New("cannot add regex scorer to a non-regex processor")
	}

	newParser := &CardsRegexEvaluator{
		r: s,
		Type: sp.ScorerType,
	}

	return setEvaluator(sp, newParser), nil
}

// WithLabelLengthScorer sets the Processor to use a label length evaluation against the given max-desired length, i. Errors, if encountered, are returned.
func (sp *CardsScoringProcessor) WithLabelLengthScorer(i int) (*CardsScoringProcessor, error) {

	if sp.ScorerType != LabelLengthScorer {
		return sp, errors.New("cannot add label length scorer to a non-length processor")
	}

	newParser := &CardsLabelLengthEvaluator{
		maxLen: i,
		Type: sp.ScorerType,
	}

	return setEvaluator(sp, newParser), nil
}

// WithMetricNameExclusionScorer sets the Processor to evaluate metric family names against a given input list.
func (sp *CardsScoringProcessor) WithMetricNameExclusionScorer(n []string) (*CardsScoringProcessor, error) {
	if sp.ScorerType != FamilyExcluderScorer {
		return sp, errors.New("cannot use family name scorer with a non-family-name processor")
	}

	newParser := &CardsFamilyNameEvaluator{
		Type: sp.ScorerType,
		excludeList: n,
	}

	return setEvaluator(sp, newParser), nil
}

//setEvaluator sets the evaluator ofa CardsScoringProcessor, logging if the operation overwrites any existing processor.
func setEvaluator(sp *CardsScoringProcessor, e CardsEvaluator) *CardsScoringProcessor {
	if sp.Evaluator != nil {
		// Allow overwriting, but log that it happened
		log.WithField("processor", sp.Name).Info("processor has evaluator, overwriting")
	}

	sp.Evaluator = e
	return sp
}

// Score runs the Processor's evaluator against the provided MetricFamily.
// If the Processor's evaluation determines the provided MetricFamily is a violation, true is returned, along with optional error
func (sp *CardsScoringProcessor) Score(mf *dto.MetricFamily) (bool, error) {
	return sp.Evaluator.Evaluate(mf)
}

// CardsEvaluator is an interface which all concrete evaluators (e.g. regex, label-length, etc.) implement
type CardsEvaluator interface {
	Evaluate(f *dto.MetricFamily) (bool, error)
}

// CardsRegexEvaluator is a struct representing a regex-utilizing evaluator
type CardsRegexEvaluator struct {
	Type ScoringType
	r string // Regex string to be evaluated
}

// Evaluate iterates through a MetricFamily's Metrics.
// Upon any Metric label name or value which matches the Evaluator's defined regex, 'true' is returned. False indicates no label names or values matched the Evaluator's internal regex. Any error encountered during processing is returned, optionally.
func (r *CardsRegexEvaluator) Evaluate(f *dto.MetricFamily) (bool, error) {
	for _, m := range f.Metric {
		matched, err := criteria.MetricLabelRegex(m, r.r)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}

// CardsLabelLengthEvaluator is an evaluator which scores based on slice length of unique label names.
// Note: This Evaluator does not account for built-in Prometheus labels job and instance, thus maxLen should generally be set to reflect this.
type CardsLabelLengthEvaluator struct {
	Type ScoringType
	maxLen int // Max length of names this evaluator accepts
}

// Evaluate uses the first metric within a MetricFamily.
// If the total length of the metric's labels exceeds the defined length, 'true' is returned. Any error encountered during processing is returned, optionally.
func (r *CardsLabelLengthEvaluator) Evaluate(f *dto.MetricFamily) (bool, error) {
	return criteria.LabelLengthCheck(f.Metric[0], r.maxLen)
}


type CardsFamilyNameEvaluator struct {
	Type ScoringType
	excludeList []string // List of MetricFamily names to exclude
}

func (r *CardsFamilyNameEvaluator) Evaluate(f *dto.MetricFamily) (bool, error) {
	return criteria.FamilyNameCheck(f, r.excludeList)
}
