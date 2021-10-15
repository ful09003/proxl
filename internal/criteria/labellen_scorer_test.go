package criteria

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

func TestLabelLengthCheck(t *testing.T) {
	mFam := &dto.MetricFamily{
		Name: proto.String("test"),
		Help: proto.String("unit test"),
		Type: dto.MetricType_GAUGE.Enum(),
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair{
					{
						Name: proto.String("firstlabel"),
						Value: proto.String("firstlabel-value"),
					},
					{
						Name: proto.String("secondlabel"),
						Value: proto.String("secondlabel-value"),
					},
					{
						Name: proto.String("thirdlabel"),
						Value: proto.String("thirdlabel-value"),
					},
				},
			},
		},
	}
	mFamWithNoLabels := &dto.MetricFamily{
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair{},
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
			name: "length triggers match",
			args: args{
				m: mFam,
				l: 3,
			},
			want: true,
			wantErr: false,
		},
		{
			name: "length does not trigger",
			args: args{
				m: mFam,
				l: 5,
			},
			want: false,
			wantErr: false,
		},
		{
			name: "no labels match as expected",
			args: args{
				m: mFamWithNoLabels,
				l: 5,
			},
			want: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LabelLengthCheck(tt.args.m.Metric[0], tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("LabelLengthCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LabelLengthCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
