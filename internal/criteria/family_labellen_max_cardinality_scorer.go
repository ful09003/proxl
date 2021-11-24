package criteria

import (
	"fmt"

	dto "github.com/prometheus/client_model/go"
)

// FamilyLabelMaxCardinality evaluates if any label in a MetricFamily has a max cardinality above the target length.
// Care should be taken that this criteria is not applied to histograms which are by nature likely to be flagged.
// This criteria must traverse all Metrics in a Family to build an accurate representation of per-label cardinality in the MetricFamily, thus we pass-through from our scorer
// BUG(mfuller): This is a very poorly-optimized function, and a great one to think about and improve <3
func FamilyLabelMaxCardinality(m *dto.MetricFamily, l int) (bool, error) {
	if m.GetType() == dto.MetricType_HISTOGRAM {
		// The caller did not heed advice; return an error here.
		return false, fmt.Errorf("invalid type for scorer (found %s), %s not supported", m.GetType(), dto.MetricType_HISTOGRAM.Enum())
	}

	labelValTermFreq := make(map[string][]string)

	for _, metric := range m.Metric {
		// Traverse Metrics, extracting label names/values for all Metric.Label
		// This is one spot where optimizations could happen?
		for _, label := range metric.Label {
			lN, lV := extractLabel(label)
			labelValTermFreq[lN] = append(labelValTermFreq[lN], lV)
		}
	}

	for _, v := range labelValTermFreq {
		if len(v) >= l {
			return true, nil
		}
	}
	return false, nil
}
