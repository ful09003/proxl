/*
Package criteria homes various cards scoring functions. These are broken up logically by filename, with 'naive' scoring methods being found in naive.go, and so on.

Scoring functions are expected to take in a pointer to a Prometheus MetricFamily, and return an int corresponding to the metricfamily "score" and optional error.
*/
package criteria

import (
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

// NaiveUntypedScorer scores a MetricFamily against the total number of metrics of an unknown type.
// This is a 'naive' scorer, as an unknown metric type is not always a signal of poor exporter quality.
func NaiveUntypedScorer(mF *dto.MetricFamily) (int, error) {
	if *mF.Type.Enum() == dto.MetricType_UNTYPED {
		log.WithField("mfam_name", mF.Name).Debug("untyped metricfamily found")
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