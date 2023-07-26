package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rtsoy/toll-calculator/types"
	"net/http"
)

type HTTPClient struct {
	Endpoint string
}

func NewClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) AggregateInvoice(distance types.Distance) error {
	b, err := json.Marshal(distance)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with non %d status code: %d", http.StatusOK, resp.StatusCode)
	}

	return nil
}
