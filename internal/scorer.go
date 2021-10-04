package internal

import (
	dto "github.com/prometheus/client_model/go"
)

type CardsScoreFn func(mFam *dto.MetricFamily) (int, error)

type CardsScoringModel struct {
	criteria []CardsScoreFn
	data     *dto.MetricFamily
}

// NewModel returns a Cards scoring model for the input MetricFamily.
func NewModel(forFam *dto.MetricFamily) *CardsScoringModel {
	return &CardsScoringModel{
		criteria: []CardsScoreFn{},
		data:     forFam,
	}
}

// WithFns adds a variable number of functions to a Cards ScoringModel
func (m *CardsScoringModel) WithFns(c ...CardsScoreFn) *CardsScoringModel {
	m.criteria = append(m.criteria, c...)

	return m
}

func (m *CardsScoringModel) Evaluate() (int, error) {
	var score int

	for _, f := range m.criteria {
		evalScore, err := f(m.data)
		if err != nil {
			return score, err
		}
		score += evalScore
	}

	return score, nil
}
