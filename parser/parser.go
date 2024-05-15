package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetNetInfoFromInterx(ctx context.Context, ip string) (*Response, error) {

	ctxWithTO, c := context.WithTimeout(ctx, time.Second*10)
	defer c()
	log.Printf("Getting net_info from: %v", ip)
	url := fmt.Sprintf("http://%v:11000/api/net_info", ip)
	req, err := http.NewRequestWithContext(ctxWithTO, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON response
	var nodeInfo Response
	err = json.Unmarshal(b, &nodeInfo)
	if err != nil {
		return nil, err
	}
	return &nodeInfo, nil
}
