package internal

import (
	"reflect"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

func TestNewModel(t *testing.T) {
	type args struct {
		forFam *dto.MetricFamily
	}
	tests := []struct {
		name string
		args args
		want *CardsScoringModel
	}{
		{
			name: "new with simple mfam",
			args: args{&dto.MetricFamily{
				Name: proto.String("blah"),
				Help: proto.String("blahhelp"),
			}},
			want: &CardsScoringModel{
				criteria: []CardsScoreFn{},
				data: &dto.MetricFamily{
					Name: proto.String("blah"),
					Help: proto.String("blahhelp"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewModel(tt.args.forFam); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsScoringModel_WithFns(t *testing.T) {
	testFn := func(mFam *dto.MetricFamily) (int, error) {
		return 1, nil
	}

	type fields struct {
		criteria []CardsScoreFn
		data     *dto.MetricFamily
	}
	type args struct {
		c []CardsScoreFn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CardsScoringModel
	}{
		{
			name: "with fns test length",
			args: args{
				c: []CardsScoreFn{testFn},
			},
			want: &CardsScoringModel{
				criteria: []CardsScoreFn{testFn},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CardsScoringModel{
				criteria: tt.fields.criteria,
				data:     tt.fields.data,
			}
			if got := m.WithFns(tt.args.c...); !reflect.DeepEqual(len(got.criteria), len(tt.want.criteria)) {
				t.Errorf("CardsScoringModel.WithFns() = %v, want %v", got, tt.want)
			}
		})
	}
}
