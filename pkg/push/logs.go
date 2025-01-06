package push

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func fetchLogs[T any](p *Push, apiURL string, transform func(T) map[string]string) ([]map[string]string, error) {
	var logs []T
	var result []map[string]string
	httpClient := http.Client{Timeout: time.Second * 10}
	cursor := ""

	for {
		httpRequest, err := http.NewRequest(http.MethodGet, apiURL, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create HTTP request: %v", err)
		}

		query := url.Values{}
		if cursor != "" {
			query.Set("cursor", cursor)
		}
		httpRequest.URL.RawQuery = query.Encode()

		httpRequest.Header.Set("x-api-key", p.apiToken)
		httpRequest.Header.Set("accept", "application/json")

		response, err := httpClient.Do(httpRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			err := response.Body.Close()
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("unexpected response code: %d", response.StatusCode)
		}

		var apiResponse struct {
			Result []T `json:"result"`
		}
		if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
			err := response.Body.Close()
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("could not unmarshal response: %v", err)
		}
		err = response.Body.Close()
		if err != nil {
			return nil, err
		}

		if p.Logger != nil {
			p.Logger.Printf("Fetched %d items from API\n", len(apiResponse.Result))
		}

		logs = append(logs, apiResponse.Result...)

		cursor = response.Header.Get("x-next-cursor")
		if cursor == "" {
			if p.Logger != nil {
				p.Logger.Println("No more pages to fetch.")
			}
			break
		} else {
			if p.Logger != nil {
				p.Logger.Printf("Next cursor: %s\n", cursor)
			}
		}
	}

	for _, log := range logs {
		transformed := transform(log)
		if len(transformed) == 0 {
			return nil, fmt.Errorf("transform function returned an empty map for log: %+v", log)
		}
		result = append(result, transformed)
	}

	return result, nil
}
