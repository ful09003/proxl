package internal

import (
	"reflect"
	"testing"

	dto "github.com/prometheus/client_model/go"
	"google.golang.org/protobuf/proto"
)

var testFn = func(mFam *dto.MetricFamily) (int, error) {
	return 1, nil
}


func TestNewScoringProcessor(t *testing.T) {
	type args struct {
		name        string
		purpose     string
		pT          ScoringType
		criticality int
	}
	tests := []struct {
		name    string
		args    args
		want    *CardsScoringProcessor
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "happy path",
			args: args{
				name:        "test",
				purpose:     "unit test",
				pT:          SimpleScorer,
				criticality: 1,
			},
			want: &CardsScoringProcessor{
				Name:        "test",
				Purpose:     "unit test",
				ScorerType:  SimpleScorer,
				Criticality: 1,
			},
		},
		{
			name: "missing name",
			args: args{
				name:    "",
				purpose: "unit test",
			},
			want:    &CardsScoringProcessor{},
			wantErr: true,
		},
		{
			name: "missing purpose",
			args: args{
				name:    "test",
				purpose: "",
			},
			want:    &CardsScoringProcessor{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewScoringProcessor(tt.args.name, tt.args.purpose, tt.args.pT, tt.args.criticality)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewScoringProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScoringProcessor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsScoringProcessor_WithRegexScorer(t *testing.T) {

	validEvaluator := &CardsRegexEvaluator{
		r:    "^.*",
		Type: RegexScorer,
	}

	anotherValidEvaluator := &CardsRegexEvaluator{
		r:    "^..*",
		Type: RegexScorer,
	}

	type fields struct {
		Name        string
		Purpose     string
		ScorerType  ScoringType
		Criticality int
		Evaluator   CardsEvaluator
	}
	type args struct {
		c string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CardsScoringProcessor
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "happy path",
			fields: fields{
				Name:        "test",
				Purpose:     "unit test",
				ScorerType:  RegexScorer,
				Criticality: 1,
			},
			args: args{
				c: "^.*",
			},
			want: &CardsScoringProcessor{
				Name:        "test",
				Purpose:     "unit test",
				ScorerType:  RegexScorer,
				Criticality: 1,
				Evaluator:   validEvaluator,
			},
		},
		{
			name: "overwriting evaluator",
			fields: fields{
				Name:        "test",
				Purpose:     "unit test",
				ScorerType:  RegexScorer,
				Criticality: 1,
				Evaluator:   validEvaluator,
			},
			args: args{
				c: anotherValidEvaluator.r,
			},
			want: &CardsScoringProcessor{
				Name:        "test",
				Purpose:     "unit test",
				ScorerType:  RegexScorer,
				Criticality: 1,
				Evaluator:   anotherValidEvaluator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &CardsScoringProcessor{
				Name:        tt.fields.Name,
				Purpose:     tt.fields.Purpose,
				ScorerType:  tt.fields.ScorerType,
				Criticality: tt.fields.Criticality,
				Evaluator:   tt.fields.Evaluator,
			}
			got, err := sp.WithRegexScorer(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardsScoringProcessor.WithRegexScorer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CardsScoringProcessor.WithRegexScorer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsRegexEvaluator_Evaluate(t *testing.T) {
	testFam := &dto.MetricFamily{
		Name: proto.String("testfam"),
		Help: proto.String("unit test struct"),
		Type: dto.MetricType_COUNTER.Enum(),
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair{
					{
						Name:  proto.String("instance"),
						Value: proto.String("localhost:9999"),
					},
					{
						Name:  proto.String("wildlabel"),
						Value: proto.String("5f8d0583-d61c-4dbd-b2f2-bc3b2a7463bd"),
					},
				},
			},
		},
	}

	type fields struct {
		Type ScoringType
		r    string
	}
	type args struct {
		f *dto.MetricFamily
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "regex matches",
			fields: fields{
				r:    "^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$",
				Type: RegexScorer,
			},
			args: args{
				f: testFam,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "regex doesn't match",
			fields: fields{
				r:    "^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{11}$",
				Type: RegexScorer,
			},
			args: args{
				f: testFam,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "regex error bubbles",
			fields: fields{
				r: "\\😎",
			},
			args: args{
				f: testFam,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardsRegexEvaluator{
				Type: tt.fields.Type,
				r:    tt.fields.r,
			}
			got, err := r.Evaluate(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardsRegexEvaluator.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CardsRegexEvaluator.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsScoringProcessor_Score(t *testing.T) {
	type fields struct {
		Name        string
		Purpose     string
		ScorerType  ScoringType
		Criticality int
		Evaluator   CardsEvaluator
	}
	type args struct {
		mf *dto.MetricFamily
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		// NOTE: We should add at least one (un-)happy path test per ScoringType here
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &CardsScoringProcessor{
				Name:        tt.fields.Name,
				Purpose:     tt.fields.Purpose,
				ScorerType:  tt.fields.ScorerType,
				Criticality: tt.fields.Criticality,
				Evaluator:   tt.fields.Evaluator,
			}
			got, err := sp.Score(tt.args.mf)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardsScoringProcessor.Score() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CardsScoringProcessor.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsLabelLengthEvaluator_Evaluate(t *testing.T) {
	mFam := &dto.MetricFamily{
		Metric: []*dto.Metric{
			{
				Label: []*dto.LabelPair{
					{
						Name:  proto.String("firstlabel"),
						Value: proto.String("firstlabel-value"),
					},
					{
						Name:  proto.String("secondlabel"),
						Value: proto.String("secondlabel-value"),
					},
				},
			},
		},
	}
	type fields struct {
		Type   ScoringType
		maxLen int
	}
	type args struct {
		f *dto.MetricFamily
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "length violates",
			fields: fields{
				maxLen: 2,
				Type:   LabelLengthScorer,
			},
			args: args{
				f: mFam,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "length does not violate",
			fields: fields{
				maxLen: 5,
				Type:   LabelLengthScorer,
			},
			args: args{
				f: mFam,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardsLabelLengthEvaluator{
				Type:   tt.fields.Type,
				maxLen: tt.fields.maxLen,
			}
			got, err := r.Evaluate(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardsLabelLengthEvaluator.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CardsLabelLengthEvaluator.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardsFamilyNameEvaluator_Evaluate(t *testing.T) {
	mFams := []*dto.MetricFamily{
		{
			Name: proto.String("test_family_1"),
			Metric: []*dto.Metric{},	
		},
		{
			Name: proto.String("test_family_2"),
			Metric: []*dto.Metric{},
		},
	}

	type fields struct {
		Type        ScoringType
		excludeList []string
	}
	type args struct {
		f *dto.MetricFamily
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "family name excluded",
			fields: fields{
				excludeList: []string{"test_family_1"},
			},
			args: args{
				f: mFams[0],
			},
			want: true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CardsFamilyNameEvaluator{
				Type:        tt.fields.Type,
				excludeList: tt.fields.excludeList,
			}
			got, err := r.Evaluate(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("CardsFamilyNameEvaluator.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CardsFamilyNameEvaluator.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
