package push

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	findingsURL = "https://api.pushsecurity.com/v1/findings"
)

type Finding struct {
	Id                  string   `json:"id"`
	Type                string   `json:"type"`
	State               string   `json:"state"`
	EmployeeId          string   `json:"employeeId"`
	PasswordId          string   `json:"passwordId"`
	AccountId           string   `json:"accountId"`
	AppType             string   `json:"appType"`
	AppId               string   `json:"appId"`
	WeakPasswordReasons []string `json:"weakPasswordReasons"`
	CreationTimestamp   int      `json:"creationTimestamp"`
}

func (p *Push) GetFindings(lookback uint32) ([]map[string]string, error) {
	currentTime := time.Now().Unix()
	lastUsedTimestampAfter := currentTime - int64(lookback)*3600
	url := fmt.Sprintf("%s?creationTimestampAfter=%d", findingsURL, lastUsedTimestampAfter)

	return fetchLogs(p, url, func(finding Finding) map[string]string {
		creationTime := time.Unix(int64(finding.CreationTimestamp), 0).UTC().Format(time.RFC3339)

		newMap := []map[string]interface{}{
			{
				"type":                finding.Type,
				"state":               finding.State,
				"employeeId":          finding.EmployeeId,
				"passwordId":          finding.PasswordId,
				"accountId":           finding.AccountId,
				"appType":             finding.AppType,
				"appId":               finding.AppId,
				"weakPasswordReasons": finding.WeakPasswordReasons,
			},
		}

		newJSON, err := json.Marshal(newMap)
		if err != nil {
			return nil
		}

		return map[string]string{
			"version":       "1",
			"id":            finding.Id,
			"TimeGenerated": creationTime,
			"category":      "Finding",
			"new":           string(newJSON),
		}
	})
}
