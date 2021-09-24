package criteria

import (
	dto "github.com/prometheus/client_model/go"
)

// NaiveUntypedScorer scores a MetricFamily against the total number of metrics of an unknown type
func NaiveUntypedScorer(mF *dto.MetricFamily) (int, error) {
	if *mF.Type.Enum() == dto.MetricType_UNTYPED {
		return 1, nil
	}

	return 0, nil
}
