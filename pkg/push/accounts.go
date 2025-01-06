package push

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	accountsURL = "https://api.pushsecurity.com/v1/accounts"
)

type Account struct {
	ID           string `json:"id"`
	EmployeeID   string `json:"employeeId"`
	AppType      string `json:"appType"`
	AppID        string `json:"appId"`
	Email        string `json:"email"`
	PasswordID   string `json:"passwordId"`
	LoginMethods struct {
		PasswordLogin  bool   `json:"passwordLogin"`
		OidcLogin      string `json:"oidcLogin"`
		SamlLogin      string `json:"samlLogin"`
		OktaSwaLogin   bool   `json:"oktaSwaLogin"`
		VendorSsoLogin string `json:"vendorSsoLogin"`
		FedCmLogin     string `json:"fedCmLogin"`
	} `json:"loginMethods"`
	CreationTimestamp int `json:"creationTimestamp"`
	LastUsedTimestamp int `json:"lastUsedTimestamp"`
}

func (p *Push) GetAccounts(lookback uint32) ([]map[string]string, error) {
	currentTime := time.Now().Unix()
	lastUsedTimestampAfter := currentTime - int64(lookback)*3600
	url := fmt.Sprintf("%s?lastUsedTimestampAfter=%d", accountsURL, lastUsedTimestampAfter)

	return fetchLogs(p, url, func(account Account) map[string]string {
		creationTime := time.Unix(int64(account.CreationTimestamp), 0).UTC().Format(time.RFC3339)

		newMap := []map[string]interface{}{
			{
				"employeeId":        account.EmployeeID,
				"accountId":         account.ID,
				"appType":           account.AppType,
				"appId":             account.AppID,
				"email":             account.Email,
				"creationTimestamp": account.CreationTimestamp,
				"lastUsedTimestamp": account.LastUsedTimestamp,
				"passwordId":        account.PasswordID,
				"loginMethods": []map[string]interface{}{
					{"type": "passwordLogin", "enabled": account.LoginMethods.PasswordLogin},
					{"type": "oidcLogin", "enabled": account.LoginMethods.OidcLogin},
					{"type": "samlLogin", "enabled": account.LoginMethods.SamlLogin},
					{"type": "oktaSwaLogin", "enabled": account.LoginMethods.OktaSwaLogin},
					{"type": "vendorSsoLogin", "enabled": account.LoginMethods.VendorSsoLogin},
					{"type": "fedCmLogin", "enabled": account.LoginMethods.FedCmLogin},
				},
			},
		}

		newJSON, err := json.Marshal(newMap)
		if err != nil {
			return nil
		}

		return map[string]string{
			"version":       "1",
			"id":            account.ID,
			"TimeGenerated": creationTime,
			"category":      "Account",
			"new":           string(newJSON),
		}
	})
}
