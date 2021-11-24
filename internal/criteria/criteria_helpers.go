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
		"val":  v,
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

// extractLabelCardinalityFromPairs iterates a slice of LabelPairs.
// Upon finding the provided label of lName, and a unique label value, an int
func extractLabelValuesForName(lName string, lPairs []*dto.LabelPair) []string {
	allValsForLabelName := make([]string, 1)

	for _, l := range lPairs {
		labelName, labelValue := extractLabel(l)
		if labelName == lName {
			// e.g. service=cool-service, where lName == service
			allValsForLabelName = append(allValsForLabelName, labelValue)
		}
	}

	return allValsForLabelName
}

func inSlice(val string, s []string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}
