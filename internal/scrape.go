package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	log "github.com/sirupsen/logrus"
)

// ReadTarget reads from an io.Reader and returns a corresponding bytes.Buffer. On error, an empty bytes.Buffer and the corresponding error are returned.
func ReadTarget(t io.Reader) (bytes.Buffer, error) {
	var rBuf bytes.Buffer

	if _, err := rBuf.ReadFrom(t); err != nil {
		return bytes.Buffer{}, err
	}

	return rBuf, nil
}

// ParseMetricFamilies parses a bytes.Buffer into a Prometheus MetricFamily. The MetricFamily (or nil MetricFamily and error) are returned.
func ParseMetricFamilies(b bytes.Buffer) (map[string]*dto.MetricFamily, error) {
	reader := bytes.NewReader(b.Bytes())

	var parser expfmt.TextParser

	mFam, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return map[string]*dto.MetricFamily{}, err
	}

	return mFam, nil

}

type CardsHttpScraper struct {
	Endpoint string // HTTP(s) endpoint representing a metrics exporter of the Prometheus text format type
	maxRetries int // Maximum retries this scraper should utilize
	b bytes.Buffer // Reusable bytes.Buffer containing the entire scrape response body
}

// NewCardsHttpScraper instantiates a new HTTP-based scraper for Cards scoring. Retries are currently ignored.
func NewCardsHttpScraper(endpoint string, retries int) (CardsHttpScraper) {
	var buf bytes.Buffer

	return CardsHttpScraper{
		Endpoint: endpoint,
		maxRetries: retries,
		b: buf,
	}
}

// ScrapeTarget performs a single HTTP scrape to the configured target, returning Prometheus MetricFamily and optional error
func (n *CardsHttpScraper) ScrapeTarget() (map[string]*dto.MetricFamily, error){
	hC := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	log.WithFields(log.Fields{
		"endpoint": n.Endpoint,
	}).Info("attempting scrape")
	
	res, err := hC.Get(n.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("scrape error: %w", err)
	}
	log.WithFields(log.Fields{
		"res_code": res.StatusCode,
		"res_content_len_bytes": res.ContentLength,
	}).Info("received response")

	parsedBody, err := ReadTarget(res.Body)
	if err != nil {
		return nil, err
	}

	return ParseMetricFamilies(parsedBody)
}
