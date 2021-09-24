package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

type FailableReader struct{}

func (f *FailableReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("ayyyyy")
}

func TestReadTarget(t *testing.T) {
	type args struct {
		t io.Reader
	}
	s := `blah`
	sR := strings.NewReader(s)
	sFR := &FailableReader{}
	sRE := bytes.NewBufferString(s)

	tests := []struct {
		name    string
		args    args
		want    bytes.Buffer
		wantErr bool
	}{
		// Happy path, string in, bytes.Buffer for the string out
		{
			name:    "blah",
			args:    args{sR},
			want:    *sRE,
			wantErr: false,
		},
		// Error from an io.Reader propagates back expectedly
		{
			name:    "error propagation",
			args:    args{sFR},
			want:    bytes.Buffer{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadTarget(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMetricFamilies(t *testing.T) {
	type args struct {
		b bytes.Buffer
	}

	badMetricsString := `#blah
blahblah
	`

	// Extracted from https://prometheus.io/docs/instrumenting/exposition_formats/#text-format-example
	goodMetricsString := `
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total counter
http_requests_total{method="post",code="200"} 1027 1395066363000
http_requests_total{method="post",code="400"}    3 1395066363000
`

	var badMetricsBuffer, goodMetricsBuffer bytes.Buffer

	badMetricsBuffer.WriteString(badMetricsString)
	goodMetricsBuffer.WriteString(goodMetricsString)

	tests := []struct {
		name    string
		args    args
		want    map[string]*dto.MetricFamily
		wantErr bool
	}{
		// Case of a bad metrics string
		{
			name:    "bad metric bytes",
			args:    args{badMetricsBuffer},
			want:    map[string]*dto.MetricFamily{},
			wantErr: true,
		},
		// Case of valid parsing metrics
		{
			name: "good metric bytes",
			args: args{goodMetricsBuffer},
			want: map[string]*dto.MetricFamily{
				"http_requests_total": &dto.MetricFamily{
					Name: proto.String("http_requests_total"),
					Help: proto.String("The total number of HTTP requests."),
					Type: dto.MetricType_COUNTER.Enum(),
					Metric: []*dto.Metric{
						{
							Label: []*dto.LabelPair{
								{
									Name:  proto.String("method"),
									Value: proto.String("post"),
								},
								{
									Name:  proto.String("code"),
									Value: proto.String("200"),
								},
							},
							Counter: &dto.Counter{
								Value: proto.Float64(1027),
							},
							TimestampMs: proto.Int64(1395066363000),
						},
						{
							Label: []*dto.LabelPair{
								{
									Name:  proto.String("method"),
									Value: proto.String("post"),
								},
								{
									Name:  proto.String("code"),
									Value: proto.String("400"),
								},
							},
							Counter: &dto.Counter{
								Value: proto.Float64(3),
							},
							TimestampMs: proto.Int64(1395066363000),
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMetricFamilies(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMetricFamilies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMetricFamilies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCardsHttpScraper(t *testing.T) {
	type args struct {
		endpoint string
		retries  int
	}
	tests := []struct {
		name string
		args args
		want CardsHttpScraper
	}{
		{
			name: "new cards scraper",
			args: args{endpoint: "my.house", retries: 1},
			want: CardsHttpScraper{
				maxRetries: 1,
				Endpoint:   "my.house",
				b:          bytes.Buffer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCardsHttpScraper(tt.args.endpoint, tt.args.retries); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCardsHttpScraper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsHttpScraper_ScrapeTarget(t *testing.T) {

	srv := newTestHTTPServer(t)
	defer srv.Close()

	type fields struct {
		Endpoint   string
		maxRetries int
		b          bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]*dto.MetricFamily
		wantErr bool
	}{
		// Happy path
		{
			name: "happy scrape",
			fields: fields{
				Endpoint: fmt.Sprintf("%s/good_case", srv.URL),
				b:        bytes.Buffer{},
			},
			want: map[string]*dto.MetricFamily{
				"joy_felt_total": &dto.MetricFamily{
					Name: proto.String("joy_felt_total"),
					Help: proto.String("A counter of joy experienced."),
					Type: dto.MetricType_COUNTER.Enum(),
					Metric: []*dto.Metric{
						{
							Label: []*dto.LabelPair{
								{
									Name:  proto.String("developer"),
									Value: proto.String("me"),
								},
							},
							Counter: &dto.Counter{
								Value: proto.Float64(9000),
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &CardsHttpScraper{
				Endpoint:   tt.fields.Endpoint,
				maxRetries: tt.fields.maxRetries,
				b:          tt.fields.b,
			}
			got, err := n.ScrapeTarget()
			if (err != nil) != tt.wantErr {
				t.Errorf("CardsHttpScraper.ScrapeTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CardsHttpScraper.ScrapeTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
