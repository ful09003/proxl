package criteria

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

func TestFamilyLabelMaxCardinality(t *testing.T) {
	testValidFlaggedMetricFamily := &dto.MetricFamily{
		Name: proto.String("test_metric"),
		Type: dto.MetricType_COUNTER.Enum(),
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair{
					{
						Name:  proto.String("label1"),
						Value: proto.String("val1"),
					},
					{
						Name:  proto.String("label1"),
						Value: proto.String("val2"),
					},
				},
			},
		},
	}
	testInvalidHistogramMetric := &dto.MetricFamily{
		Name: proto.String("test_metric"),
		Type: dto.MetricType_HISTOGRAM.Enum(),
		Metric: []*dto.Metric{
			{
				Histogram: &dto.Histogram{
					SampleCount: proto.Uint64(10),
					SampleSum:   proto.Float64(10),
					Bucket: []*dto.Bucket{
						{
							CumulativeCount: proto.Uint64(4),
						},
					},
				},
				Label: []*dto.LabelPair{
					{
						Name:  proto.String("label1"),
						Value: proto.String("val1"),
					},
					{
						Name:  proto.String("label1"),
						Value: proto.String("val2"),
					},
				},
			},
		},
	}
	type args struct {
		m *dto.MetricFamily
		l int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "happy path",
			args: args{
				m: testValidFlaggedMetricFamily,
				l: 2,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "histograms fail",
			args: args{
				m: testInvalidHistogramMetric,
				l: 1,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FamilyLabelMaxCardinality(tt.args.m, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("FamilyLabelMaxCardinality() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FamilyLabelMaxCardinality() = %v, want %v", got, tt.want)
			}
		})
	}
}
