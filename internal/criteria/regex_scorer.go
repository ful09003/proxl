package criteria

import (
	"regexp"

	dto "github.com/prometheus/client_model/go"
)

// MetricLabelRegex takes a Metric pointer as well as string pattern, and evaluates a regex against all label names *and* values. Matches are represented by a return boolean, along with any errors.
func MetricLabelRegex(m *dto.Metric, p string) (bool, error) {
	exp, err := regexp.Compile(p)

	if err != nil {
		return false, err
	}

	for _, l := range m.Label {
		lName, lVal := extractLabel(l)
		if exp.MatchString(lVal) || exp.MatchString(lName) {
			return true, nil
		}
	}

	return false, nil
}
