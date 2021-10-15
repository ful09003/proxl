package criteria

import (
	dto "github.com/prometheus/client_model/go"
)

//LabelLengthCheck checks the first Metric of a MetricFamily and returns true if the length of the Metric Labels exceeds the desired length, along with any errors encountered.
func LabelLengthCheck(m *dto.Metric, l int) (bool, error) {

	labels := extractLabelPair(m.Label)
	return len(labels) >= l, nil
}