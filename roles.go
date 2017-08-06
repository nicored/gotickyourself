package tickspot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Role struct {
	SubscriptionID int     `json:"subscription_id,omitempty" yaml:"subscription_id"`
	Company        string  `json:"company,omitempty"  yaml:"company"`
	APIToken       string  `json:"api_token,omitempty"  yaml:"api_token"`
	IsDefault      bool    `json:"-"  yaml:"is_default"`
	WeeklyHours    float64 `json:"-" yaml:"weekly_hours"`
}

// curl -u "email_address:password" \
// -H 'User-Agent: MyCoolApp (me@example.com)' \
// https://www.tickspot.com/api/v2/roles.json
func (t *Tick) GetRoles() ([]*Role, error) {
	url := fmt.Sprintf("%s/api/v2/roles.json", t.BaseUrl)

	req, err := http.NewRequest(MethodGET, url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(t.Client.Username, t.Client.Password)
	req.Header.Add("User-Agent", fmt.Sprintf("GoTickYourself (%s)", t.Client.Username))

	resp, err := t.sendRequest(req)

	roles := []*Role{}
	err = json.Unmarshal(resp.Body, &roles)
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, errors.New(string(resp.Body))
	}

	return roles, nil
}
