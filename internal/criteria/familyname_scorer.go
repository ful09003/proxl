package criteria

import (
	dto "github.com/prometheus/client_model/go"
)

//FamilyNameCheck checks the name of the MetricFamily and returns true if the name exists in the input string slice, along with any errors encountered.
func FamilyNameCheck(m *dto.MetricFamily, n []string) (bool, error) {
	s := m.GetName()

	for _, excludeName := range n {
		if s == excludeName {
			return true, nil
		}
	}

	return false, nil
}