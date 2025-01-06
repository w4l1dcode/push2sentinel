package push

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	appsURL = "https://api.pushsecurity.com/v1/apps"
)

type App struct {
	ID                string        `json:"id"`
	Type              string        `json:"type"`
	ApprovalStatus    interface{}   `json:"approvalStatus"`
	SensitivityLevel  interface{}   `json:"sensitivityLevel"`
	OwnerID           interface{}   `json:"ownerId"`
	Notes             string        `json:"notes"`
	Website           string        `json:"website"`
	Description       string        `json:"description"`
	FriendlyName      string        `json:"friendlyName"`
	Labels            []interface{} `json:"labels"`
	CreationTimestamp int           `json:"creationTimestamp"`
}

func (p *Push) GetApps(lookback uint32) ([]map[string]string, error) {
	currentTime := time.Now().Unix()
	lastUsedTimestampAfter := currentTime - int64(lookback)*3600
	url := fmt.Sprintf("%s?creationTimestampAfter=%d", appsURL, lastUsedTimestampAfter)

	return fetchLogs(p, url, func(app App) map[string]string {
		creationTime := time.Unix(int64(app.CreationTimestamp), 0).UTC().Format(time.RFC3339)

		newMap := []map[string]interface{}{
			{
				"type":             app.Type,
				"approvalStatus":   app.ApprovalStatus,
				"sensitivityLevel": app.SensitivityLevel,
				"ownerId":          app.OwnerID,
				"notes":            app.Notes,
				"website":          app.Website,
				"description":      app.Description,
				"friendlyName":     app.FriendlyName,
				"labels":           app.Labels,
			},
		}

		newJSON, err := json.Marshal(newMap)
		if err != nil {
			return nil
		}

		return map[string]string{
			"version":       "1",
			"id":            app.ID,
			"TimeGenerated": creationTime,
			"category":      "App",
			"new":           string(newJSON),
		}
	})
}
