package criteria

import (
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

// NaiveUntypedScorer scores a MetricFamily against the total number of metrics of an unknown type
func NaiveUntypedScorer(mF *dto.MetricFamily) (int, error) {
	if *mF.Type.Enum() == dto.MetricType_UNTYPED {
		return 1, nil
	}

	return 0, nil
}

// NaiveLabelScorer scores a MetricFamily based on total number of labels.
// The function accomplishes this by summing all unique label values across the Metrics in the MetricFamily.
func NaiveLabelScorer(mF *dto.MetricFamily) (int, error) {
	var score int

	r := make(map[string]int)

	for _, m := range mF.Metric {
		vals := extractLabelPair(m.Label)
		for lName := range vals {
			r[lName] += 1
		}
	}

	for n, v := range r {
		log.WithFields(log.Fields{"label_name": n, "label_cardinality": v}).Debug("adding score")
		score += v
	}

	return score, nil
}

func extractLabelPair(lPairs []*dto.LabelPair) map[string]int {
	rVal := make(map[string]int)

	for _, lP := range lPairs {
		rVal[lP.GetName()] += 1
	}

	return rVal
}
