package push

import (
	"encoding/json"
	"time"
)

const (
	employeesURL = "https://api.pushsecurity.com/v1/employees"
)

type Employee struct {
	ID                string   `json:"id"`
	Email             string   `json:"email"`
	FirstName         string   `json:"firstName"`
	LastName          string   `json:"lastName"`
	Department        string   `json:"department"`
	Location          string   `json:"location"`
	Licensed          bool     `json:"licensed"`
	ChatopsEnabled    bool     `json:"chatopsEnabled"`
	Groups            []string `json:"groups"`
	CreationTimestamp int      `json:"creationTimestamp"`
}

func (p *Push) GetEmployees(lookBackHours uint32) ([]map[string]string, error) {
	return fetchLogs(p, lookBackHours, "creationTimestampAfter", employeesURL, func(employee Employee) map[string]string {
		creationTime := time.Unix(int64(employee.CreationTimestamp), 0).UTC().Format(time.RFC3339)

		newMap := []map[string]interface{}{
			{
				"email":          employee.Email,
				"firstName":      employee.FirstName,
				"lastName":       employee.LastName,
				"department":     employee.Department,
				"location":       employee.Location,
				"licensed":       employee.Licensed,
				"chatOpsEnabled": employee.ChatopsEnabled,
				"groups":         employee.Groups,
			},
		}

		newJSON, err := json.Marshal(newMap)
		if err != nil {
			return nil
		}

		return map[string]string{
			"version":       "1",
			"id":            employee.ID,
			"TimeGenerated": creationTime,
			"category":      "Employee",
			"new":           string(newJSON),
		}
	})
}
