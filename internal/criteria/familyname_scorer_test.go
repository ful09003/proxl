package criteria

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

func TestFamilyNameCheck(t *testing.T) {
	type args struct {
		m *dto.MetricFamily
		n []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "matching name found",
			args: args{
				m: &dto.MetricFamily{
					Name: proto.String("matching_family_name"),
				},
				n: []string{"matching_family_name", "some_other_name"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "no matching name found",
			args: args{
				m: &dto.MetricFamily{
					Name: proto.String("no_matching_family_name"),
				},
				n: []string{"matching_family_name", "some_other_name"},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FamilyNameCheck(tt.args.m, tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("FamilyNameCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FamilyNameCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
