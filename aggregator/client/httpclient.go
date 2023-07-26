package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rtsoy/toll-calculator/types"
	"net/http"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) Client {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	//invReq := &types.GetInvoiceRequest{
	//	ObuID: int32(id),
	//}
	//
	//b, err := json.Marshal(invReq)
	//if err != nil {
	//	return nil, err
	//}

	targetURL := c.Endpoint + fmt.Sprintf("/invoice?obuID=%d", id)
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responded with non %d status code: %d", http.StatusOK, resp.StatusCode)
	}

	var inv types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}

	return &inv, nil
}

func (c *HTTPClient) Aggregate(ctx context.Context, request *types.AggregateRequest) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.Endpoint+"/aggregate", bytes.NewReader(b))
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
