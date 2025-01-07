package push

import (
	"encoding/json"
	"time"
)

const (
	browsersURL = "https://api.pushsecurity.com/v1/browsers"
)

type Browser struct {
	ID                  string `json:"id"`
	EmployeeID          string `json:"employeeId"`
	Email               string `json:"email"`
	Version             string `json:"version"`
	TokenType           string `json:"tokenType"`
	IsActive            bool   `json:"isActive"`
	Browser             string `json:"Browser"`
	Os                  string `json:"os"`
	ExtensionVersion    string `json:"extensionVersion"`
	CreationTimestamp   int    `json:"creationTimestamp"`
	LastOnlineTimestamp int    `json:"lastOnlineTimestamp"`
}

func (p *Push) GetBrowsers(lookBackHours uint32) ([]map[string]string, error) {
	return fetchLogs(p, lookBackHours, "lastOnlineTimestampAfter", browsersURL, func(browser Browser) map[string]string {
		creationTime := time.Unix(int64(browser.CreationTimestamp), 0).UTC().Format(time.RFC3339)

		newMap := []map[string]interface{}{
			{
				"employeeID":          browser.EmployeeID,
				"email":               browser.Email,
				"tokenType":           browser.TokenType,
				"isActive":            browser.IsActive,
				"Browser":             browser.Browser,
				"os":                  browser.Os,
				"extensionVersion":    browser.ExtensionVersion,
				"LastOnlineTimestamp": browser.LastOnlineTimestamp,
			},
		}

		newJSON, err := json.Marshal(newMap)
		if err != nil {
			return nil
		}

		return map[string]string{
			"version":       "1",
			"id":            browser.ID,
			"TimeGenerated": creationTime,
			"category":      "Browser",
			"new":           string(newJSON),
		}
	})
}
