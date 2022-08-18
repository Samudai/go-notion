package notion

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type AuthResponse struct {
	AccessToken   string  `json:"access_token"`
	TokenType     string  `json:"token_type"`
	BotID         string  `json:"bot_id"`
	WorkspaceName *string `json:"workspace_name"`
	WorkspaceIcon *string `json:"workspace_icon"`
	WorkspaceID   string  `json:"workspace_id"`
	Owner         Owner   `json:"owner"`
}

type Owner struct {
	Type string `json:"type"`
	User User   `json:"user"`
}

type NotionError struct {
	Error string `json:"error"`
}

// AuthUser authenticates with notion and returns an user access token.
func AuthUser(code string) (*AuthResponse, error) {
	url := baseURL + "/oauth/token"

	clientID := os.Getenv("NOTION_CLIENT_ID")
	clientSecret := os.Getenv("NOTION_CLIENT_SECRET")
	redirectURI := os.Getenv("NOTION_REDIRECT_URI")

	payload := strings.NewReader(`{
    "grant_type": "authorization_code",
    "code": "` + code + `",
    "redirect_uri": "` + redirectURI + `"
}`)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	auth := base64.URLEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		var resp NotionError
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", resp.Error)
	}

	var resp AuthResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
