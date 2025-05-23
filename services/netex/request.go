package netex

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/plugimt/transat-backend/models"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

var mu sync.Mutex
var lastRequestTime time.Time

func CallStopMonitoringRequest(stops []string) (*models.SIRI, error) {
	mu.Lock()
	defer mu.Unlock()

	if time.Since(lastRequestTime) < 1*time.Second {
		time.Sleep(1 * time.Second)
	}

	content, err := GenerateStopMonitoringRequest(stops)
	if err != nil {
		return nil, err
	}

	url := url.URL{
		Scheme: APIScheme,
		Host:   APIHost,
		Path:   APIPath,
	}

	request := http.Request{
		Method: "POST",
		URL:    &url,
		Body:   io.NopCloser(bytes.NewBufferString(content)),
		Header: http.Header{
			"Content-Type": []string{"application/xml"},
			"datasetId":    []string{datasetId},
		},
	}

	resp, err := httpClient.Do(&request)
	if err != nil {
		return nil, err
	}

	lastRequestTime = time.Now()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var siri models.SIRI
	err = xml.Unmarshal(body, &siri)
	if err != nil {
		return nil, err
	}

	return &siri, nil
}
