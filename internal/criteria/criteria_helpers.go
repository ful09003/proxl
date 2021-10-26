package criteria

import (
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

// extractLabel returns two strings from a LabelPair, representing the label name and value
func extractLabel(lPair *dto.LabelPair) (string, string) {
	n, v := lPair.GetName(), lPair.GetValue()
	
	log.WithFields(log.Fields{
		"name": n,
		"val": v,
	}).Debug("extracted lpair")
	
	return n, v
}

func extractLabelPair(lPairs []*dto.LabelPair) map[string]int {
	rVal := make(map[string]int)

	for _, lP := range lPairs {
		rVal[lP.GetName()] += 1
	}

	return rVal
}
