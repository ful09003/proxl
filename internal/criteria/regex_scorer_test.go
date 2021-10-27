package criteria

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

func TestMetricLabelRegex(t *testing.T) {
	type args struct {
		m *dto.Metric
		p string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// Matching regex (UUID label value)
		{
			name: "matching regex (uuid label)",
			args: args{
				m: &dto.Metric{
					Label: []*dto.LabelPair{
						{
							Name:  proto.String("plain-name"),
							Value: proto.String("5f8d0583-d61c-4dbd-b2f2-bc3b2a7463bd"),
						},
						{
							Name:  proto.String("another-plain-name"),
							Value: proto.String("4444444"),
						},
					},
				},
				p: "^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$",
			},
			want:    true,
			wantErr: false,
		},
		// No matching regex
		{
			name: "no matching regex (uuid label)",
			args: args{
				m: &dto.Metric{
					Label: []*dto.LabelPair{
						{
							Name:  proto.String("plain-name"),
							Value: proto.String("plain-value"),
						},
						{
							Name:  proto.String("another-plain-name"),
							Value: proto.String("another-plain-value"),
						},
					},
				},
				p: "^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$",
			},
			want:    false,
			wantErr: false,
		},
		// Matches label names also
		{
			name: "label name matching",
			args: args{
				m: &dto.Metric{
					Label: []*dto.LabelPair{
						{
							Name:  proto.String("5f8d0583-d61c-4dbd-b2f2-bc3b2a7463bd"),
							Value: proto.String("groovy"),
						},
					},
				},
				p: "^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$",
			},
			want:    true,
			wantErr: false,
		},
		// Non-compiling regex, at least in current implementation :P
		{
			name: "non-compiling regex",
			args: args{
				m: &dto.Metric{},
				p: "\\😎",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MetricLabelRegex(tt.args.m, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricLabelRegex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricLabelRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
